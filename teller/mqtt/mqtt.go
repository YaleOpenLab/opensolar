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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	flags "github.com/jessevdk/go-flags"
)

func publishMessage(mqttopts *mqtt.ClientOptions) error {
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

func subscribeMessage(mqttopts *mqtt.ClientOptions) error {
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

	token := client.Subscribe(opts.Topic, byte(opts.Qos), nil)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for receiveCount < opts.Num {
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
	Id        string `long:"id" description:"the clientid" default:"id"`
	Cleansess bool   `long:"cleansess" description:"set clean seession"`
	Qos       int    `long:"qos" description:"quality of service"`
	Num       int    `long:"num" description:"number of messages to subscribe to" default:"1"`
	Message   string `long:"message" description:"message text to publish"`
	Action    string `long:"action" description:"pub/sub" required:"true"`
	Store     string `long:"store" description:"store directory" default:":memory:"`
	Https     bool   `long:"secure" description:"start https"`
}

// ./mqtt --broker 0.0.0.0:1883 --action pub --message cool --topic topic  --id pub --user publisher
// ./mqtt --action sub --topic topic --broker 0.0.0.0:1883 --id sub --user subscriber

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
	log.Printf("\tclientid:  %s\n", opts.Id)
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
	mqttopts.SetClientID(opts.Id)
	mqttopts.SetUsername(opts.User)
	mqttopts.SetPassword(opts.Password)
	mqttopts.SetCleanSession(opts.Cleansess)
	if opts.Store != ":memory:" {
		mqttopts.SetStore(mqtt.NewFileStore(opts.Store))
	}

	// openssl req -new -newkey rsa:2048 -nodes -keyout server.key -out server.csr
	// openssl x509 -req -sha256 -days 365 -in server.csr -signkey server.key -out server.crt

	if opts.Https {
		log.Println("creating root cert chain to work with https broker")
		certFile := "server.crt"

		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		certs, err := ioutil.ReadFile(certFile)
		if err != nil {
			log.Fatalf("Failed to append", err)
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Println("No certs appended, using system certs only")
		}

		config := &tls.Config{
			RootCAs: rootCAs,
		}

		mqttopts.SetTLSConfig(config)
	}

	if opts.Action == "pub" {
		err = publishMessage(mqttopts)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = subscribeMessage(mqttopts)
		if err != nil {
			log.Fatal(err)
		}
	}
}
