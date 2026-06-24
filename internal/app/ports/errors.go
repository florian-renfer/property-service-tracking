package ports

import "errors"

// ErrPropertyNotFound indicates no property exists for the requested ID.
var ErrPropertyNotFound = errors.New("property not found")
