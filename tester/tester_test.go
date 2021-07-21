package tester

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/smithy-go/middleware"
	"testing"
)

// Mock that satisfies the commandSenderLister interface
type mockClient struct {
	listCommandInvocationRetryIndex int
	mockSendCommand                 func(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (*ssm.SendCommandOutput, error)
	mockListCommandInvocations      []func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error)
}

func (m *mockClient) SendCommand(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (*ssm.SendCommandOutput, error) {
	return m.mockSendCommand(ctx, params, optFns...)
}

func (m *mockClient) ListCommandInvocations(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (*ssm.ListCommandInvocationsOutput, error) {
	// Todo this logic of mocking multiple return values for the same function is a bit convoluted, consider improving
	if m.listCommandInvocationRetryIndex < len(m.mockListCommandInvocations) {
		m.listCommandInvocationRetryIndex = m.listCommandInvocationRetryIndex + 1
	}
	return m.mockListCommandInvocations[m.listCommandInvocationRetryIndex-1](ctx, params, optFns...)
}

func TestTcpConnectionTest(t *testing.T) {
	var cases = []struct {
		caseName      string
		client        func(t *testing.T) *mockClient
		tagName       string
		endpoint      string
		port          string
		expected      interface{}
		expectedError error
	}{
		{
			caseName: "Should return false, error with appropriate message if no invocations are found",
			client: func(t *testing.T) *mockClient {
				return &mockClient{
					mockSendCommand: mockSendCommandHelper(t),
					mockListCommandInvocations: []func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error){
						func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error) {
							// CommandId should be should be set correctly
							if e, a := "dummyCommandId", *params.CommandId; e != a {
								t.Errorf("Expected CommandId to be %s, got %s", e, a)
							}
							return &ssm.ListCommandInvocationsOutput{}, nil
						}},
				}
			},
			tagName:       "ec2NameTag",
			endpoint:      "dummyEndpoint",
			port:          "dummyPort",
			expected:      false,
			expectedError: maxRetriesExceededError{underlying: noInvocationFoundError{}},
		},
		{
			caseName: "Should return true if command invocation completes successfully",
			client: func(t *testing.T) *mockClient {
				return &mockClient{
					mockSendCommand: mockSendCommandHelper(t),
					mockListCommandInvocations: []func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error){
						func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error) {
							return &ssm.ListCommandInvocationsOutput{
								CommandInvocations: []types.CommandInvocation{
									{
										Status: types.CommandInvocationStatusSuccess,
									},
									{
										Status: types.CommandInvocationStatusSuccess,
									},
								},
							}, nil

						}},
				}
			},
			tagName:       "ec2NameTag",
			endpoint:      "dummyEndpoint",
			port:          "dummyPort",
			expected:      true,
			expectedError: nil,
		},
		{
			caseName: "Success should be returned if command invocation completes successfully after retry",
			client: func(t *testing.T) *mockClient {
				return &mockClient{
					mockSendCommand: mockSendCommandHelper(t),
					mockListCommandInvocations: []func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error){
						func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error) {
							return &ssm.ListCommandInvocationsOutput{
								CommandInvocations: []types.CommandInvocation{
									{
										Status: types.CommandInvocationStatusSuccess,
									},
									{
										Status: types.CommandInvocationStatusPending,
									},
								},
							}, nil
						},
						func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error) {
							return &ssm.ListCommandInvocationsOutput{
								CommandInvocations: []types.CommandInvocation{
									{
										Status: types.CommandInvocationStatusSuccess,
									},
									{
										Status: types.CommandInvocationStatusSuccess,
									},
								},
							}, nil
						},
					},
				}
			},
			tagName:       "ec2NameTag",
			endpoint:      "dummyEndpoint",
			port:          "dummyPort",
			expected:      true,
			expectedError: nil,
		},
		{
			caseName: "Should Fail if command invocation completes with Failure after retry",
			client: func(t *testing.T) *mockClient {
				return &mockClient{
					mockSendCommand: mockSendCommandHelper(t),
					mockListCommandInvocations: []func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error){
						func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error) {
							return &ssm.ListCommandInvocationsOutput{
								CommandInvocations: []types.CommandInvocation{
									{
										Status: types.CommandInvocationStatusSuccess,
									},
									{
										Status: types.CommandInvocationStatusPending,
									},
								},
							}, nil
						},
						func(ctx context.Context, params *ssm.ListCommandInvocationsInput, optFns ...func(*ssm.Options)) (output *ssm.ListCommandInvocationsOutput, e error) {
							return &ssm.ListCommandInvocationsOutput{
								CommandInvocations: []types.CommandInvocation{
									{
										InstanceId: stringPointer("dummyInstanceId"),
										Status:     types.CommandInvocationStatusFailed,
									},
									{
										InstanceId: stringPointer("dummyInstanceId2"),
										Status:     types.CommandInvocationStatusSuccess,
									},
								},
							}, nil
						},
					},
				}
			},
			tagName:       "ec2NameTag",
			endpoint:      "dummyEndpoint",
			port:          "dummyPort",
			expected:      false,
			expectedError: FatalError{Underlying: failedForInstanceIdError{instanceId: "dummyInstanceId"}},
		},
	}
	for _, c := range cases {
		t.Run(c.caseName, func(t *testing.T) {
			actual, actualErr := TcpConnectionTestWithNameTagE(t, c.client(t), c.tagName, c.endpoint, c.port)
			if (c.expectedError != nil && actualErr == nil) || (c.expectedError == nil && actualErr != nil) {
				t.Errorf("Expected error %v, but got %v", c.expectedError, actualErr)
			}
			if c.expectedError != nil && actualErr != nil {
				if c.expectedError.Error() != actualErr.Error() {
					t.Errorf("Expected error message to be %s, but got %s", c.expectedError.Error(), actualErr.Error())
				}
			}
			if actual != c.expected {
				t.Errorf("Expected %v, but got %v", c.expected, actual)
			}
		})
	}
}

func mockSendCommandHelper(t *testing.T) func(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (output *ssm.SendCommandOutput, e error) {
	return func(ctx context.Context, params *ssm.SendCommandInput, optFns ...func(*ssm.Options)) (output *ssm.SendCommandOutput, e error) {
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
			Command: &types.Command{
				CommandId: stringPointer("dummyCommandId"),
			},
			ResultMetadata: middleware.Metadata{},
		}, nil
	}
}
