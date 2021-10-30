package handlers

import (
	"log"
	"net/http"
	"src/authentification"
	"src/database"
	"src/structs"
	"src/utils"
	"text/template"
)

// HomePage : Handler for home page
func CategoriesPage(w http.ResponseWriter, r *http.Request) {
	login := "templates/html/pieces/guestLogin.html"

	loggedIn, user, err := authentification.IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
	}
	if loggedIn {
		login = "templates/html/pieces/userMenu.html"
	}

	files := []string{"templates/html/pages/categoriesPage.html", login, "templates/html/pieces/flexFilter.html", "templates/html/pieces/header.html", "templates/html/pieces/footer.html", "templates/html/pieces/base.html"}

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

	if r.Method == "POST" {
		posts, err = database.GetPostsByCategory(db, r.FormValue("categorieName"))
		if err != nil {
			log.Fatal(err)
		}

		isNoPost := false
		var homePageReturn []structs.HomePageReturn
		for i := 0; i < len(posts); i++ {
			var tmp structs.HomePageReturn
			tmp.Likes, tmp.Dislikes, err = database.GetLikeNbByPost(db, posts[i].ID)
			if err != nil {
				log.Fatal(err)
			}

			postUser, err := database.GetUser(db, posts[i].UserID)
			if err != nil {
				log.Fatal(err)
			}

			tmp.Name = postUser.Name
			tmp.Post, err = database.GetFullPost(db, posts[i])
			if err != nil {
				log.Fatal(err)
			}

			homePageReturn = append(homePageReturn, tmp)
		}
		if len(homePageReturn) == 0 {
			isNoPost = true
		}

		categories, err := database.GetAllCategories(db)
		if err != nil {
			log.Fatal(err)
		}

		var catName string
		catName, err = database.GetCategoryNameById(db, r.FormValue("categorieName"))
		if err != nil {
			log.Fatal(err)
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}

		isLike := r.FormValue("isHerelikes")
		isDate := r.FormValue("isHeredate")
		if isLike == "hidden" {
			if isDate == "decroi" {
				homePageReturn = utils.SortHomePageReturnByDateDecroissant(homePageReturn)
			}
		} else {
			if isLike == "croi" {
				homePageReturn = utils.SortHomePageReturnByLikesCroissant(homePageReturn)
			} else {
				homePageReturn = utils.SortHomePageReturnByLikesDecroissant(homePageReturn)
			}
		}

		data := struct {
			User       structs.User
			Data       []structs.HomePageReturn
			Categories []structs.Category
			CatName    string
			CatID      string
			IsNoPost   bool
		}{user, homePageReturn, categories, catName, r.FormValue("categorieName"), isNoPost}

		err = ts.Execute(w, data)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
