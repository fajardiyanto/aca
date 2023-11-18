package models

type Department struct {
	ID          uint      	`gorm:"primaryKey" json:"id"`
	Name 		string 		`gorm:"column:name" json:"name"`
}

func (*Department) TableName() string {
	return "department"
}