package core

import (
	"fmt"
	"log"

	mqttlib "github.com/YaleOpenLab/opensolar/teller/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	// "github.com/sparrc/go-ping"
)

func TrackProject(projIndex int, brokerurl string, topic string, errs chan error){
	go func() {
		err := trackProject(projIndex, brokerurl, topic)
		if err != nil {
			errs <- fmt.Errorf("track project error: %s", err.Error())
		}
		close(errs)
	}()
}

func trackProject(projIndex int, brokerurl string, topic string) error {
	/*
	pinger, err := ping.NewPinger(brokerurl)
	if err != nil {
		return err
	}

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
	*/
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return err
	}
	log.Println(project.Index)
	// start a subscriber connected to a specific topic

	mqttopts := mqtt.NewClientOptions()
	mqttopts.AddBroker(brokerurl)
	mqttopts.SetClientID("platformID")
	mqttopts.SetUsername("platform")
	qos := 0
	num := 1

	err = mqttlib.SubscribeMessage(mqttopts, topic, qos, num)
	if err != nil {
		return err
	}

	return nil
}
