package gotry

import (
	"time"
)

// Try to execute the given function f and repeat if f returns an error.
//
// There are multiple RetryOptions available:
//
// * Delay            = 0s   How long to wait inbetween calls to f.
// * MaxTries:        = 5    How many errors returned from f should be tolerated before giving up.
// * Timeout:         = 5s   How long overall the retries are allowed to take before aborting.
// * AfterRetry:      = nil  A function that is called with the resulting error after every failing call to f.
// * AfterRetryLimit: = nil  A function that is called with the latest error after reaching MaxTries.
// * AfterTimeout:    = nil  A function that is called with the latest error after reaching Timeout.
// * ReturnChannel:   = nil  A channel that replaces the otherwise created channel to which the resulting value in case of success is sent.
//
// Its return values are:
// * a chan RetryResult: The channel to which the result of a successful call to f is sent.
// * an error:           If the call to Try was misconfigured.
func Try(f func() (interface{}, error), options ...RetryOption) (<-chan *RetryResult, error) {
	if f == nil {
		return nil, ErrFIsMissing
	}
	retryOptions := NewRetryOptionsWithDefault(options...)
	resultChannel := make(chan *RetryResult)
	go func() {
		tryCount := 0
		var timeout <-chan time.Time
		var lastError error

		if retryOptions.Timeout > 0 {
			timeout = time.After(retryOptions.Timeout)
		}

		for {
			if tryCount >= retryOptions.MaxTries {
				if retryOptions.AfterRetryLimit != nil {
					retryOptions.AfterRetryLimit(lastError)
				}
				resultChannel <- &RetryResult{
					StopReason: ErrMaxTriesReached,
					LastError:  lastError,
				}
				close(resultChannel)
				return
			}

			select {
			case <-timeout:
				if retryOptions.AfterTimeout != nil {
					retryOptions.AfterTimeout(lastError)
				}
				resultChannel <- &RetryResult{
					StopReason: ErrTimeout,
					LastError:  lastError,
				}
				close(resultChannel)
				return
			default:
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
			close(resultChannel)
			return
		}
	}()
	return resultChannel, nil
}
