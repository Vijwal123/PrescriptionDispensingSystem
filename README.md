# PrescriptionDispensingSystem

A backend service built with GoLang and PostgreSQL for managing users, medicines, and prescriptions in a hospital or pharmacy environment. The API supports role-based access control with JWT authentication and ensures safe stock management, including concurrent updates and atomic operations.

# Key Features:
User Management: Roles – Doctor, Pharmacist, Admin.
Admin can manage the medicine catalog.
Doctors can issue prescriptions.
Pharmacists can dispense medicine and update stock.

# Medicine Catalog & Stock Management:
Add, view, and update medicines.
Stock updates are atomic to prevent race conditions.

# Prescription Handling:
Doctors issue prescriptions with patient name, medicine, and quantity.
Dispensing prescriptions decrements stock automatically.

# Security: JWT-based authentication and authorization for all endpoints.

# Testing: Includes unit and concurrency tests for stock operations.

# Documentation: Swagger/OpenAPI documentation included.

# Tech Stack:
Backend: GoLang + Fiber
Database: PostgreSQL
Authentication: JWT tokens
Testing: Unit tests + concurrency tests
Version Control: Git + GitHub
