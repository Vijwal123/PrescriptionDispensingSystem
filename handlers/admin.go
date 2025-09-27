package handlers

import (
	"context"
	"prescription/db"
	"prescription/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// Medicine godoc
// @Summary Add a new medicine
// @Description Admin adds medicine to stock
// @Tags Medicine
// @Accept json
// @Produce json
// @Param medicine body models.Medicine true "Medicine data"
// @Success 200 {object} models.Medicine
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /medicine [post]
func Medicine(fi *fiber.Ctx) error {
	var medicine models.Medicine
	if err := fi.BodyParser(&medicine); err != nil {
		return fi.JSON("Invalid Input")
	}

	ctx := context.Background()
	tx, err := db.Postdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fi.Status(500).JSON(err.Error())
	}
	defer tx.Rollback(ctx)

	var currentStock int
	err = tx.QueryRow(ctx, `
        SELECT stock_quantity 
        FROM medicine 
        WHERE medicine_name=$1 AND dosage_form=$2
        FOR UPDATE
    `, medicine.Medicine_Name, medicine.Dosage_Form).Scan(&currentStock)

	if err != nil {
		if err.Error() == "no rows in result set" {
	
			err = tx.QueryRow(ctx, `
                INSERT INTO medicine (medicine_name, dosage_form, stock_quantity)
                VALUES ($1, $2, $3)
                RETURNING stock_quantity
            `, medicine.Medicine_Name, medicine.Dosage_Form, medicine.Stock_Quantity).Scan(&medicine.Stock_Quantity)
			if err != nil {
				return fi.Status(500).JSON(err.Error())
			}
		} else {
			return fi.Status(500).JSON(err.Error())
		}
	} else {

		newStock := currentStock + medicine.Stock_Quantity
		err = tx.QueryRow(ctx, `
            UPDATE medicine 
            SET stock_quantity=$1
            WHERE medicine_name=$2 AND dosage_form=$3
            RETURNING stock_quantity
        `, newStock, medicine.Medicine_Name, medicine.Dosage_Form).Scan(&medicine.Stock_Quantity)
		if err != nil {
			return fi.Status(500).JSON(err.Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fi.Status(500).JSON(err.Error())
	}

	return fi.JSON(medicine)
}

// RemoveMedicine godoc
// @Summary Delete a medicine
// @Description Admin removes medicine from stock
// @Tags Medicine
// @Produce json
// @Param medicine_name path string true "Medicine name"
// @Success 200 {string} string "Medicine Deleted"
// @Failure 404 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /medicine/{medicine_name} [delete]
func RemoveMedicine(fi *fiber.Ctx) error {

	medname := fi.Params("medicine_name")
	res, err := db.Postdb.Exec(context.Background(), `DELETE FROM medicine WHERE medicine_name=$1`, medname)
	if err != nil {
		return fi.Status(404).JSON("Failed to delete medicine")
	}

	if res.RowsAffected() == 0 {
		return fi.Status(404).JSON("medicine not found")
	}

	return fi.JSON("Medicine Deleted")
}
