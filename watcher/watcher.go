package main

import (
	"encoding/json"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

// ParticlePingResponse is a structure to parse returned particle.io data
type ParticlePingResponse struct {
	Online bool `json:"online"`
	Ok     bool `json:"ok"`
}

func main() {
	var err error
	// read from config file. Doesn't check for config file. If you don;t have one, create
	// one and then start this notifier daemon
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	email1 := viper.Get("pa1").(string)              // pa = platform admin
	email2 := viper.Get("pa2").(string)              // platform admin 2
	accessToken := viper.Get("accessToken").(string) // the access token to access the particle io interface
	deviceID := viper.Get("deviceId").(string)       // the device id associated with the IoT hub

	body := "https://api.particle.io/v1/devices/" + deviceID + "/ping"

	for {
		payload := strings.NewReader("access_token=" + accessToken)
		data, err := erpc.PutRequest(body, payload)
		if err != nil {
			log.Println("did not receive success response", err)
			return
		}
		var x ParticlePingResponse
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println("did not unmarshal json", err)
			return
		}
		if !x.Ok || !x.Online {
			// the platform is not online, so we need to send an email to the platform admins alerting them of the same
			// read config from the config file
			// read from config.yaml in the working directory
			log.Println("SENDING ALERT EMAIL TO: ", email1, "AND:", email2)

			err = SendIoTHubDownEmail("S.U.Pasto School, Puerto Rico", email1, email2)
			if err != nil {
				log.Println("Failed to send notification, quitting!")
				return
			}
		}
		time.Sleep(2 * 24 * time.Hour) // check every hour whether the IoT Hub is up or not
	}
}

// SendIoTHubDownEmail is an email to the platform notifying that the IoT device for a particular project is down.
func SendIoTHubDownEmail(location string, email1 string, email2 string) error {
	body := "Greetings from your remote notifier! \n\nWe're writing to let you know that your remote IoT Hub in: " + location +
		" has not been responding to pings for a while. The timestamp of this alert is: " + utils.Timestamp() + " Please take action at the earliest." + "\n\n\n" +
		"Have a nice day! \nYour Friendly Notifier"

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	from := viper.Get("email").(string)    // interface to string
	pass := viper.Get("password").(string) // interface to string
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	// to can also be an array of addresses if needed
	msg := "From: " + from + "\n" +
		"To: " + email1 + "\n" +
		"Subject: OpenSolar IoT Hub DOWN:\n\n" + body

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email1}, []byte(msg))
	if err != nil {
		return errors.Wrap(err, "smtp error")
	}

	msg = "From: " + from + "\n" +
		"To: " + email1 + "\n" +
		"Subject: OpenSolar IoT Hub DOWN:\n\n" + body

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{email2}, []byte(msg))
	if err != nil {
		return errors.Wrap(err, "smtp error")
	}

	return nil
}
