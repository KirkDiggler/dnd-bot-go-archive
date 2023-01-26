package dnderr

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
