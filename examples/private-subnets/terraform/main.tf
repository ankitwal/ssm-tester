locals {
  project = "private-subnets-example-test"
  region  = data.aws_region.aws_current_region.name
}
data "aws_region" "aws_current_region" {
}

module "aws_vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "3.2.0"
  name    = "simple-example"
  cidr    = "10.0.0.0/16"

  azs = ["${local.region}a", "${local.region}b", "${local.region}c"]

  # Db Subnets to host RDS
  database_subnets = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  database_subnet_tags = {
    Name = "${local.project}-db-subnet"
  }
  database_dedicated_network_acl = true

  # Private Subnets to host EC2 Instances
  intra_subnets               = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  intra_dedicated_network_acl = true
  intra_subnet_tags = {
    Name = "${local.project}-intra-subnet"
  }

  vpc_tags = {
    Name = local.project
  }
  enable_dns_hostnames = true
  enable_dns_support   = true
}

module "endpoints" {
  source  = "terraform-aws-modules/vpc/aws//modules/vpc-endpoints"
  version = "3.2.0"

  vpc_id = module.aws_vpc.vpc_id

  # Endpoints requirements for SSM https://docs.aws.amazon.com/systems-manager/latest/userguide/setup-create-vpc.html
  endpoints = {
    ssm = {
      service             = "ssm"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "ssm-vpc-endpoint" }
    },
    ssmmessages = {
      service             = "ssmmessages"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "ssmmessages-vpc-endpoint" }
    },
    ec2messages = {
      service             = "ec2messages"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "ec2messages-vpc-endpoint" }
    },
    kms = {
      service             = "kms"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "kms-vpc-endpoint" }
    },
    ec2 = {
      service             = "ec2"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "ec2-vpc-endpoint" }
    },
    s3 = {
      service         = "s3"
      service_type    = "Gateway"
      route_table_ids = [module.aws_vpc.default_route_table_id]
      tags            = { Name = "s3-vpc-endpoint" }
    },

    ## endpoints for application
    # Logging
    cloudwatch_logs = {
      service             = "logs"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "cloudwatchlog-vpc-endpoint" }
    },
    # Monitoring
    cloudwatch_monitoring = {
      service             = "monitoring"
      private_dns_enabled = true
      security_group_ids  = [module.aws_vpc.default_security_group_id]
      subnet_ids          = module.aws_vpc.intra_subnets
      tags                = { Name = "cloudwatchlog-vpc-endpoint" }
    }
  }
}
