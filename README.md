# ssm-tester
[![CircleCI](https://circleci.com/gh/ankitwal/ssm-tester/tree/main.svg?style=svg)](https://circleci.com/gh/ankitwal/ssm-tester/tree/main)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/ankitwal/ssm-tester/main/LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/ankitwal/ssm-tester/tester?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/ankitwal/ssm-tester)](https://goreportcard.com/report/github.com/ankitwal/ssm-tester)


Infrastructure testing helper for AWS Resources that uses AWS SSM to remotely execute commands on EC2 machines.
## Demo


## Quick Start 

### Requirements 

ssm-tester requires the the EC2 instances to be integrated with [AWS Systems Manager](https://aws.amazon.com/systems-manager/). 
This requires your EC2 instances to all ssm-agent installed( installed by default on Amazon Linux 2), and certain AWS SSM related resources provisioned. 
You may see this [example](./examples/private-subnets/terraform/main.tf) of a minimum AWS Systems Manager integrated infra, and 
[this AWS Documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-setting-up.html) for a more comprehensive guide.  
If you already use AWS Systems Manager in your AWS Infrastructure then you should be able to use this out of box. Alternatively you may
consider layering on AWS SSM required resources 

### Using tester
1. Import ssm-tester/tester
    ```go
    	import "github.com/ankitwal/ssm-tester/tester"
    ```
2. Initialise the ssm service client
    ```go
    	// Initialise AWS SSM client service.
    	ssmClient := tester.NewSSMClientWithDefaultConfig(t)
    ```
3. 
```go
	// Initialise AWS SSM client service.
	ssmClient := tester.NewSSMClientWithDefaultConfig(t)

	// create retry configuration for
	// the tester the number of times the tester should retry polling for the result of the test command
	maxRetriesToPollResult := 5
	waitBetweenRetries := 3 * time.Second

	t.Run("TestAppInstanceConnectivityToDatabase", func(t *testing.T) {
		// get the required resource values using terratest's terraform module
		dbEndpoint := terraform.Output(t, terraformOptions, "database_endpoint")
		dbPort := terraform.Output(t, terraformOptions, "database_port")
		tag := terraform.Output(t, terraformOptions, "instance_name_tag")
		tester.TcpConnectionTestWithTagName(t, ssmClient, tag, dbEndpoint, dbPort, maxRetriesToPollResult, waitBetweenRetries)
	})
```
