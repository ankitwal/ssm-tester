package tester

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type commandSender interface {
	SendCommand(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (*ssm.SendCommandOutput, error)
}

type commandLister interface {
	ListCommandInvocations(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error)
}

type commandSenderLister interface {
	commandSender
	commandLister
}
