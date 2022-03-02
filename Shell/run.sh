#!/bin/bash
# Run script for Packagebird-Server on Debian Linux 
# Must be executed as super-user
echo "Running Packagebird Server for Linux\n"

# Check if MongoDB already installed
if ! type mongo &> /dev/null ; then
    echo "MongoDB not installed! Please install MongoDB either manually 'https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/' or by running the Packagebird-Server 'install.sh' script\n"
    exit 1
fi

# Start MongoDB processes
echo "Starting MongoDB if not already running"
sudo systemctl daemon-reload
sudo systemctl start mongod

# Ensure file is executable
chmod +x $HOME/bin/Packagebird-Server/Packagebird-Server

# Check if Packagebird-Server already installed
if ! type $HOME/bin/Packagebird-Server/Packagebird-Server &> /dev/null ; then
    echo "Packagebird not installed! Please install using the 'install.sh' script!"
    exit 1
fi

# Start Packagebird process
echo "Starting Packagebird Server process if not already running"
sudo $HOME/bin/Packagebird-Server/Packagebird-Server
