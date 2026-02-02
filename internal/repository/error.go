package repository

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
)
