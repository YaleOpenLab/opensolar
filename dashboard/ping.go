package main

import (
	"encoding/json"
	"log"

	erpc "github.com/Varunram/essentials/rpc"
)

func opensPing() bool {
	data, err := erpc.GetRequest(platformURL + "/ping")
	if err != nil {
		log.Println(err)
		return false
	}

	var x erpc.StatusResponse

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return false
	}

	return x.Code == 200
}

func openxPing() bool {
	data, err := erpc.GetRequest("https://api.openx.solar/ping")
	if err != nil {
		log.Println(err)
		return false
	}

	var x erpc.StatusResponse

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return false
	}

	return x.Code == 200
}

func buildsPing() bool {
	data, err := erpc.GetRequest("https://builds.openx.solar/ping")
	if err != nil {
		log.Println(err)
		return false
	}

	var x erpc.StatusResponse

	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println(err)
		return false
	}

	return x.Code == 200
}

func websitePing() bool {
	data, err := erpc.GetRequest("https://openx.solar")
	if err != nil {
		log.Println(err)
		return false
	}

	return string(data)[2:14] == "doctype html"
}
