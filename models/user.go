package models

import (
	"gorm.io/gorm"
)

// available balance = balance - pending bets amounts
type User struct {
	gorm.Model
	Login    string  `gorm:"index:idx_login,unique" json:"login"` // we'll use email
	Password string  `json:"password"`                            // we store it encrypted
	Balance  float64 `json:"balance"`
	IsAI     bool    `json:"is_ai"`
}
