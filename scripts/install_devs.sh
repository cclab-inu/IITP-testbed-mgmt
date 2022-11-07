#!/bin/bash

# Ubuntu Upgrade
sudo apt-get update && sudo apt-get upgrade -y
sudo apt-get install wget -y

# net-tools install
sudo apt-get install net-tools -y

# vim install
sudo apt-get install vim -y

# openSSH install
sudo apt-get install openssh-server -y

# temporarily using build files directory
mkdir -p /tmp/tools
cd /tmp/tools

# Installation Go Binaries
goBinary=$(curl -s https://go.dev/dl/ | grep linux | head -n 1 | cut -d'"' -f4 | cut -d"/" -f3)
wget --quiet https://dl.google.com/go/$goBinary -O /tmp/build/$goBinary
sudo tar -C /usr/local -xzf /tmp/build/$goBinary

# Environment var settings
echo "export GOPATH=\$HOME/go" >> ~/.profile
echo "export GOROOT=/usr/local/go" >> ~/.profile
echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> ~/.profile

# Installation Apparmor
sudo apt-get install -y apparmor apparmor-utils

sudo rm -rf /tmp/tools