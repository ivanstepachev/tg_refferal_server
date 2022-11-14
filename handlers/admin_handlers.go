package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/ivanstepachev/tg_refferal/store/models"
	"github.com/gorilla/mux"
)

func (h *handlers) UsersListHandler(w http.ResponseWriter, r *http.Request) {
	var Users []models.User
	cond := r.URL.Query().Get("cond")

	if cond == "ref" {
		result := h.DB.Order("id").Not("beneficiary_id = ?", "landing").Not("beneficiary_id = ?", "directly").Find(&Users)
		if result.Error != nil {
			log.Println(result.Error)
		}
	} else if cond == "paid" {
		result := h.DB.Order("id").Where("is_paid = ?", true).Find(&Users)
		if result.Error != nil {
			log.Println(result.Error)
		}
	} else if cond == "notpaid" {
		result := h.DB.Order("id").Where("is_paid = ?", false).Find(&Users)
		if result.Error != nil {
			log.Println(result.Error)
		}
	} else {
		result := h.DB.Order("id").Find(&Users)
		if result.Error != nil {
			log.Println(result.Error)
		}
	}
	
	tmpl, err := template.ParseFiles("static/users.html")
	if err != nil {
		// Write error to template
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	if err := tmpl.Execute(w, Users); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (h *handlers) UserAddHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := template.ParseFiles("static/add_user.html")
		if err != nil {
			// Write error to template
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
		}
		telegramId := r.PostForm.Get("telegram_id")
		isPaidString := r.PostForm.Get("is_paid")
		refCode := r.PostForm.Get("beneficiary_id")
		balance, err := strconv.Atoi(r.PostForm.Get("balance"))
		if err != nil {
			log.Println(err.Error())
		}
		isPaid := false
		if isPaidString == "on" {
			isPaid = true
		}
		args := newUserArgs{
			telegramId: telegramId,
			isPaid: isPaid,
			Balance: balance,
			refCode: refCode,
		}
		_, err = h.createNewUser(&args)
		if err != nil {
			log.Println(err.Error())
		}
		http.Redirect(w, r, "/admin/users", http.StatusMovedPermanently)
	}
}

func (h *handlers) UserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err.Error())
	}
	var user = models.User{Id: userId}
	tx := h.DB.First(&user)
	if tx.Error != nil {
		log.Println(tx.Error.Error())
	}
	if tx.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		
		tmpl, err := template.ParseFiles("static/user.html")
		if err != nil {
			// Write error to template
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())
		}
		telegramId := r.PostForm.Get("telegram_id")
		isPaidString := r.PostForm.Get("is_paid")
		refCode := r.PostForm.Get("beneficiary_id")
		balance, err := strconv.Atoi(r.PostForm.Get("balance"))
		if err != nil {
			log.Println(err.Error())
		}

		isPaid := false
		if isPaidString == "on" {
			isPaid = true
		}

		tx := h.DB.Model(&user).Updates(map[string]interface{}{
			"telegram_id": telegramId, 
			"telegram_link": fmt.Sprintf("tg://user?id=%v", telegramId), 
			"is_paid": isPaid,
			"balance": balance,
			"beneficiary_id": refCode,

		})
		if tx.Error != nil {
			log.Println(tx.Error.Error())
	
		}
		url := fmt.Sprintf("/admin/users/%v", userId)
		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}

func (h *handlers) UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
	}
	id, err := strconv.Atoi(r.PostForm.Get("id"))
	if err != nil {
		log.Println(err.Error())
	}
	h.DB.Delete(&models.User{}, id)
	http.Redirect(w, r, "/admin/users", http.StatusMovedPermanently)
}

func (h *handlers) TransactionsListHandler(w http.ResponseWriter, r *http.Request) {
	var Transactions []models.Transaction
	result := h.DB.Order("id").Find(&Transactions)
	if result.Error != nil {
		log.Println(result.Error)
	}
	
	tmpl, err := template.ParseFiles("static/transactions.html")
	if err != nil {
		// Write error to template
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	if err := tmpl.Execute(w, Transactions); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
