package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"prescription/db"
	"prescription/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	db.ConnectPostDb()
	app := fiber.New()
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Post("/medicine", handlers.Medicine)
	app.Get("/medicines", handlers.GetallMeds)
	app.Post("/updatemeds", handlers.DispenseStock)
	return app
}

func TestRegisterAndLogin(t *testing.T) {
	app := setupApp()

	// Register
	req := httptest.NewRequest("POST", "/register",
		strings.NewReader(`{"name":"TestAdmin","email":"admin@test.com","password":"pass123","role":"Admin"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)

	// Login
	req2 := httptest.NewRequest("POST", "/login",
		strings.NewReader(`{"email":"admin@test.com","password":"pass123"}`))
	req2.Header.Set("Content-Type", "application/json")
	resp2, _ := app.Test(req2, -1)
	assert.Equal(t, 200, resp2.StatusCode)

	// Optional: decode token
	var result map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&result)
	assert.NotEmpty(t, result["token"])
}

// Add Medicine
func TestAddMedicine(t *testing.T) {
	app := setupApp()

	req := httptest.NewRequest("POST", "/medicine",
		strings.NewReader(`{"medicine_name":"Paracetamol","dosage_form":"Tablet","stock_quantity":10}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)

	var med map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&med)
	assert.Equal(t, "Paracetamol", med["medicine_name"])
}

// Dispese Medicine
func TestDispenseMedicine(t *testing.T) {
	app := setupApp()

	ctx := context.Background()
	var currentStock int
	err := db.Postdb.QueryRow(ctx, `SELECT stock_quantity FROM medicine WHERE medicine_name='Paracetamol'`).Scan(&currentStock)
	assert.NoError(t, err)

	dispenseQty := 5
	req := httptest.NewRequest("POST", "/updatemeds",
		strings.NewReader(`{"medicine_name":"Paracetamol","dosage_form":"Tablet","stock_quantity":`+fmt.Sprint(dispenseQty)+`}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)
	assert.Equal(t, 200, resp.StatusCode)

	var newStock int
	err = db.Postdb.QueryRow(ctx, `SELECT stock_quantity FROM medicine WHERE medicine_name='Paracetamol'`).Scan(&newStock)
	assert.NoError(t, err)

	assert.Equal(t, currentStock-dispenseQty, newStock, "Stock after dispense mismatch")
}

// Concurrency test for Dispense
func TestConcurrentDispense(t *testing.T) {
	app := setupApp()
	ctx := context.Background()

	_, _ = db.Postdb.Exec(ctx, `UPDATE medicine SET stock_quantity=10 WHERE medicine_name='Paracetamol'`)

	var wg sync.WaitGroup
	errors := make(chan error, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("POST", "/updatemeds",
				strings.NewReader(`{"medicine_name":"Paracetamol","dosage_form":"Tablet","stock_quantity":7}`))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				errors <- err
				return
			}
			if resp.StatusCode != 200 {
				errors <- assert.AnError
			}
		}()
	}

	wg.Wait()
	close(errors)

	failCount := 0
	for e := range errors {
		if e != nil {
			failCount++
		}
	}

	assert.Equal(t, 1, failCount)
}
