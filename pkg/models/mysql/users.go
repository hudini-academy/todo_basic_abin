package mysql

import (
	"TODO/pkg/models"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email string, hashed_password []byte) error {
	stmt := `INSERT INTO users(name,email,hashed_password,created)
	VALUES(?,?,?,utc_timestamp())`
	_, err := m.DB.Exec(stmt, name, email, hashed_password)
	if err != nil {
		return err
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	log.Println(email, password)
	row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}
	log.Println("dadasd")
	// Check whether the hashed password and plain-text password provided match
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}

// stmt := `select id from users where email = ? and hashed_password =?`
// rows, err := m.DB.Query(stmt, email, password)
// if err != nil {
// 	log.Println(err)
// 	return false, err
// }
// defer rows.Close()
// return true, nil

func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
