#!/bin/bash
# This script can be used to verify the hash chain and retrieve past history
# from a single hash
echo "Enter the hash received from the teller: "
read input
hash=$input # get this value from the user
echo "\n"
while [ "$hash" != "" ]
do
  temp=$(ipfs cat $hash |  grep -a IPFSHASHCHAIN)
  set -- junk $temp ; shift # to split the above into separate fields organized by space
  echo "found hash: $2"
  hash=$2
done
