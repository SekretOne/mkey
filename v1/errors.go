package mkey

import "fmt"

// MarshalError returns when a problem occurs when converting to the multi-field string format.
//
// Causes include:
//   - bad tagging instructions
//   - exported subfield has no marshalling instructions.
type MarshalError struct {
	message string
	source  error
}

func (m MarshalError) Error() string {
	return m.message
}

func newMarshalError(msg string, srcerr error) MarshalError {
	return MarshalError{message: msg, source: srcerr}
}

func (m MarshalError) Is(err error) bool {
	return m.source == err || m.message == err.Error()
}

func (m MarshalError) Unwrap() error {
	return m.source
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

func newUnmarshalError(msg string, srcerr error) UnmarshalError {
	return UnmarshalError{message: msg, source: srcerr}
}
