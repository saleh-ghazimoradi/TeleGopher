package repository

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrDuplicateEmail       = errors.New("duplicate email")
	ErrPrivateNotFound      = errors.New("private conversation not found")
	ErrPrivateAlreadyExists = errors.New("private conversation already exists between these users")
	ErrSameUser             = errors.New("cannot create private chat with yourself")
)
