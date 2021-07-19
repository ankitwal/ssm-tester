terraform {
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
