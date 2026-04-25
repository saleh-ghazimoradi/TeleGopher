package repository

import "errors"

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrEmailExists          = errors.New("email already exists")
	ErrSameUser             = errors.New("cannot create private conversation with the same user")
	ErrPrivateAlreadyExists = errors.New("private conversation already exists")
)
