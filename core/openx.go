package core

import (
	"encoding/json"
	"log"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"

	openx "github.com/YaleOpenLab/openx/database"
)

// this file handles interactions with openx. Privileged access due to access code

// RetrieveUser retrieves a user from openx's database
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

// ValidateUser validates a user with openx's database
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

// NewUser creates a new user in openx's database
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

// CheckUsernameCollision checks for username collisions in openx's database
func CheckUsernameCollision(name string) bool {
	body := consts.OpenxURL + "/platform/user/collision?code=" + consts.TopSecretCode + "&name=" + name
	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		return true
	}

	return data[0] == byte(1) // 0 means no collision, 1 means a collison was found
}
