package chatservice

import "errors"

// This function defines custom errors for this package
var (
	errorExistsButNotConvertible = errors.New("there exists an error but it is not convertible to string")
)
