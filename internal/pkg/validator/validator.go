package validator

import (
	"log/slog"
	"sync"

	validator "github.com/go-playground/validator/v10"
)

const (
	pin = "pin"
)

//nolint:gochecknoglobals // Singleton
var (
	once     sync.Once
	validate *validator.Validate
)

func GetInstance() *validator.Validate {
	if validate == nil {
		once.Do(func() {
			// register custom
			err := validator.New().RegisterValidation(pin, validatePin)
			if err != nil {
				slog.Error("error registering custom validation:", "error", err.Error())
				return
			}

			validate = validator.New()
		})
	}

	return validate
}
