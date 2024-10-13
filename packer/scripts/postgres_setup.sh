#!/bin/bash

# Ensure that required environment variables are set
if [ -z "$DB_USER" ]; then
  echo "Error: DB_USER is not set."
  exit 1
fi

if [ -z "$DB_PASSWORD" ]; then
  echo "Error: DB_PASSWORD is not set."
  exit 1
fi

if [ -z "$DB_NAME" ]; then
  echo "Error: DB_NAME is not set."
  exit 1
fi


sudo apt install postgresql postgresql-contrib -y

# Create PostgreSQL user and database
sudo -u postgres psql -c "CREATE USER $DB_USER WITH LOGIN CREATEDB PASSWORD '$DB_PASSWORD';"
sudo -u postgres psql -c "CREATE DATABASE $DB_NAME WITH OWNER = $DB_USER;"