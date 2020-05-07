#!/bin/bash

echo "welcome to openx / opensolar"
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "OS: Linux"
    sudo apt-get -y upgrade
    GO=$(which go)
    if [[ "$GO" != "" ]] ; then
        echo "go installed"
    else 
        echo "go not installed, installing go"
        wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
        sudo tar -xvf go1.13.5.linux-amd64.tar.gz
        sudo mv go /usr/bin/
        echo 'GOROOT=/usr/bin/go' >> ~/.profile
        echo 'GOPATH=$HOME' >> ~/.profile
        echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
        source ~/.profile
    fi
    sudo apt-get update
    sudo apt-get install build-essential
elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "OS: macOS"
    BREW=$(which brew)
    echo "BREW: " $BREW
    if [[ "$BREW" != "" ]] ; then
        echo "brew installed"
    else 
        echo "brew not installed, installing brew"
        BREWIN=$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)
        echo $BREWIN
    fi
    which brew
    GO=$(which go)
    mkdir $GOPATH
    if [[ "$GO" != "" ]] ; then
        echo "go installed"
    else 
        echo "go not installed, installing go"
        brew install golang
    fi
fi
which go
go version
mkdir $GOPATH
go get -v github.com/YaleOpenLab/opensolar
go get -v github.com/YaleOpenLab/openx
cd $GOPATH/src/github.com/YaleOpenLab/opensolar
go get -v ./...
cd ../openx/
go get -v ./...
mv dummyconfig.yaml config.yaml
cd ../opensolar/
mv dummyconfig.yaml config.yaml
cp bootstrap/platform.sh ../openx/
cp opensolar ../openx/
cd ../openx/
chmod +x platform.sh
./openx && ./platform.sh && ./opensolar