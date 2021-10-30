package handlers

import (
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"text/template"
)

// HomePage : Handler for home page
func HomePage(w http.ResponseWriter, r *http.Request) {
	login := "templates/html/pieces/guestLogin.html"

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if loggedIn {
		login = "templates/html/pieces/userMenu.html"
	}

	files := []string{"templates/html/pages/homePage.html", login, "templates/html/pieces/blank.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error : check if template file exists", 500)
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}
	posts, err := database.GetAllPosts(db)
	if err != nil {
		log.Fatal(err)
	}
	isNoPost := false

	fullPosts, err := database.GetFullPosts(db, posts)
	if err != nil {
		log.Fatal(err)
	}

	if len(fullPosts) == 0 {
		isNoPost = true
	}

	categories, err := database.GetAllCategories(db)
	if err != nil {
		log.Fatal(err)
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	if len(fullPosts) > 6 {
		fullPosts = fullPosts[:5]
	}

	data := struct {
		User       structs.User
		Posts      []structs.FullPost
		Categories []structs.Category
		IsNoPost   bool
	}{user, fullPosts, categories, isNoPost}

	err = ts.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
