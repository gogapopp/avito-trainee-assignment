package models

import "time"

type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city" validate:"required,oneof=Москва Санкт-Петербург Казань"`
}

type PVZWithReceptions struct {
	PVZ        PVZ                 `json:"pvz"`
	Receptions []ReceptionProducts `json:"receptions"`
}
