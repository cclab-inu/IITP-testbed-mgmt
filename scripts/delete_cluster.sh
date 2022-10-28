#!/bin/bash

sudo kubeadm reset -f
kubectl config delete-context kubernetes-admin@kubernetes
sudo rm -rf $HOME/.kube
sudo rm $HOME/k8s_init.log
sudo /usr/local/bin/cilium uninstall