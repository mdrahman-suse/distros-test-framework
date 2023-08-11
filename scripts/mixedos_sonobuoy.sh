#!/bin/bash

# Usage: sh mixedos_sonobouy.sh install 0.56.16 amd64
# Usage: sh mixedos_sonobouy.sh delete

if [ -z "$2" ]
then
    version="0.56.16"
else
    version=$2
fi

if [ -z "$3" ]
then
    arch="amd64"
else
    arch=$3
fi

installation(){
    echo "Installing sonobuoy version ${version} for mixedos validation"
    if [ ! -d "my-sonobuoy-plugins" ]; 
    then
        echo "Cloning repo: https://github.com/phillipsj/my-sonobuoy-plugins.git"
        git clone https://github.com/phillipsj/my-sonobuoy-plugins.git
    fi
    wait
    if [ ! -f "sonobuoy_${version}_linux_${arch}.tar.gz" ];
    then
        echo "Downloading sonobouy installer"
        wget -q https://github.com/vmware-tanzu/sonobuoy/releases/download/v${version}/sonobuoy_${version}_linux_${arch}.tar.gz
    fi
    wait
    tar -xvf sonobuoy_${version}_linux_${arch}.tar.gz
    sudo chmod +x sonobuoy && sudo mv sonobuoy /usr/local/bin/sonobuoy
}

deletion(){
    echo "Deleting sonobuoy installer"
    sudo rm -rf my-sonobuoy-plugins
    sudo rm -rf sonobuoy_*
    sudo rm -rf /usr/local/bin/sonobuoy
}

if [ "$1" == "install" ];
then
    installation
elif [ "$1" == "delete" ];
then
    deletion
else
    echo "Invalid argument, please pass required arg [install or delete]"
fi


