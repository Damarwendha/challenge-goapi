package models

type LaundryDetail struct {
	Id          int `json:"id"`
	Service_Id  int `json:"service_id"`
	Quantity    int `json:"quantity"`
	Total_Price int `json:"total_price"`
}
