// Package tester enables execution of test commands on target EC2 instances via the AWS Systems Manager(SSM).
// tester will poll AWS SSM API for the success or failure of the command sent to target VMs and report success of failures accordingly.
package tester

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"testing"
	"time"
)

// Todo Documentation
func RunTestCaseForTarget(t *testing.T, client commandSenderLister, testCase commandParameterBuilder, target targetParamBuilder,
	maxRetries int, waitBetweenRetries time.Duration) {
	_, err := RunTestCaseForTargetE(t, client, testCase, target, maxRetries, waitBetweenRetries)
	if err != nil {
		t.Error(err)
	}
}

// Todo Documentation
// RunTestCaseForTargetE accepts a testCase, target and some retry related configuration and uses AWS SSM to run
// It returns false and an error if no instances are found to match the Name tag.
// It returns false and an error if any one of the instances cannot run the command successfully or within timeout.
// It returns false and error for any other error.
func RunTestCaseForTargetE(t *testing.T, client commandSenderLister, testCase commandParameterBuilder, target targetParamBuilder,
	maxRetries int, waitBetweenRetries time.Duration) (bool, error) {
	// send test command
	sendCommandInput := newSendCommandInput(testCase, target)
	sendCommandOutput, err := client.SendCommand(context.Background(), sendCommandInput)
	if err != nil {
		return false, err
	}
	// poll for test command execution results
	retryAction := getListCommandAction(t, client, *sendCommandOutput.Command.CommandId)
	result, err := retry(t, "Poll For Invocation Results", maxRetries, waitBetweenRetries, retryAction)
	if err != nil {
		return false, err
	}
	return result.(bool), nil
}

func newSendCommandInput(testCase commandParameterBuilder, target targetParamBuilder) *ssm.SendCommandInput {
	return &ssm.SendCommandInput{
		DocumentName:    stringPointer(testCase.documentName()),
		DocumentVersion: stringPointer(testCase.documentVersion()),
		Parameters:      testCase.buildCommandParameters(),
		Targets:         target.buildTargetParameters(),
	}
}

func buildListCommandInput(commandId string) *ssm.ListCommandInvocationsInput {
	return &ssm.ListCommandInvocationsInput{
		CommandId: &commandId,
	}
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
			return false, invocationsIncompleteError{}
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

type invocationsIncompleteError struct {
}

func (err invocationsIncompleteError) Error() string {
	return "command invocations pending, inprogress or delayed"
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
