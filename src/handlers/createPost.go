package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"text/template"
	"time"
)

// CreatePost : Hnadler for create post page
func CreatePost(w http.ResponseWriter, r *http.Request) {

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if r.URL.Path != "/createPost" {
		log.Fatal(errors.New("URL.Path incorrect (CreatePost)"))
	}

	type CreatePostErrorStruct struct {
		Title, Categories string
	}
	var errStruct CreatePostErrorStruct
	isErrors := false

	var postData structs.Post

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "POST" {
		r.ParseForm()
		postData.Title = r.FormValue("title")
		postData.UserID = user.ID
		postData.Content = r.FormValue("content")
		postData.Categories = r.Form["Categories"]
		postData.PostedTime = time.Now()

		if len(postData.Title) > 254 {
			errStruct.Title = "Title is too long (255 caractères maximum) !"
			isErrors = true
		} else if len(postData.Title) == 0 {
			errStruct.Title = "Please enter a title !"
			isErrors = true
		}

		if len(postData.Categories) == 0 {
			errStruct.Categories = "You need to select at least 1 category."
			isErrors = true
		}

		if !isErrors {
			err = database.AddPost(db, postData)
			if err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	categories, err := database.GetAllCategories(db)
	if err != nil {
		log.Fatal(err)
	}

	files := []string{"templates/html/pages/createPost.html", "templates/html/pieces/userMenu.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		User       structs.User
		Categories []structs.Category
		Errors     CreatePostErrorStruct
	}{user, categories, errStruct}

	err = ts.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}
