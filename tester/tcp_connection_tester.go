package tester

import (
	"fmt"
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
func TcpConnectionTestWithNameTagE(t *testing.T, client commandSenderLister, tagName, endpoint string, port string, maxRetries int, waitBetweenRetries time.Duration) (bool, error) {
	// build target using tagName
	target := NewTagNameTarget(tagName)
	// build testCase with bash command to check tcp connectivity
	testCase := NewShellTestCase(tcpConnectionTestShellCommand(3, endpoint, port), true)
	return UseThisToTestE(t, client, testCase, target, maxRetries, waitBetweenRetries)
}

func tcpConnectionTestShellCommand(timeoutInSeconds int, endpoint string, port string) string {
	return fmt.Sprintf("timeout %d bash -c '</dev/tcp/%s/%s'", timeoutInSeconds, endpoint, port)
}
