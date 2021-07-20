package tester

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/smithy-go/middleware"
	"testing"
)

// Mock that satisfies the commandSenderLister interface
type mockClient struct {
	mockSendCommand            func(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (*ssm.SendCommandOutput, error)
	mockListCommandInvocations func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error)
}

func (m *mockClient) SendCommand(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (*ssm.SendCommandOutput, error) {
	return m.mockSendCommand(ctx, params, optFns...)
}

func (m *mockClient) ListCommandInvocations(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error) {
	return m.mockListCommandInvocations(ctx, params, optFns...)
}

func TestTcpConnectionTest(t *testing.T) {
	cases := []struct {
		client   func(t *testing.T) *mockClient
		tagName  string
		endpoint string
		port     string
		expected interface{}
	}{
		{client: func(t *testing.T) *mockClient {
			return &mockClient{
				mockSendCommand: func(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (output *ssm.SendCommandOutput, e error) {
					// Name tag should be set correctly
					if e, a := "ec2NameTag", params.Targets[0].Values[0]; e != a {
						t.Errorf("Expected Name tag to be %s, got %s", e, a)
					}
					// Document Name should be set correctly
					if e, a := "AWS-RunShellScript", *params.DocumentName; e != a {
						t.Errorf("Expected DocumentName to be set to %s, got %s", e, a)
					}
					// Parameters Command should be set correctly
					if e, a := `bash -c '</dev/tcp/dummyEndpoint/dummyPort'`, params.Parameters["command"][0]; e != a {
						t.Errorf("Expected DocumentName to be set to %s, got %s", e, a)
					}
					return &ssm.SendCommandOutput{
						Command:        nil,
						ResultMetadata: middleware.Metadata{},
					}, nil
				},
				mockListCommandInvocations: nil,
			}
		},
			tagName:  "ec2NameTag",
			endpoint: "dummyEndpoint",
			port:     "dummyPort",
		},
	}

	for _, c := range cases {
		_, err := TcpConnectionTestWithNameTag(c.client(t), c.tagName, c.endpoint, c.port)
		if err != nil {
			t.Errorf("Failed with \n%s", err.Error())
		}
	}
}
