# PrescriptionDispensingSystem

Backend Framework: GoLang
Database: PostgreSQL (transactional DB required)
Optional: Redis (for locks or queues)
Documentation: Swagger/OpenAPI
Testing: Unit tests, including concurrent access tests
Version Control: Git + GitHub

Description

A backend system to manage prescriptions and medicine stock in hospitals or pharmacies. The system supports role-based access control with Admin, Doctor, and Pharmacist roles, ensuring proper authorization for each action.

Features

User Management

Roles: Doctor, Pharmacist, Admin

Authentication: JWT tokens

Authorization:

Doctors can issue prescriptions

Pharmacists can dispense medicine and update stock

Admin can manage medicine catalog

Medicine Catalog & Stock

Admin can add new medicines with details: name, dosage form, stock quantity

Pharmacists can view stock levels

Prescription Handling

Doctors can issue a prescription for a patient (patient name, medicine, quantity)

When a prescription is dispensed, the stock decrements automatically

Atomic stock updates to prevent race conditions

Concurrency & Transactions

Multiple concurrent prescriptions handled safely

Stock updates are atomic to avoid negative stock

Race conditions are prevented when multiple pharmacists dispense the same medicine
