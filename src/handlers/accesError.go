package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"text/template"
)

func AccesError(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/acessError" {
		log.Fatal(errors.New("URL.Path incorrect (acessError)"))
	}

	login := "templates/html/pieces/guestLogin.html"

	loggedIn, currentUser, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if loggedIn {
		login = "templates/html/pieces/userMenu.html"
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	files := []string{"templates/html/erreurs/erreurAcces.html", login, "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	ownerData, _ := database.GetUser(db, currentUser.ID)

	data := struct {
		User structs.User
		Name string
	}{currentUser, ownerData.Name}

	err = ts.Execute(w, data)

	if err != nil {
		log.Fatal(err)
	}
	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}
}
