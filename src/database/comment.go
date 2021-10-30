package database

import (
	"database/sql"
	"log"
	"src/structs"
	"strconv"
	"strings"
	"time"
)

// AddComment : Add a comment to database (database, comment structure)
func AddComment(db *sql.DB, comment structs.Comment) error {
	_, err := db.Exec("INSERT INTO comments (userID, parentID, postID, content, postedTime) VALUES (?, ?, ?, ?, ?)",
		comment.UserID, comment.ParentID, comment.PostID, comment.Content, comment.PostedTime)

	if err != nil {
		return err
	}
	return nil
}

// RemoveComment : Remove comment from database (database, comment id)
func RemoveComment(db *sql.DB, ID int) error {
	_, err := db.Exec("DELETE FROM comments WHERE commentID = ?", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetComment : Get comment in database (database, comment id)
func GetComment(db *sql.DB, ID int) (structs.Comment, error) {
	var comment structs.Comment
	var empty structs.Comment
	rows, err := db.Query("select commentId, userID, parentID, postID, content, postedTime from comments where commentId = ?", ID)
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&comment.ID, &comment.UserID, &comment.ParentID, &comment.PostID, &comment.Content, &comment.PostedTime)
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
	return comment, nil
}

// GetCommentsWithUserID : Return an array of comment structure with all comments from an user (database, user id)
func GetCommentsWithUserID(db *sql.DB, ID int) ([]structs.Comment, error) {
	var posts []structs.Comment
	rows, err := db.Query("select commentID from users where userID = ?", ID)
	if err != nil {
		log.Fatal(err)
		return make([]structs.Comment, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Comment, 0), err
		}

		temp2 := strings.Split(temp, ",")
		for i := 0; i < len(temp2); i++ {
			var post structs.Comment
			var value int

			post, err = GetComment(db, value)
			if err != nil {
				log.Fatal(err)
				return make([]structs.Comment, 0), err
			}

			posts = append(posts, post)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.Comment, 0), err
	}
	return posts, nil
}

// GetCommentsWithPost : Return an array of comment structure with all comments from an user (database, post id)
func GetCommentsWithPost(db *sql.DB, ID int) ([]structs.Comment, error) {
	var comments []structs.Comment
	rows, err := db.Query("select commentID from comments where postID = ?", ID)
	if err != nil {
		log.Fatal(err)
		return make([]structs.Comment, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Comment, 0), err
		}

		temp2 := strings.Split(temp, ",")
		for i := 0; i < len(temp2); i++ {
			var comment structs.Comment

			value, err := strconv.Atoi(temp)
			if err != nil {
				log.Fatal(err)
				return make([]structs.Comment, 0), err
			}

			comment, err = GetComment(db, value)
			if err != nil {
				log.Fatal(err)
				return make([]structs.Comment, 0), err
			}

			comments = append(comments, comment)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.Comment, 0), err
	}
	return comments, nil
}

func GetFullComments(db *sql.DB, comments []structs.Comment) ([]structs.FullComment, error) {
	var fullComment []structs.FullComment
	for i := range comments {
		var tmp structs.FullComment
		tmp.ID = comments[i].ID
		tmp.UserID = comments[i].UserID
		tmp.ParentID = comments[i].ParentID
		tmp.PostID = comments[i].PostID
		tmp.Content = comments[i].Content
		tmp.PostedTime = comments[i].PostedTime

		ownerData, err := GetUser(db, comments[i].UserID)
		if err != nil {
			log.Fatal(err)
			return fullComment, err
		}
		tmp.Owner = ownerData.Name
		rawDate := comments[i].PostedTime.Format(time.RFC1123)
		Date := rawDate[5:16]
		hour := rawDate[17:25]
		tmp.Date = Date + "  à  " + hour

		fullComment = append(fullComment, tmp)
	}
	return fullComment, nil
}
