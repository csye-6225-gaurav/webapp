#!/bin/bash
# Disabling SELinux
echo "SELINUX=permissive" | sudo tee /etc/selinux/config
sudo setenforce 0
sudo mv /tmp/webapp.service /etc/systemd/system/webapp.service
sudo mkdir -p /var/log/webapp/
sudo chown csye6225:csye6225 /var/log/webapp/

sudo systemctl daemon-reload
sudo systemctl enable webapp.service
sudo systemctl start webapp.service