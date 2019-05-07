package gotry

import "errors"

var (
	// ErrTimeout is used when the retrying process reaches the Timeout without
	// succeeding.
	ErrTimeout = errors.New("timeout occurred")
	// ErrMaxTriesReached is used when the retrying process does not succeed
	// before MaxTries is reached.
	ErrMaxTriesReached = errors.New("max tries reached")
)

// IsTimeout returns true if the cause of the given error is a ErrTimeout.
func IsTimeout(err error) bool {
	return err == ErrTimeout
}

// IsMaxTriesReached returns true if the cause of the given error is a ErrMaxTriesReached.
func IsMaxTriesReached(err error) bool {
	return err == ErrMaxTriesReached
}
