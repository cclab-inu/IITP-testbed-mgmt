#!/bin/bash

docker stop $(docker ps -q)
docker rm $(docker ps -a -q)
docker rmi $(docker images -q)

systemctl stop docker.service
systemctl stop containerd.service

sudo apt list --installed | grep docker

#editing