package models

import "time"

type User struct {
	Id int `json:"id" gorm:"primaryKey"`
	TelegramId string `json:"telegram_id" gorm:"index"`
	TelegramLink string `json:"telegram_link"`
	IsPaid bool `json:"is_paid"`
	Balance int `json:"balance"`
	// Beneficiary is TelegramId who recieved bonus for recommendation
	BeneficiaryId string `json:"beneficiary_id"` 
	CreatedAt time.Time `json:"created_at"`
}
