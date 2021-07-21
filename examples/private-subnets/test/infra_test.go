package test

import (
	"context"
	"github.com/ankitwal/ssm-tester/tester"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"testing"
)

func TestInfra(t *testing.T) {
	config, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Error(err)
	}
	_, err = tester.TcpConnectionTestWithNameTagE(t, ssm.NewFromConfig(config), "private-subnets-example-test",
		"private-subnets-example-test.cvcfzdalq7qz.ap-southeast-1.rds.amazonaws.com", "3306")
	if err != nil {
		t.Error(err)
	}
}
