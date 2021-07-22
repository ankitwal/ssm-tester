package test

import (
	"context"
	"github.com/ankitwal/ssm-tester/tester"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"time"
)

func TestInfra(t *testing.T) {

	// this example uses terratest to initialise the terraform stack and get output value
	// please see here for some 'how to use terratest' basics: https://terratest.gruntwork.io/examples/
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
	})

	// init and apply terraform stack ensuring clean up
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// get the required resource values using terratest's terraform module
	dbEndpoint := terraform.Output(t, terraformOptions, "database_endpoint")
	dbPort := terraform.Output(t, terraformOptions, "database_port")
	tag := terraform.Output(t, terraformOptions, "instance_name_tag")

	// Initialise AWS SSM client service.
	config, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Error(err)
	}

	// the number of times you want the tester to retry polling for the result of the test command
	maxRetriesToPollResult := 5
	// the number of time to sleep between retries
	waitBetweenRetries := 3 * time.Second
	// run the test and check for errors. If no errors are returned the command is success ful for all instances
	_, err = tester.TcpConnectionTestWithNameTagE(t, ssm.NewFromConfig(config), tag, dbEndpoint, dbPort, maxRetriesToPollResult, waitBetweenRetries)
	if err != nil {
		t.Error(err)
	}
}
