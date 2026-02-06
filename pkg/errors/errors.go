package errors

import (
	"errors"
	"fmt"
	"strings"
)

type GeelatoError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Err     error  `json:"-"`
}

func (e *GeelatoError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *GeelatoError) Unwrap() error {
	return e.Err
}

func (e *GeelatoError) WithDetails(details string) *GeelatoError {
	e.Details = details
	return e
}

func (e *GeelatoError) WithErr(err error) *GeelatoError {
	e.Err = err
	return e
}

func New(code Code, details ...string) *GeelatoError {
	msg := Message(code)
	var detail string
	if len(details) > 0 {
		detail = strings.Join(details, ": ")
	}
	return &GeelatoError{
		Code:    code,
		Message: msg,
		Details: detail,
	}
}

func Wrap(err error, code Code, details ...string) *GeelatoError {
	if ge, ok := err.(*GeelatoError); ok {
		return ge
	}
	msg := Message(code)
	var detail string
	if len(details) > 0 {
		detail = strings.Join(details, ": ")
	}
	return &GeelatoError{
		Code:    code,
		Message: msg,
		Details: detail,
		Err:     err,
	}
}

func Is(err error, code Code) bool {
	var ge *GeelatoError
	if errors.As(err, &ge) {
		return ge.Code == code
	}
	return false
}

func Equal(err error, target error) bool {
	return errors.Is(err, target)
}

func FromStd(err error) *GeelatoError {
	return &GeelatoError{
		Code:    ErrUnknown,
		Message: err.Error(),
		Err:     err,
	}
}
