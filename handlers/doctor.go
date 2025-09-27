package handlers

import (
	"context"
	"prescription/db"
	"prescription/models"

	"github.com/gofiber/fiber/v2"
)

// MakePresc godoc
// @Summary Create prescription
// @Description Doctor creates a prescription for a patient
// @Tags Prescription
// @Accept json
// @Produce json
// @Param prescription body models.Prescription true "Prescription data"
// @Success 200 {object} models.Prescription
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /presc [post]
func MakePresc(fi *fiber.Ctx) error {
	var prescription models.Prescription

	err := fi.BodyParser(&prescription)
	if err != nil {
		return err
	}

	b := `INSERT INTO prescriptions (patient_name, medicine_name, dosage_form, quantity) VALUES ($1,$2,$3,$4) RETURNING id`
	err = db.Postdb.QueryRow(context.Background(), b, prescription.Patient_Name, prescription.Medicine_Name, prescription.Dosage_Form, prescription.Quantity).Scan(&prescription.Id)
	if err != nil {
		return err
	}

	return fi.JSON(prescription)
}
