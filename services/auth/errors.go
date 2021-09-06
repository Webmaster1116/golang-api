package auth

import "errors"

var ErrNoTokenProvided = errors.New("token is not provided")
var ErrTokenInvalid = errors.New("token is invalid")
var ErrTokenClaimsInvalid = errors.New("token claims are invalid")
var ErrTokenUserNotFound = errors.New("token user isn't found")
