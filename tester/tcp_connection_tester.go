package tester

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"testing"
)

// TcpConnectionTestWithNameTagE return true if instances are found and all instance can run connection command successfully.
// It returns false and an error if no instances are found to match the Name tag.
// It returns false and an error if any one of the instances cannot run the command successfully or within timeout.
// It Returns false and error for any other error.
func TcpConnectionTestWithNameTagE(t *testing.T, client commandSenderLister, tagName string, endpoint string, port string) (bool, error) {
	// build command
	tcpConnectionTestSendCommandInput := buildSendCommandInputForTcpConnectionWithNameTag(endpoint, port, tagName)
	fmt.Println(tcpConnectionTestSendCommandInput)
	// send command and poll for results
	return sendCommandAndPollResults(t, client, tcpConnectionTestSendCommandInput)
}

func buildSendCommandInputForTcpConnectionWithNameTag(endpoint string, port string, tagName string) *ssm.SendCommandInput {
	return buildSendCommandInput(buildParametersForTcpConnectionWithDefaultTimeout(endpoint, port), buildTargetsFromNameTag(tagName))
}

func buildParametersForTcpConnectionWithDefaultTimeout(endpoint string, port string) map[string][]string {
	const defaultTimeout= "5"
	command := []string{fmt.Sprintf("bash -c '</dev/tcp/%s/%s'", endpoint, port)}
	return buildParametersForCommand(command, defaultTimeout)
}

func buildTargetsFromNameTag(tagName string) []types.Target {
	values := []string{tagName}
	return []types.Target{{Key: stringPointer("tag:Name"), Values: values}}
}
