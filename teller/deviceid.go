package main

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"

	consts "github.com/YaleOpenLab/openx/consts"
)

// deviceid sets the deviceid and stores it in a retrievable location

// GenerateRandomString generates a random string of length _n_
func GenerateRandomString(n int) (string, error) {
	// generate a crypto secure random string
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "could not read random bytes to generate random string")
	}

	return hex.EncodeToString(b), nil
}

// GenerateDeviceID generates a random 16 character device ID
func GenerateDeviceID() (string, error) {
	rs, err := GenerateRandomString(16)
	if err != nil {
		return "", errors.Wrap(err, "could not generate random string")
	}
	upperCase := strings.ToUpper(rs)
	return upperCase, nil
}

// CheckDeviceID checks the device's ID against a locally saved copy
func CheckDeviceID() error {
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
		deviceId, err := GenerateDeviceID()
		if err != nil {
			return errors.Wrap(err, "could not generate device id")
		}
		ColorOutput("GENERATED UNIQUE DEVICE ID: "+deviceId, GreenColor)
		_, err = file.Write([]byte(deviceId))
		if err != nil {
			return errors.Wrap(err, "could not write device id to file")
		}
		file.Close()
		err = SetDeviceId(LocalRecipient.U.Username, LocalRecipient.U.Pwhash, deviceId)
		if err != nil {
			return errors.Wrap(err, "could not store device id in remote platform")
		}
	}
	return nil
}

// GetDeviceID retrieves the deviceId from storage
func GetDeviceID() (string, error) {
	path := consts.TellerHomeDir + "/deviceid.hex"
	file, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "could not open teller home path")
	}
	// read the hex string from the file
	data := make([]byte, 32)
	numInt, err := file.Read(data)
	if err != nil {
		return "", errors.Wrap(err, "could not read from file")
	}
	if numInt != 32 {
		return "", errors.New("Length of strings doesn't match, quitting!")
	}
	return string(data), nil
}
