package gotry

import (
	"time"
)

const (
	defaultMaxTries = 5
	defaultDelay    = time.Duration(0)
	defaultTimeout  = time.Duration(5) * time.Second
)

// RetryOptions is the struct containing all relevant options for the Try func-
// tion. This is exported for convenience but you should not need to use it. Use
// the `NewRetryOptionsWithDefaults()` function instead.
type RetryOptions struct {
	Delay           time.Duration
	MaxTries        int
	Timeout         time.Duration
	AfterRetry      func(error)
	AfterRetryLimit func(error)
	AfterTimeout    func(error)
	ReturnChannel   chan *RetryResult
}

// RetryOption is a configuration wrapper for the Try function. It performs some
// operation on RetryOptions, most likely setting a value.
type RetryOption func(options *RetryOptions)

// RetryResult is the data type that is passed back from the Try function once
// the retried operation succeeds and returns a value.
type RetryResult struct {
	Value interface{}
	// Why the retrying was stopped. Either too many tries or timeout.
	StopReason error
	// The last error returned by the operation.
	LastError error
}

// Delay sets the Delay option and determines how long is slept between execu-
// tions of the retried operation.
func Delay(t time.Duration) RetryOption {
	return func(options *RetryOptions) {
		options.Delay = t
	}
}

// MaxTries sets the MaxTries options and determines the maximum number of fai-
// ling executions of the retried operation before aborting and sending an error
// to the channel.
func MaxTries(t int) RetryOption {
	return func(options *RetryOptions) {
		options.MaxTries = t
	}
}

// Timeout sets the Timeout option and determines for how long the operation
// will be retried before aborting.
func Timeout(t time.Duration) RetryOption {
	return func(options *RetryOptions) {
		options.Timeout = t
	}
}

// AfterRetry sets the AfterRetry function which is called with the last
// occurring error after every retry.
func AfterRetry(afterRetry func(err error)) RetryOption {
	return func(options *RetryOptions) {
		options.AfterRetry = afterRetry
	}
}

// AfterRetryLimit sets the AfterRetryLimit function which is called with the
// last occurring error after the retry limit has been reached.
func AfterRetryLimit(afterRetryLimit func(err error)) RetryOption {
	return func(options *RetryOptions) {
		options.AfterRetryLimit = afterRetryLimit
	}
}

// AfterTimeout sets the AfterTimeout function which is called once the retry
// operation times out.
func AfterTimeout(afterTimeout func(err error)) RetryOption {
	return func(options *RetryOptions) {
		options.AfterTimeout = afterTimeout
	}
}

// ReturnChannel sets the ReturnChannel option to which the results of the retry
// process are sent. Use this if you want to explicitly create your own channel.
// Otherwise one will be created.
func ReturnChannel(returnChannel chan *RetryResult) RetryOption {
	return func(options *RetryOptions) {
		options.ReturnChannel = returnChannel
	}
}

// NewRetryOptionsWithDefault builds a RetryOptions struct with default values
// and applies all given RetryOptions to it.
func NewRetryOptionsWithDefault(options ...RetryOption) *RetryOptions {
	retryOptions := &RetryOptions{
		Delay:           defaultDelay,
		MaxTries:        defaultMaxTries,
		Timeout:         defaultTimeout,
		AfterRetry:      nil,
		AfterRetryLimit: nil,
		AfterTimeout:    nil,
		ReturnChannel:   make(chan *RetryResult),
	}

	for _, option := range options {
		option(retryOptions)
	}

	return retryOptions
}
