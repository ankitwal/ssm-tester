# ssm-tester
[![CircleCI](https://circleci.com/gh/ankitwal/ssm-tester/tree/main.svg?style=svg)](https://circleci.com/gh/ankitwal/ssm-tester/tree/main)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/ankitwal/ssm-tester/main/LICENSE)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/mod/github.com/ankitwal/ssm-tester/tester?tab=overview)
[![Go Report Card](https://goreportcard.com/badge/github.com/ankitwal/ssm-tester)](https://goreportcard.com/report/github.com/ankitwal/ssm-tester)


Infrastructure testing helper for AWS Resources that uses AWS SSM to remotely execute commands on EC2 machines, to enable infrastructure engineering teams to write 
tests that validate *behaviour*.

## Demo
![automate infrastructure behaviour testing](https://i.ibb.co/cgbthDN/terminal-screenshot.png)

## Why 

### Validating infrastructure without ssh

* Modern cloud architecture is moving away from ssh, and without access to instance it is hard to validate certain behaviour.
* Manually running commands to validate infrastructure correctness is slow and unreliable.
* This means some behaviour may only get tested when we run an application on the provisioned infrastructure app on it, eg. connectivity to database, connectivity to required internet endpoint.
This slows the feedback loop in turn means lower quality. The application delivery team is a consumer/customer of the the infrastructure delivery team. 
An infrastructure delivery team should not have to rely on it's customers to validate to its code.
* Additionally some behaviour is hard to validate, and wont get immediately feedback with even application running correctly on the provisioned infra. Example:
    * Broken connectivity to logging endpoints/service may only get detected if a team member notices missing logs, often these are not even being looked at in lower environments. Or worse it may only be
    detected when instance in production start falling over since their disks have gone to full from failing to flush logs to a remote logging service.

ssm-tester allows infrastructure delivery teams to write tests than can execute custom commands on ec2 instances and hence validate for otherwise hard to test behaviour.

### Testing Behaviour over Configuration 
When we write infrastructure as code - we want to not only test against the configuration we create but also test the our infrastructure for *behaviour*!
Specially when we write infrastructure code in declarative tooling like terraform, tests that validate configuration may have limited value.  
For example, validating for **configuration**:  
* does my security group have outgoing allowed to the RDS Security group
* does my application subnet network ACL have rules allowing outgoing to the RDS Subnet
* does my application subnet network ACL have rules allowing ephemeral ports open for return traffic from RDS subnets
* does my application subnet have a route table attached with routes to the database subnet 

These tests may essentially be a repeat of the configuration specified in our Infrastructure declarative code
and does not validate the behaviour we want to guarantee in our infrastructure.  
Instead it would be better if we could write tests to validate **behaviour**: 
* does my provisioned infrastructure allow my application EC2 instances to connect via TCP to my RDS endpoint. 
    * this would ideally validate that the configuration for security groups, subnets, NACLs, route tables cumulatively allows this behaviour.  
* can my provisioned instance pull a required secret from secrets manager
    * this would validate that required networking configuration + IAM Instance Profile + Role configuration cumulatively allows for this behaviour 

ssm-tester enables users to write automated tests that validate behaviour so infrastructure engineering teams do not have to wait for application teams to report
broken infrastructure or worse wait for incidents in production. 

## Quick Start 

### Requirements 

ssm-tester requires the the EC2 instances to be integrated with [AWS Systems Manager](https://aws.amazon.com/systems-manager/). 
This requires your EC2 instances to all ssm-agent installed( installed by default on Amazon Linux 2), and certain AWS SSM related resources provisioned. 
You may see this [example](examples/simple-example/terraform/main.tf) of a minimum AWS Systems Manager integrated infra, and 
[this AWS Documentation](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-setting-up.html) for a more comprehensive guide.  
If you already use AWS Systems Manager in your AWS Infrastructure then you should be able to use this out of box. Alternatively you may
consider layering on AWS SSM required resources 

### Using ssm-tester/tester

1. Import ssm-tester/tester
    ```go
    	import "github.com/ankitwal/ssm-tester/tester"
    ```
2. Initialise the ssm service client - this is used to the AWS SSM API
    ```go
    	// Initialise AWS SSM client service.
    	ssmClient := tester.NewSSMClientWithDefaultConfig(t)
    ```
3. Initialise retry config - this is used to manage to polling for the test result
    ```go
    	// create retry configuration for
	    retryConfig := tester.NewRetryDefaultConfig()
    ```
4. Write some tests  
    ```go
       t.Run("TestAppInstanceCanConnectToImportantEndpoint", func(t *testing.T) {   
           // 4.1 create a new test with a custom test command
           // this example uses curl, it relies on your target instances having curl installed
           testCase := tester.NewShellTestCase("curl https://www.importantendpoint.com --max-time=2", true)
   
           // 4.2 specify the ec2 instance to target for the test
           target := tester.NewTagNameTarget(terraform.Output(t, terraformOptions, "app_instance_name_tag"))
   
           // 4.3 run the test 
           tester.RunTestCaseForTarget(t, ssmClient, testCase, target, retryConfig)   
       })
    ```
### More examples 

Write some tests with built in [TcpConnectionTestWithTagName](https://pkg.go.dev/github.com/ankitwal/ssm-tester/tester#TcpConnectionTestWithNameTag) helper 
 
```go
    t.Run("TestAppInstanceConnectivityToDatabase", func(t *testing.T) {
        dbEndpoint := "mydb.privatedns" 
        dbPort := "3306" 
        tag := "app_instance_name_tag" 
   
        // run the test
        tester.TcpConnectionTestWithTagName(t, ssmClient, tag, dbEndpoint, dbPort, retryConfig)
    })
```
   
Write some tests with using [terratest](https://terratest.gruntwork.io), please see [examples](examples/simple-example/) for working examples 

```go
	// this example uses terratest to initialise the terraform stack and get output value
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../terraform",
	})

	// init and apply terraform stack ensuring clean up
	t.Cleanup(func() { terraform.Destroy(t, terraformOptions) })
	terraform.InitAndApply(t, terraformOptions)

    t.Run("TestAppInstanceConnectivityToDatabase", func(t *testing.T) {
        // get the required resource values using terratest's terraform module
        dbEndpoint := terraform.Output(t, terraformOptions, "database_endpoint")
        dbPort := terraform.Output(t, terraformOptions, "database_port")
        tag := terraform.Output(t, terraformOptions, "instance_name_tag")
   
        // run the test 
        tester.TcpConnectionTestWithTagName(t, ssmClient, tag, dbEndpoint, dbPort, retryConfig)    
    })
``` 

Write negative tests
```go
    t.Run("TestAppInstanceShouldNOTHaveConnectivityToPublicInternet", func(t *testing.T) {
          // build a tcp connectivity test case with public endpoint and port, 
          // with condition false, i.e the tests passes if the command fails on all target instances
    	testCase := tester.NewShellTestCase(fmt.Sprintf("timeout 2 bash -c '</dev/tcp/%s/%s'", "www.example.com", "443"), false)
   
          // specify the ec2 instance to target for the test
          target := tester.NewTagNameTarget(terraform.Output(t, terraformOptions, "instance_name_tag"))
   
          // run the test
          tester.RunTestCaseForTarget(t, ssmClient, testCase, target, retryConfig)
    })
```
   
Write tests to app instances have the required IAM, and networking configuration to be able to pull secrets that might be required by app
```go
    t.Run("TestAppInstanceShouldNOTHaveConnectivityToPublicInternet", func(t *testing.T) {
          // build a testCase command that validates that the instance has networking and IAM access to a secret that will be required by the application 
          // this relies on aws cli being installed on the instance(AMI) being targeted.  
    	testCase := tester.NewShellTestCase(`aws secretsmanager list-secret-version-ids --secret-id "secret-required-by-app"`), true)
   
          // specify the ec2 instance to target for the test
          target := tester.NewTagNameTarget(terraform.Output(t, terraformOptions, "instance_name_tag"))
   
          // run the test
          tester.RunTestCaseForTarget(t, ssmClient, testCase, target, retryConfig)
    })
```
