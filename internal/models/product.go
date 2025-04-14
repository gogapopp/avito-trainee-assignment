package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type" validate:"required,oneof=электроника одежда обувь"`
	ReceptionID string    `json:"receptionId"`
}

type NewProductRequest struct {
	Type  string `json:"type" validate:"required,oneof=электроника одежда обувь"`
	PVZID string `json:"pvzId" validate:"required,uuid"`
}
