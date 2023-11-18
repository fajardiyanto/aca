package models

type Simper struct {
	ID          uint      	`gorm:"primaryKey" json:"id"`
	Valid       string    	`gorm:"column:valid" json:"valid"`
	Type		string 		`gorm:"column:type" json:"type"`
	BloodType	string		`gorm:"column:blood_type" json:"blood_type"`
	Violation	string 		`gorm:"column:violation" json:"violation"`
	Vehicle 	string 		`gorm:"column:vehicle" json:"vehicle"`
}

type Vehicle struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

func (*Simper) TableName() string {
	return "simper"
}