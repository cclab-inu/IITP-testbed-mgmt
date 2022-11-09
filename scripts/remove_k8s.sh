#!/bin/bash
kubeadm reset

systemctl stop kubelet
systemctl stop docker

rm -rf /var/lib/cni /var/lib/kubelet/* /var/lib/etcd /run/flannel /etc/cni /etc/kubernetes ~/.kub
sudo apt-get purge -y kubeadm kubectl kubelet kubernetes-cni kube*
sudo apt-get autoremove -y kubeadm kubectl kubelet kubernetes-cni kube*