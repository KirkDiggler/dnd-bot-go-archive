package dnderr

import "fmt"

type AlreadyExistsError struct {
	Msg string
}

func (e *AlreadyExistsError) Error() string {
	return e.Msg
}

func NewAlreadyExistsError(msg string) error {
	return &AlreadyExistsError{
		Msg: msg,
	}
}
type InvalidOperationError struct {
	Msg string
}

func (e *InvalidOperationError) Error() string {
	return e.Msg
}

func NewInvalidOperationError(msg string) error {
	return &InvalidOperationError{
		Msg: msg,
	}
}

type MissingParameterError struct {
	Param string
}

func (e *MissingParameterError) Error() string {
	return "Missing parameter: " + e.Param
}

func NewMissingParameterError(param string) error {
	return &MissingParameterError{
		Param: param,
	}
}

type NotFoundError struct {
	Msg string
}

func (e *NotFoundError) Error() string {
	return e.Msg
}

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		Msg: msg,
	}
}

type ResourceExhaustedError struct {
	Msg string
}

func (e *ResourceExhaustedError) Error() string {
	return e.Msg
}

func NewResourceExhaustedError(msg string) error {
	return &ResourceExhaustedError{
		Msg: msg,
	}
}

type InvalidEntityError struct {
	Msg string
}

func (e *InvalidEntityError) Error() string {
	return e.Msg
}

func NewInvalidEntityError(msg string) error {
	return &InvalidEntityError{
		Msg: msg,
	}
}

type InvalidParameterError struct {
	Param string
	Msg   string
}

func (e *InvalidParameterError) Error() string {
	return "Invalid parameter: " + e.Param + " - " + e.Msg
}

func NewInvalidParameterError(param string, msg interface{}) error {
	return &InvalidParameterError{
		Param: param,
		Msg:   fmt.Sprintf("%+v", msg),
	}
}
