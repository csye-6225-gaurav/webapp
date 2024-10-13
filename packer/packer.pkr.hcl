packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "profile" {
  type        = string
  description = "The AWS CLI profile to use"
}

variable "region" {
  type        = string
  description = "The AWS region"
}

variable "source_ami" {
  type        = string
  description = "The base AMI to use for creating the new image"
}

variable "instance_type" {
  type        = string
  description = "The instance type for the build"
}

variable "subnet_id" {
  type        = string
  description = "The Subnet ID in which the instance will be created"
}

variable "ssh_username" {
  type        = string
  description = "The SSH username for the instance"
  default     = "ubuntu"
}

variable "db_user" {
  type        = string
  description = "PostgreSQL database user"
}

variable "db_password" {
  type        = string
  description = "PostgreSQL database password"
}

variable "db_name" {
  type        = string
  description = "PostgreSQL database name"
}

variable "ami_name" {
  type        = string
  description = "Prefix for the AMI name"
}


source "amazon-ebs" "ubuntu" {
  profile       = var.profile
  region        = var.region
  source_ami    = var.source_ami
  ami_name      = var.ami_name
  instance_type = var.instance_type
  subnet_id     = var.subnet_id
  ssh_username  = var.ssh_username
}

build {
  name    = "test-ami-builder"
  sources = ["source.amazon-ebs.ubuntu"]


  provisioner "shell" {
    script = "./scripts/create_user.sh"
  }
  # Create PostgreSQL user and database
  provisioner "shell" {
    environment_vars = [
      "DB_USER= ${var.db_user}",
      "DB_PASSWORD=${var.db_password}",
      "DB_NAME=${var.db_name}"
    ]
    script = "./scripts/postgres_setup.sh"
  }
  provisioner "file" {
    source      = "webapp"
    destination = "/tmp/webapp"
  }
  provisioner "file" {
    source      = "../.env"
    destination = "/tmp/.env"
  }
  provisioner "shell" {
    script = "./scripts/binary_env_setup.sh"
  }
  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }
  provisioner "shell" {
    script = "./scripts/systemd_conf.sh"
  }
}
