package models

type Position struct {
	ID          uint      	`gorm:"primaryKey" json:"id"`
	Name 		string 		`gorm:"column:name" json:"name"`
}

func (*Position) TableName() string {
	return "position"
}