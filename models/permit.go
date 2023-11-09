package models

import "time"

type Permit struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"column:name" json:"name"`
	PermitID    string    `gorm:"column:permit_id" json:"permit_id"`
	Region      string    `gorm:"column:region" json:"region"`
	NIK         string    `gorm:"column:nik" json:"nik"`
	Company     string    `gorm:"column:company" json:"company"`
	Departement string    `gorm:"column:departement" json:"departement"`
	Position    string    `gorm:"column:position" json:"position"`
	Image       string    `gorm:"column:image" json:"image"`
	Valid       string    `gorm:"column:valid" json:"valid"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (*Permit) TableName() string {
	return "permit"
}
