package test

import (
	"context"
	"fmt"
	"github.com/ankitwal/ssm-tester/tester"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"time"
)

func TestInfra(t *testing.T) {

	t.Parallel()
	// this example uses terratest to initialise the terraform stack and get output value
	// please see here for some 'how to use terratest' basics: https://terratest.gruntwork.io/examples/
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
	})

	// init and apply terraform stack ensuring clean up
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Initialise AWS SSM client service.
	config, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		t.Error(err)
	}
	ssmClient := ssm.NewFromConfig(config)

	// create retry configuration for
	// the tester the number of times the tester should retry polling for the result of the test command
	maxRetriesToPollResult := 5
	waitBetweenRetries := 3 * time.Second

	t.Run("TestAppInstanceConnectivityToDatabase", func(t *testing.T) {
		t.Parallel()
		// get the required resource values using terratest's terraform module
		dbEndpoint := terraform.Output(t, terraformOptions, "database_endpoint")
		dbPort := terraform.Output(t, terraformOptions, "database_port")
		tag := terraform.Output(t, terraformOptions, "instance_name_tag")
		tester.TcpConnectionTestWithTagName(t, ssmClient, tag, dbEndpoint, dbPort, maxRetriesToPollResult, waitBetweenRetries)
	})
	t.Run("TestAppInstanceConnectivityToLoggingService", func(t *testing.T) {
		t.Parallel()
		// get the required resource values using terratest's terraform module
		loggingEndpoint := terraform.Output(t, terraformOptions, "logging_endpoint")
		loggingPort := "443"
		tag := terraform.Output(t, terraformOptions, "instance_name_tag")
		tester.TcpConnectionTestWithTagName(t, ssmClient, tag, loggingEndpoint, loggingPort, maxRetriesToPollResult, waitBetweenRetries)
	})
	t.Run("TestAppInstanceConnectivityToMonitoringService", func(t *testing.T) {
		t.Parallel()
		// get the required resource values using terratest's terraform module
		monitoringEndpoint := terraform.Output(t, terraformOptions, "monitoring_endpoint")
		monitoringPort := "443"
		tag := terraform.Output(t, terraformOptions, "instance_name_tag")
		tester.TcpConnectionTestWithTagName(t, ssmClient, tag, monitoringEndpoint, monitoringPort, maxRetriesToPollResult, waitBetweenRetries)
	})
	t.Run("TestAppInstanceShouldNotHaveConnectivityToPublicInternet", func(t *testing.T) {
		t.Parallel()
		// build a tcp connectivity test case with public endpoint and port
		testCase := tester.NewShellTestCase(fmt.Sprintf("timeout 2 bash -c '</dev/tcp/%s/%s'", "www.example.com", "443"), false)
		target := tester.NewTagNameTarget(terraform.Output(t, terraformOptions, "instance_name_tag"))
		tester.RunTestCaseForTarget(t, ssmClient, testCase, target, maxRetriesToPollResult, waitBetweenRetries)
	})

	t.Run("TestEchoCommandViaCustomTestCase", func(t *testing.T) {
		t.Parallel()
		testCase := tester.NewShellTestCase("aws", false)
		tagNameTarget := tester.NewTagNameTarget(terraform.Output(t, terraformOptions, "instance_name_tag"))
		tester.RunTestCaseForTarget(t, ssmClient, testCase, tagNameTarget, maxRetriesToPollResult, waitBetweenRetries)
	})

}
