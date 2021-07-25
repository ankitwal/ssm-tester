output "vpc_id" {
  value = module.aws_vpc.vpc_id
}
output "database_subnet_group_name" {
  value = module.aws_vpc.database_subnet_group_name
}
output "database_subnet_group" {
  value = module.aws_vpc.database_subnet_group
}
output "instance_name_tag" {
  value = local.project
}
output "database_port" {
  value = module.db.db_instance_port
}
output "database_endpoint" {
  value = module.db.db_instance_address
}
output "monitoring_endpoint" {
  value = format("monitoring.%s.amazonaws.com", local.region)
}
output "logging_endpoint" {
  value = format("logs.%s.amazonaws.com", local.region)
}
