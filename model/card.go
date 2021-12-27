package model

type Card struct {
	CommonTableCols
	CardID      string `json:"cardID" gorm:"type:char(36); unique; not null"`
	ColumnID    string `json:"columnID" gorm:"type:char(36); foreign_key; not null"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Order       int    `json:"order" gorm:"not null"`
	Status      *bool  `json:"status"`
}
