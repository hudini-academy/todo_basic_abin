package mysql

import (
	"TODO/pkg/models"
	"database/sql"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, hashed_password string) error {
	stmt := `INSERT INTO users(name,email,hashed_password,created)
	VALUES(?,?,?,utc_timestamp())`
	_, err := m.DB.Exec(stmt, name, email, hashed_password)
	if err != nil {
		return err
	}
	return nil
}
func (m *UserModel) Authenticate(email, password string) (bool, error) {
	stmt := `select id from users where email = ? and hashed_password =?`
	rows, err := m.DB.Query(stmt, email, password)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer rows.Close()
	return true, nil

}
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
