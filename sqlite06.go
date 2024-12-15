package sqlite06

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var (
	Filename = ""
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

func openConnection() (*sql.DB, error) {
	db, e := sql.Open("sqlite3", Filename)
	if e != nil {
		return nil, e
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, e := openConnection()
	if e != nil {
		fmt.Println(e)
		return -1
	}
	defer db.Close()

	userID := -1
	stmt := fmt.Sprintf(`SELECT ID FROM users WHERE username = '%s'`, username)
	rows, e := db.Query(stmt)
	defer rows.Close()

	for rows.Next() {
		var id int
		e = rows.Scan(&id)
		if e != nil {
			fmt.Println("exists() Scan", e)
			return -1
		}
		userID = id
	}
	return userID
}

func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)

	db, e := openConnection()
	if e != nil {
		fmt.Println(e)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("user already exists:", d.Username)
		return -1
	}

	stmt := `INSERT INTO users VALUES (NULL, ?)`
	_, e = db.Exec(stmt, d.Username)
	if e != nil {
		fmt.Println(e)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	stmt = `INSERT INTO Userdata VALUES (?, ?, ?, ?)`
	_, e = db.Exec(stmt, userID, d.Name, d.Username, d.Description)
	if e != nil {
		fmt.Println("db.Exec()", e)
		return -1
	}

	return userID
}

func DeleteUser(id int) error {
	db, e := openConnection()
	if e != nil {
		return e
	}
	defer db.Close()

	stmt := fmt.Sprintf(`SELECT Username FROM users WHERE ID = %d`, id)
	rows, e := db.Query(stmt)
	defer rows.Close()

	var username string
	for rows.Next() {
		e = rows.Scan(&username)
		if e != nil {
			return e
		}
	}

	if exists(username) != id {
		return fmt.Errorf("user with ID %d does not exist", id)
	}

	stmt = `DELETE FROM userdata WHERE userid = ?`
	_, e = db.Exec(stmt, id)
	if e != nil {
		return e
	}

	stmt = `DELETE FROM users WHERE id = ?`
	_, e = db.Exec(stmt, id)
	if e != nil {
		return e
	}

	return nil
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, e := openConnection()
	if e != nil {
		return nil, e
	}
	defer db.Close()

	rows, e := db.Query(`SELECT id, username, name, surname, description FROM users, userdata WHERE users.id = userdata.userid`)
	if e != nil {
		return Data, e
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		e = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
		if e != nil {
			return nil, e
		}
	}
	return Data, nil
}

func UpdateUser(d Userdata) error {
	db, e := openConnection()
	if e != nil {
		return e
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("user does not exist")
	}
	d.ID = userID
	stmt := `UPDATE userdata SET name = ?, Surname = ?, Description = ? WHERE userid = ?`
	_, e = db.Exec(stmt, d.Name, d.Surname, d.Description, d.ID)
	if e != nil {
		return e
	}

	return nil
}
