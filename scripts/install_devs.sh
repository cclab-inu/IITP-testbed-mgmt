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

# Installation Go Binaries (if you want recent version of golang, uncomment under lines and comment fixed version download lines)
# goBinary=$(curl -s https://go.dev/dl/ | grep linux | head -n 1 | cut -d'"' -f4 | cut -d"/" -f3)
# wget --quiet https://dl.google.com/go/$goBinary -O /tmp/tools/$goBinary
# sudo tar -C /usr/local -xzf /tmp/tools/$goBinary

# Installation Go Binaries (Fixed Verision : 1.18.8)
wget --quiet https://go.dev/dl/go1.18.8.linux-amd64.tar.gz -O /tmp/tools/go1.18.8.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf /tmp/tools/go1.18.8.linux-amd64.tar.gz

# Environment var settings
echo "export GOPATH=\$HOME/go" >> ~/.profile
echo "export GOROOT=/usr/local/go" >> ~/.profile
echo "export PATH=\$PATH:/usr/local/go/bin:\$HOME/go/bin" >> ~/.profile
. ~/.profile

# Installation Apparmor
sudo apt-get install -y apparmor apparmor-utils

sudo rm -rf /tmp/tools