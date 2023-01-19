package dnderr

type missingParameterError struct {
	Param string
}

func (e *missingParameterError) Error() string {
	return "Missing parameter: " + e.Param
}

func NewMissingParameterError(param string) error {
	return &missingParameterError{
		Param: param,
	}
}

type notFoundError struct {
	Msg string
}

func (e *notFoundError) Error() string {
	return e.Msg
}

func NewNotFoundError(msg string) error {
	return &notFoundError{
		Msg: msg,
	}
}
