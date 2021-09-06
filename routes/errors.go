package routes

import "errors"

var ErrVerificationTemplateNotFound = errors.New("cannot find verification email template")
var ErrResetPasswordTemplateNotFound = errors.New("cannot find reset password email template")
