package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"strconv"
)

// ModifyModoPage : Handler for modifying user perms
func ModifyModoPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/supprAccountPage" {
		log.Fatal(errors.New("URL.Path incorrect (supprAccountPage)"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if user.Perms != 2 {
		http.Redirect(w, r, "/acessError", http.StatusFound)
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	var userId int
	if r.FormValue("promote") != "" {
		userId, err = strconv.Atoi(r.FormValue("promote"))
		if err != nil {
			log.Fatal(err)
		}
	} else if r.FormValue("demote") != "" {
		userId, err = strconv.Atoi(r.FormValue("demote"))
		if err != nil {
			log.Fatal(err)
		}
	}

	userToModify, err := database.GetUser(db, userId)
	if err != nil {
		log.Fatal(err)
	}

	if r.FormValue("promote") != "" {
		userToModify.Perms = 1
	} else if r.FormValue("demote") != "" {
		userToModify.Perms = 0
	}

	err = database.UpdateUser(db, userToModify)
	if err != nil {
		log.Fatal(err)
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/bossPage", http.StatusFound)
}
