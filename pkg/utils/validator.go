package utils

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"regexp"
	"strings"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errors[e.Field()] = getErrorMessage(e)
		}
	}

	return errors
}

func getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	case "gte":
		return e.Field() + " must be greater than or equal to " + e.Param()
	case "lte":
		return e.Field() + " must be less than or equal to " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}

// CleanUsername removes special characters and spaces from username, keeping only alphanumeric and underscores
func CleanUsername(username string) string {
	// Convert to lowercase
	username = strings.ToLower(username)

	// Remove all non-alphanumeric characters except underscores
	reg := regexp.MustCompile("[^a-z0-9_]+")
	username = reg.ReplaceAllString(username, "")

	// Remove leading/trailing underscores
	username = strings.Trim(username, "_")

	// If username is empty after cleaning, return a default
	if username == "" {
		return "user"
	}

	return username
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}