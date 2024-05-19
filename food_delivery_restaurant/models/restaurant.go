package models

type Restaurant struct {
	ID      int    `json:"id"`
	Name    string `json:"name" validate:"required"`
	OwnerID int    `json:"owner_id" validate:"required"`
}
