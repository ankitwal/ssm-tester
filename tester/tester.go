// Package tester enables execution of test commands on target EC2 instances via the AWS Systems Manager(SSM).
// tester will poll AWS SSM API for the success or failure of the command sent to target VMs and report success of failures accordingly.
package tester

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"testing"
	"time"
)

func sendCommandAndPollResults(t *testing.T, client commandSenderLister, sendCommandInput *ssm.SendCommandInput, maxRetries int, waitBetweenRetries time.Duration) (bool, error) {
	// Command Sender
	sendCommandOutput, err := client.SendCommand(context.Background(), sendCommandInput)
	if err != nil {
		return false, err
	}
	// Command Result Poller
	retryAction := getListCommandAction(t, client, *sendCommandOutput.Command.CommandId)
	result, err := retry(t, "Poll For Invocation Results", maxRetries, waitBetweenRetries, retryAction)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
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

func buildSendCommandInput(parameters map[string][]string, targets []types.Target) *ssm.SendCommandInput {
	const (
		awsShellScript  = "AWS-RunShellScript"
		documentVersion = "$LATEST"
	)
	return &ssm.SendCommandInput{
		DocumentName:    stringPointer(awsShellScript),
		DocumentVersion: stringPointer(documentVersion),
		Parameters:      parameters,
		Targets:         targets,
	}
}

func buildParametersForCommand(command []string, timeout string) map[string][]string {
	const (
		commands         = "commands"
		executionTimeout = "executionTimeout"
	)
	parameters := map[string][]string{}
	parameters[commands] = command
	parameters[executionTimeout] = []string{timeout}
	return parameters
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
// Returns an error, signalling to the retry function to try again in the case of pending, in progress or delayed invocation
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
			return false, fatalError{Underlying: failedForInstanceIdError{instanceId: *v.InstanceId}}
		}
	}
	// In the case that all the invocations were Successful
	return true, nil
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
