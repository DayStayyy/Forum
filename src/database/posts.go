package database

import (
	"database/sql"
	"log"
	"math/rand"
	"src/structs"
	"strconv"
	"strings"
	"time"
)

// AddPost : Add a new Post to database (database, post structure)
func AddPost(db *sql.DB, post structs.Post) error {
	_, err := db.Exec("INSERT INTO posts (userID, title, content, categories, postedTime) VALUES (?, ?, ?, ?, ?)", post.UserID, post.Title, post.Content, strings.Join(post.Categories, ","), post.PostedTime)
	if err != nil {
		return err
	}
	return nil
}

// RemovePost : Delete a post from database (database, post id)
func RemovePost(db *sql.DB, ID int) error {
	_, err := db.Exec("DELETE FROM posts WHERE postId = ?", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetPost : Return a post structure with all post informations (database, post id)
func GetPost(db *sql.DB, ID int) (structs.Post, error) {
	var post structs.Post
	var empty structs.Post
	rows, err := db.Query("select postID, userID, title, content, categories, postedTime from posts where postId = ?", ID)
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		var tempCategories string
		var tempTime string
		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &tempCategories, &tempTime)
		if err != nil {
			log.Fatal(err)
			return empty, err
		}

		post.PostedTime, err = time.Parse(time.RFC3339, tempTime)
		if err != nil {
			log.Fatal(err)
			return empty, err
		}

		post.Categories = strings.Split(tempCategories, ",")
		if err != nil {
			log.Fatal(err)
			return empty, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	return post, nil
}

// GetPostsWithUserID : Return an array of post structure with all posts from an user (database, user id)
func GetPostsWithUserID(db *sql.DB, ID int) ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := db.Query("select postID from posts where userID = ?", ID)
	if err != nil {
		log.Fatal(err)
		return make([]structs.Post, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)

		if err != nil {
			log.Fatal(err)
			return make([]structs.Post, 0), err
		}
		temp2 := strings.Split(temp, ",")
		for i := 0; i < len(temp2); i++ {
			var post structs.Post
			value, err := strconv.Atoi(temp)

			post, err = GetPost(db, value)
			if err != nil {
				log.Fatal(err)
				return make([]structs.Post, 0), err
			}
			posts = append(posts, post)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.Post, 0), err
	}
	return posts, nil
}

// GetAllPosts : Return an array of post structure with all posts in database (db)
func GetAllPosts(db *sql.DB) ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := db.Query("select postId from posts")
	if err != nil {
		log.Fatal(err)
		return make([]structs.Post, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Post, 0), err
		}
		var post structs.Post
		value, err := strconv.Atoi(temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Post, 0), err
		}

		post, err = GetPost(db, value)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Post, 0), err
		}
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.Post, 0), err
	}

	return sortPostsByDate(posts), nil
}

func GetPostsByCategory(db *sql.DB, ID string) ([]structs.Post, error) {
	var posts []structs.Post
	rows, err := db.Query("select postId, categories from posts")
	if err != nil {
		log.Fatal(err)
		return make([]structs.Post, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		var postID int

		err = rows.Scan(&postID, &temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Post, 0), err
		}

		if isInStringArray(strings.Split(temp, ","), ID) {
			var post structs.Post
			post, err = GetPost(db, postID)
			if err != nil {
				log.Fatal(err)
				return make([]structs.Post, 0), err
			}
			posts = append(posts, post)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.Post, 0), err
	}

	return sortPostsByDate(posts), nil
}

func stringAToIntA(array []string) ([]int, error) {
	var result []int
	if len(array) == 1 {
		return result, nil
	}
	for i := 0; i < len(array); i++ {
		val, err := strconv.Atoi(array[i])
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}

func sortPostsByDate(posts []structs.Post) []structs.Post {
	if len(posts) < 2 {
		return posts
	}

	left, right := 0, len(posts)-1
	pivot := rand.Intn(len(posts))
	posts[pivot], posts[right] = posts[right], posts[pivot]

	for i := 0; i < len(posts); i++ {
		if posts[i].PostedTime.After(posts[right].PostedTime) {
			posts[left], posts[i] = posts[i], posts[left]
			left++
		}
	}

	posts[left], posts[right] = posts[right], posts[left]

	sortPostsByDate(posts[:left])
	sortPostsByDate(posts[left+1:])

	return posts
}

func isInStringArray(array []string, str string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == str {
			return true
		}
	}
	return false
}

func GetFullPosts(db *sql.DB, posts []structs.Post) ([]structs.FullPost, error) {
	var fullPost []structs.FullPost
	var err error
	for i := range posts {
		var tmp structs.FullPost
		tmp.Likes, tmp.Dislikes, err = GetLikeNbByPost(db, posts[i].ID)
		if err != nil {
			log.Fatal(err)
			return fullPost, err
		}

		tmp.ID = posts[i].ID
		tmp.UserID = posts[i].UserID
		tmp.Title = posts[i].Title
		tmp.Content = posts[i].Content

		for j := 0; j < len(posts[i].Categories); j++ {
			name, err := GetCategoryNameById(db, posts[i].Categories[j])
			if err != nil {
				log.Fatal(err)
			}
			tmp.Categories = append(tmp.Categories, name)
		}

		ownerData, err := GetUser(db, posts[i].UserID)
		if err != nil {
			log.Fatal(err)
			return fullPost, err
		}
		tmp.Owner = ownerData.Name

		rawDate := posts[i].PostedTime.Format(time.RFC1123)
		Date := rawDate[5:16]
		hour := rawDate[17:25]
		tmp.Date = Date + "  à  " + hour
		fullPost = append(fullPost, tmp)
	}
	return fullPost, nil
}

func GetFullPost(db *sql.DB, post structs.Post) (structs.FullPost, error) {
	var fullPost structs.FullPost
	var err error

	fullPost.Likes, fullPost.Dislikes, err = GetLikeNbByPost(db, post.ID)
	if err != nil {
		log.Fatal(err)
		return fullPost, err
	}

	fullPost.ID = post.ID
	fullPost.UserID = post.UserID
	fullPost.Title = post.Title
	fullPost.Content = post.Content

	for j := 0; j < len(post.Categories); j++ {
		name, err := GetCategoryNameById(db, post.Categories[j])
		if err != nil {
			log.Fatal(err)
		}
		fullPost.Categories = append(fullPost.Categories, name)
	}

	ownerData, err := GetUser(db, post.UserID)
	if err != nil {
		log.Fatal(err)
		return fullPost, err
	}
	fullPost.Owner = ownerData.Name

	rawDate := post.PostedTime.Format(time.RFC1123)
	Date := rawDate[5:16]
	hour := rawDate[17:25]
	fullPost.Date = Date + "  à  " + hour

	return fullPost, nil
}
