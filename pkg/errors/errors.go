package errors

import "net/http"

type AppError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	IsDefault bool   `json:"is_default"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrBadRequest         = &AppError{Code: http.StatusBadRequest, Message: "Invalid request data"}
	ErrNotFound           = &AppError{Code: http.StatusNotFound, Message: "Resource not found"}
	ErrInternalServer     = &AppError{Code: http.StatusInternalServerError, Message: "Internal server error"}
	ErrInvalidIINLength   = &AppError{Code: http.StatusBadRequest, Message: "IIN must be exactly 12 digits"}
	ErrInvalidIINFormat   = &AppError{Code: http.StatusBadRequest, Message: "IIN must contain only numeric digits"}
	ErrInvalidIINChecksum = &AppError{Code: http.StatusBadRequest, Message: "Invalid IIN checksum"}
	ErrInvalidDateOfBirth = &AppError{Code: http.StatusBadRequest, Message: "Invalid date of birth in IIN"}
	ErrInvalidCenturyCode = &AppError{Code: http.StatusBadRequest, Message: "Invalid 7th digit in IIN"}
)
