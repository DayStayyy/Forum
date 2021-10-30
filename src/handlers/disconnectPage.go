package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
)

// Disconnect : Handler for disconnecting user (blank page)
func Disconnect(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/disconnect" {
		log.Fatal(errors.New("URL.Path incorrect (Disconnect)"))
	}

	cookie, err := r.Cookie("remember_token")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}
	user, err := database.GetUserByCookie(db, cookie.Value)
	if err != nil {
		log.Fatal(err)
	}
	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	var connect structs.ConnectInfos
	connect.Name = user.Mail
	connect.Password = ""

	err = authentification.DisconnectUser(w, r, connect)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
