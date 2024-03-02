# Providing a reference to our default VPC
resource "aws_default_vpc" "default_vpc" {
}

# A single URL provided by our load balancer that, behind the scenes, will redirect our traffic to our underlying containers
resource "aws_alb" "application_load_balancer" {
  name               = "foodtruck-lb" # Naming our load balancer
  load_balancer_type = "application"
  subnets            = var.subnets # multiple subnets # aws_default_subnet.default_subnet.*.id
  # Referencing the security group
  security_groups = ["${aws_security_group.load_balancer_security_group.id}"]
}

# Creating a security group for the load balancer:
# A security group acts as a firewall that controls the traffic allowed to and from the resources in your virtual private cloud (VPC). 
# Specify the ports and protocols to allow for inbound traffic and for outbound traffic.
resource "aws_security_group" "load_balancer_security_group" {
  ingress {
    from_port   = 80 # Allowing traffic in from port 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # Allowing traffic in from all sources
  }

  egress {
    from_port   = 0             # Allowing any incoming port
    to_port     = 0             # Allowing any outgoing port
    protocol    = "-1"          # Allowing any outgoing protocol 
    cidr_blocks = ["0.0.0.0/0"] # Allowing traffic out to all IP addresses
  }
}

# Target groups route requests to one or more registered targets, such as EC2 instances, using the protocol and port number that you specify
resource "aws_lb_target_group" "target_group" {
  name        = "target-group"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_default_vpc.default_vpc.id # Referencing the default VPC.
  health_check {
    healthy_threshold   = "2"
    unhealthy_threshold = "6"
    interval            = "30"
    matcher             = "200,301,302"
    path                = "/"
    protocol            = "HTTP"
    timeout             = "5"
  }
}

# A listener is a process that checks for connection requests, using the protocol and port that you configure. 
# The rules defined for a listener determine how the load balancer routes requests to the targets in one or more target groups.
resource "aws_lb_listener" "listener" {
  load_balancer_arn = aws_alb.application_load_balancer.arn # Referencing our load balancer
  port              = "80"
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.target_group.arn # Referencing our target group
  }
}
