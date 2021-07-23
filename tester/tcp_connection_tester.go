package tester

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"testing"
	"time"
)

// TcpConnectionTestWithNameTag return true if instances are found and all instance can run connection command successfully.
// It fails the test if no instances are found to match the Name tag.
// It fails the test if any one of the instances cannot run the command successfully or within timeout, or any other error.
// It passes the test if all found instances for tag Name run the command successfully in the given timeouts.
//
// tagName should be the string value of the tag:Name of the target EC2 instances.
// endpoint, and port should be the network endpoint to validate tcp connectivity to.
//
// maxRetries specifies the number of times the test should poll AWS API for results of the command sent to the target EC2 VMs.
// waitBetweenRetires specifies the duration in time.Seconds to wait between each retry.
// these values may need to be adjusted for the total number of ec2 instances that are expected to run the test command.
func TcpConnectionTestWithNameTag(t *testing.T, client commandSenderLister, tagName string, endpoint string, port string, maxRetries int, waitBetweenRetries time.Duration) {
	_, err := TcpConnectionTestWithNameTagE(t, client, tagName, endpoint, port, maxRetries, waitBetweenRetries)
	if err != nil {
		t.Error(err)
	}
}

// TcpConnectionTestWithNameTagE return true if instances are found and all instance can run connection command successfully.
// It returns false and an error if no instances are found to match the Name tag.
// It returns false and an error if any one of the instances cannot run the command successfully or within timeout.
// It returns false and error for any other error.
// It returns true and nil error if all found instances for tag Name run the command successfully in the given timeouts.
func TcpConnectionTestWithNameTagE(t *testing.T, client commandSenderLister, tagName string, endpoint string, port string, maxRetries int, waitBetweenRetries time.Duration) (bool, error) {
	// build command
	tcpConnectionTestSendCommandInput := buildSendCommandInputForTcpConnectionWithNameTag(endpoint, port, tagName)
	fmt.Println(tcpConnectionTestSendCommandInput)
	// send command and poll for results
	return sendCommandAndPollResults(t, client, tcpConnectionTestSendCommandInput, maxRetries, waitBetweenRetries)
}

func buildSendCommandInputForTcpConnectionWithNameTag(endpoint string, port string, tagName string) *ssm.SendCommandInput {
	return buildSendCommandInput(buildParametersForTcpConnectionWithDefaultTimeout(endpoint, port), buildTargetsFromNameTag(tagName))
}

func buildParametersForTcpConnectionWithDefaultTimeout(endpoint string, port string) map[string][]string {
	// ball park execution timeout that should be okay for testing tcp connectivity
	const defaultTimeout = "10"
	command := []string{fmt.Sprintf("bash -c '</dev/tcp/%s/%s'", endpoint, port)}
	return buildParametersForCommand(command, defaultTimeout)
}

func buildTargetsFromNameTag(tagName string) []types.Target {
	values := []string{tagName}
	return []types.Target{{Key: stringPointer("tag:Name"), Values: values}}
}
