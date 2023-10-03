package user

import (
	"crypto/sha512"
	"encoding/base64"
	"notes-manager/src/pkg/generator"
	"strings"
	"time"

	"github.com/google/uuid"
)

// getUserLoginKey generates a cache key for the user based on the login.
func getUserLoginKey(login string) string {
	return userLoginKeyBase + strings.ToLower(login)
}

// getUserKey generates a cache key for the user based on the user's ID.
func getUserKey(userId uuid.UUID) string {
	return userKeyBase + userId.String()
}

// generateHashedPassword returns the hashed version of the password using SHA-512 and the provided salt.
func generateHashedPassword(password, salt string) (string, error) {
	hasher := sha512.New()
	_, err := hasher.Write([]byte(password + salt))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

// New creates a new user instance with the provided login and password.
// It generates a salt, hashes the password with the salt, and sets default values for the user.
func New(login, password string) (*User, error) {
	// Generate a salt for password hashing
	salt := generator.GenerateString(8)

	// Hash the password using SHA-512
	hashedPassword, err := generateHashedPassword(password, salt)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now().UTC()

	return &User{
		Id:           uuid.New(),
		Login:        strings.ToLower(login),
		Password:     hashedPassword,
		Salt:         salt,
		Role:         DefaultRole,
		RegisteredAt: currentTime,
		LastLoginAt:  currentTime,
	}, nil
}

// TableName defines the table name for the User model when working with an ORM.
func (User) TableName() string {
	return "ns_users"
}

// ValidatePassword checks if the provided password matches the hashed password of the user.
func (u User) ValidatePassword(password string) (bool, error) {
	hashedPassword, err := generateHashedPassword(password, u.Salt)
	if err != nil {
		return false, err
	}
	return u.Password == hashedPassword, nil
}
