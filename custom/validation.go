package custom

import "github.com/go-playground/validator/v10"

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than" + " " + fe.Param() + " characters"
	case "gte":
		return "Should be more than" + " " + fe.Param() + " characters"
	case "min":
		return "Should be more than" + " " + fe.Param() + " characters"
	case "max":
		return "Should not be more than" + " " + fe.Param() + " characters"
	case "email":
		return "This field must be a valid E-mail address"
	default:
		return "Unknown error"
	}
}
