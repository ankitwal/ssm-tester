package tester

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// retry runs the specified action. If it returns a value, return that value. If it returns a fatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, return a MaxRetriesExceeded error.
// based on terratest RetryWithInterface - needed some custom code to wrap and propagate underlying errors for MaxRetriesExceeded
func retry(t *testing.T, actionDescription string, maxRetries int, waitBetweenRetries time.Duration, action func() (interface{}, error)) (interface{}, error) {
	var output interface{}
	var err error

	for i := 0; i <= maxRetries; i++ {

		output, err = action()
		if err == nil {
			return output, nil
		}

		if _, isFatalErr := err.(fatalError); isFatalErr {
			log.Printf("Returning due to fatal error: %v", err)
			return output, err
		}

		log.Printf("%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), waitBetweenRetries)
		time.Sleep(waitBetweenRetries)
	}

	return output, maxRetriesExceededError{underlying: err}
}

// fatalError is a marker interface for errors that should not be retried.
type fatalError struct {
	Underlying error
}

func (err fatalError) Error() string {
	return fmt.Sprintf("fatalError stopped immediately - underlying error: %v}", err.Underlying)
}

type maxRetriesExceededError struct {
	underlying error
}

func (m maxRetriesExceededError) Error() string {
	return fmt.Sprintf("max retires exceeded - last underlying error: %s", m.underlying.Error())
}
