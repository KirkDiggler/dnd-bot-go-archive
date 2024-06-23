package internal

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrAlreadyExists  = errors.New("record already exists")
	ErrTabEmpty       = errors.New("tab is empty, go home you're drunk")
)
