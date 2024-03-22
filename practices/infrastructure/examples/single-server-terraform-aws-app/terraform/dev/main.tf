
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

  key_name = "my-aws-keypair" # Ensure this key pair exists in your AWS account

  # Security groups and networking
  vpc_security_group_ids = ["sg-12345678"] # Replace with your security group ID
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
  name        = "app-sg"
  
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
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
