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
	"time"
)

// CreateCommentPage : Hnadler for create comment page
func CreateCommentPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createComment" {
		log.Fatal(errors.New("URL.Path incorrect (CreateComment)"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if r.Method == "POST" {
		files := []string{"templates/html/loading/postCommentary.html", "templates/html/pieces/userMenu.html", "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		db, err := database.ConnectDatabase("./database.db")
		if err != nil {
			log.Fatal(err)
		}

		var postId int
		var parentId int
		content := r.FormValue("content")

		if r.FormValue("postID") != "" {
			postId, err = strconv.Atoi(r.FormValue("postID"))
			if err != nil {
				log.Fatal(err)
			}
		}
		if r.FormValue("parentID") != "" {
			parentId, err = strconv.Atoi(r.FormValue("parentID"))
			if err != nil {
				log.Fatal(err)
			}
		}

		var comment structs.Comment
		comment.UserID = user.ID
		comment.PostID = postId
		comment.ParentID = parentId
		comment.Content = content
		comment.PostedTime = time.Now()
		err = database.AddComment(db, comment)
		if err != nil {
			log.Fatal(err)
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}

		data := struct {
			User   structs.User
			PostID int
		}{user, postId}

		err = ts.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
