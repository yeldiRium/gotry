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
		if reason := shouldAbort(tryCount, timeout, retryOptions); reason != nil {
			switch reason {
			case ErrTimeout:
				if retryOptions.AfterTimeout != nil {
					retryOptions.AfterTimeout(lastError)
				}
			case ErrMaxTriesReached:
				if retryOptions.AfterRetryLimit != nil {
					retryOptions.AfterRetryLimit(lastError)
				}
			}
			resultChannel <- &RetryResult{
				StopReason: reason,
				LastError:  lastError,
			}
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

// Checks whether MaxTries or Timeout was reached and returns fitting error or nil, if all is well.
func shouldAbort(tryCount int, timeout <-chan time.Time, options *RetryOptions) error {
	if tryCount >= options.MaxTries {
		return ErrMaxTriesReached
	}

	select {
	case <-timeout:
		return ErrTimeout
	default:
		return nil
	}
}
