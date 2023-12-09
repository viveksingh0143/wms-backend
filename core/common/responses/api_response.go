package responses

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"star-wms/core/utils"
	"star-wms/core/validation"
)

type PaginationMeta struct {
	CurrentPage  int   `json:"currentPage"`
	PageSize     int   `json:"pageSize"`
	TotalRecords int64 `json:"totalRecords"`
	TotalPages   int   `json:"totalPages"`
}

type APIResponse struct {
	Status      int             `json:"status"`
	Message     string          `json:"message"`                // e.g., "Operation successful" or "An error occurred"
	Data        interface{}     `json:"data"`                   // This will hold the actual response data, e.g., a list of permissions or a single permission
	Metadata    *PaginationMeta `json:"metadata,omitempty"`     // This can be nil if there's no pagination info to provide
	Errors      []string        `json:"errors"`                 // A slice to hold validation or other errors
	FieldErrors []*IError       `json:"field_errors,omitempty"` // A slice to hold validation field errors
	Statistics  interface{}     `json:"statistics"`             // This will hold the statistics data
}

func NewSuccessResponse(status int, message string) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
	}
}

func NewSuccessDataResponse(status int, message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(status int, message string, err error) APIResponse {
	var iError *IError
	if errors.As(err, &iError) {
		return APIResponse{
			Status:      http.StatusBadRequest,
			Message:     "Please provide valid/required values",
			FieldErrors: []*IError{iError},
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return APIResponse{
			Status:  http.StatusNotFound,
			Message: "No record found...",
			Errors:  []string{err.Error()},
		}
	}
	return APIResponse{
		Status:  status,
		Message: message,
		Errors:  []string{err.Error()},
	}
}

func NewValidationErrorResponse(err error, anyType interface{}) APIResponse {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fieldErrors := GetAllErrors(err, anyType)
		return APIResponse{
			Status:      http.StatusBadRequest,
			Message:     "Please provide valid/required values",
			FieldErrors: fieldErrors,
		}
	}
	return NewErrorResponse(http.StatusBadRequest, "Failed to bind data", err)
}

func NewPageResponse(data interface{}, totalRecords int64, currentPage int, pageSize int) APIResponse {
	var totalPages int64 = 0
	if pageSize > 0 {
		totalPages = (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	}
	return APIResponse{
		Status:  http.StatusOK,
		Message: "Successfully records fetched",
		Data:    data,
		Metadata: &PaginationMeta{
			CurrentPage:  currentPage,
			PageSize:     pageSize,
			TotalRecords: totalRecords,
			TotalPages:   int(totalPages),
		},
	}
}

func GetAllErrors(err error, anyType interface{}) []*IError {
	t := reflect.TypeOf(anyType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	var errs []*IError
	for _, err := range err.(validator.ValidationErrors) {
		field, _ := t.FieldByName(err.Field())
		var errMsg string
		tag := err.Tag()

		if template, ok := validation.TagMessages[tag]; ok {
			if utils.HasFormatVerbs(template) {
				errMsg = fmt.Sprintf(template, err.Param())
			} else {
				errMsg = template
			}
		} else {
			errMsg = err.Error()
		}
		var el = &IError{}
		validationTag := field.Tag.Get("validationTag")
		if validationTag != "" {
			el.Field = validationTag
		} else {
			el.Field = field.Tag.Get("json")
		}
		el.Message = errMsg
		el.Value = err.Param()
		errs = append(errs, el)
	}
	return errs
}
