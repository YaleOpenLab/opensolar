package main

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
)

func getToken(username, password string) (string, error) {
	form := url.Values{}
	form.Add("username", username)
	form.Add("pwhash", utils.SHA3hash(password))

	retdata, err := erpc.PostForm(platformURL+"/token", form)
	if err != nil {
		log.Println(err)
		return "", err
	}

	type tokenResponse struct {
		Token string `json:"Token"`
	}

	var x tokenResponse

	err = json.Unmarshal(retdata, &x)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return x.Token, nil
}

func getUToken(wg *sync.WaitGroup, username string) {
	defer wg.Done()
	var err error

	Token, err = getToken(username, "password")
	if err != nil {
		log.Fatal("error while fetching recipient token: ", err)
	}

	var wg1 sync.WaitGroup
	wg1.Add(1)
	go validateRecp(&wg1, username, Token)
	wg1.Wait()
}

func getAToken(wg *sync.WaitGroup) {
	defer wg.Done()
	var err error

	AdminToken, err = getToken("admin", "password")
	if err != nil {
		log.Fatal("error while fetching recipient token: ", err)
	}

	var wg1 sync.WaitGroup
	wg1.Add(1)
	go adminTokenHandler(&wg1, 1)
	wg1.Wait()
}
