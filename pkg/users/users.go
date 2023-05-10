package users

import (
	"database/sql"
	"errors"
	"time"
)

// User struct representing the user data
type User struct {
	ID          int
	Username    string
	Password    string
	AnsiEnabled bool
	Active      bool
	Level       int
	CreatedOn   time.Time
	LastLogin   time.Time
	RulesRead   bool
	Echo        bool
	Email       string
	AllowEmail  bool
}

// GetUserByName retrieves a user from the database by their username
func GetUserByName(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, password, ansi_enabled, active, level, created_on, last_login, rules_read, echo, email, allow_email FROM users WHERE username = ?`
	row := db.QueryRow(query, username)

	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.AnsiEnabled,
		&user.Active,
		&user.Level,
		&user.CreatedOn,
		&user.LastLogin,
		&user.RulesRead,
		&user.Echo,
		&user.Email,
		&user.AllowEmail,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// ToggleAnsi toggles the ANSI setting for the specified user
func ToggleAnsi(db *sql.DB, username string) error {
	user, err := GetUserByName(db, username)
	if err != nil {
		return err
	}

	user.AnsiEnabled = !user.AnsiEnabled

	return UpdateUser(db, user)
}

// UpdateUser updates the user information in the database
func UpdateUser(db *sql.DB, user *User) error {
	query := `UPDATE users SET ansi_enabled = ?, active = ?, level = ?, last_login = ?, rules_read = ?, echo = ?, email = ?, allow_email = ? WHERE id = ?`
	_, err := db.Exec(query, user.AnsiEnabled, user.Active, user.Level, user.LastLogin, user.RulesRead, user.Echo, user.Email, user.AllowEmail, user.ID)
	return err
}

// Add other functions as needed
