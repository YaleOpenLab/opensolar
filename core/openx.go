package core

import (
	"encoding/json"
	"log"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"

	openx "github.com/YaleOpenLab/openx/database"
)

// this file handles everything related to openx interaction. Privileged access due to access code

func RetrieveUser(key int) (openx.User, error) {
	var user openx.User
	keyString, err := utils.ToString(key)
	if err != nil {
		return user, err
	}
	body := consts.OpenxURL + "/platform/user/retrieve?code=" + consts.TopSecretCode + "&key=" + keyString
	data, err := erpc.GetRequest(body)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func ValidateUser(name string, pwhash string) (openx.User, error) {
	var user openx.User
	body := consts.OpenxURL + "/platform/user/validate?code=" + consts.TopSecretCode + "&name=" + name + "&pwhash=" + pwhash
	data, err := erpc.GetRequest(body)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func NewUser(name string, pwhash string, seedpwd string, realname string) (openx.User, error) {
	var user openx.User
	body := consts.OpenxURL + "/platform/user/new?code=" + consts.TopSecretCode + "&name=" + name + "&pwhash=" + pwhash +
		"&seedpwd=" + seedpwd + "&realname=" + realname

	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func CheckUsernameCollision(name string) bool {
	body := consts.OpenxURL + "/platform/user/collision?code=" + consts.TopSecretCode + "&name=" + name
	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		return true
	}

	log.Println("DATA: ", data)
	return data[0] == byte(1) // 0 means no collision, 1 means a collison was found
}
