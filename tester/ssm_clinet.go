package tester

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"testing"
)

// NewSSMClientWithDefaultConfig initialises and returns and instance of AWS SSM Service Client
// with default config.
func NewSSMClientWithDefaultConfig(t *testing.T) *ssm.Client {
	config, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Error(err)
	}
	return ssm.NewFromConfig(config)
}

