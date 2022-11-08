#!/bin/bash

# Must execute install_docker.sh first.

# apt-repository update
sudo apt-get update

# Check the Docker Status
systemctl status docker

# Swap-off & Disable Firewall
sudo swapoff -a && sudo sed -i '/swap/s/^/#/' /etc/fstab
ufw disable
# sudo systemctl stop firewalld
# sudo systemctl disable firewalld

# Settings IPTable Proxy
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
EOF
sudo sysctl --system

#@ref : https://kubernetes.io/ko/docs/setup/production-environment/tools/kubeadm/install-kubeadm/
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF
apt-get update

# installatino kubelet, kubeadm, kubectl
apt-get -y install kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl
systemctl daemon-reload
systemctl restart kubelet
systemctl status kubelet