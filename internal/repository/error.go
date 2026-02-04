package repository

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrDuplicateEmail       = errors.New("duplicate email")
	ErrPrivateAlreadyExists = errors.New("private chat already exists")
	ErrSameUser             = errors.New("cannot create private chat with same user")
)
