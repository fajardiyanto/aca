package models

import "time"

type Auth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Password  string    `gorm:"column:password" json:"password"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
}

type LoginModel struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (*Auth) TableName() string {
	return "users"
}
