# SSM-Tester Example

This example project is intended to demonstrate the use of ssm-tester to write automated infrastructure tests that 
validate *behaviour*. 

## How to use
```shell script
go test -v -timeout 30m ./test 
```
> **WARNING**: This example tests creates actual AWS resources and may incur costs. The tests should clean up
resources it creates after itself. In the case of an unexpected failure to clean up resources please ensure 
that you remove any unwanted AWS resources to avoid incurring additional costs.  

## Contents
```
examples/private-subnets
├── README.md                       # this file  
├── terraform                       # sample terraform project
│   ├── ec2-asg.tf
│   ├── init.tf
│   ├── main.tf
│   ├── outputs.tf
│   ├── rds.tf
│   ├── terraform.tfstate
│   └── terraform.tfstate.backup
└── test                            # directory for infrastructure tests
    └── infra_test.go               # go test file with infra tests for the above terraform provisioned infra
```

This is a simple example of infrastructure 
* An autoscaling group deployed in a private subnet
* A RDS database deployed in an internal subnet
* A monitoring endpoint for the application instances in the autoscaling group to forward metrics
* A logging endpoint for the application instances in the autoscaling group to forward logs

The tests in [infra_test.go](./test/infra_test.go) defines automated infrastructure tests using ssm-tester to validate 
some *behaviour* of the above provisioned infra. 
* do the application have tcp connectivity to the database
* do the application instances have tcp connectivity to the looging, monitoring endpoint
* ensure that application instances do not have connectivity to a public internet endpoint

