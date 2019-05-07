package gotry

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestSucceedingTryWithDefaultValues(t *testing.T) {
	opCallCount := 0
	retryCount := 0

	op := func() (interface{}, error) {
		opCallCount++

		if opCallCount < 4 {
			return nil, fmt.Errorf("this is error no. %v", opCallCount)
		}
		return true, nil
	}

	afterRetry := func(err error) {
		retryCount++
	}

	resultChannel := make(chan *RetryResult)
	go Try(
		op,
		resultChannel,
		AfterRetry(afterRetry),
	)

	result := <-resultChannel
	if result.StopReason != nil {
		t.Errorf("Retry stopped although it should not have. StopReason: %v", result.StopReason)
	}
	if result.Value != true {
		t.Errorf("The result value was not as expected. Expected %v, got %v.", true, result.Value)
	}
}

func TestFailingWithTooManyTries(t *testing.T) {
	opError := errors.New("lul an error")
	op := func() (interface{}, error) {
		return nil, opError
	}

	retryLimitCallbackCalled := false
	retryLimitCallback := func(err error) {
		retryLimitCallbackCalled = true
	}

	resultChannel := make(chan *RetryResult)
	go Try(
		op,
		resultChannel,
		AfterRetryLimit(retryLimitCallback),
	)
	result := <-resultChannel
	if result.Value != nil {
		t.Errorf("Retry returned a value but should not have: %v", result.Value)
	}
	if !IsMaxTriesReached(result.StopReason) {
		t.Errorf("Retry failed with wrong reason: %v. Should have been MaxTriesReached.", result.StopReason)
	}
	if result.LastError != opError {
		t.Errorf("Retry reported an unexpected LastError: %v. Should have been %v.", result.LastError, opError)
	}
	if !retryLimitCallbackCalled {
		t.Errorf("AfterRetryLimit callback was not called.")
	}
}

func TestFailingWithTimeout(t *testing.T) {
	opError := errors.New("lul an error")
	op := func() (interface{}, error) {
		return nil, opError
	}

	timeoutCallbackCalled := false
	timeoutCallback := func(err error) {
		timeoutCallbackCalled = true
		if err != opError {
			t.Errorf("TimoutCallback was called with wrong error. Was %v, should have been %v.", err, opError)
		}
	}

	resultChannel := make(chan *RetryResult)
	go Try(
		op,
		resultChannel,
		Timeout(time.Duration(1)*time.Second),
		Delay(time.Duration(500)*time.Millisecond),
		AfterTimeout(timeoutCallback),
	)

	result := <-resultChannel
	if result.Value != nil {
		t.Errorf("ReRetry returned a value but should not have: %v", result.Value)
	}
	if !IsTimeout(result.StopReason) {
		t.Errorf("Retry failed with wrong reason: %v. Should have been Timeout.", result.StopReason)
	}
	if result.LastError != opError {
		t.Errorf("Retry reported an unexpected LastError: %v. Should have been %v.", result.LastError, opError)
	}
	if !timeoutCallbackCalled {
		t.Errorf("AfterTimeout callback was not called.")
	}
}

func TestRetryCallback(t *testing.T) {
	opCallCount := 0
	opError := errors.New("lul an error")
	op := func() (interface{}, error) {
		opCallCount++
		if opCallCount < 5 {
			return nil, opError
		}
		return true, nil
	}

	retryCallbackCount := 0
	retryCallback := func(err error) {
		retryCallbackCount++
		if err != opError {
			t.Errorf("Retry function was not called with the correct error. Got %v, should have been %v.", err, opError)
		}
	}

	resultChannel := make(chan *RetryResult)
	go Try(
		op,
		resultChannel,
		MaxTries(5),
		AfterRetry(retryCallback),
	)

	<-resultChannel
	if retryCallbackCount != 4 {
		t.Errorf("RetryCallback was not called the correct amount of times. It was called %v times but should have been called %v times.", retryCallbackCount, 4)
	}
}

func TestTypecastingAfterSucceeding(t *testing.T) {
	type MockStruct struct {
		Thingy bool
		Thang  int
	}

	returnStruct := &MockStruct{
		Thingy: false,
		Thang:  3,
	}
	op := func() (interface{}, error) {
		return returnStruct, nil
	}

	resultChannel := make(chan *RetryResult)
	go Try(op, resultChannel)

	result := <-resultChannel
	mockStructValue := result.Value.(*MockStruct)
	if mockStructValue.Thingy != false {
		t.Errorf("Result was not as expected. Got %v, expected %v.", mockStructValue.Thingy, returnStruct.Thingy)
	}
	if mockStructValue.Thang != 3 {
		t.Errorf("Result was not as expected. Got %v, expected %v.", mockStructValue.Thang, returnStruct.Thang)
	}
}

/*
TODO: make this test pass. this library needs to be able to handle long,
blocking operations. it is goTRY and should absolutely be usable for single-exe-
cution scenarios.

func TestFailingWithLongBlockingOperation(t *testing.T) {
	op := func() (interface{}, error) {
		time.Sleep(time.Duration(2) * time.Second)
		return nil, errors.New("this error should not be returned since retry should abort earlier")
	}

	resultChannel := make(chan *RetryResult)
	go Try(
		op,
		resultChannel,
		Timeout(time.Duration(1)*time.Second),
	)

	result := <-resultChannel
	if result.Value != nil {
		t.Errorf("Retry returned a value but should not have: %v", result.Value)
	}
	if !IsTimeout(result.StopReason) {
		t.Errorf("Retry failed with wrong reason: %v. Should have been Timeout.", result.StopReason)
	}
	if result.LastError != nil {
		t.Errorf("Retry reported an unexpected LastError: %v. Should have been %v.", result.LastError, nil)
	}
}
*/
