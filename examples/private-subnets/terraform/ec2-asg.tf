data "aws_ami" "amazon_latest_al2_ami_id" {
  most_recent = true
  owners = [
  "amazon"]

  filter {
    name = "name"
    values = [
    "amzn*"]
  }
}

locals {
  capacity = length(module.aws_vpc.intra_subnets)
}

module "asg" {
  source  = "terraform-aws-modules/autoscaling/aws"
  version = "~> 4.0"

  # Autoscaling group
  name = local.project

  min_size            = local.capacity
  max_size            = local.capacity
  desired_capacity    = local.capacity
  health_check_type   = "EC2"
  vpc_zone_identifier = module.aws_vpc.intra_subnets

  # Ensure ASG Instances are refreshed on change in launch template
  instance_refresh = {
    strategy = "Rolling"
  }

  # Launch template
  lt_name                = local.project
  update_default_version = true

  use_lt    = true
  create_lt = true

  image_id      = data.aws_ami.amazon_latest_al2_ami_id.id
  instance_type = "t3.micro"
  ebs_optimized = true

  # Maps to vpc_security_group_ids
  security_groups = []
  tags = [
    {
      key                 = "Environment"
      value               = "test"
      propagate_at_launch = true
    },
    {
      key                 = "Project"
      value               = local.project
      propagate_at_launch = true
    },
  ]
  iam_instance_profile_arn = aws_iam_instance_profile.ec2_instance.arn
}

# IAM Policy Required for SSM
data "aws_iam_policy" "ssm" {
  arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}
data "aws_iam_policy_document" "ssm_s3" {
  statement {
    effect = "Allow"
    actions = [
      "s3:GetObject"
    ]
    resources = [
      "arn:aws:s3:::aws-ssm-${data.aws_region.aws_current_region.name}/*",
      "arn:aws:s3:::aws-windows-downloads-${data.aws_region.aws_current_region.name}/*",
      "arn:aws:s3:::amazon-ssm-${data.aws_region.aws_current_region.name}/*",
      "arn:aws:s3:::amazon-ssm-packages-${data.aws_region.aws_current_region.name}/*",
      "arn:aws:s3:::${data.aws_region.aws_current_region.name}-birdwatcher-prod/*",
      "arn:aws:s3:::aws-ssm-distributor-file-${data.aws_region.aws_current_region.name}/*",
      "arn:aws:s3:::aws-ssm-document-attachments-${data.aws_region.aws_current_region.name}/*",
      "arn:aws:s3:::patch-baseline-snapshot-${data.aws_region.aws_current_region.name}/*"
    ]
  }
}

resource "aws_iam_policy" "ssm_s3" {
  name = "${local.project}-ssm-s3"
  policy = data.aws_iam_policy_document.ssm_s3.json
}

data "aws_iam_policy_document" "instance_assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}
resource "aws_iam_role" "ec2_instance" {
  name               = "${local.project}-instance-role"
  assume_role_policy = data.aws_iam_policy_document.instance_assume_role_policy.json
}
resource "aws_iam_role_policy_attachment" "ec2_instance_ssm" {
  policy_arn = data.aws_iam_policy.ssm.arn
  role       = aws_iam_role.ec2_instance.name
}
resource "aws_iam_role_policy_attachment" "ec2_instance_ssm_s3" {
  policy_arn = aws_iam_policy.ssm_s3.arn
  role = aws_iam_role.ec2_instance.name
}
resource "aws_iam_instance_profile" "ec2_instance" {
  role = aws_iam_role.ec2_instance.id
  name = "${local.project}-instance-profile"
}
