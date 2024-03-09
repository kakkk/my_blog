package consts

import (
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrHasExist       = errors.New("key has exist")
	ErrExpired        = errors.New("expired")
)
