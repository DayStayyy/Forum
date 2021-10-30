package main

import (
	"fmt"
	"net/http"
	"src/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Listening on port : 8080")
	fs := http.FileServer(http.Dir("./fileServer"))
	http.Handle("/fileServer/", http.StripPrefix("/fileServer/", fs))
	css := http.FileServer(http.Dir("./css"))
	http.Handle("/css/", http.StripPrefix("/css/", css))
	fmt.Println("Listening on port : 8080")

	http.HandleFunc("/", handlers.HomePage)
	http.HandleFunc("/registration", handlers.Registration)
	http.HandleFunc("/userPage", handlers.UserPage)
	http.HandleFunc("/createPost", handlers.CreatePost)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/postPage", handlers.PostPage)
	http.HandleFunc("/disconnect", handlers.Disconnect)
	http.HandleFunc("/postLike", handlers.PostLike)
	http.HandleFunc("/postDislike", handlers.PostDislike)
	http.HandleFunc("/supprAccountPage", handlers.SupprUser)
	http.HandleFunc("/supprPostPage", handlers.SupprPost)
	http.HandleFunc("/supprPostPage2", handlers.SupprPostFromPost)
	http.HandleFunc("/bossPage", handlers.BossPage)
	http.HandleFunc("/createComment", handlers.CreateCommentPage)
	http.HandleFunc("/categoriesPage", handlers.CategoriesPage)
	http.HandleFunc("/profilSettings", handlers.ModifyPage)
	http.HandleFunc("/modifier", handlers.ModifyProfilInfo)
	http.HandleFunc("/modifyModo", handlers.ModifyProfilInfo)
	http.HandleFunc("/acessError", handlers.AccesError)

	fmt.Println("Listening on port : 8080")
	http.ListenAndServe(":8080", nil)
}
