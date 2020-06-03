package main

import (
	"bufio"
	"net/url"

	//"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"

	//	rpc "github.com/YaleOpenLab/openx/rpc"
	erpc "github.com/Varunram/essentials/rpc"
	consts "github.com/YaleOpenLab/opensolar/consts"
	oracle "github.com/YaleOpenLab/opensolar/oracle"
)

// refreshLogin runs once every 5 minutes in order to fetch the latest recipient details
func refreshLogin(username string, pwhash string) error {
	var err error
	for {
		time.Sleep(consts.LoginRefreshInterval)
		err = login(username, pwhash)
		if err != nil {
			colorOutput(CyanColor, err)
		}
		LocalProject, err = getLocalProjectDetails(loginProjIndex)
		if err != nil {
			colorOutput(CyanColor, err)
		}
	}
}

// EndHandler runs when the teller shuts down. Records the start time and location of the
// device in ipfs and commits it as two transactions to the Stellar blockchain
func endHandler() error {
	colorOutput(CyanColor, "Gracefully shutting down, please do not press any button in the process")
	var err error

	NowHash, err = getLatestBlockHash()
	if err != nil {
		colorOutput(RedColor, err)
	}

	hashString := "Device Shutting down. Info: " + DeviceInfo + " Device Location: " + DeviceLocation +
		" Device Unique ID: " + DeviceID + " " + "Start hash: " + StartHash + " Now hash: " + NowHash +
		"Ipfs HashChainHeader: " + HashChainHeader
	// note that we don't commit the latest hash chain header's hash here because this gives us a tighter timeline
	// to audit what really happened
	ipfsHash, err := storeDataInIpfs(hashString)
	if err != nil {
		colorOutput(RedColor, err)
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

	colorOutput(CyanColor, "sent device shutdown notice")
	commitDataShutdown()
	// save last known state of the system in the recipient's list of known hashes. Call this
	// last since there would still be data that we want to measure when the above commands are running
	return nil
	// have a return because we don't want to sigint while we send emails and stuff
}

func splitAndSend2Tx(memo string) (string, string, error) {
	// 10 padding chars + 46 (ipfs hash length) characters
	firstHalf := memo[:28]
	secondHalf := memo[28:]
	tx1, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, 1, firstHalf)
	if err != nil {
		tx1, err = sendXLM(LocalRecipient.U.StellarWallet.PublicKey, 1, firstHalf)
		if err != nil {
			return "", "", err
		}
	}

	time.Sleep(5 * time.Second)

	tx2, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, 1, secondHalf)
	if err != nil {
		tx2, err = sendXLM(LocalRecipient.U.StellarWallet.PublicKey, 1, secondHalf)
		if err != nil {
			return "", "", err
		}
	}

	log.Printf("tx hash: %s, tx2 hash: %s", tx1, tx2)
	return tx1, tx2, nil
}

func checkPayback() {
	for {
		colorOutput(CyanColor, "Payback interval reached. Paying back automatically")
		assetName := LocalProject.DebtAssetCode
		amount := float64(EnergyValue)*oracle.MonthlyBill()/1000000 + 1
		refreshLogin(loginUsername, loginPwhash)
		err := projectPayback(assetName, amount)
		if err != nil {
			colorOutput(RedColor, "Error while paying back", err, "trying again")
			time.Sleep(5 * time.Second)
			refreshLogin(loginUsername, loginPwhash)
			err = projectPayback(assetName, amount)
			if err != nil {
				sendDevicePaybackFailedEmail()
			}
		}
		time.Sleep(time.Duration(LocalProject.PaybackPeriod) * consts.OneWeekInSecond)
	}
}

// updateState stores the current state of the teller in ipfs and commits the ipfs hash to the blockchain
func updateState(trigger bool) {
	for {
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			colorOutput(RedColor, "error while trying to read data file")
			time.Sleep(consts.TellerPollInterval)
		}
		subcommand := string(data)
		refreshLogin(loginUsername, loginPwhash)
		ipfsHash, err := storeDataInIpfs("Device ID: " + DeviceID + " UPDATESTATE" + subcommand)
		if err != nil {
			colorOutput(RedColor, "Error while fetching ipfs hash", err)
			// time.Sleep(consts.TellerPollInterval)
		}

		ipfsHash = "STATUPD: " + ipfsHash
		// Stellar allows one to send as many stroops as desired to the same account, so
		// send timestamp stroops to ourselves.

		// memo field restricted to 28 bytes - AAAAAAAAAAAAAAAAAAAAAAAAAAAA

		// don't use platform RPCs for interacting with the blockchain

		hash1, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, float64(utils.Unix()), ipfsHash[:28])
		if err != nil {
			colorOutput(RedColor, err)
		}

		log.Println(hash1)
		time.Sleep(5 * time.Second)

		hash2, err := sendXLM(LocalRecipient.U.StellarWallet.PublicKey, float64(utils.Unix()), ipfsHash[29:])
		if err != nil {
			colorOutput(RedColor, err)
		}

		log.Println(hash2)
		// we updated state as hash1 and hash2
		colorOutput(MagentaColor, "Updated State: "+hash1+" "+hash2)
		if trigger {
			break // we trigerred this manually, don't want to keep doing this
		}

		time.Sleep(consts.TellerPollInterval)
	}
}

func storeDataInIpfs(data string) (string, error) {
	form := url.Values{}
	form.Add("username", LocalRecipient.U.Username)
	form.Add("token", Token)
	form.Add("data", data)

	var retdata []byte
	var err error
	// retdata, err := erpc.PostForm(APIURL+"/ipfs/putdata", form)
	if strings.Contains(APIURL, "localhost") {
		// connect to openx which runs on http instead of opensolar
		retdata, err = erpc.PostForm("http://localhost:8080/ipfs/putdata", form)
	} else {
		retdata, err = erpc.PostForm(APIURL+"/ipfs/putdata", form)
	}

	if err != nil {
		colorOutput(RedColor, err)
		return "", err
	}

	retdata = retdata[1:47]
	colorOutput(CyanColor, "IPFS HASH: ", string(retdata))

	if len(string(retdata)) != 46 { // 46 is the length of an ipfs hash
		return "", errors.New("ipfs hash storage failed")
	}

	return string(retdata), nil
}

// commitDataShutdown is called when the teller errors out and goes down
func commitDataShutdown() {
	// retrieve data from local storage
	colorOutput(CyanColor, "printing data before shutdown")
	path := consts.TellerHomeDir + "/data.txt"

	fileHash, err := storeDataInIpfs("TELLER SHUTDOWN + " + utils.Timestamp())
	if err != nil {
		colorOutput(RedColor, "Couldn't store data in file: ", err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		colorOutput(RedColor, "error while opening file", err)
		return
	}

	file.Write([]byte(fileHash))
	err = file.Close()
	if err != nil {
		log.Println("error while closing file: ", err)
	}

	err = storeStateHistory(fileHash)
	if err != nil {
		return
	}
}

type statusResponse struct {
	Code   int
	Status string
}

// generateDeviceID generates a random 16 character device ID
func generateDeviceID() (string, error) {
	rs := utils.GetRandomString(consts.TellerDeviceIDLen)
	upperCase := strings.ToUpper(rs)
	return upperCase, nil
}

// checkDeviceID checks the device's ID against a locally saved copy
func checkDeviceID() error {
	// checks whether we've set device id beforehand
	if _, err := os.Stat(consts.TellerHomeDir); os.IsNotExist(err) {
		// directory does not exist, create a device id
		colorOutput(CyanColor, "Creating home directory for teller")
		os.MkdirAll(consts.TellerHomeDir, os.ModePerm)
		path := consts.TellerHomeDir + "/deviceid.hex"
		file, err := os.Create(path)
		if err != nil {
			return errors.Wrap(err, "could not create device id file")
		}
		deviceID, err := generateDeviceID()
		if err != nil {
			return errors.Wrap(err, "could not generate device id")
		}
		colorOutput(GreenColor, "GENERATED UNIQUE DEVICE ID: "+deviceID)
		_, err = file.Write([]byte(deviceID))
		if err != nil {
			return errors.Wrap(err, "could not write device id to file")
		}
		file.Close()
		err = setDeviceID(LocalRecipient.U.Username, deviceID)
		if err != nil {
			return errors.Wrap(err, "could not store device id in remote platform")
		}
	}
	return nil
}

// getDeviceID retrieves the deviceId from storage
func getDeviceID() (string, error) {
	path := consts.TellerHomeDir + "/deviceid.hex"
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "could not open teller home path")
	}

	defer func() {
		if ferr := file.Close(); ferr != nil {
			log.Println(ferr)
			err = ferr
		}
	}()
	// read the hex string from the file
	data := make([]byte, consts.TellerDeviceIDLen)
	readBytes, err := file.Read(data)
	if err != nil {
		return "", errors.Wrap(err, "could not read from file")
	}
	if readBytes != consts.TellerDeviceIDLen {
		return "", errors.New("length of strings doesn't match, quitting")
	}
	return string(data), nil
}

type energyStruct struct {
	EnergyTimestamp string `json:"energy_timestamp"`
	Unit            string `json:"unit"`
	Value           uint32 `json:"value"`
	OwnerID         string `json:"owner_id"`
	AssetID         string `json:"asset_id"`
}

func updateEnergyData() error {
	EnergyValue = 0

	origPath := "data.txt"
	hcPath := consts.TellerHomeDir + "/data.txt"

	presentData, err := ioutil.ReadFile(origPath)
	if err != nil {
		return errors.Wrap(err, "could not open data file for reading")
	}

	hc, err := os.OpenFile(hcPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "could not open data file for reading")
	}

	defer hc.Close()

	_, err = hc.Write(presentData)
	if err != nil {
		return errors.Wrap(err, "could not write data to the hc file")
	}

	// read size of the updated file
	hcData, err := ioutil.ReadFile(hcPath)
	if err != nil {
		return errors.Wrap(err, "could not open data file for reading")
	}

	for {
		size, err := hc.Stat()
		if err != nil {
			colorOutput(RedColor, err)
			break
		}
		colorOutput(CyanColor, "File size is: ", size.Size())
		if size.Size() >= int64(consts.TellerMaxLocalStorageSize) {
			colorOutput(CyanColor, "flushing data to ipfs")
			// close the file, store in ipfs, get hash, delete file and create same file again
			// with the previous file's hash (so people can verify) as the first line
			err = hc.Close()
			if err != nil {
				colorOutput(RedColor, "couldn't close file, trying again")
				break
			}
			fileHash, err := storeDataInIpfs(string(hcData))
			if err != nil {
				colorOutput(RedColor, "Couldn't hash file: ", err)
			}
			HashChainHeader = fileHash
			fileHash = "IPFSHASHCHAIN: " + fileHash + "\n" // the header of the ipfs hashchain that we form
			// colorOutput(CyanColor, "HashChainHeader: ", HashChainHeader)
			os.Remove(hcPath)
			_, err = os.Create(hcPath)
			if err != nil {
				colorOutput(RedColor, "error while opening file", err)
				continue
			}
			hc, err := os.OpenFile(hcPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				colorOutput(RedColor, "error while opening file", err)
				continue
			}
			defer hc.Close()
			hc.Write([]byte(fileHash))
		}
		break
	}

	// now that the hash chain is done, take care of accumulating data
	f, err := os.Open(origPath)
	if err != nil {
		return errors.Wrap(err, "could not open data file for reading")
	}

	defer func() {
		if ferr := f.Close(); ferr != nil {
			err = ferr
		}
	}()

	reader := bufio.NewReader(f)

	for {
		var data []byte

		for i := 0; i < 7; i++ { // formatted according to the responses received from the lumen unit
			// which is further read by the subscriber
			line, _, err := reader.ReadLine()
			if err != nil {
				colorOutput(RedColor, "reached end of file")
				err = os.Remove("data.txt")
				if err != nil {
					return err
				}
				file, err := os.Create("data.txt") // create a new file to reset TellerEnergy
				if err != nil {
					return err
				}
				err = file.Close()
				if err != nil {
					return err
				}
				return nil
			}
			data = append(data, line...)
		}

		var x energyStruct
		err = json.Unmarshal(data, &x)
		if err != nil {
			return errors.Wrap(err, "could not unmarshal json data struct")
		}

		EnergyValue += x.Value
	}
}

// readEnergyData reads energy data from a local file and stores it in the remote opensolar instance
func readEnergyData() {
	for {
		time.Sleep(LocalProject.PaybackPeriod * consts.OneWeekInSecond / 2)
		colorOutput(CyanColor, "reading energy data from file")
		refreshLogin(loginUsername, loginPwhash)
		err := updateEnergyData()
		if err != nil {
			colorOutput(RedColor, "error while reading energy data: ", err)
			err := updateEnergyData()
			if err != nil {
				continue
			}
		}

		// need to update remote with the energy data
		colorOutput(CyanColor, "storing energy data on opensolar")
		refreshLogin(loginUsername, loginPwhash)
		data, err := putEnergy(EnergyValue)
		if err != nil {
			colorOutput(RedColor, err)
			data, err = putEnergy(EnergyValue)
			if err != nil {
				continue
			}
		}
		colorOutput(CyanColor, string(data))
	}
}
