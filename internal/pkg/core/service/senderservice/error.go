package senderservice

import "errors"

// This group of constants define the error strings for this package, be it formatted or not
var (
	ErrorExistsButNotConvertible = errors.New("there exists an error but it is not convertible to string")
)
