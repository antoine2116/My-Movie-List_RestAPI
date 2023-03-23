package users

import "errors"

var (
	ErrUserAlreadyExists   = errors.New("an user alreay exists with the same email")
	ErrBadCredentials      = errors.New("invalid email or password")
	ErrGoogleAuthFailed    = errors.New("an error occured during Google authentication")
	ErrGitHubAuthFailed    = errors.New("an error occured during GitHub authentication")
	ErrMissingOAuthCode    = errors.New("missing code")
	ErrInvalidToken        = errors.New("missing token or invalid format")
	ErrMismatchedPasswords = errors.New("passwords do not match")
)
