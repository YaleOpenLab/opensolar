sudo apt-get -y upgrade
wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
sudo tar -xvf go1.13.5.linux-amd64.tar.gz
sudo mv go /usr/bin/
echo 'GOROOT=/usr/bin/go' >> ~/.profile
echo 'GOPATH=$HOME' >> ~/.profile
echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
source ~/.profile
which go
go version
sudo apt-get update
sudo apt-get install build-essential
mkdir go
go get -v github.com/YaleOpenLab/opensolar
cd ~/go/src/github.com/YaleOpenLab/opensolar
go get -v ./...
cd teller
go get -v ./...
go build