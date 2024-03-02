output "docker_repository_url" {
  description = "docker repository url"
  value       = "${aws_ecr_repository.stream_aggregator.repository_url}:8080"
}

output "prod_ecr_tag" {
  description = "tag for image to push to aws ecr"
  value       = local.prod_ecr_tag
}

output "aws_availability_zones" {
  value = data.aws_availability_zones.available_zones.names
}

output "app_url" {
  value = module.load_balancer.app_url
}
