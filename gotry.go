package gotry

import (
	"time"
)

// Try to execute the given function f and repeat if f returns an error.
// Results and Stop Errors are sent to the given resultChannel.
// Decide for yourself whether you want to run Try in a goroutine.
//
// There are multiple RetryOptions available:
//
// * Delay            = 0s   How long to wait between calls to f.
// * MaxTries:        = 5    How many errors returned from f should be tolerated before giving up.
// * Timeout:         = 5s   How long overall the retries are allowed to take before aborting.
// * AfterRetry:      = nil  A function that is called with the resulting error after every failing call to f.
// * AfterRetryLimit: = nil  A function that is called with the latest error after reaching MaxTries.
// * AfterTimeout:    = nil  A function that is called with the latest error after reaching Timeout.
func Try(
	f func() (interface{}, error),
	resultChannel chan *RetryResult,
	options ...RetryOption,
) {
	retryOptions := NewRetryOptionsWithDefault(options...)
	tryCount := 0
	var lastError error
	var timeout <-chan time.Time

	if retryOptions.Timeout > 0 {
		timeout = time.After(retryOptions.Timeout)
	}

	for {
		if isMaxTriesReached(retryOptions, tryCount) {
			handleMaxTries(retryOptions, resultChannel, lastError)
			return
		}
		if isTimedOut(retryOptions, timeout) {
			handleTimeout(retryOptions, resultChannel, lastError)
			return
		}

		tryCount++
		value, err := f()
		if err != nil {
			lastError = err
			if retryOptions.AfterRetry != nil {
				retryOptions.AfterRetry(lastError)
			}
			time.Sleep(retryOptions.Delay)
			continue
		}
		resultChannel <- &RetryResult{
			Value: value,
		}
		return
	}
}

func isTimedOut(options *RetryOptions, timeout <-chan time.Time) bool {
	select {
	case <-timeout:
		return true
	default:
		return false
	}
}

// HandleTimeout by calling the callback if it exists and then sending an error
// to the resultChannel.
func handleTimeout(options *RetryOptions, resultChannel chan *RetryResult, lastError error) {
	if options.AfterTimeout != nil {
		options.AfterTimeout(lastError)
	}
	sendError(resultChannel, lastError, ErrTimeout)
}

func isMaxTriesReached(options *RetryOptions, tryCount int) bool {
	return tryCount >= options.MaxTries
}

// HandleMaxTries by calling the callback if it exists and then sending an error
// to the resultChannel.
func handleMaxTries(options *RetryOptions, resultChannel chan *RetryResult, lastError error) {
	if options.AfterRetryLimit != nil {
		options.AfterRetryLimit(lastError)
	}
	sendError(resultChannel, lastError, ErrMaxTriesReached)
}

// SendError sends an error in a RetryResult to the resultChannel. The error
// contains information about why the retrying was stopped and what error the
// retried operation last returned.
func sendError(resultChannel chan *RetryResult, lastError error, stopReason error) {
	resultChannel <- &RetryResult{
		StopReason: stopReason,
		LastError:  lastError,
	}
}
