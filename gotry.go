package gotry

import (
	"time"
)

// Try ...
func Try(f func() (interface{}, error), options ...RetryOption) (<-chan *RetryResult, error) {
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
