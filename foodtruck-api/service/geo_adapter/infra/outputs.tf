output "zookeeper_connect_string" {
  value = aws_msk_cluster.kafka.zookeeper_connect_string
}

output "bootstrap_brokers_tls" {
  description = "TLS connection host:port pairs"
  value       = aws_msk_cluster.kafka.bootstrap_brokers_tls
}

output "docker_repository_url" {
  description = "docker repository url"
  value       = "${aws_ecr_repository.geo_adapter.repository_url}:8080"
}

output "app_url" {
  value = module.load_balancer.app_url
}
