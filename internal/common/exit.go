package common

import "errors"

type Exit struct {
	ID  string
	Err error
}

var ErrFatal = errors.New("fatal")
