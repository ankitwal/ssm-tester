package tester

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// create tiny interfaces to enable DI with mocks in unit tests
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

// Interface type to abstract testCases from the tester
type commandParameterBuilder interface {
	documentName() string
	documentVersion() string
	buildCommandParameters() map[string][]string
}

// Interface type to abstract concrete type from tester
type targetParamBuilder interface {
	buildTargetParameters() []types.Target
}

