package model

type Column struct {
	CommonTableCols
	ColumnID string `json:"columnID" gorm:"type:char(36); primary_key; not null"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
}
