package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"strconv"
)

// SupprUser : Handler for deleting a user in bossPage
func SupprUser(w http.ResponseWriter, r *http.Request) {
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

	userId, err := strconv.Atoi(r.FormValue("supprUser"))
	if err != nil {
		log.Fatal(err)
	}

	err = database.RemoveUser(db, userId)
	if err != nil {
		log.Fatal(err)
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, "/bossPage", http.StatusFound)
}
