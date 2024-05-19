package models

type Meal struct {
	ID      int     `json:"id,omitempty"`
	ResName string  `json:"res_name"`
	Name    string  `json:"name"`
	Price   float64 `json:"price"`
	Amount  int     `json:"amount"`
	ResID   int     `json:"res_id,omitempty"`
}
