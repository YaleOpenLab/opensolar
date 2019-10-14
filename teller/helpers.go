package main

import (
	"bufio"
	//"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	//"encoding/json"

	ipfs "github.com/Varunram/essentials/ipfs"
	utils "github.com/Varunram/essentials/utils"
	//	rpc "github.com/YaleOpenLab/openx/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	oracle "github.com/YaleOpenLab/opensolar/oracle"
)

// refreshLogin runs once every 5 minutes in order to fetch the latest recipient details
// for eg, if the recipient loads his balance on the platform, we need it to be reflected on
// the teller
func refreshLogin(username string, pwhash string) error {
	var err error
	for {
		time.Sleep(consts.TellerPollInterval)
		err = login(username, pwhash)
		if err != nil {
			log.Println(err)
		}
	}
}

// EndHandler runs when the teller shuts down. Records the start time and location of the
// device in ipfs and commits it as two transactions to the blockchain
func endHandler() error {
	log.Println("Gracefully shutting down, please do not press any button in the process")
	var err error

	NowHash, err = getLatestBlockHash()
	if err != nil {
		log.Println(err)
	}

	hashString := "Device Shutting down. Info: " + DeviceInfo + " Device Location: " + DeviceLocation +
		" Device Unique ID: " + DeviceId + " " + "Start hash: " + StartHash + " Now hash: " + NowHash +
		"Ipfs HashChainHeader: " + HashChainHeader
	// note that we don't commit the latest hash chain header's hash here because this gives us a tighter timeline
	// to audit what really happened
	ipfsHash, err := ipfs.IpfsAddString(hashString)
	if err != nil {
		log.Println(err)
	}
	memo := "IPFSHASH: " + ipfsHash

	tx1, tx2, err := splitAndSend2Tx(memo)
	if err != nil {
		log.Fatal("could not split and send 2tx: ", err)
	}

	err = sendDeviceShutdownEmail(tx1, tx2)
	if err != nil {
		log.Fatal("could not send device shutdown email: ", err)
	}

	log.Println("sent device shutdown notice")
	commitDataShutdown()
	// save last known state of the system in the recipient's list of known hashes
	// Call this last since there would still be data that we want ot measure when the above commands
	// are still running
	return nil
	// have a return because we don't want to sigint while we send emails and stuff
}

func splitAndSend2Tx(memo string) (string, string, error) {
	// 10 padding chars + 46 (ipfs hash length) characters
	firstHalf := memo[:28]
	secondHalf := memo[28:]
	tx1, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, 1, firstHalf)
	if err != nil {
		return "", "", err
	}
	time.Sleep(2 * time.Second)
	tx2, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, 1, secondHalf)
	if err != nil {
		return "", "", err
	}
	log.Printf("tx hash: %s, tx2 hash: %s", tx1, tx2)
	return tx1, tx2, nil
}

// so the teller will be run on the hub and has some data that the platform might need
// The teller must serve some data to other entities as well. So we need to run a server for that
// and it must be over tls for preventing mitm attacks
func checkPayback() {
	for {
		log.Println("Paybck interval reached. Paying back automatically")
		assetName := LocalProject.DebtAssetCode
		amount := oracle.MonthlyBill() // TODO: consumption data must be accumulated from zigbee in the future

		err := projectPayback(assetName, amount)
		if err != nil {
			log.Println("Error while paying amount back", err)
			sendDevicePaybackFailedEmail()
		}
		time.Sleep(time.Duration(time.Duration(LocalProject.PaybackPeriod) * consts.OneWeekInSecond))
	}
}

// updateState hashes the current state of the teller into ipfs and commits the ipfs hash
// to the blockchain
func updateState(trigger bool) {
	for {
		subcommand := "Energy production data for this cycle: " + "100" + "W"
		// no spaces since this won't allow us to send in a requerst which has strings in it
		// TODO: replace this with real data rather than fake data that we have here
		// use rest api for ipfs since this may be too heavy to load on a pi. If not, we can shift
		// this to the pi as well to achieve a s tate of good decentralization of information.
		ipfsHash, err := ipfs.IpfsAddString("Device ID: " + DeviceId + " UPDATESTATE" + subcommand)
		if err != nil {
			log.Println("Error while fetching ipfs hash", err)
			time.Sleep(consts.TellerPollInterval)
		}

		ipfsHash = "STATUPD: " + ipfsHash
		// send _timestamp_ stroops to ourselves, we just pay the network fee of 100 stroops
		// this gives us 10**5 updates per xlm, which is pretty nice, considering that we
		// do about 288 updates a day, this amounts to 347 days' worth updates with 1 XLM
		// memo field restricted to 28 bytes - AAAAAAAAAAAAAAAAAAAAAAAAAAAA
		// we could ideally send the smallest amount of 1 stroop but stellar allows you to
		// send yourself as much money as you want, so we can have any number here
		// we could also time this amount to be the state update number itself.
		// TODO: is this an ideal solution?

		// don't use platform RPCs for interacting with the blockchain
		// But we do need to track this somehow, so maybe hash the device id and "STATUPS: "
		// so we can track if but others viewing the blockchain can't (since the deviceId is assumed
		// to be unique)
		hash1, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, float64(utils.Unix()), ipfsHash[:28])
		if err != nil {
			log.Println(err)
		}

		hash2, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, float64(utils.Unix()), ipfsHash[29:])
		if err != nil {
			log.Println(err)
		}

		// we updated state as hash1 and hash2
		// send email to the platform for this?  maybe overkill
		// TODO: Define structures on the backend that would keep track of this state change
		colorOutput("Updated State: "+hash1+" "+hash2, MagentaColor)
		if trigger {
			break // we trigerred this manually, don't want to keep doing this
		}
		time.Sleep(consts.TellerPollInterval)
	}
}

// TODO and MWTODO: think upon this problem and arrive at a solution. Might be useful to do
// we don't want all data to be public - figure out which parts need to be private and which public
// stream data from the pilot particle instance and write to a file name data.txt
// one can run verify.sh to get a list of all the ipfs hashes in the hashchain (we can make this hashchain header
// public or available to all t he entities involved in the workflow)

// storeDataLocal stores the data we observe in real time to a file and st ores the hashchain header
func storeDataLocal() {
	log.Println("storing a local copy of data")
	path := consts.TellerHomeDir + "/data.txt"

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport}

	body := "https://api.particle.io/v1/devices/events?access_token=3f7d69aa99956fd77c5466f3f52eb6132f500210"
	resp, err := client.Get(body)
	if err != nil {
		log.Println("error while reading from streaming endpoint: ", err)
		return
	}

	defer func() {
		if ferr := resp.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	reader := bufio.NewReader(resp.Body)
	x := make([]byte, 200)
	// open and write to file
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(path)
		if err != nil {
			log.Println("error while opening file", err)
			return
		}
	} else {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			// don't start the teller if we can't read the last known hash since this would break continuity
			log.Println(err)
			return
		}
		HashChainHeader = string(data)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("error while opening file", err)
		return
	}

	log.Println("starting to stream data from particle board: ")
	// this loop waits for inputs (in this case from the particle API) and continually
	// writes it to a data stream
	for {
		_, err = reader.Read(x)
		if err != nil {
			log.Println(err)
			continue
		}
		//log.Println("streaming data from particle board: ", string(x))
		_, err = file.Write(x)
		if err != nil {
			log.Println("error while writing to file", err)
			continue
		}
		size, err := file.Stat()
		if err != nil {
			log.Println(err)
			continue
		}
		// comment since this would fill console out and we can't read anything
		// log.Println("File size is: ", size.Size())
		if size.Size() >= int64(consts.TellerMaxLocalStorageSize) {
			log.Println("flushing data to ipfs")
			// close the file, store in ipfs, get hash, delete file and create same file again
			// with the previous file's hash (so people can verify)
			// we need to store this in ipfs, delete this file and then commit the ipfs hash as
			// the first line in a new file. This whole construction is like a blockchain so we could say
			// we have a blockchain within a blockchain
			// log.Println("size limit reached, taking action")
			err = file.Close()
			if err != nil {
				log.Println("couldn't close file, trying again")
				time.Sleep(2 * time.Second)
				continue
			}
			fileHash, err := ipfs.IpfsAddBytes([]byte(path))
			if err != nil {
				log.Println("Couldn't hash file: ", err)
			}
			HashChainHeader = fileHash
			fileHash = "IPFSHASHCHAIN: " + fileHash + "\n" // the header of the ipfs hashchain that we form
			// log.Println("HashChainHeader: ", HashChainHeader)
			os.Remove(path)
			_, err = os.Create(path)
			if err != nil {
				log.Println("error while opening file", err)
				time.Sleep(2 * time.Second)
				continue
			}
			file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				log.Println("error while opening file", err)
				time.Sleep(2 * time.Second)
				continue
			}
			file.Write([]byte(fileHash))
		}
	}
}

// commitDataShutdown is called when the teller errors out and goes down
func commitDataShutdown() {
	// retrieve data from local storage
	log.Println("printing data before shutdown")
	path := consts.TellerHomeDir + "/data.txt"

	fileHash, err := ipfs.IpfsAddBytes([]byte(path))
	if err != nil {
		log.Println("Couldn't hash file: ", err)
	}

	defer func() {
		if ferr := os.Remove(path); ferr != nil {
			ferr = err
		}
	}()

	_, err = os.Create(path)
	if err != nil {
		log.Println("error while opening file", err)
		return
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("error while opening file", err)
		return
	}

	file.Write([]byte(fileHash))
	err = file.Close()
	if err != nil {
		defer func() {
			if ferr := os.Remove(path); ferr != nil {
				ferr = err
			}
		}()
	}

	err = storeStateHistory(fileHash)
	if err != nil {
		return
	}
}

const tellerUrl = "https://localhost"

type statusResponse struct {
	Code   int
	Status string
}

// GenerateDeviceID generates a random 16 character device ID
func generateDeviceID() (string, error) {
	rs := utils.GetRandomString(16)
	upperCase := strings.ToUpper(rs)
	return upperCase, nil
}

// CheckDeviceID checks the device's ID against a locally saved copy
func checkDeviceID() error {
	// checks whether there is a device id set on this device beforehand
	if _, err := os.Stat(consts.TellerHomeDir); os.IsNotExist(err) {
		// directory does not exist, create a device id
		log.Println("Creating home directory for teller")
		os.MkdirAll(consts.TellerHomeDir, os.ModePerm)
		path := consts.TellerHomeDir + "/deviceid.hex"
		file, err := os.Create(path)
		if err != nil {
			return errors.Wrap(err, "could not create device id file")
		}
		deviceId, err := generateDeviceID()
		if err != nil {
			return errors.Wrap(err, "could not generate device id")
		}
		colorOutput("GENERATED UNIQUE DEVICE ID: "+deviceId, GreenColor)
		_, err = file.Write([]byte(deviceId))
		if err != nil {
			return errors.Wrap(err, "could not write device id to file")
		}
		file.Close()
		err = setDeviceId(LocalRecipient.U.Username, deviceId)
		if err != nil {
			return errors.Wrap(err, "could not store device id in remote platform")
		}
	}
	return nil
}

// GetDeviceID retrieves the deviceId from storage
func getDeviceID() (string, error) {
	path := consts.TellerHomeDir + "/deviceid.hex"
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "could not open teller home path")
	}

	defer func() {
		if ferr := file.Close(); ferr != nil {
			err = ferr
		}
	}()
	// read the hex string from the file
	data := make([]byte, 32)
	readBytes, err := file.Read(data)
	if err != nil {
		return "", errors.Wrap(err, "could not read from file")
	}
	if readBytes != 32 {
		return "", errors.New("Length of strings doesn't match, quitting!")
	}
	return string(data), nil
}
