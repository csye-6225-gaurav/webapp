#!/bin/bash


echo set debconf to Noninteractive
echo 'debconf debconf/frontend select Noninteractive' | sudo debconf-set-selections
sudo apt-get update -y && sudo apt-get upgrade -y
echo "Creating the csye6225 group..."
sudo groupadd csye6225
sudo useradd -s /sbin/nologin -M -g csye6225 csye6225