package models

import (
	"time"
)

type Transaction struct {
	Id int `json:"id" gorm:"primaryKey"`
	Amount int `json:"amount"`
	Card string `json:"card"`
	Done bool `json:"done" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	ChangedAt time.Time `json:"changed_at"`
	TelegramId string `json:"telegram_id"`
}