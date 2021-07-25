# SSM-Tester Example

This is a simple example of infrastructure 
* An autoscaling group deployed in a private subnet
* A RDS database deployed in an internal subnet
* A monitoring endpoint for the application instances in the autoscaling group to forward metrics
* A logging endpoint for the application instances in the autoscaling group to forward logs 

## How to use
```shell script
go test -v ./test 
```
> **WARNING**: This example tests creates actual AWS resources and may incur costs. The tests should clean up
resources it creates after itself. In the case of an unexpected failure to clean up resources please ensure 
that you remove any unwanted AWS resources to avoid incurring additional costs.  
