terraform {
  required_version = ">= 0.15.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.50"
    }
  }
}

provider "aws" {
  region = "ap-southeast-1"
  default_tags {
    tags = {
      Environment = "private-subnets-example-test"
      Project     = "ssm-tester"
      Terraform   = "true"
    }
  }
}
