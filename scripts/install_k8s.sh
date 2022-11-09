#!/bin/bash

# Must execute install_docker.sh first.

# apt-repository update
sudo apt-get update

# Swap-off & Disable Firewall
sudo swapoff -a && sudo sed -i '/swap/s/^/#/' /etc/fstab
ufw disable
sudo systemctl stop firewalld
sudo systemctl disable firewalld

# Setting Network Time Protocol
sudo apt install ntp -y
sudo service ntp restart
sudo ntpq -p

# Settings IPTable Proxy
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sudo sysctl --system

## make directory : /etc/docker 
sudo mkdir /etc/docker

# Docker Demon Settings
cat <<EOF | sudo tee /etc/docker/daemon.json
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2"
}
EOF

sudo mkdir -p /etc/systemd/system/docker.service.d
sudo systemctl daemon-reload
sudo systemctl restart docker
sudo systemctl enable docker

# installatino kubelet, kubeadm, kubectl
#@ref : https://kubernetes.io/ko/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
sudo apt-get update

sudo apt-get -y install kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
sudo systemctl daemon-reload
sudo systemctl restart kubelet
sudo systemctl status kubelet

# kubernetes cheat sheet : https://kubernetes.io/ko/docs/reference/kubectl/cheatsheet/