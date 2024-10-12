package mkey

import "fmt"

func (m MarshalError) Error() string {
	if m.source == nil {
		return fmt.Sprintf("cannot marshal %T: %v", m.message, m.message)
	}
	return m.message
}

func newMarshalError(msg string, source error) MarshalError {
	return MarshalError{message: msg, source: source}
}

type UnmarshalError struct {
	message string
	source  error
}

func (u UnmarshalError) Error() string {
	if u.source == nil {
		return fmt.Sprintf("%v: %v", u.message, u.message)
	}
	return u.message
}

func newUnmarshalError(msg string, source error) UnmarshalError {
	return UnmarshalError{message: msg, source: source}
}
