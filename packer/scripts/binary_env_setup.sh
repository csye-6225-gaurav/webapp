#!/bin/bash

sudo mkdir -p /usr/bin/
sudo mv /tmp/webapp /usr/bin/
# Change the ownership of the webapp file to csye6225
sudo chown csye6225:csye6225 /usr/bin/webapp
sudo chmod +x /usr/bin/webapp
sudo mv /tmp/.env /usr/bin/