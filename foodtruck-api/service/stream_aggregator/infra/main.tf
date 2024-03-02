terraform {
  required_version = ">= 1.4.6"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
    archive = {
      source = "hashicorp/archive"
    }
    null = {
      source = "hashicorp/null"
    }
  }
}

provider "aws" {
  region = local.region

  default_tags {
    tags = {
      app = "foodtruck-stream-aggregator"
    }
  }
}

locals {
  monorepo_root_path = "./../../../"
  region             = "ap-southeast-1"
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name               = "ecs_task_execution_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}


/**
1. Dockerize the app.
2. Create an image repository on AWS ECR and push the image.
3. Create an AWS ECS cluster.
4. Create an AWS ECS task.
5. Create an AWS ECS service.
6. Create a load balancer.
*/

# 1. Build the local docker image
resource "null_resource" "run_docker_img" {
  provisioner "local-exec" {
    command = "docker compose --file ${local.monorepo_root_path}/docker-compose.yml up --detach foodtruck_stream_aggregator"
  }
}

# 2. Push docker image to ECR 
resource "aws_ecr_repository" "stream_aggregator" {
  name = "foodtruck-stream-aggregator"
}

// get the access to the effective Account ID, User ID, and ARN in which Terraform is authorized.
data "aws_caller_identity" "current" {}

locals {
  prod_ecr_tag = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com/${var.image_name}:latest"
}

resource "null_resource" "publish_to_ecr" {
  provisioner "local-exec" {
    # rmb to remove --endpoint-url http://localhost:4566
    # docker tag ${var.container_name} ${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com/${var.container_name}:latest
    command = <<EOF
     	aws ecr get-login-password --region ap-southeast-1 --endpoint-url http://localhost:4566
      | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com
      docker tag ${var.image_name} ${aws_ecr_repository.stream_aggregator.repository_url}/:latest
      docker push ${aws_ecr_repository.stream_aggregator.repository_url}/:latest
     EOF
  }
  depends_on = [null_resource.run_docker_img]
}

# 3. Create ECS cluster
resource "aws_ecs_cluster" "foodtruck_cluster" {
  name = "foodtruck-cluster"
}

# 4. Create EC2 task
resource "aws_ecs_task_definition" "foodtruck_stream_aggregator_task" {
  family                   = "foodtruck-stream-aggregator-task" # Naming our first task
  container_definitions    = <<DEFINITION
  [
    {
      "name": "foodtruck-stream-aggregator-task",
      "image": "${aws_ecr_repository.stream_aggregator.repository_url}",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 8080,
          "hostPort": 8080
        }
      ],
      "memory": 512,
      "cpu": 256
    }
  ]
  DEFINITION
  requires_compatibilities = ["FARGATE"] # Stating that we are using ECS Fargate
  network_mode             = "awsvpc"    # Using awsvpc as our network mode as this is required for Fargate
  memory                   = 512         # Specifying the memory our container requires
  cpu                      = 256         # Specifying the CPU our container requires
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
}

# 5. Configure VPC networking
data "aws_availability_zones" "available_zones" {
  state = "available"
}

# Providing a reference to our default subnets
resource "aws_default_subnet" "default_subnet" {
  count             = 3
  availability_zone = data.aws_availability_zones.available_zones.names[count.index]
}

#6. Create load balancer
module "load_balancer" {
  source  = "./../../../infra/modules/load_balancer"
  name    = "foodtruck_stream_aggregator_lb"
  subnets = aws_default_subnet.default_subnet.*.id
}

# 7. Create ECS service
resource "aws_ecs_service" "foodtruck_stream_aggregator_service" {
  name            = "stream-aggregator-service"                                  # Naming our first service
  cluster         = aws_ecs_cluster.foodtruck_cluster.id                         # Referencing our created Cluster
  task_definition = aws_ecs_task_definition.foodtruck_stream_aggregator_task.arn # Referencing the task our service will spin up
  launch_type     = "FARGATE"
  desired_count   = 3 # Setting the number of containers we want deployed to 3

  load_balancer {
    target_group_arn = module.load_balancer.aws_lb_target_group_arn # Referencing our target group
    container_name   = aws_ecs_task_definition.foodtruck_stream_aggregator_task.family
    container_port   = 8080 # Specifying the container port
  }

  network_configuration {
    subnets          = aws_default_subnet.default_subnet.*.id              # multiple subnets
    assign_public_ip = true                                                # Providing our containers with public IPs
    security_groups  = ["${aws_security_group.service_security_group.id}"] # Setting the security group
  }
}

# Creating a security group for ECS cluster, to make sure only traffic from load balancer is allowed:
resource "aws_security_group" "service_security_group" {
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"
    # Only allowing traffic in from the load balancer security group
    security_groups = ["${module.load_balancer.aws_security_group_id}"]
  }

  egress {
    from_port   = 0             # Allowing any incoming port
    to_port     = 0             # Allowing any outgoing port
    protocol    = "-1"          # Allowing any outgoing protocol 
    cidr_blocks = ["0.0.0.0/0"] # Allowing traffic out to all IP addresses
  }
}
