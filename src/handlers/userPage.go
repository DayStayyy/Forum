package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"strconv"
	"text/template"
)

// UserPage : Handler for user page
func UserPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/userPage" {
		log.Fatal(errors.New("URL.Path incorrect (UserPage)"))
	}

	loggedIn, currentUser, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "GET" {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	if r.Method == "POST" {
		files := []string{"templates/html/pages/userPage.html", "templates/html/pieces/userMenu.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		ID, err := strconv.Atoi(r.FormValue("user"))
		if err != nil {
			log.Fatal(err)
		}

		userData, err := database.GetUser(db, ID)
		if err != nil {
			log.Fatal(err)
		}

		posts, err := database.GetPostsWithUserID(db, ID)
		if err != nil {
			log.Fatal(err)
		}
		isNoPost := false

		fullPosts, err := database.GetFullPosts(db, posts)
		if err != nil {
			log.Fatal(err)
		} else if len(fullPosts) == 0 {
			isNoPost = true
		}

		if userData.Description == "" {
			userData.Description = "<br>" + userData.Name + " n'a pas de description<br><br>"
		}

		posts, err = database.GetPostsWithUserID(db, ID)
		if err != nil {
			log.Fatal(err)
		}

		if len(fullPosts) > 5 {
			fullPosts = fullPosts[:5]
		}

		data := struct {
			User     structs.User
			UserLook structs.User
			Posts    []structs.FullPost
			IsPost   bool
		}{currentUser, userData, fullPosts, isNoPost}

		err = ts.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
	}
}
