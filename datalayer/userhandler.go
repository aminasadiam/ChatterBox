package datalayer

import (
	"database/sql"
	"time"
)

type User struct {
	Id           uint
	Username     string
	Email        string
	Password     string
	ActiveCode   string
	RegisterDate time.Time
	IsActive     bool
	IsDelete     bool
	IsAdmin      bool
}

func (handler *SQLhandler) GetUsers() ([]User, error) {
	rows, err := handler.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		u, err := getRowsDataUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (handler *SQLhandler) GetUserById(userId int) (User, error) {
	row := handler.db.QueryRow("SELECT * FROM users WHERE id = ?", userId)
	return getRowDataUser(row)
}

func (handler *SQLhandler) GetUserByEmail(email string) (User, error) {
	row := handler.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	return getRowDataUser(row)
}

func (handler *SQLhandler) AddUser(user User) error {
	_, err := handler.db.Exec("INSERT INTO users (username, email, password, activecode, registerdate, isactive, isdelete, isadmin) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.Username, user.Email, user.Password, user.ActiveCode, user.RegisterDate, user.IsActive, user.IsDelete, user.IsAdmin)
	return err
}

func (handler *SQLhandler) ActiveUser(activeCode string) error {
	_, err := handler.db.Exec("UPDATE users SET isactive = 1 WHERE activecode = ?", activeCode)
	return err
}

func getRowDataUser(row *sql.Row) (User, error) {
	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.ActiveCode, &user.RegisterDate, &user.IsActive, &user.IsDelete, &user.IsAdmin)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func getRowsDataUser(rows *sql.Rows) (User, error) {
	var user User
	err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.ActiveCode, &user.RegisterDate, &user.IsActive, &user.IsDelete, &user.IsAdmin)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
