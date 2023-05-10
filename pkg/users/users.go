package users

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/iammegalith/tvbbs/pkg/utils"
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

// Login handles the login workflow
func Login(conn net.Conn, db *sql.DB) (*User, error) {
	utils.ShowText(conn, "ascii/prelog.txt")
	var username string
	var password string
	var user *User
	var err error

	for i := 0; i < 3; i++ {
		fmt.Fprint(conn, utils.Strings("prelog_prompt"))
		fmt.Fscanln(conn, &username)

		if strings.ToLower(username) == "new" {
			return RegisterUser(conn, db)
		}

		user, err = GetUserByName(db, username)
		if err != nil {
			fmt.Fprintln(conn, "User not found. Try again.")
			continue
		}

		fmt.Fprint(conn, "Password: ")
		fmt.Fscanln(conn, &password)

		if password == user.Password {
			return user, nil
		}
		fmt.Fprintln(conn, "Password incorrect. Try again.")
	}

	utils.ShowText(conn, "ansi/later.txt")
	return nil, errors.New("failed login attempts")
}

// RegisterUser handles the user registration workflow
func RegisterUser(conn net.Conn, db *sql.DB) (*User, error) {
	var username, password, confirmPassword string

	for {
		fmt.Fprint(conn, "Enter a new username: ")
		fmt.Fscanln(conn, &username)

		_, err := GetUserByName(db, username)
		if err == nil {
			fmt.Fprintln(conn, "Username already exists. Try another one.")
			continue
		}

		break
	}

	for {
		fmt.Fprint(conn, "Enter a password: ")
		fmt.Fscanln(conn, &password)

		fmt.Fprint(conn, "Confirm password: ")
		fmt.Fscanln(conn, &confirmPassword)

		if password == confirmPassword {
			break
		}

		fmt.Fprintln(conn, "Passwords do not match. Try again.")
	}

	// Add user to the database and return the new user
	newUser := &User{
		Username:    username,
		Password:    password,
		AnsiEnabled: true,
		Active:      true,
		Level:       1, // Set the appropriate level for new users
		CreatedOn:   time.Now(),
		LastLogin:   time.Now(),
		RulesRead:   false,
		Echo:        false,
		Email:       "",
		AllowEmail:  false,
	}

	err := AddUser(db, newUser)
	if err != nil {
		return nil, err
	}

	// Assuming menus.MainMenu exists and takes a net.Conn and *sql.DB as parameters
	// menus.MainMenu(conn, db)

	return newUser, nil
}

// AddUser adds a new user to the database
func AddUser(db *sql.DB, user *User) error {
	query := `INSERT INTO users (username, password, ansi_enabled, active, level, created_on, last_login, rules_read, echo, email, allow_email) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, user.Username, user.Password, user.AnsiEnabled, user.Active, user.Level, user.CreatedOn, user.LastLogin, user.RulesRead, user.Echo, user.Email, user.AllowEmail)
	return err
}

