package database

import (
	"database/sql"
	"log"
	"src/structs"
)

// AddLikeToPost : Add a new post like to database (database, like structure)
func AddLikeToPost(db *sql.DB, like structs.Like) error {
	_, err := db.Exec("INSERT INTO likes (userID, postID, commentID, isPositive, postedTime) VALUES (?, ?, ?, ?, ?)", like.UserID, like.PostID, like.CommentID, like.IsPositive, like.PostedTime)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// AddLikeToComment : Add a new comment like to database (database, like structure)
func AddLikeToComment(db *sql.DB, like structs.Like) error {
	_, err := db.Exec("INSERT INTO likes (userID, postID, commentID, isPositive, postedTime) VALUES (?, ?, ?, ?, ?)", like.UserID, like.PostID, like.CommentID, like.IsPositive, like.PostedTime)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// RemoveLikeFromPost : Delete a post like from database (database, like id)
func RemoveLikeFromPost(db *sql.DB, postID int, userID int) error {
	_, err := db.Exec("DELETE FROM likes WHERE postID = ? AND userId = ?", postID, userID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// RemoveLikeFromComment : Delete a comment like from database (database, like id)
func RemoveLikeFromComment(db *sql.DB, ID int) error {
	_, err := db.Exec("DELETE FROM likes WHERE likeID = ?", ID)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// GetLike : Return a like structure with all like informations (database, like id)
func GetLike(db *sql.DB, ID int) (structs.Like, error) {
	var like structs.Like
	var empty structs.Like
	rows, err := db.Query("select likeID, userID, postID, commentID, isPositive, postedTime from likes where likeId = ?", ID)
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.IsPositive, &like.PostedTime)
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
	return like, nil
}

// GetLikeByUserID : Return a like structure with all like informations (database, like id)
func IfLikeIsInDB(db *sql.DB, userID int, postID int) (bool, bool, error) {
	row := db.QueryRow("select isPositive from likes where userID = ? AND postID = ?", userID, postID)

	var boolean int
	err := row.Scan(&boolean)
	if err != nil {
		return false, false, nil
	}

	err = row.Err()
	if err != nil {
		log.Fatal(err)
		return false, false, err
	}
	if boolean == 1 {
		return true, true, nil
	}
	return true, false, nil
}

func GetLikeNbByPost(db *sql.DB, ID int) (int, int, error) {
	var likes int
	var dislikes int
	rows, err := db.Query("select isPositive from likes where postId = ?", ID)
	if err != nil {
		log.Fatal(err)
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var boolean int
		err = rows.Scan(&boolean)
		if err != nil {
			log.Fatal(err)
			return 0, 0, err
		}
		if boolean == 1 {
			likes++
		} else {
			dislikes++
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return 0, 0, err
	}
	return likes, dislikes, nil
}
