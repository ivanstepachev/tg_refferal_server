package handlers

import (
	"github.com/ivanstepachev/tg_refferal/store/models"
)

// startMessage handle /start message with and without ref code
func startMessage(h *handlers, telegramId string, message string) (models.User, error) {
	// Check if first message contains ref code /start ore /start 111111
	if message[:6] == "/start" && len(message) > 8 {
		if message[7:] == landingRef {
			user, err := h.createNewUser(&newUserArgs{
				telegramId: telegramId,
				isPaid:     true,
				Balance: 0,
				refCode:    "landing",
			})
			if err != nil {
				return models.User{}, err
			}
			return user, nil
			// TODO send answer
		} else {
			// refCode is Benecficiary Telegram ID, not paid yet
			refCode := message[7:]
			user, err := h.createNewUser(&newUserArgs{
				telegramId: telegramId,
				isPaid:     false,
				Balance: 0,
				refCode:    refCode,
			})
			if err != nil {
				return models.User{}, err
			}
			return user, nil
			// TODO send answer
		}
		// /start without ref code, higher price than with ref code
	} else {
		user, err := h.createNewUser(&newUserArgs{
			telegramId: telegramId,
			isPaid:     false,
			Balance: 0,
			refCode:    "directly",
		})
		if err != nil {
			return models.User{}, err
		}
		return user, nil
		// TODO send answer
	}
}
