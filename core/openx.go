package core

import (
	"encoding/json"
	"log"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"
	"github.com/pkg/errors"

	openx "github.com/YaleOpenLab/openx/database"
)

// RetrieveUser retrieves a user from openx's database
func RetrieveUser(key int) (openx.User, error) {
	var user openx.User
	keyString, err := utils.ToString(key)
	if err != nil {
		return user, err
	}
	body := consts.OpenxURL + "/platform/user/retrieve?code=" +
		consts.TopSecretCode + "&key=" + keyString

	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}

	if user.Index == 0 {
		log.Println(string(data))
		return user, errors.New("problem with retrieving user")
	}

	return user, nil
}

// ValidateUser validates a user with openx's database
func ValidateUser(name string, token string) (openx.User, error) {
	var user openx.User
	body := consts.OpenxURL + "/platform/user/validate?code=" +
		consts.TopSecretCode + "&username=" + name + "&token=" + token

	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		log.Println(err)
		return user, err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		log.Println(err)
		return user, err
	}

	if user.Index == 0 {
		log.Println(string(data))
		return user, errors.New("problem with user validation")
	}

	return user, nil
}

// NewUser creates a new user on openx
func NewUser(name string, pwhash string, seedpwd string, email string) (openx.User, error) {
	var user openx.User
	body := consts.OpenxURL + "/platform/user/new?code=" + consts.TopSecretCode + "&username=" + name + "&pwhash=" + pwhash +
		"&seedpwd=" + seedpwd + "&email=" + email

	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		return user, err
	}

	if user.Index == 0 {
		log.Println(string(data))
		return user, errors.New("problem with user creation")
	}

	return user, nil
}

// CheckUsernameCollision checks for username collisions while creating a new user
func CheckUsernameCollision(name string) bool {
	body := consts.OpenxURL + "/platform/user/collision?code=" +
		consts.TopSecretCode + "&username=" + name

	log.Println(body)
	data, err := erpc.GetRequest(body)
	if err != nil {
		return true
	}

	return data[0] == byte(1) // 0 means no collision, 1 means a collison was found
}
