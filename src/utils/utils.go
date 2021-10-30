package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"src/database"
	"src/structs"
	"strings"
)

func isAValidEmail(mail string) bool {
	if strings.Count(mail, "@") != 1 {
		return false
	}
	mail = mail[strings.Index(mail, "@")+1:]
	if strings.Count(mail, ".") != 1 {
		return false
	}
	return true
}

// Hash : Return token hash for database storage (string)
func Hash(input string) (string, error) {
	h := hmac.New(sha256.New, []byte("anEncryptionKey"))
	h.Reset()
	_, err := h.Write([]byte(input))
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	b := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(b), nil
}

// CheckAuthForm : Checks if registration form is valid and return errors
func CheckAuthForm(user structs.User, firstName string, name string, passwd string) (bool, struct{ Name, Mail, Passwd string }) {
	var errorStruct struct {
		Name, Mail, Passwd string
	}
	boolean := false

	if len(firstName) == 0 || len(name) == 0 {
		errorStruct.Name = "Vous devez entrer un prénom ou nom valide."
		boolean = true
	}
	if len(user.Mail) == 0 || !isAValidEmail(user.Mail) {
		errorStruct.Mail = "Vous devez entrer une adresse email valide."
		boolean = true
	} else {
		db, err := database.ConnectDatabase("./database.db")
		if err != nil {
			log.Fatal(err)
		}
		emails, err := database.GetUserEmails(db)
		if err != nil {
			log.Fatal(err)
		}
		err = database.DisconnectDatabase(db)
		if err != nil {
			log.Fatal(err)
		}
		if IsInStringArray(emails, user.Mail) {
			errorStruct.Mail = "Email déja utilisé."
			boolean = true
		}
	}
	if len(passwd) < 9 {
		errorStruct.Passwd = "Vous devez entrer un mot de passe valide. (8 caractères minimum)"
		boolean = true
	}
	return boolean, errorStruct
}

func IsInStringArray(array []string, str string) bool {
	for i := 0; i < len(array); i++ {
		if array[i] == str {
			return true
		}
	}
	return false
}

func insertInArray(array []structs.FullComment, element structs.FullComment, index int) []structs.FullComment {
	var result []structs.FullComment
	for i := 0; i < index; i++ {
		result = append(result, array[i])
	}
	result = append(result, element)
	for i := index; i < len(array); i++ {
		result = append(result, array[i])
	}
	return result
}
