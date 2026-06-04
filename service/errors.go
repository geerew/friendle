package service

import "errors"

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

var (
	ErrRoundNotFound       = errors.New("round not found")
	ErrRoundNotActive      = errors.New("round not active")
	ErrRoundNotAwaitingWord = errors.New("round not awaiting word")
	ErrRoundNotFinished    = errors.New("round not finished")
	ErrNotPicker           = errors.New("not the picker")
	ErrPickerCannotGuess   = errors.New("picker cannot guess")
	ErrWordNotSet          = errors.New("word not set")
	ErrAlreadyFinished     = errors.New("already finished")
	ErrNoAttemptsLeft      = errors.New("no attempts left")
	ErrInvalidWordLength   = errors.New("word must be 5 letters")
	ErrInvalidAnswerWord   = errors.New("not a valid answer word")
	ErrInvalidGuessWord    = errors.New("not in word list")
	ErrNotSiteAdmin        = errors.New("not a site admin")
	ErrUserNotFound        = errors.New("user not found")
	ErrUsernameTaken       = errors.New("username already exists")
	ErrNoUpdateData        = errors.New("no data to update")
	ErrCredentialsRequired = errors.New("username and password required")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong     = errors.New("password must be no more than 128 characters")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidCurrentPassword = errors.New("invalid current password")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrLastAdmin           = errors.New("unable to delete the last admin user")
)
