package main

import (
	"log"
	"net/http"


	"github.com/ivanstepachev/tg_refferal/handlers"
	"github.com/ivanstepachev/tg_refferal/store/db"
	"github.com/gorilla/mux"
)

func main() {
	DB := db.Init()
	
	
	h := handlers.New(DB)

	// http Serve Mux for admin dashboard
	router := mux.NewRouter()

	router.HandleFunc("/refferal", h.TelegramApiHandler).Methods("POST")
	router.HandleFunc("/payment", h.PaymentApiHandler).Methods("POST")
	
	router.HandleFunc("/admin/users", h.UsersListHandler).Methods("GET")
	router.HandleFunc("/admin/users/delete", h.UserDeleteHandler).Methods("POST")

	router.HandleFunc("/admin/users/add", h.UserAddHandler)
	router.HandleFunc("/admin/users/{id:[0-9]+}", h.UserUpdateHandler)

	router.HandleFunc("/admin/transactions", h.TransactionsListHandler).Methods("GET")


	log.Println("Server starts on port :80")
	http.ListenAndServe(":80", router)
}