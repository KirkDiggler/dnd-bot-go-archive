package internal

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrTabEmpty       = errors.New("tab is empty, go home you're drunk")
)
