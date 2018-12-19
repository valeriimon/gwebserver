package mainservice

import (
	"encoding/json"
	"errors"

	uuid "github.com/satori/go.uuid"
)

// User : user model
type User struct {
	Firstname    string        `json:"firstname,omitempty"`
	Lastname     string        `json:"lastname,omitempty"`
	Email        string        `json:"email,omitempty"`
	Password     string        `json:"password,omitempty"`
	Applications []Application `json:"applications,omitempty"`
	ID           string        `json:"id,omitempty"`
}

// Create : create new user
func (user *User) Create() error {
	var database *Database
	uID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	user.ID = uID.String()
	user.Password = EncryptValueByHex(user.Password)

	userBuf, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	err = user.checkDuplicateEmail()
	if err != nil {
		return err
	}

	database.Save("users", user.Email, userBuf)

	return nil
}

func (user *User) checkDuplicateEmail() error {
	var database *Database
	result, err := database.Get("users", user.Email)

	// Ignore error if bucket is undefined
	if err != nil && err.Error() != "Bucket Not Found" {
		return err
	}
	if result == nil {
		return nil
	}

	return errors.New("Error duplicate email")
}

// Authenticate : Authenticate user via email and password
func (user *User) Authenticate() (*User, error) {
	var database *Database

	user.Password = EncryptValueByHex(user.Password)

	u, err := database.Auth("users", user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return u, nil
}
