package database

import (
	"database/sql"
	"log"
	"src/structs"
	"strconv"
)

// AddCategory : Add a category to database (database, category structure)
func AddCategory(db *sql.DB, category structs.Category) error {
	_, err := db.Exec("INSERT INTO category (categoryID, name, color) VALUES (?, ?)", category.ID, category.Name, category.Color)
	if err != nil {
		return err
	}
	return nil
}

// RemoveCategory : Remove category from database (database, category id)
func RemoveCategory(db *sql.DB, ID int) error {
	_, err := db.Exec("DELETE FROM category WHERE categoryID = ?", ID)
	if err != nil {
		return err
	}
	return nil
}

// GetCategory : Get comment in database (database, comment id)
func GetCategory(db *sql.DB, ID int) (structs.Category, error) {
	var category structs.Category
	var empty structs.Category
	rows, err := db.Query("select categoryID, name, color from categories where categoryID = ?", ID)
	if err != nil {
		log.Fatal(err)
		return empty, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&category.ID, &category.Name, &category.Color)
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
	return category, nil
}

// GetAllCategories : Return an array of category structure with all categories in database (db)
func GetAllCategories(db *sql.DB) ([]structs.Category, error) {
	var categories []structs.Category
	rows, err := db.Query("select categoryID from categories")
	if err != nil {
		log.Fatal(err)
		return make([]structs.Category, 0), err
	}
	defer rows.Close()

	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Category, 0), err
		}

		var category structs.Category
		value, err := strconv.Atoi(temp)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Category, 0), err
		}

		category, err = GetCategory(db, value)
		if err != nil {
			log.Fatal(err)
			return make([]structs.Category, 0), err
		}
		categories = append(categories, category)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return make([]structs.Category, 0), err
	}

	return categories, nil
}

func GetCategoryNameById(db *sql.DB, id string) (string, error) {
	rows, err := db.Query("select categoryID, name from categories")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var ID int
		err = rows.Scan(&ID, &name)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

		testID, err := strconv.Atoi(id)
		if err != nil {
			log.Fatal(err)
			return "", err
		}

		if testID == ID {
			return name, nil
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return "Categories", nil
}
