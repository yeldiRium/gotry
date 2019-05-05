package gotry

import (
	"testing"
	"time"
)

var setupDefaultOptions = func() *RetryOptions {
	return &RetryOptions{
		Delay:           time.Duration(50) * time.Millisecond,
		Timeout:         time.Duration(10) * time.Second,
		MaxTries:        10,
		AfterRetry:      nil,
		AfterRetryLimit: nil,
		ReturnChannel:   make(chan *RetryResult),
	}
}

func TestDelayOption(t *testing.T) {
	options := setupDefaultOptions()
	newDelay := time.Duration(500) * time.Second
	Delay(newDelay)(options)

	if options.Delay != newDelay {
		t.Errorf("Delay() did not set the delay on RetryOptions correcty. Expected %v, got %v.", newDelay, options.Delay)
	}
}

func TestMaxTriesOption(t *testing.T) {
	options := setupDefaultOptions()
	maxTries := 20
	MaxTries(maxTries)(options)

	if options.MaxTries != maxTries {
		t.Errorf("MaxTries() did not set the max tries on RetryOptions correcty. Expected %v, got %v.", maxTries, options.MaxTries)
	}
}

func TestTimeoutOption(t *testing.T) {
	options := setupDefaultOptions()
	timeout := time.Duration(500) * time.Second
	Timeout(timeout)(options)

	if options.Timeout != timeout {
		t.Errorf("Timeout() did not set the timeout on RetryOptions correcty. Expected %v, got %v.", timeout, options.Timeout)
	}
}

func TestAfterRetryOption(t *testing.T) {
	options := setupDefaultOptions()
	afterRetry := func(err error) {}
	AfterRetry(afterRetry)(options)

	if options.AfterRetry == nil {
		t.Errorf("AfterRetry() did not set the after retry function.")
	}
}

func TestAfterRetryLimitOption(t *testing.T) {
	options := setupDefaultOptions()
	afterRetryLimit := func(err error) {}
	AfterRetryLimit(afterRetryLimit)(options)

	if options.AfterRetryLimit == nil {
		t.Errorf("AfterRetryLimit() did not set the after retry function.")
	}
}

func TestReturnChannelOption(t *testing.T) {
	options := setupDefaultOptions()
	returnChannel := make(chan *RetryResult)
	ReturnChannel(returnChannel)(options)

	if options.ReturnChannel != returnChannel {
		t.Errorf("ReturnChannel() did not set the return channel on RetryOptions correcty. Expected %v, got %v.", returnChannel, options.ReturnChannel)
	}
}

func TestNewRetryOptionsWithDefault(t *testing.T) {
	options := newRetryOptionsWithDefault()

	if options.Delay != defaultDelay {
		t.Errorf("The default delay was not set correctly. Expected %v, got %v", defaultDelay, options.Delay)
	}
	if options.MaxTries != defaultMaxTries {
		t.Errorf("The default max tries was not set correctly. Expected %v, got %v", defaultMaxTries, options.MaxTries)
	}
	if options.Timeout != defaultTimeout {
		t.Errorf("The default timeout was not set correctly. Expected %v, got %v", defaultTimeout, options.Timeout)
	}
	if options.AfterRetry != nil {
		t.Error("The default AfterRetry function was not set correctly.")
	}
	if options.AfterRetryLimit != nil {
		t.Error("The default AfterRetryLimit function was not set correctly.")
	}
	if options.ReturnChannel == nil {
		t.Errorf("The default ReturnChannel was not set.")
	}
}

func TestNewRetryOptionsWithDefaultAndAdditionalOptions(t *testing.T) {
	delay := time.Duration(500) * time.Second
	maxTries := 20
	timeout := time.Duration(20) * time.Hour
	afterRetry := func(err error) {}
	afterRetryLimit := func(err error) {}
	returnChannel := make(chan *RetryResult)

	options := newRetryOptionsWithDefault(
		ReturnChannel(returnChannel),
		Delay(delay),
		AfterRetry(afterRetry),
		MaxTries(maxTries),
		Timeout(timeout),
		AfterRetryLimit(afterRetryLimit),
	)

	if options.Delay != delay {
		t.Errorf("The default delay was not set correctly. Expected %v, got %v", delay, options.Delay)
	}
	if options.MaxTries != maxTries {
		t.Errorf("The default max tries was not set correctly. Expected %v, got %v", defaultMaxTries, options.MaxTries)
	}
	if options.Timeout != timeout {
		t.Errorf("The default timeout was not set correctly. Expected %v, got %v", timeout, options.Timeout)
	}
	if options.AfterRetry == nil {
		t.Error("The default AfterRetry function was not set correctly.")
	}
	if options.AfterRetryLimit == nil {
		t.Error("The default AfterRetryLimit function was not set correctly.")
	}
	if options.ReturnChannel != returnChannel {
		t.Errorf("The default ReturnChannel was not set.")
	}
}
