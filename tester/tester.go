package tester

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gruntwork-io/terratest/modules/retry"
	"testing"
	"time"
)

// Return true if connection is possible
// Returns false if tcp connection is not established
// Returns an error if there is an error.
func TcpConnectionTestWithNameTag(t *testing.T, client commandSenderLister, tagName string, endpoint string, port string) (bool, error) {
	sendCommandOutput, err := client.SendCommand(context.Background(), buildSendCommandInputForTcpConnectionWithNameTag(endpoint, port, tagName))
	if err != nil {
		t.Error(err)
	}
	retryAction := getListCommandAction(t, client, *sendCommandOutput.Command.CommandId)
	result, err := retry.DoWithRetryInterfaceE(t, "Poll For Invocation Results", 5, 3*time.Second, retryAction)
	switch err.(type) {
	case retry.FatalError:
		// if FatalError Fail immediately
		t.Errorf(err.Error())
	case retry.MaxRetriesExceeded:
		// return false and error
		return false, errors.New(err.(retry.MaxRetriesExceeded).Error())
	case nil:
		return result.(bool), nil
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
			return nil, errors.New("No invocations found yet, force retry")
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
			return false, errors.New("Force retry")
		case types.CommandInvocationStatusFailed,
			types.CommandInvocationStatusCancelled,
			types.CommandInvocationStatusCancelling,
			types.CommandInvocationStatusTimedOut:
			// Todo log message to help debug failure
			// return false and nil error to signal retry to stop and to return false to the user
			return false, nil
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
		DocumentName:           stringPointer("AWS-RunShellScript"),
		CloudWatchOutputConfig: nil,
		Comment:                nil,
		DocumentHash:           nil,
		DocumentHashType:       "",
		DocumentVersion:        nil,
		InstanceIds:            nil,
		MaxConcurrency:         nil,
		MaxErrors:              nil,
		NotificationConfig:     nil,
		OutputS3BucketName:     nil,
		OutputS3KeyPrefix:      nil,
		OutputS3Region:         nil,
		Parameters:             buildParametersForTcpConnection(endpoint, port),
		ServiceRoleArn:         nil,
		Targets:                buildTargetsFromNameTag(tagName),
		TimeoutSeconds:         0,
	}
}
func buildParametersForTcpConnection(endpoint string, port string) map[string][]string {
	command := []string{fmt.Sprintf("bash -c '</dev/tcp/%s/%s'", endpoint, port)}
	return buildParametersForCommand(command)
}
func buildParametersForCommand(command []string) map[string][]string {
	parameters := map[string][]string{}
	parameters["command"] = command
	parameters["executionTimeout"] = []string{"5"}
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
