# Installing the teller

The teller is and endpoint which interacts directly with IoT devices to stream data and trigger the smart contract run.

### Prerequisites

1. Golang
2. Ipfs

### Building from Source

` go build` inside this repo should ideally work in most scenarios. If it doesn't, please open an issue with the error logs.

### Downloading a prebuilt version

[The builds website](https://builds.openx.solar/fe) has daily builds for opensolar, openx and the teller. Running them should be as simple as running the executable.

### Setting up config params

Duplicate [dummyconfig.yaml](dummyconfig.yaml), rename to config.yaml and replace the relevant values with those relevant to the project.
