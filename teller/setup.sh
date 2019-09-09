#!/bin/bash
whoami
MACHINE_TYPE=`uname -m`
sudo apt-get -y upgrade
sudo apt-get install build-essential

if [ ${MACHINE_TYPE} == 'x86_64' ]; then
  # need to install the 64 bit version of go
  wget https://dl.google.com/go/go1.13.linux-amd64.tar.gz
elif [ ${MACHINE_TYPE} == 'armv6l' ] ; then
  wget https://dl.google.com/go/go1.13.linux-armv6l.tar.gz
else
  wget https://dl.google.com/go/go1.13.linux-386.tar.gz
fi

sudo tar -xvf go1.12.4.linux-amd64.tar.gz
sudo mv go /usr/local
echo 'GOROOT=/usr/local/go' >> ~/.profile
echo 'GOPATH=$HOME' >> ~/.profile
echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
source ~/.profile
which go
go version
mkdir go
go get -v github.com/YaleOpenLab/openx
cd ~/go/src/github.com/YaleOpenLab/openx
go get -v ./...
go build
cp dummyconfig.yaml config.yaml
alias openx = ./$PWD/openx
go get -v github.com/YaleOpenLab/opensolar
cd ~/go/src/github.com/YaleOpenLab/opensolar
go get -v ./...
go build
alias opensolar = ./$PWD/opensolar
cd teller
go get -v ./...
go build
cp dummyconfig.yaml config.yaml
alias teller = ./$PWD/teller
## start teller, openx and opensolar
# env GOOS=linux GOARCH=arm GOARM=5 go build
