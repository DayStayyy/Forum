package database

import (
	"database/sql"
	"log"
	"src/structs"
	"strconv"
)

// AddUser : Add a user to database (database, user structure)
func AddUser(db *sql.DB, user structs.User) error {
	_, err := db.Exec("INSERT INTO users (userName, userPassword, userMail, userImage, nationality, description, perms, rememberToken) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Name, user.Password, user.Mail, user.Image, user.Nationality, user.Description, user.Perms, user.RememberToken)

	if err != nil {
		return err
	}
	return nil
}

// RemoveUser : Remove user from database (database, user id)
func RemoveUser(db *sql.DB, ID int) error {
	_, err := db.Exec("DELETE FROM users WHERE userID = ?", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetUser : Get user in database (database, user id)
func GetUser(db *sql.DB, ID int) (structs.User, error) {
	var user structs.User
	var empty structs.User
	rows, err := db.Query("select userID, userName, userPassword, userMail, userImage, nationality, description, perms, rememberToken from users where userId = ?", ID)
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Password, &user.Mail, &user.Image, &user.Nationality, &user.Description, &user.Perms, &user.RememberToken)
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
	return user, nil
}

// UpdateUser : Update user with a new user sructure
func UpdateUser(db *sql.DB, user structs.User) error {
	_, err := db.Exec("UPDATE users set userName=?, userPassword=?, userMail=?, userImage=?, nationality=?, description=?, perms=?, rememberToken=? WHERE userID =? ",
		user.Name, user.Password, user.Mail, user.Image, user.Nationality, user.Description, user.Perms, user.RememberToken, user.ID)

	if err != nil {
		return err
	}
	return nil
}

// GetNbOfUsers : Get nb of users in database (db)
func GetNbOfUsers(db *sql.DB) (int, error) {
	rows, err := db.Query("select userID from users")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer rows.Close()

	counter := 0
	for rows.Next() {
		counter++
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return counter, nil
}

// GetUserByCookie : Get user in database from his cookie (database, token string)
func GetUserByCookie(db *sql.DB, token string) (structs.User, error) {
	var empty structs.User
	var user structs.User

	rows, err := db.Query("select userID, userName, userPassword, userMail, userImage, nationality, description, perms, rememberToken from users where rememberToken = ?", token)
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Password, &user.Mail, &user.Image, &user.Nationality, &user.Description, &user.Perms, &user.RememberToken)
		if err != nil {
			log.Fatal(err)
			return empty, err
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
			return empty, err
		}
	}
	return user, nil
}

// GetUserEmails : Get emails of all users in database (db)
func GetUserEmails(db *sql.DB) ([]string, error) {
	var empty []string
	var results []string
	rows, err := db.Query("select userMail from users")
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		rows.Scan(&temp)
		results = append(results, temp)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	return results, nil
}

// GetAllUsers : Get all users in database (db)
func GetAllUsers(db *sql.DB) ([]structs.User, error) {
	var users []structs.User
	rows, err := db.Query("select userId from users")
	if err != nil {
		log.Fatal(err)
		return make([]structs.User, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.User, 0), err
		}
		var user structs.User
		value, err := strconv.Atoi(temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.User, 0), err
		}

		user, err = GetUser(db, value)
		if err != nil {
			log.Fatal(err)
			return make([]structs.User, 0), err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.User, 0), err
	}
	return users, nil
}
