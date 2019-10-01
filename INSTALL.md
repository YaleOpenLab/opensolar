# Installing opensolar

Opensolar is built on top of the base [openx](https://github.com/YaleOpenLab/openx) platform and uses the Stellar blockchain as the basis for its smart contracts. Opensolar uses IPFS to store legal documents and anything that requires persistent storage and verifiability. Encrypted private keys, project parameters and all other data are stored on [boltdb](https://github.com/boltdb/bolt) on the same platform instance as [openx.solar](https://openx.solar).

### Operating system

Opensolar has been tested on Ubuntu 16.04 LTS and macOS 10.13+. The build status on other operating systems is unknown. There are no plans to test on alternate operating systems at the moment.

### Prerequisites

1. Golang

To Save you from a few clicks, use the following script to download and install golang.

Linux:
```
sudo apt-get -y upgrade
wget https://dl.google.com/go/go1.12.4.linux-amd64.tar.gz
sudo tar -xvf go1.12.4.linux-amd64.tar.gz
sudo mv go /usr/local
echo 'GOROOT=/usr/local/go' >> ~/.profile
echo 'GOPATH=$HOME' >> ~/.profile
echo 'PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.profile
source ~/.profile
which go
```

MacOS:
```
brew install golang
```
If you don't have brew installed, its highly recommended:
```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

2. IPFS

To install IPFS, please follow this excellent installation guide from the IPFS team: [guide](https://docs.ipfs.io/guides/guides/install/)

### Building from Source

Assuming that you have GOPATH set from installing go in step 1, please run this following script:

```
cd $GOPATH/src/github.com/YaleOpenLab/opensolar
go get -v ./...
go build
```
This gets the necessary dependencies for opensolar and builds the opensolar executable. You're now ready to go provided you have your openx instance ready. For more instructions on how to do that, please refer the openx docs.

### Downloading a prebuilt version

[The builds website](https://builds.openx.solar/fe) has daily builds for opensolar, openx and the teller. Running them should be as simple as running the executable.
