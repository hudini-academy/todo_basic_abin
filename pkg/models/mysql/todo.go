package mysql

import (
	"TODO/pkg/models"
	"database/sql"
	"strings"
)

// Define a TodoModel type which wraps a sql.DB connection pool.
type ToDoModel struct {
	DB *sql.DB
}

type SpclToDoModel struct {
	DB *sql.DB
}

// This will insert a new todo into the database.

func (m *ToDoModel) Insert(title string) (int, error) {
	// Write the SQL statement we want to execute.
	// Entering into both databases
	stmt := `INSERT INTO todo (title, created, expires) 
    VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, 365)
	if err != nil {
		return 0, nil
	}
	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SpclToDoModel) Insert(title string) (int, error) {
	if strings.HasPrefix(title, "Special:") {
		// Write the SQL statement we want to execute.
		// Entering into both databases
		stmt := `INSERT INTO spcl_todo (title, created, expires) 
    VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
		result, err := m.DB.Exec(stmt, title, 365)
		if err != nil {
			return 0, nil
		}
		// Use the LastInsertId() method on the result object to get the ID of our
		// newly inserted record in the snippets table.
		id, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		return int(id), nil
	}
	return 0, nil
}

// This will return a specific todo based on its id.
func (m *ToDoModel) Delete(id int) error {
	stmt := `Delete from todo where id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// This will return all the created todo.
func (m *ToDoModel) GetAll() ([]*models.ToDo, error) {
	stmt := `SELECT * from todo`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the models.ToDo objects.
	todo := []*models.ToDo{}
	for rows.Next() {
		// Create a pointer to a new zeroed ToDo struct.
		s := &models.ToDo{}
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of todo.
		todo = append(todo, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the todo slice.
	return todo, nil
}
func (m *ToDoModel) Update(id int, title string) error {
	//Update the DB using sql queries
	stmt := `Update todo set title =? WHERE id=?`
	_, err := m.DB.Exec(stmt, title, id)
	if err != nil {
		return err
	}
	return nil
}

// For Special task
func (m *SpclToDoModel) SpclTaskAdd(title string) (int, error) {
	// Write the SQL statement we want to execute.
	stmt := `INSERT INTO spcl_todo (heading, born, dead) 
    VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, 365)
	if err != nil {
		return 0, nil
	}
	// Use the LastInsertId
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}
func (m *SpclToDoModel) SpclTaskDelete(id int) error {
	stmt := `Delete from spcl_todo where id = ?`
	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}
func (m *SpclToDoModel) SpclTaskGetAll() ([]*models.ToDo, error) {
	stmt := `SELECT * from spcl_todo`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the models.ToDo objects.
	spcltodo := []*models.ToDo{}
	for rows.Next() {
		// Create a pointer to a new zeroed ToDo struct.
		s := &models.ToDo{}
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of todo.
		spcltodo = append(spcltodo, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// Return the spcltodo slice.
	return spcltodo, nil
}
func (m *SpclToDoModel) SpclGetAll() ([]*models.ToDo, error) {
	stmt := `SELECT * from spcl_todo`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the models.ToDo objects.
	todo := []*models.ToDo{}
	for rows.Next() {
		// Create a pointer to a new zeroed ToDo struct.
		s := &models.ToDo{}
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		todo = append(todo, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todo, nil
}
