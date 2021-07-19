output "vpc_id" {
  value = module.aws_vpc.vpc_id
}
output "database_subnet_group_name" {
  value = module.aws_vpc.database_subnet_group_name
}
output "database_subnet_group" {
  value = module.aws_vpc.database_subnet_group
}
output "database_endpoint" {
  value = module.db.db_instance_endpoint
}
