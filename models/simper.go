package models

type Simper struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PermitID  string `gorm:"permit_id" json:"permit_id"`
	Valid     string `gorm:"column:valid" json:"simper_valid"`
	Simpol 	  string `gorm:"column:simpol" json:"simpol"`
	NoSimpol  string `gorm:"column:no_simpol" json:"no_simpol"`
	Type      string `gorm:"column:type" json:"simper_type"`
	BloodType string `gorm:"column:blood_type" json:"blood_type"`
	Vehicle   string `gorm:"column:vehicle" json:"vehicle"`
}

type Vehicle struct {
	Number 	string `gorm:"number" json:"number"`
	Type 	string `gorm:"type_vehicle" json:"type_vehicle"`
	Name 	string `gorm:"name_vehicle" json:"name_vehicle"`
}

type SimperResponse struct {
	ID        uint    `json:"id"`
	PermitID  string  `json:"permit_id"`
	Valid     string  `json:"simper_valid"`
	Simpol 	  string  `json:"simpol"`
	NoSimpol  string  `json:"no_simpol"`
	Type      string  `json:"simper_type"`
	BloodType string  `json:"blood_type"`
	Vehicle   []Vehicle `json:"vehicle"`
}

func (*Simper) TableName() string {
	return "simper"
}
