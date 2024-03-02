
# 1. Build the local docker image
resource "null_resource" "run_docker_img" {
  provisioner "local-exec" {
    command = "docker compose --file ${local.monorepo_root_path}/docker-compose.yml up --detach kafka zookeeper foodtruck_geo_adapter"
  }
}

# 2. Push docker image to ECR 
resource "aws_ecr_repository" "geo_adapter" {
  name = "foodtruck-geo-adapter"
}

resource "null_resource" "publish_to_ecr" {
  provisioner "local-exec" {
    # rmb to remove --endpoint-url http://localhost:4566
    # docker tag ${var.container_name} ${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com/${var.container_name}:latest
    command = <<EOF
     	aws ecr get-login-password --region ap-southeast-1 --endpoint-url http://localhost:4566
      | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${local.region}.amazonaws.com
      docker tag ${var.image_name} ${aws_ecr_repository.geo_adapter.repository_url}/:latest
      docker push ${aws_ecr_repository.geo_adapter.repository_url}/:latest
     EOF
  }
  depends_on = [null_resource.run_docker_img]
}


# 3. Create ECS cluster
resource "aws_ecs_cluster" "foodtruck_cluster" {
  name = "foodtruck-cluster"
}

# 4. Create EC2 task
resource "aws_ecs_task_definition" "foodtruck_geo_adapter_task" {
  family                   = "foodtruck-geo-adapter-task" # Naming our first task
  container_definitions    = <<DEFINITION
  [
    {
      "name": "foodtruck-geo-adapter-task",
      "image": "${aws_ecr_repository.geo_adapter.repository_url}",
      "essential": true,
      "portMappings": [
        {
          "containerPort": 6000,
          "hostPort": 6000
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

#5. Create load balancer
module "load_balancer" {
  source  = "./../../../infra/modules/load_balancer"
  name    = "foodtruck_stream_aggregator_lb"
  subnets = aws_subnet.private_subnet.*.id
}

resource "aws_ecs_service" "geo_adapter_service" {
  name            = "geo-adapter-service"                                  # Naming our first service
  cluster         = aws_ecs_cluster.foodtruck_cluster.id                   # Referencing our created Cluster
  task_definition = aws_ecs_task_definition.foodtruck_geo_adapter_task.arn # Referencing the task our service will spin up
  launch_type     = "FARGATE"
  desired_count   = 3 # Setting the number of containers we want deployed to 3

  load_balancer {
    target_group_arn = module.load_balancer.aws_lb_target_group_arn # Referencing our target group
    container_name   = aws_ecs_task_definition.foodtruck_geo_adapter_task.family
    container_port   = 8080 # Specifying the container port
  }

  network_configuration {
    subnets          = aws_subnet.private_subnet.*.id           # multiple subnets
    assign_public_ip = true                                     # Providing our containers with public IPs
    security_groups  = ["${aws_security_group.geo_adapter.id}"] # Setting the security group
  }
}
