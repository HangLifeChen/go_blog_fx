package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func BindValidator() {
	// bind validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("max_decimals", maxDecimals)
		v.RegisterValidation("password", validatePassword)
	}
}
