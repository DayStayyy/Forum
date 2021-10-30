package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"strconv"
)

// SupprPost : Handler for deleting a post in bossPage
func SupprPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/supprPostPage" {
		log.Fatal(errors.New("URL.Path incorrect (supprPostPage)"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
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

	postId, err := strconv.Atoi(r.FormValue("supprPost"))
	if err != nil {
		log.Fatal(err)
	}

	if user.Perms < 1 {
		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/acessError", http.StatusFound)
	} else {
		err = database.RemovePost(db, postId)
		if err != nil {
			log.Fatal(err)
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.Redirect(w, r, "/bossPage", http.StatusFound)
}

// SupprPostFromPost : Handler for deleting a post in the postPage
func SupprPostFromPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/supprPostPage2" {
		log.Fatal(errors.New("URL.Path incorrect (supprPostPage2)"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
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

	postId, err := strconv.Atoi(r.FormValue("supprPost"))
	if err != nil {
		log.Fatal(err)
	}

	post, err := database.GetPost(db, postId)
	if err != nil {
		log.Fatal(err)
	}

	if user.ID != post.UserID {
		fmt.Println("toto")
		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/acessError", http.StatusFound)
	} else {
		err = database.RemovePost(db, postId)
		if err != nil {
			log.Fatal(err)
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
