package models

import "gorm.io/gorm"

// @Description Bet represents a betting record
type Bet struct {
	gorm.Model
	UserID            uint        `json:"user_id"`
	Amount            float64     `json:"amount"`
	ItemValue         uint        `json:"item_value"`
	AcceptorItemValue *uint       `json:"acceptor_item_value"`
	Result            *ResultType `json:"result"` // it's nil on status Pending
	Status            BetStatus   `json:"status"`
	AcceptorID        uint        `json:"acceptor_id"` // it's another user id
}

type BetStatus uint

const (
	Pending  BetStatus = iota
	Declined           // can be declined by acceptor
	Accepted           // = played
)
