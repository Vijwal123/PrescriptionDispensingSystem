package handlers

import (
	"context"
	"prescription/db"
	"prescription/models"
	"github.com/jackc/pgx/v5" 
	"github.com/gofiber/fiber/v2"
)

// GetallMeds godoc
// @Summary Get all medicines
// @Description Pharmacist views stock levels
// @Tags Medicine
// @Produce json
// @Success 200 {array} models.Medicine
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /medicines [get]
func GetallMeds(fi *fiber.Ctx) error {
	res, err := db.Postdb.Query(context.Background(), `SELECT medicine_name, dosage_form, stock_quantity FROM medicine`)
	if err != nil {
		return err
	}
	defer res.Close()

	var medicines []models.Medicine
	for res.Next() {
		var a models.Medicine
		res.Scan(&a.Medicine_Name, &a.Dosage_Form, &a.Stock_Quantity)

		medicines = append(medicines, a)
	}

	return fi.JSON(medicines)
}


// DispenseStock godoc
// @Summary Dispense medicine stock
// @Description Pharmacist dispenses medicine and updates stock
// @Tags Medicine
// @Accept json
// @Produce json
// @Param dispense body models.Medicine true "Dispense data"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /updatemeds [post]	
func DispenseStock(fi *fiber.Ctx) error {
	var dis models.Medicine
	if err := fi.BodyParser(&dis); err != nil {
		return fi.Status(400).JSON("invalid input")
	}

	query := `
		UPDATE medicine
		SET stock_quantity = stock_quantity - $1
		WHERE medicine_name = $2 AND dosage_form = $3 AND stock_quantity >= $1
		RETURNING stock_quantity
	`

	var newStock int
	err := db.Postdb.QueryRow(context.Background(), query,
		dis.Stock_Quantity, dis.Medicine_Name, dis.Dosage_Form).Scan(&newStock)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fi.Status(400).JSON("Medicine not found or insufficient stock")
		}
		return fi.Status(500).JSON(err.Error())
	}

	dis.Stock_Quantity = newStock
	return fi.JSON(dis)
}
