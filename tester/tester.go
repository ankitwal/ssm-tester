package tester

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"log"
)

func TcpConnectionTestWithNameTag(client commandSenderLister, tagName string, endpoint string, port string) (interface{}, error) {
	sendCommandOutput, err := client.SendCommand(context.Background(), buildSendCommandInputForTcpConnectionWithNameTag(endpoint, port, tagName))
	if err != nil {
		log.Fatal(err.Error())
	}
	return sendCommandOutput, err
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
