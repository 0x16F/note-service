package responses

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
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

func CustomErrorHandler() func(c *fiber.Ctx, err error) error {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		switch err.(type) {
		case *fiber.Error:
			code = fiber.StatusNotFound
			return c.Status(code).SendString(New(code, "Page not found", nil).Error())
		default:
			cErr := Error{}
			if err := json.Unmarshal([]byte(err.Error()), &cErr); err != nil {
				return c.Status(code).SendString(err.Error())
			}

			code = cErr.Code

			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			return c.Status(code).SendString(err.Error())
		}
	}
}
