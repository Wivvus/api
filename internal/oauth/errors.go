package oauth

import "errors"

var (
	ErrNoUserInSession = errors.New("user not found in session")
)
