
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14"
}

provider "aws" {
  region = "us-east-1" # or your preferred region
}

resource "aws_instance" "app_server" {
  ami           = "ami-1234567890abcdef0" # replace with a valid AMI for your region
  instance_type = "t3.xlarge"

  # Go to AWS EC2 key-pair listing
  key_name = "key-0dd9b3a1cbe164237" # Ensure this key pair exists in your AWS account

  # Security groups and networking

  # Security group as defined below, use as is if you want to create a new security group every time
  # Or replace with your preexisting security group ID
  vpc_security_group_ids = aws_security_group.app_sg.id
  subnet_id              = "subnet-12345678" # Replace with your subnet ID

  user_data = <<-EOF
              #!/bin/bash
              sudo apt-get install -y git docker.io docker-compose
              git clone https://github.com/yourusername/yourrepo.git /home/ubuntu/yourrepo
              cd /home/ubuntu/yourrepo
              sudo docker-compose up -d
              EOF

  tags = {
    Name = "DockerComposeServer"
  }
}

resource "aws_security_group" "app_sg" {
  # Display name in console on AWS
  name        = "app-sg"
  
  # use individual ingress blocks 
  ingress {
    # Beginning bound of port range
    from_port   = 80
    # end bound of port range
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # wildcard, lock down to a specific cidr range, anything that can get to it is allowed to.
    # 
  }
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # wildcard, lock down to a specific cidr range, anything that can get to it is allowed to.
  }
  # hard code my home & work ip, or pull out of a parameter store
  # Or add all IPS to a security group.
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # wildcard, lock down to a specific cidr range, anything that can get to it is allowed to.
  }
}


resource "aws_ebs_volume" "db_volume" {
  availability_zone = "us-east-1"
  size              = 100
  type              = "gp3"
  iops              = 3000 # Only necessary for io1 or io2
  throughput        = 125  # Only applicable for gp3
  encrypted         = true
  kms_key_id        = "arn:aws:kms:us-west-2:123456789012:key/abcd1234-a123-456a-a12b-a123b4cd56ef"

  tags = {
    Name = "MyDatabaseVolume"
  }
}
