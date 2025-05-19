package err

import "errors"

var (
	ErrBodyParse    = errors.New("body parsing error")
	ErrParamsParse  = errors.New("parameter parsing error")
	ErrQueryParse   = errors.New("query parameter parsing error")
	ErrValidate     = errors.New("data validation error")
	ErrUnauthorized = errors.New("Invalid or expired key")
)
