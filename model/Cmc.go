package model

import "time"

type Cmc struct {
	CommonTableCols
	CmcID       string    `json:"cmcID" gorm:"type:char(36); unique; not null"`
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}
