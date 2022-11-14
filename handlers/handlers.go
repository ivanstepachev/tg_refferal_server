package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ivanstepachev/tg_refferal/store/models"
	"github.com/jinzhu/gorm"
	"github.com/tidwall/gjson"
)

// Path to Telegram_ID and Message keys of TG api JSON
const tgIdKey = "telegram_id"
const tgMessageKey = "text"

const paymentStatusKey = "status"
const paymentIdKey = "payment_id"

// Ref code to define that user from landing page
const landingRef = "landing"

type handlers struct {
	DB *gorm.DB
}

type newUserArgs struct {
	telegramId string
	isPaid     bool
	Balance int
	refCode    string
}

// createNewUser write new user to DB
func (h *handlers) createNewUser(u *newUserArgs) (models.User, error) {
	user := models.User{
		TelegramId:    u.telegramId,
		TelegramLink:  fmt.Sprintf("tg://user?id=%v", u.telegramId),
		IsPaid:        u.isPaid,
		Balance:       u.Balance,
		BeneficiaryId: u.refCode,
	}
	if result := h.DB.Create(&user); result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// Select User by telegram id
func (h *handlers) selectUser(telegramId string) (models.User, error) {
	var user models.User
	err := h.DB.Where("telegram_id = ?", telegramId).First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// addBonus add bonus to benificiary
func (h *handlers) addBonus(bonus int, user *models.User) (models.User, error) {
	var beneficiary models.User
	refCode := user.BeneficiaryId
	err := h.DB.Where("telegram_id = ?", refCode).First(&beneficiary).Error
	if err != nil {
		return models.User{}, err
	}
	beneficiary.Balance += bonus
	err = h.DB.Save(&beneficiary).Error
	if err != nil {
		return models.User{}, err
	}
	// TODO send message about bonus
	return beneficiary, nil
}

// New initializing new db storage
func New(DB *gorm.DB) handlers {
	return handlers{
		DB: DB,
	}
}

// logRequest writes request JSON logs to txt file
func logRequests(jsonString string, path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	currentTime := time.Now().Format(time.Stamp)
	toLog := fmt.Sprintf("%v - %v\n", currentTime, jsonString)
	_, err = f.WriteString(toLog)
	if err != nil {
		return err
	}
	return nil
}

// TelegramApiHandler handles requests from Telegram bot
func (h *handlers) TelegramApiHandler(w http.ResponseWriter, r *http.Request) {
	req, _ := ioutil.ReadAll(r.Body)
	// Recieve JSON string from byte array parsed from body request
	jsonString := string(req[:])
	err := logRequests(jsonString, "telegram_requests.txt")
	if err != nil {
		log.Println(err.Error())
	}
	message := gjson.Get(jsonString, tgMessageKey).Str
	telegramId := gjson.Get(jsonString, tgIdKey).Str
	// Handle only correct requests, check if this key:value does not exist
	if len(telegramId) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if not exist than create new user
	user, err := h.selectUser(telegramId)
	if err != nil {
		log.Println(err.Error())
	}
	if (user == models.User{}) {
		_, err := startMessage(h, telegramId, message)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// PaymentApiHandler handles information about success payments from users by telegram_id
func (h *handlers) PaymentApiHandler(w http.ResponseWriter, r *http.Request) {
	// var user models.User
	req, _ := ioutil.ReadAll(r.Body)
	// Recieve JSON string from byte array parsed from body request
	status := gjson.Get(string(req[:]), paymentStatusKey).Str
	telegramId := gjson.Get(string(req[:]), paymentIdKey).Str
	err := logRequests(string(req[:]), "payment_requests.txt")
	if err != nil {
		log.Println(err.Error())
	}
	// Stop execution if payment status is not success
	if status != "success" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	user, err := h.selectUser(telegramId)
	if err != nil {
		log.Println(err.Error())
	}
	if (user == models.User{}) {
		w.WriteHeader(http.StatusAccepted)
		return
	}
	user.IsPaid = true
	err = h.DB.Save(&user).Error
	if err != nil {
		log.Println(err.Error())
	}
	// Check Beneficiary of the current payment and add balance
	if user.BeneficiaryId != "landing" && user.BeneficiaryId != "directly" {
		beneficiary, err := h.addBonus(250, &user)
		if err != nil {
			log.Println(err.Error())
		}
		// Check Grand Beneficiary of the current payment and add balance
		// addBonus can return an empty value of User, for this reason we compare with empty struct
		if beneficiary.BeneficiaryId != "landing" && beneficiary.BeneficiaryId != "directly" && (beneficiary != models.User{}) {
			_, err = h.addBonus(50, &beneficiary)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
	w.WriteHeader(http.StatusOK)
}

// {'update_id': 117378454, 'message': {'message_id': 49, 'from': {'id': 5104671785, 'is_bot': False, 'first_name': 'Ivan', 'username': 'ivan_stepachev', 'language_code': 'ru'}, 'chat': {'id': 5104671785, 'first_name': 'Ivan', 'username': 'ivan_stepachev', 'type': 'private'}, 'date': 1661203162, 'text': '/start 896205315', 'entities': [{'offset': 0, 'length': 6, 'type': 'bot_command'}]}}
