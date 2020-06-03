/*
 * Copyright (c) 2013 IBM Corp.
 *
 * All rights reserved. This program and the accompanying materials
 * are made available under the terms of the Eclipse Public License v1.0
 * which accompanies this distribution, and is available at
 * http://www.eclipse.org/legal/epl-v10.html
 *
 * Contributors:
 *    Seth Hoenig
 *    Allan Stockdill-Mander
 *    Mike Robertson
 */

package main

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	flags "github.com/jessevdk/go-flags"
)

// PublishMessage is a helper used to publish a message to an MQTT broker
func PublishMessage(mqttopts *mqtt.ClientOptions) error {
	client := mqtt.NewClient(mqttopts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil { // token.Wait() returns a bool
		return token.Error()
	}

	for i := 0; i < opts.Num; i++ {
		log.Println("Publishing message: ", opts.Message)
		token := client.Publish(opts.Topic, byte(opts.Qos), false, opts.Message)
		token.Wait()
	}

	client.Disconnect(250)
	log.Println("Publisher Disconnected")
	return nil
}

// SubscribeMessage is a helper used to subscribe to an MQTT broker
func SubscribeMessage(mqttopts *mqtt.ClientOptions, topic string, qos int, num int) error {
	receiveCount := 0
	receiver := make(chan [2]string)
	var messages []string

	mqttopts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		receiver <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(mqttopts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	token := client.Subscribe(topic, byte(qos), nil)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for receiveCount < num {
		incoming := <-receiver
		messages = append(messages, incoming[1])
		log.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}

	client.Disconnect(250)
	log.Println("Subscriber Disconnected")
	log.Println("MESSAGES: ", messages)
	return nil
}

var opts struct {
	Topic     string `long:"topic" description:"The topic name to/from which to publish/subscribe"`
	Broker    string `long:"broker" description:"The broker url" default:"localhost:1883"`
	Password  string `long:"password" description:"The password"`
	User      string `long:"user" description:"the user" default:"username"`
	ID        string `long:"id" description:"the clientid" default:"id"`
	Cleansess bool   `long:"cleansess" description:"set clean seession"`
	Qos       int    `long:"qos" description:"quality of service"`
	Num       int    `long:"num" description:"number of messages to subscribe to" default:"1"`
	Message   string `long:"message" description:"message text to publish"`
	Action    string `long:"action" description:"pub/sub" required:"true"`
	Store     string `long:"store" description:"store directory" default:":memory:"`
}

// ./mqtt --broker 0.0.0.0:1883 --action pub --message cool --topic test --id pub --user publisher
// ./mqtt --action sub --topic test --broker 0.0.0.0:1883 --id sub --user subscriber

// set id to username when connecting

func main() {
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal("Failed to parse arguments / Help command")
	}

	if !(opts.Action == "pub" || opts.Action == "sub") {
		log.Println("Invalid setting for -action, must be pub or sub")
		return
	}

	if opts.Topic == "" {
		log.Println("Invalid setting for -topic, must not be empty")
		return
	}

	log.Printf("\taction:    %s\n", opts.Action)
	log.Printf("\tbroker:    %s\n", opts.Broker)
	log.Printf("\tclientid:  %s\n", opts.ID)
	log.Printf("\tuser:      %s\n", opts.User)
	log.Printf("\tpassword:  %s\n", opts.Password)
	log.Printf("\ttopic:     %s\n", opts.Topic)
	log.Printf("\tmessage:   %s\n", opts.Message)
	log.Printf("\tqos:       %d\n", opts.Qos)
	log.Printf("\tcleansess: %v\n", opts.Cleansess)
	log.Printf("\tnum:       %d\n", opts.Num)
	log.Printf("\tstore:     %s\n", opts.Store)

	mqttopts := mqtt.NewClientOptions()
	mqttopts.AddBroker(opts.Broker)
	mqttopts.SetClientID(opts.ID)
	mqttopts.SetUsername(opts.User)
	mqttopts.SetPassword(opts.Password)
	mqttopts.SetCleanSession(opts.Cleansess)
	if opts.Store != ":memory:" {
		mqttopts.SetStore(mqtt.NewFileStore(opts.Store))
	}

	// openssl req -new -newkey rsa:2048 -nodes -keyout server.key -out server.csr
	// openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt

	if opts.Action == "pub" {
		err = PublishMessage(mqttopts)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = SubscribeMessage(mqttopts, opts.Topic, opts.Qos, opts.Num)
		if err != nil {
			log.Fatal(err)
		}
	}
}
