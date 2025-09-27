package models

type UserDetails struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Medicine struct {
	Medicine_Name  string `json:"medicine_name"`
	Dosage_Form    string `json:"dosage_form"`
	Stock_Quantity int    `json:"stock_quantity"`
}

type Prescription struct {
	Id            int    `json:"id"`
	Patient_Name  string `json:"patient_name"`
	Medicine_Name string `json:"medicine_name"`
	Dosage_Form   string `json:"dosage_form"`
	Quantity      int    `json:"quantity"`
}

type LoginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}