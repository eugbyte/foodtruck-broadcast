output "aws_lb_target_group_arn" {
  description = "arn of aws_lb_target_group"
  value       = aws_lb_target_group.target_group.arn
}

output "app_url" {
  description = "the url to access the application"
  value       = aws_alb.application_load_balancer.dns_name
}

output "aws_security_group_id" {
  description = "id of aws_security_group for the load balancer"
  value       = aws_security_group.load_balancer_security_group.id
}
