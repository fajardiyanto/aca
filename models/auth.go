package models

import "time"

type Auth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	Name      string    `gorm:"column:name" json:"name"`
	Email 	  string 	`gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	Role      string    `gorm:"column:role" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
}

type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type RegisterModel struct {
	Email    string `json:"email"`
	Name	 string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type TokenModel struct {
	Authorized bool   `json:"authorized"`
	Exp        string `json:"exp"`
	Role       string `json:"role"`
	UserID     string `json:"user_id"`
}

func (*Auth) TableName() string {
	return "users"
}
