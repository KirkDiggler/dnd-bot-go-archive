package errors

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
