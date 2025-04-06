package types

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	UserID   string `gorm:"column:user_id;" json:"user_id"`
	Username string `gorm:"column:username;" json:"username"`
	Password string `gorm:"column:password;" json:"password"`
	Email    string `gorm:"column:email;" json:"email"`
	Phone    string `gorm:"column:phone;" json:"phone"`
	Role     string `gorm:"column:role;" json:"role"`
}
