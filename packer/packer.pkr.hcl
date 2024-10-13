packer {
  required_plugins {
    amazon = {
      version = ">= 1.2.8"
      source  = "github.com/hashicorp/amazon"
    }
  }
}


source "amazon-ebs" "test" {
  profile       = "dev"
  region        = "us-east-1"
  source_ami    = "ami-0866a3c8686eaeeba"
  ami_name      = "test-ami-{{timestamp}}"
  instance_type = "t2.micro"
  subnet_id     = "subnet-08d00738e387c298a"
  ssh_username  = "ubuntu"
}

build {
  name    = "test-ami-builder"
  sources = ["source.amazon-ebs.test"]


  provisioner "shell" {
    inline = [
      "echo set debconf to Noninteractive",
      "echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections",
      "sudo apt-get update -y && sudo apt-get upgrade -y",
      "sudo apt install postgresql postgresql-contrib -y",
    ]
  }
  # Create PostgreSQL user and database
  provisioner "shell" {
    inline = [
      "sudo -u postgres psql -c \"CREATE USER clouduser WITH LOGIN CREATEDB PASSWORD 'CloudUser123';\"",
      "sudo -u postgres psql -c \"CREATE DATABASE cloud_db WITH OWNER = clouduser;\""
    ]
  }
  provisioner "file" {
    source      = "webapp"
    destination = "/tmp/webapp"
  }
  provisioner "shell" {
    inline = [
      "sudo mkdir -p /usr/bin/",
      "sudo mv /tmp/webapp /usr/bin/",
      "sudo chmod +x /usr/bin/webapp"
    ]
  }
}
