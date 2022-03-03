#!/bin/bash
# Install script for Packagebird-Server component for Packagebird 
# MongoDB dependency requires Ubuntu Bionic Version 18.04 to be used
# Script must be executed as super-user
echo "Installing Packagebird Server for Linux\n"

# Housing-cleaning if running inside a Docker container
apt-get update && apt-get install -y sudo
sudo apt-get install systemctl
sudo apt-get install wget
sudo apt-get install curl

# Create directory for Server binary
mkdir -p $HOME/bin/Packagebird-Server

# Download file from host, changes wtih each version
curl -L https://github.com/wguner/Dependency-Management/releases/download/v0.1.0-beta/Packagebird-Server-LinuxAMD64.gz > $HOME/bin/Packagebird-Server/Packagebird-Server.gz

# Extract archive contents
gzip -d $HOME/bin/Packagebird-Server/Packagebird-Server.gz

# Check if MongoDB already installed
if type mongo &> /dev/null; then
    echo "MongoDB already installed.\nTo start the server, run the 'run.sh' script"
    exit 0
fi

# Install MongoDB if not already present
# Based on instructions provided by MongoDB
echo "Installing MongoDB"

# Modify keys
sudo apt-get install gnupg
wget -qO - https://www.mongodb.org/static/pgp/server-5.0.asc | sudo apt-key add -

# Echo contents to the terminal and modify list
echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/5.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-5.0.list

# Update local package database
sudo apt-get update

# Finally install the database
sudo apt-get install -y mongodb-org

# Not
echo "MongoDB installed. To start the server, run the 'run.sh' script"
exit 0