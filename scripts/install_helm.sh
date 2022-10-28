#!/bin/bash

wget https://get.helm.sh/helm-v3.8.2-linux-amd64.tar.gz
tar xvf helm-*-linux-amd64.tar.gz
sudo mv linux-amd64/helm /usr/local/bin
rm helm-*-linux-amd64.tar.gz
rm -r linux-amd64
