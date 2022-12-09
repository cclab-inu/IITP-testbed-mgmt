#!/bin/bash

# set default
if [ "$CNI" == "" ]; then
    CNI=cilium
fi

# use docker as default CRI
if [ "$CRI_SOCKET" == "" ]; then
    # if docker, let kubeadm figure it out
    if [ -S /var/run/docker.sock ]; then
        CRI_SOCKET=""
    elif [ -S /var/run/containerd/containerd.sock ]; then
        CRI_SOCKET=unix:///var/run/containerd/containerd.sock
    elif [ -S /var/run/crio/crio.sock ]; then
        CRI_SOCKET=unix:///var/run/crio/crio.sock
    fi
fi

# check supported CNI
if [ "$CNI" != "flannel" ] && [ "$CNI" != "weave" ] && [ "$CNI" != "calico" ] && [ "$CNI" != "cilium" ]; then
    echo "Usage: CNI={flannel|weave|calico|cilium} CRI_SOCKET=unix:///path/to/socket_file MASTER={true|false} $0"
    exit
fi

# turn off swap
sudo swapoff -a

# activate br_netfilter
sudo modprobe br_netfilter
sudo bash -c "echo '1' > /proc/sys/net/bridge/bridge-nf-call-iptables"
sysctl=`sudo cat /etc/sysctl.conf`
if [[ "${sysctl}" != *"net.bridge.bridge-nf-call-iptables=1"* ]];then
    sudo bash -c "echo 'net.bridge.bridge-nf-call-iptables=1' >> /etc/sysctl.conf"
fi

# initialize the master node
if [ "$CNI" == "calico" ]; then
    sudo kubeadm init --cri-socket=$CRI_SOCKET --pod-network-cidr=192.168.0.0/16 | tee -a $HOME/k8s_init.log
else # weave, flannel, cilium
    sudo kubeadm init --cri-socket=$CRI_SOCKET --pod-network-cidr=10.244.0.0/16 | tee -a $HOME/k8s_init.log
fi

# make kubectl work for non-root user
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $USER:$USER $HOME/.kube/config
export KUBECONFIG=$HOME/.kube/config
bashrc=`cat $HOME/.bashrc`
if [[ "${bashrc}" != *"export KUBECONFIG=$HOME/.kube/config"* ]];then
    echo "export KUBECONFIG=$HOME/.kube/config" | tee -a $HOME/.bashrc
fi

if [ "$CNI" == "flannel" ]; then
    # install a pod network (flannel)
    kubectl apply -f https://raw.githubusercontent.com/flannel-io/flannel/v0.17.0/Documentation/kube-flannel.yml
elif [ "$CNI" == "weave" ]; then
    # install a pod network (weave)
    export kubever=$(kubectl version | base64 | tr -d '\n')
    kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$kubever"
elif [ "$CNI" == "calico" ]; then
    # install a pod network (calico)
    kubectl apply -f https://projectcalico.docs.tigera.io/manifests/calico.yaml
elif [ "$CNI" == "cilium" ]; then
    # install a pod network (cilium)
    curl -L --remote-name-all https://github.com/cilium/cilium-cli/releases/latest/download/cilium-linux-amd64.tar.gz{,.sha256sum}
    sha256sum --check cilium-linux-amd64.tar.gz.sha256sum
    sudo tar xzvfC cilium-linux-amd64.tar.gz /usr/local/bin
    rm cilium-linux-amd64.tar.gz{,.sha256sum}
    /usr/local/bin/cilium install
fi

# enable hubble 
cilium hubble enable
