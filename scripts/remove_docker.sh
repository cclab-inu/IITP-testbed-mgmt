#!/bin/bash

# Docker image stop & remove
docker stop $(docker ps -q)
docker rm -f $(docker ps -a -q)
docker rmi -f $(docker images -a -q)
docker network prune -f

systemctl stop docker.service
systemctl stop containerd.service

# Docker uninstallation
# Check the docker is installed
dpkg -l | grep -i docker

# Docker uninstall
sudo apt-get purge -y docker-engine docker docker.io docker-ce docker-ce-cli
sudo apt-get autoremove -y --purge docker-engine docker docker.io docker-ce

# Remove the files that related docker
sudo rm -rf /var/lib/docker /etc/docker
sudo rm /etc/apparmor.d/docker
sudo groupdel docker
sudo rm -rf /var/run/docker.sock

# All the Docker misc files remove
# Caution about general files that named 'docker'
# sudo find / -name "*docker*" -exec `rm -rf` {} +