#!/bin/bash

# apt-repository update
sudo apt-get update

# add a packages
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common -y

# get a GPG Key & repository setting (docker) 
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
sudo apt-get update

# check docker is installed & get a docker packages info
sudo apt-cache policy docker-ce

# install docker
sudo apt-get install -y docker-ce docker-ce-cli containerd.io

# start docker
sudo systemctl enable docker
sudo systemctl restart docker

# docker permission error solution
# check with if statements the user use docker non-root or root
# but, should to perceive security warnings that the docker group grants privileges equivalent to the root user.
# @ref from : https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user
sudo groupadd docker
sudo usermod -a -G docker $USER
newgrp docker

# edit done