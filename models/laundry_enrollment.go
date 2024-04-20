package models

type LaundryEnrollment struct {
	Id          int `json:"id"`
	Customer_Id int `json:"customer_id"`
	Service_Id  int `json:"service_id"`
	Quantity    int `json:"quantity"`
}
