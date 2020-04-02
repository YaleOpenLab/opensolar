# Teller

Teller runs on the IoT hub in a particular project (in case of the opensolar platform). The teller is responsible for assimilating data from the Zigbee devices installed on the solar panels, filtering out noise and transmitting resulting data to external third party providers who may be interested to partner on providing various services. This would also be used to get data about the amount that the recipient owes per month and is used to trigger payback.

We need the teller to be secure and tamper resistant, because it is going to face the recipient of a particular project. Broadly, what the teller needs is:

1. The code that is running on the device is same as the one we want to run on the device.
2. The teller doesn't randomly stop and hence reduce the amount owed by the recipient.

## Working

### Start

On start, a couple metrics are recorded for later use by the teller. It also performs a series of functions including authentication.

- Read config from the config.yaml file to get the username, pwhash, seedpwd and other parameters. Error out if some are missing
- Call the platform API to get the project's index
- Logon to the platform using the given credentials
- Runs a refresh login routine in the background in order to continuously update the recipient
- Decrypt the seed using the given seed pwd. Error out if seedpwd can not be derived
- Get Project details and check whether provided device id matches with the id stored on the server
- Store start time and location of the teller
- Get email of the platform so we can use it in API calls to send emergency emails when required.

### Routines

The teller runs some routines in the background to be able to automatically achieve some functions that are deemed ideal for the teller. The list of go routines include:

- Payback: The teller automatically checks whether the recipient should pay back  towards an order and if so, proceeds to pay the required amount with the help of the oracle. If payback fails, it sends an email to the recipient and the platform and depending on severity emails the guarantor and investors.

- Hash Chain - The teller manages to pull in data from from the zigbee device(s) and write(s) it to the `data.txt` file open in RAM. This acts as the handler for the hashchain described below

- Update State - The teller also updates the state of the teller in parallel to updating the hashchain.  It hashes the deviceId and the power consumption data over an interval and commits it to ipfs. It also propagates two transactions on the blockchain with the ipfs hash (along with some padding to distinguish from spam) in the memo fields

- Start Server - The teller also serves a ping endpoint and the hh endpoint for the investor or recipient to check if the teller is alive. This ip should not be ideally exposed to the public since the IoT Hubs are especially vulnerable to DoS attacks.

### Hashchain

The primary function of the teller is to assimilate data on the IoT hub and provide a mechanism to attest that the data that is being reported is correct. This can be done in multiple ways but we take the approach of a hash chain to accomplish what we want to do. The working of the hashchain can be explained in the below steps:

- Open a file named `data.txt` in the home directory of the teller
- Pipe data from the zigbee devices installed to this file
- Monitor the size of the file in parallel to the write operations. Have a threshold for how big this can go (since this is a pi, this can't store 1GB or something in RAM). This is tored in RAM, so that's to be taken into consideration
- Once the file size hits the marked threshold, hash the contents of the file into ipfs and get the header. Replace all the contents of the given file with this dingle header (along with minimal padding in order to be able to identify that this is not a random set of bytes)
- Continue writing the data

### Shutdown

A shutdown maybe triggered due to a bug in the teller or due to manual intervention in case of emergencies. In either case, we would like to notify the platform about this and take a set of steps to ensure that the teller shuts down gracefully.

- Record the blockstamp when the shutdown occurred. This way, no one can shutdown the teller in advance and claim that it had shutdown earlier.
- Commit the blockstamp, device info, device location, start hash, end hash and the hashchain header to ipfs.
- Propagate two transactions which contain the ipfs hash in their memo fields. The first one has a padding "IPFSHASH:" to denote that this isn't a random memo
- Send an email to the recipient with the two transactions and the device Id to inform them about the shutdown so they can contact help in case they did not trigger this.
- Update the hashchain header as described above. The difference in hashchain headers helps us identify when exactly the shutdown occurred (along with the blockcstamp) and helps us filter logs using the `verify.sh` script.

### Daemon Mode

The teller can also be run in daemon mode in case one does not wish to use the CLI interface that is provided. There are some cases in which this makes sense as the developer might not want the recipient or involved entities to meddle with the functioning of the teller. The daemon mode is also preferable when the IoT device does not have a screen attached to it (although the platform / developer might want the CLI to be able to query some information later).

### CLI Function

The CLI that the teller offers is a bit limited in comparison to the emulator since its aimed at recipients instead of investors. Still, the CLI aims to offer critical functions which may be useful for a developer to be able to debug remotely.  This may be expanded in the future depending on feedback from recipients on the usefulness of this function.


## What the teller aims to achieve

The aim of the teller package is to be as light as possible and at the same time collect data and ensure that it can facilitate the triggering of smart contracts which is essential to the project flow. One can also think of the teller as a mini emulator (although it can function as a full emulator if need be) which has some functions automated and not within direct control of users.
