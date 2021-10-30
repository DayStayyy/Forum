package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/fileUpload"
	"src/structs"
	"src/utils"
	"strings"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

// Registration : Handler for registration page
func Registration(w http.ResponseWriter, r *http.Request) {

	var ifErrors bool
	var errorStruct struct {
		Name, Mail, Passwd string
	}
	errorStruct.Name = ""
	errorStruct.Mail = ""
	errorStruct.Passwd = ""

	if r.URL.Path != "/registration" {
		log.Fatal(errors.New("URL.Path incorrect (Registration)"))
	}

	if r.Method == "POST" {
		passwd, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 5)
		if err != nil {
			log.Fatal(err)
		}

		var user structs.User
		user.ID = 0
		user.Name = strings.Title(strings.ToLower(r.FormValue("prenom"))) + " " + strings.ToUpper((r.FormValue("nom")))
		user.Password = string(passwd)
		user.Mail = r.FormValue("email")
		user.Description = r.FormValue("description")
		user.Nationality = r.FormValue("nationality")
		user.Perms = 0

		ifErrors, errorStruct = utils.CheckAuthForm(user, r.FormValue("prenom"), r.FormValue("nom"), r.FormValue("password"))
		if !ifErrors {
			user.Image = fileUpload.UploadFile(r, user.Name)

			db, err := database.ConnectDatabase("./database.db")
			if err != nil {
				log.Fatal(err)
			}
			err = database.AddUser(db, user)
			if err != nil {
				log.Fatal(err)
			}
			err = database.DisconnectDatabase(db)
			if err != nil {
				log.Fatal(err)
			}

			var connect structs.ConnectInfos
			connect.Name = r.FormValue("email")
			connect.Password = r.FormValue("password")
			err = authentification.ConnectUser(w, connect)
			if err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, "/userPage", http.StatusFound)
		}
	}

	files := []string{"templates/html/pages/registrationPage.html", "templates/html/pieces/guestLogin.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error : check if template file exists", 500)
	}

	var user structs.User

	data := struct {
		Error structs.RegistrationErrors
		User  structs.User
	}{errorStruct, user}

	err = ts.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
