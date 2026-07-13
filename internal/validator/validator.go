package validator

import (
	"os"

	"github.com/go-playground/validator/v10"
)

var Instance *validator.Validate

func init() {
	Instance = validator.New()

	Instance.RegisterValidation("dir_exists", validateDirExists)
}

func validateDirExists(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
