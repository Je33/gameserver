package domain

import "github.com/pkg/errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrConversion  = errors.New("conversion error")
	ErrNoDocuments = errors.New("no documents")
	ErrConfig      = errors.New("config error")
	ErrSignature   = errors.New("signature error")
)
