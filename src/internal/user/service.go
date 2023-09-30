package user

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"notes-manager/src/pkg/generator"
	"strings"
	"time"

	"github.com/google/uuid"
)

func getUserLoginKey(login string) string {
	return fmt.Sprintf("ns:users:login:%s", login)
}

func getUserKey(userId uuid.UUID) string {
	return fmt.Sprintf("ns:users:%s", userId)
}

func New(login, password string) *User {
	// Генерируем соль для пароля
	salt := generator.GenerateString(8)

	// Хешируем пароль
	hasher := sha512.New()
	hasher.Write([]byte(password + salt))
	hashedPassword := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	currentTime := time.Now().UTC()

	return &User{
		Id:           uuid.New(),
		Login:        strings.ToLower(login),
		Password:     string(hashedPassword),
		Salt:         salt,
		Role:         DEFAULT_ROLE,
		RegisteredAt: currentTime,
		LastLoginAt:  currentTime,
	}
}

func (User) TableName() string {
	return "ns_users"
}

func (u User) ValidatePassword(password string) bool {
	hasher := sha512.New()
	hasher.Write([]byte(password + u.Salt))
	hashedPassword := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return u.Password == hashedPassword
}
