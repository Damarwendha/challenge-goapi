package models

import "time"

type TrxLaundry struct {
	Id                int       `json:"id"`
	Customer_Id       int       `json:"customer_id"`
	Entry_Date        time.Time `json:"entry_date"`
	Finish_Date       time.Time `json:"finish_date"`
	Laundry_Detail_Id int       `json:"laundry_detail_id"`
}
