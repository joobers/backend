package errors

import errorslib "errors"

var (
	InternalServerError = errorslib.New("Internal server error")
	InvalidToken        = errorslib.New("Invalid token")
)
