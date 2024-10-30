#!/bin/bash
sudo mv /tmp/cloudwatch-config.json /opt/cloudwatch-config.json
wget https://amazoncloudwatch-agent.s3.amazonaws.com/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i -E ./amazon-cloudwatch-agent.deb
sudo systemctl enable amazon-cloudwatch-agent