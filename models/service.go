package models

type Service struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Unit_type_id int    `json:"unit_type_id"`
}
