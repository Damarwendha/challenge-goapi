package models

import "time"

type Laundry struct {
	Id             int       `json:"id"`
	Customer_Name  string    `json:"customer_name"`
	Customer_Phone string    `json:"customer_phone"`
	Entry_Date     time.Time `json:"entry_date"`
	Finish_Date    time.Time `json:"finish_date"`
	Service_Name   string    `json:"service_name"`
	Quantity       int       `json:"quantity"`
	Total_Price    int       `json:"total_price"`
}
