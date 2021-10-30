package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"src/utils"
	"strconv"
	"text/template"
)

// PostPage : Handler for post page
func PostPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/postPage" {
		log.Fatal(errors.New("URL.Path incorrect (PostPage)"))
	}

	login := "templates/html/pieces/guestLogin.html"

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if loggedIn {
		login = "templates/html/pieces/userMenu.html"
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == "GET" {
		files := []string{"templates/html/erreurs/erreurAcces.html", login, "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		ownerData, _ := database.GetUser(db, user.ID)
		data := struct {
			User structs.User
			Name string
		}{user, ownerData.Name}

		err = ts.Execute(w, data)
		if err != nil {
			log.Fatal(err)
		}
	}

	if r.Method == "POST" {
		files := []string{"templates/html/pages/postPage.html", login, "templates/html/pieces/blank.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Fatal(err)
		}

		postID, err := strconv.Atoi(r.FormValue("postID"))
		if err != nil {
			log.Fatal(err)
		}

		post, err := database.GetPost(db, postID)
		if err != nil {
			log.Fatal(err)
		}

		comments, err := database.GetCommentsWithPost(db, post.ID)
		if err != nil {
			log.Fatal(err)
		}

		var IsNoComments bool

		if len(comments) == 0 {
			IsNoComments = true
			var comment structs.Comment
			comment.Content = "Aucun commentaire n'a été publié."
			comments = append(comments, comment)
		}

		var posts []structs.Post
		posts = append(posts, post)

		fullPosts, err := database.GetFullPosts(db, posts)
		if err != nil {
			log.Fatal(err)
		}

		var fullPost structs.FullPost
		if len(fullPosts) == 1 {
			fullPost = fullPosts[0]
		}

		fullComment, err := database.GetFullComments(db, comments)
		if err != nil {
			log.Fatal(err)
		}

		fullComment = utils.CommentSort(fullComment)

		data := struct {
			User         structs.User
			Post         structs.FullPost
			Comments     []structs.FullComment
			IsNoComments bool
		}{user, fullPost, fullComment, IsNoComments}

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
