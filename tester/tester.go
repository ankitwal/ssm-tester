package tester

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"log"
	"testing"
	"time"
)

// Return true if instances are found and all instance can run connection command successfully
// Returns false and an error if no instances are found to match the Name tag
// Returns false and an error if any one of the instances cannot run the command successfully or within timeout.
// Returns false and error for any other error
func TcpConnectionTestWithNameTagE(t *testing.T, client commandSenderLister, tagName string, endpoint string, port string) (bool, error) {
	sendCommandOutput, err := client.SendCommand(context.Background(), buildSendCommandInputForTcpConnectionWithNameTag(endpoint, port, tagName))
	if err != nil {
		return false, err
	}
	retryAction := getListCommandAction(t, client, *sendCommandOutput.Command.CommandId)
	result, err := retry(t, "Poll For Invocation Results", 5, 3*time.Second, retryAction)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

func getListCommandAction(t *testing.T, client commandLister, commandId string) func() (interface{}, error) {
	return func() (interface{}, error) {
		listCommandOutput, err := client.ListCommandInvocations(context.Background(), buildListCommandInput(commandId))
		if err != nil {
			t.Error(err)
		}
		if len(listCommandOutput.CommandInvocations) == 0 {
			// Todo Log message to help debugging
			return nil, noInvocationFoundError{}
		}

		//Check status of all found invocations
		result, err := checkAllInvocationForStatus(listCommandOutput)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

// Returns true if all the invocations have succeeded.
// Returns false if any of the invocations has failed
// Returns an error, signalling to the retry function to try again in the case of pending, inprogress or delayed invocaiton
func checkAllInvocationForStatus(listCommandOutput *ssm.ListCommandInvocationsOutput) (bool, error) {
	for _, v := range listCommandOutput.CommandInvocations {
		switch v.Status {
		case types.CommandInvocationStatusPending,
			types.CommandInvocationStatusInProgress,
			types.CommandInvocationStatusDelayed:
			// Todo log message about why retry
			return false, errors.New("command invocations pending or delayed")
		case types.CommandInvocationStatusFailed,
			types.CommandInvocationStatusCancelled,
			types.CommandInvocationStatusCancelling,
			types.CommandInvocationStatusTimedOut:
			// Todo log message to help debug failure
			// return false and nil error to signal retry to stop and to return false to the user
			return false, FatalError{Underlying: failedForInstanceIdError{instanceId: *v.InstanceId}}
		}
	}
	// In the case that all the invocations were Successful
	return true, nil
}

func buildListCommandInput(commandId string) *ssm.ListCommandInvocationsInput {
	return &ssm.ListCommandInvocationsInput{
		CommandId:  &commandId,
		Details:    false,
		Filters:    nil,
		InstanceId: nil,
		MaxResults: 0,
		NextToken:  nil,
	}
}

func buildSendCommandInputForTcpConnectionWithNameTag(endpoint string, port string, tagName string) *ssm.SendCommandInput {
	return &ssm.SendCommandInput{
		DocumentName:    stringPointer("AWS-RunShellScript"),
		DocumentVersion: stringPointer("$LATEST"),
		Parameters:      buildParametersForTcpConnection(endpoint, port),
		Targets:         buildTargetsFromNameTag(tagName),
	}
}
func buildParametersForTcpConnection(endpoint string, port string) map[string][]string {
	command := []string{fmt.Sprintf("bash -c '</dev/tcp/%s/%s'", endpoint, port)}
	return buildParametersForCommand(command)
}
func buildParametersForCommand(command []string) map[string][]string {
	const (
		commands         = "commands"
		executionTimeout = "executionTimeout"
	)
	parameters := map[string][]string{}
	parameters[commands] = command
	parameters[executionTimeout] = []string{"5"}
	return parameters
}
func buildTargetsFromNameTag(tagName string) []types.Target {
	values := []string{tagName}
	return []types.Target{{Key: stringPointer("tag:Name"), Values: values}}
}
func stringPointer(s string) *string {
	temp := s
	return &temp
}

// retry runs the specified action. If it returns a value, return that value. If it returns a FatalError, return that error
// immediately. If it returns any other type of error, sleep for sleepBetweenRetries and try again, up to a maximum of
// maxRetries retries. If maxRetries is exceeded, return a MaxRetriesExceeded error.
func retry(t *testing.T, actionDescription string, maxRetries int, sleepBetweenRetries time.Duration, action func() (interface{}, error)) (interface{}, error) {
	var output interface{}
	var err error

	for i := 0; i <= maxRetries; i++ {

		output, err = action()
		if err == nil {
			return output, nil
		}

		if _, isFatalErr := err.(FatalError); isFatalErr {
			log.Printf("Returning due to fatal error: %v", err)
			return output, err
		}

		log.Printf("%s returned an error: %s. Sleeping for %s and will try again.", actionDescription, err.Error(), sleepBetweenRetries)
		time.Sleep(sleepBetweenRetries)
	}

	return output, maxRetriesExceededError{underlying: err}
}

// FatalError is a marker interface for errors that should not be retried.
type FatalError struct {
	Underlying error
}

func (err FatalError) Error() string {
	return fmt.Sprintf("FatalError{Underlying: %v}", err.Underlying)
}

type maxRetriesExceededError struct {
	underlying error
}

func (m maxRetriesExceededError) Error() string {
	return fmt.Sprintf("max retires exceeded - last underlying error: %s", m.underlying.Error())
}

type noInvocationFoundError struct {
}

func (err noInvocationFoundError) Error() string {
	return "no invocations found"
}

type failedForInstanceIdError struct {
	instanceId string
}

func (err failedForInstanceIdError) Error() string {
	return fmt.Sprintf("command invocations failed for instanceId %s", err.instanceId)
}
