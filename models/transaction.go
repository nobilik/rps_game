package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID uint            `json:"user_id"`
	Amount float64         `json:"amount"`
	Type   TransactionType `json:"type"`
	BetID  *uint           `json:"bet_id"` // for win & loss
}

type TransactionType uint

const (
	TopUp TransactionType = iota
	Withdrawal
	Win
	Lose
)
