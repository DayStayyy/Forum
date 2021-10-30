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

// PostPage : Handler for post page
func BossPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/bossPage" {
		log.Fatal(errors.New("URL.Path incorrect (bosspage)"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else if user.Perms == 0 {
		http.Redirect(w, r, "/acessError", http.StatusFound)
	}
	files := []string{"templates/html/pages/bossPage.html", "templates/html/pieces/userMenu.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error : check if template file exists", 500)
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	users, err := database.GetAllUsers(db)
	if err != nil {
		log.Fatal(err)
	}

	tmpPosts, err := database.GetAllPosts(db)
	if err != nil {
		log.Fatal(err)
	}

	posts, err := database.GetFullPosts(db, tmpPosts)
	if err != nil {
		log.Fatal(err)
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}

	isNoPosts := false
	if len(posts) == 0 {
		isNoPosts = true
	}

	data := struct {
		User      structs.User
		Users     []structs.User
		Posts     []structs.FullPost
		IsNoPosts bool
	}{user, users, posts, isNoPosts}

	err = ts.Execute(w, data)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", 500)
	}
}
