package handlers

import "errors"

const (
	StatusOK  = "OK"
	StatusERR = "ERROR"
)

var (
	ErrNoError  = errors.New("No error")
	ErrDecode   = errors.New("failed to decode request body")
	ErrValidate = errors.New("failed to validate body")
)
