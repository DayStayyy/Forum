package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/structs"
	"text/template"
)

// Login : Handler for login page
func Login(w http.ResponseWriter, r *http.Request) {
	var mail, passwd string

	if r.URL.Path != "/login" {
		log.Fatal(errors.New("URL.Path incorrect (Login)"))
	}

	if r.Method == "POST" {
		var connect structs.ConnectInfos
		connect.Name = r.FormValue("email")
		connect.Password = r.FormValue("password")

		if len(connect.Name) == 0 {
			mail = "Vous devez entrer une adresse email valide."
		} else if len(connect.Password) < 9 {
			passwd = "Vous devez entrer un mot de passe valide (8 caractères minimum)."
		} else {
			err := authentification.ConnectUser(w, connect)
			if err == nil {
				http.Redirect(w, r, "/", http.StatusFound)
			} else {
				if passwd == "" {
					passwd = "Email ou mot de passe invalide."
				}
			}
		}
	}

	files := []string{"templates/html/pages/loginPage.html", "templates/html/pieces/guestLogin.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error : check if template file exists", 500)
	}

	var user structs.User

	data := struct {
		User   structs.User
		Mail   string
		Passwd string
	}{user, mail, passwd}

	err = ts.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
