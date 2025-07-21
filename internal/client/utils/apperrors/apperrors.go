package apperrors

import (
	"errors"
)

var ErrFileNotFound = errors.New("file not found")
var ErrTokenNotFound = errors.New("token not found")
