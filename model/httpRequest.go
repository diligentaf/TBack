package model

import "time"

type GetCmcHTTPRequest struct {
	CmcID       string    `json:"cmcID" gorm:"type:char(36); unique; not null"`
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}

type GetBySymbolHTTPRequest struct {
	Symbol   string `json:"symbol"`
	FromDate string `json:"from_date" form:"from_date" query:"from_date"`
	ToDate   string `json:"to_date" form:"to_date" query:"to_date"`
}
