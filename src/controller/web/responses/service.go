package responses

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

func New(code int, message, developer interface{}) error {
	err := Error{
		Code:      code,
		Message:   message,
		Developer: developer,
	}

	encoded, _ := json.Marshal(&err)
	return errors.New(string(encoded))
}

func BadRequest(message, developer interface{}) error {
	return New(http.StatusBadRequest, message, developer)
}

func System(message, developer interface{}) error {
	return New(http.StatusInternalServerError, internalErrorMessage, developer)
}

func NotAuthorized() error {
	return New(http.StatusUnauthorized, "Вам необходимо авторизоваться", "")
}

func Permissions(scope interface{}) error {
	return New(http.StatusForbidden, fmt.Sprintf("У вас недостаточно прав для того, чтобы %s", scope), "")
}
