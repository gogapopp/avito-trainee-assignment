package models

import "time"

type Reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId" validate:"required,uuid"`
	Status   string    `json:"status"`
}

type ReceptionProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}

type NewReceptionRequest struct {
	PVZID string `json:"pvzId" validate:"required,uuid"`
}
