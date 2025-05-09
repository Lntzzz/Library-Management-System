package xerror

import (
	"errors"
	"fmt"
	"net/http"
)

type Error interface {
	Code() int
	Message() string
	Error() string // as Status
	Details() interface{}
	SetCode(httpCode int) Error
	SetMessage(msg string) Error
	SetDetails(details interface{}) Error
	IsOK() bool
	Is3xx() bool
	Is4xx() bool
	Is5xx() bool
}

// all predictable errors

var (
	ErrUNKNOWN             Error = New5xxErr(http.StatusInternalServerError, "Unknown internal error")
	ErrInvalidArgument     Error = New4xxErr(http.StatusBadRequest, "Invalid argument")
	ErrForbidden           Error = New4xxErr(http.StatusForbidden, "Forbidden")
	ErrInvalidPin          Error = New4xxErr(http.StatusBadRequest, "Request header x-jdcloud-pin is empty")
	ErrPinNotFound         Error = New4xxErr(http.StatusNotFound, "Pin not found")
	ErrXJcloudPinIsEmpty   Error = New4xxErr(http.StatusBadRequest, "Request header x-jcloud-pin is empty.")
	ErrDecodeXJcloudPin    Error = New4xxErr(http.StatusBadRequest, "Decode request header x-jcloud-pin failed.")
	ErrInvalidUserInHeader Error = New4xxErr(http.StatusBadRequest, "User in request header invalid.")
	ErrUnsupportedRegion   Error = New4xxErr(http.StatusBadRequest, "Unsupported region")

	ErrBookNotFound     Error = New4xxErr(http.StatusNotFound, "Book not found")
	ErrBookConflict     Error = New4xxErr(http.StatusConflict, "Book conflict")
	ErrAddBookFailed    Error = New5xxErr(http.StatusInternalServerError, "Add book failed")
	ErrQueryBookFailed  Error = New5xxErr(http.StatusInternalServerError, "Query book failed")
	ErrDeleteBookFailed Error = New5xxErr(http.StatusInternalServerError, "Delete book failed")
	ErrUpdateBookFailed Error = New5xxErr(http.StatusInternalServerError, "Update book failed")
)

func Wrap(err error, message string) Error {
	if err == nil {
		return nil
	}
	return ErrUNKNOWN.SetMessage(fmt.Sprintf("%v:%v", message, err.Error()))
}

func WrapOpen(err OpenError, message string) Error {
	var old *OpenErrObj
	if !errors.As(err, &old) {
		return ErrUNKNOWN
	}
	if old.M != "" {
		message = fmt.Sprintf("%v:%v", message, old.M)
	} else {
		message = fmt.Sprintf("%v", message)
	}
	ne := NE(old.R, old.T, old.P, message, old.C, old.D)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}

func Unwrap2Open(err error) OpenError {
	if err == nil {
		return nil
	}
	if ec, ok := err.(ErrCode); ok {
		return ec
	}
	var old *OpenErrObj
	if errors.As(err, &old) {
		return old
	}
	return nil
}

func Is(err error, code ErrCode) bool {
	if err == nil {
		return false
	}
	if ec, ok := err.(ErrCode); ok {
		return ec.Code() == code.Code()
	}
	var old *OpenErrObj
	return errors.As(err, &old) && old.Code() == code.Code()
}

func Is4xx(err error) bool {
	if err == nil {
		return false
	}
	if ec, ok := err.(ErrCode); ok {
		return ec.Code() >= 400 && ec.Code() < 500
	}
	var old *OpenErrObj
	return errors.As(err, &old) && old.Code() >= 400 && old.Code() < 500
}

func NewErr(code int, message string) Error {
	if code >= 400 && code < 500 {
		return New4xxErr(code, message)
	}
	if code >= 500 && code < 600 {
		return New5xxErr(code, message)
	}
	ne := NE("", ErrNull, "", message, code, nil)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}

// New4xxErr creates a new 4xx (HTTP 4xx) error
func New4xxErr(httpStatusCode int, message string) Error {
	t := ErrInvalidArguments
	switch httpStatusCode {
	case http.StatusBadRequest:
		t = ErrInvalidArguments
	case http.StatusNotFound:
		t = ErrNotFound
	case http.StatusConflict:
		t = ErrConflict
	case http.StatusForbidden:
		t = ErrPermissionDenied
	}
	ne := NE("", t, "", message, httpStatusCode, nil)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}

// New5xxErr creates a new 5xx (HTTP 5xx) error
func New5xxErr(httpStatusCode int, message string) Error {
	ne := NE("", ErrInternal, "", message, httpStatusCode, nil)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}

type ErrCode struct {
	oe *OpenErrObj
}

func (e ErrCode) Error() string {
	return e.oe.Error()
}

func (e ErrCode) Code() int {
	return e.oe.Code()
}

func (e ErrCode) SetCode(code int) Error {
	var old *OpenErrObj
	errors.As(e.oe, &old)
	ne := NE(old.R, old.T, old.P, old.M, code, old.D)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}

func (e ErrCode) IsOK() bool {
	return e.oe.Code() == 200
}

func (e ErrCode) Is3xx() bool {
	return e.oe.Code() >= 300 && e.oe.Code() < 400
}

func (e ErrCode) Is4xx() bool {
	return e.oe.Code() >= 400 && e.oe.Code() < 500
}

func (e ErrCode) Is5xx() bool {
	return e.oe.Code() >= 500 && e.oe.Code() < 600
}

func (e ErrCode) Message() string {
	return e.oe.Message()
}

func (e ErrCode) SetMessage(msg string) Error {
	var old *OpenErrObj
	errors.As(e.oe, &old)
	ne := NE(old.R, old.T, old.P, msg, old.C, old.D)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}

func (e ErrCode) Details() interface{} {
	return e.oe.Details()
}

func (e ErrCode) SetDetails(details interface{}) Error {
	var old *OpenErrObj
	errors.As(e.oe, &old)
	ne := NE(old.R, old.T, old.P, old.M, old.C, details)
	var oe *OpenErrObj
	errors.As(ne, &oe)
	return ErrCode{oe: oe}
}
