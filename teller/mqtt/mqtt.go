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
	"crypto/x509"
	"io/ioutil"
	"crypto/tls"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	flags "github.com/jessevdk/go-flags"
)

func sub(mqttopts *mqtt.ClientOptions) {
	client := mqtt.NewClient(mqttopts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Sample Publisher Started")
	for i := 0; i < opts.Num; i++ {
		log.Println("---- doing publish ----")
		token := client.Publish(opts.Topic, byte(opts.Qos), false, opts.Payload)
		token.Wait()
	}

	client.Disconnect(250)
	log.Println("Sample Publisher Disconnected")
}

var opts struct {
	Topic     string `long:"topic" description:"The topic name to/from which to publish/subscribe"`
	Broker    string `long:"broker" description:"The broker URI" default:"tls://localhost:8883"`
	Password  string `long:"password" description:"The password"`
	User      string `long:"user" description:"the user" default:"username"`
	Id        string `long:"id" description:"the clientid" default:"id"`
	Cleansess bool   `long:"cleansess" description:"set clean seession"`
	Qos       int    `long:"qos" description:"quality of service"`
	Num       int    `long:"num" description:"number of messages to subscribe to" default:"1"`
	Payload   string `long:"message" description:"message text to publish"`
	Action    string `long:"action" description:"pub/sub" required:"true"`
	Store     string `long:"store" description:"store directory" default:":memory:"`
}

// ./mqtt --broker 0.0.0.0:1883 --action pub --message cool --topic topic --user user --id 1
// ./mqtt --action sub --topic topic --broker 0.0.0.0:1883 --user user --id blah

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

	log.Printf("Sample Info:\n")
	log.Printf("\taction:    %s\n", opts.Action)
	log.Printf("\tbroker:    %s\n", opts.Broker)
	log.Printf("\tclientid:  %s\n", opts.Id)
	log.Printf("\tuser:      %s\n", opts.User)
	log.Printf("\tpassword:  %s\n", opts.Password)
	log.Printf("\ttopic:     %s\n", opts.Topic)
	log.Printf("\tmessage:   %s\n", opts.Payload)
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
	if opts.Action == "pub" {
		sub(mqttopts)
	} else {
		receiveCount := 0
		choke := make(chan [2]string)

		mqttopts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		})

		client := mqtt.NewClient(mqttopts)
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		token = client.Subscribe(opts.Topic, byte(opts.Qos), nil)
		if token.Wait() && token.Error() != nil {
			log.Println(token.Error())
			os.Exit(1)
		}

		for receiveCount < opts.Num {
			incoming := <-choke
			log.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			receiveCount++
		}

		client.Disconnect(250)
		log.Println("Sample Subscriber Disconnected")
	}
}
