package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// custom validation function to ensure the number of decimal
// places is not greater than the specified number
func maxDecimals(fl validator.FieldLevel) bool {
	decimalPlaces, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}

	var field string
	switch fl.Field().Kind() {
	case reflect.String:
		field = fl.Field().String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Float32, reflect.Float64:
		// convert to string
		field = fmt.Sprintf("%."+strconv.Itoa(decimalPlaces+1)+"f", fl.Field().Float())
	default:
		return false
	}
	// check decimal places
	parts := strings.Split(field, ".")
	if len(parts) < 2 {
		return true
	}
	return len(parts[1]) <= decimalPlaces
}

// coustom validation function to ensure the password is valid
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) == 0 {
		return true
	}
	// check password length is between 8 and 20
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	// check if password contains at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// check if password contains at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// check if password contains at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasNumber
}
