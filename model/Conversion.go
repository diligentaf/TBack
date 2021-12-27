package model

type Conversion struct {
	From  string  `json:"from"`
	To    string  `json:"to"`
	Price float64 `json:"price"`
}
