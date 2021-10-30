package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/fileUpload"
	"src/structs"
	"strconv"
	"text/template"
)

func ModifyProfilInfo(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/modifier" {
		log.Fatal(errors.New("URL.Path incorrect ()"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if r.Method == "POST" {

		files := []string{"templates/html/loading/pageModifier.html", "templates/html/pieces/userMenu.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		name := r.FormValue("name")
		description := r.FormValue("description")

		userID, err := strconv.Atoi(r.FormValue("userID"))
		if err != nil {
			log.Fatal(err)
		}

		db, err := database.ConnectDatabase("./database.db")
		if err != nil {
			log.Fatal(err)
		}

		userToModify, err := database.GetUser(db, userID)
		if err != nil {
			log.Fatal(err)
		}

		if name != "" {
			userToModify.Name = name
		} else if description != "" {
			userToModify.Description = description
		} else {
			userToModify.Image = fileUpload.UploadFile(r, userToModify.Name)
		}

		err = database.UpdateUser(db, userToModify)
		if err != nil {
			log.Fatal(err)
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
		data := struct {
			User structs.User
		}{user}
		err = ts.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
