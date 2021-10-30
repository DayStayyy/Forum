package handlers

import (
	"errors"
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"strconv"
	"time"
)

// PostDislike : Handler for adding a dislike to a post (blank page)
func PostDislike(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/postDislike" {
		log.Fatal(errors.New("URL.Path incorrect (postDislike)"))
	}

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if loggedIn {
		db, err := database.ConnectDatabase("./database.db")
		if err != nil {
			log.Fatal(err)
		}

		postId, err := strconv.Atoi(r.FormValue("PouceRouge"))
		if err != nil {
			log.Fatal(err)
		}

		likes, isPositive, err := database.IfLikeIsInDB(db, user.ID, postId)
		if likes && !isPositive {
			database.RemoveLikeFromPost(db, postId, user.ID)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			if isPositive {
				err = database.RemoveLikeFromPost(db, postId, user.ID)
				if err != nil {
					log.Fatal(err)
				}
			}
			var newLike structs.Like
			newLike.UserID = user.ID
			newLike.PostID = postId
			newLike.PostedTime = time.Now()
			newLike.IsPositive = 0

			err = database.AddLikeToPost(db, newLike)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}
