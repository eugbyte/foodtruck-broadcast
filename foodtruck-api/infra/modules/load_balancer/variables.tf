variable "subnets" {
  description = "List of default subnet ids."
  type        = list(string)
}

variable "name" {
  description = "name of load balancer"
  type        = string
}
