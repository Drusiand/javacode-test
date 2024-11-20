package app

import "errors"

var (
	ErrApply     = errors.New("failed to apply operation")
	ErrGetAmount = errors.New("failed to get amount")
)
