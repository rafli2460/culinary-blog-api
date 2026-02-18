package logger

import (
	"errors"

	"github.com/rs/zerolog/log"
)

// LogError logs the error with zerolog and returns the original error.
// Use this for infrastructure/database errors that need to be logged but bubbled up.
func LogError(err error, message string) error {
	log.Error().Err(err).Msg(message)
	return err
}

// LogErrorWithFields logs the error with additional fields and returns the original error.
func LogErrorWithFields(err error, message string, fields map[string]interface{}) error {
	event := log.Error().Err(err)
	for k, v := range fields {
		event.Interface(k, v)
	}
	event.Msg(message)
	return err
}

// ValidationError logs a warning and returns a new error with the message.
// Use this for business logic validation failures (e.g. invalid input).
func ValidationError(message string) error {
	log.Warn().Msg(message)
	return errors.New(message)
}

// SystemError logs an error and returns a new error with the message.
// Use this when you want to create a new error from scratch but log it as an error.
func SystemError(message string) error {
	log.Error().Msg(message)
	return errors.New(message)
}
