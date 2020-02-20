package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	consts "github.com/YaleOpenLab/opensolar/consts"
)

// storeParticleDataLocal stores the data we observe in real time to a file
func storeParticleDataLocal() {
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

		// log.Println("File size is: ", size.Size())
		if size.Size() >= int64(consts.TellerMaxLocalStorageSize) {
			log.Println("flushing data to ipfs")
			// close the file, store in ipfs, get hash, delete file and create same file again
			// with the previous file's hash (so people can verify) as the first line
			err = file.Close()
			if err != nil {
				log.Println("couldn't close file, trying again")
				time.Sleep(2 * time.Second)
				continue
			}
			fileHash, err := storeDataInIpfs(path)
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
