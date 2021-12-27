package model

import "time"

// CommonTableCols ...
type CommonTableCols struct {
	Seq       uint `gorm:"type:int unsigned auto_increment;unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
