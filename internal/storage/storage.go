package storage

import "errors"

var (
	ErrBadOperation = errors.New("wrong operation type")
	ErrBalance      = errors.New("not enough money")
)
