package errors

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation failed")
	ErrConflict     = errors.New("resource already exists")
	ErrInternal     = errors.New("internal server error")
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user already exists")
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrProjectExists   = errors.New("project already exists")
)
