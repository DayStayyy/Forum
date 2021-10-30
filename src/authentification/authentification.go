package authentification

import (
	"errors"
	"log"
	"net/http"
	"src/database"
	"src/structs"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// SignIn : User sign in function (responseWriter, user structure)
func SignIn(w http.ResponseWriter, user structs.User) error {
	if user.RememberToken != "" {
		token, err := generateRememberToken()
		if err != nil {
			log.Fatal(err)
			return err
		}
		user.RememberToken = token

		db, err := database.ConnectDatabase("./database.db")
		if err != nil {
			log.Fatal(err)
			return err
		}
		err = database.UpdateUser(db, user)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.RememberToken,
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}

// ConnectUser : Function to conenct user (response writer, connect infos)
func ConnectUser(w http.ResponseWriter, infos structs.ConnectInfos) error {
	userExist, userID, err := checkIfUserExists(infos)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if !userExist {
		err = errors.New("User mail=" + infos.Name + " do not exist in database (authentification.ConnectUser)")
		return err
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
		return err
	}

	user, err := database.GetUser(db, userID)
	if err != nil {
		log.Fatal(err)
		return err
	}

	_, err = checkPassword(user, infos)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			err2 := database.DisconnectDatabase(db)
			if err2 != nil {
				log.Fatal(err2)
				return err2
			}
			return err
		}
		err2 := database.DisconnectDatabase(db)
		if err2 != nil {
			log.Fatal(err2)
			return err2
		}
		log.Fatal(err)
		return err
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = SignIn(w, user)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// DisconnectUser : Function to disconenct user(response writer, user structure)
func DisconnectUser(w http.ResponseWriter, r *http.Request, infos structs.ConnectInfos) error {
	isLogged, user, err := IsLoggedIn(r)
	if err != nil {
		log.Fatal(err)
		return err
	}
	if !isLogged {
		err = errors.New("Can't disconnect an unconndcted user (authentification.DisconnectUser)")
		log.Fatal(err)
		return err
	}

	c := http.Cookie{
		Name:    "remember_token",
		Value:   "Unused",
		Expires: time.Unix(1414414788, 1414414788000)}
	http.SetCookie(w, &c)

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
		return err
	}

	user.RememberToken = c.Expires.String()
	err = database.UpdateUser(db, user)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//IsLoggedIn : Return true and user struct if user is connected
func IsLoggedIn(r *http.Request) (bool, structs.User, error) {
	var empty structs.User
	cookie, err := r.Cookie("remember_token")
	if err == http.ErrNoCookie {
		return false, empty, nil
	} else if err != nil {
		log.Fatal(err)
		return false, empty, err
	}

	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
		return false, empty, err
	}

	user, err := database.GetUserByCookie(db, cookie.Value)
	if err != nil {
		log.Fatal(err)
		return false, empty, err
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
		return false, empty, err
	}
	return true, user, nil
}

func checkIfUserExists(infos structs.ConnectInfos) (bool, int, error) {
	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
		return false, 0, err
	}

	rows, err := db.Query("SELECT userID FROM users WHERE userMail = ?", infos.Name)
	if err != nil {
		log.Fatal(err)
		return false, 0, err
	}

	count := 0
	id := 0
	for rows.Next() {
		count++
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
			return false, 0, err
		}
	}

	err = database.DisconnectDatabase(db)
	if err != nil {
		log.Fatal(err)
		return false, 0, err
	}

	if count == 1 {
		return true, id, nil
	}
	return false, 0, nil
}

func checkPassword(user structs.User, infos structs.ConnectInfos) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(infos.Password))
	if err != nil {
		return false, err
	}
	return true, nil
}

func CheckLoginForm(infos structs.ConnectInfos) (bool, string) {
	db, err := database.ConnectDatabase("./database.db")
	if err != nil {
		log.Fatal(err)
	}

	exist, id, err := checkIfUserExists(infos)
	if err != nil {
		log.Fatal(err)
	} else if !exist {
		return false, "Login or password incorrect."
	}

	user, err := database.GetUser(db, id)
	if err != nil {
		log.Fatal(err)
	}

	passwdValid, err := checkPassword(user, infos)
	if !passwdValid {
		return false, "Login or password incorrect."
	}

	return true, ""
}
