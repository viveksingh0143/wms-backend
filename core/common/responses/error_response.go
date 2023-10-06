package responses

import "fmt"

type IError struct {
	Field   string      `json:"field,omitempty"`
	Message string      `json:"message"`
	Value   interface{} `json:"value,omitempty"`
}

func (e *IError) Error() string {
	return fmt.Sprintf("Validation error: %s", e.Message)
}

func NewInputError(field string, msg string, value interface{}) error {
	return &IError{
		Field:   field,
		Message: msg,
		Value:   value,
	}
}
