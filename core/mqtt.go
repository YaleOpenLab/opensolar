package core

import (
	"fmt"
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	consts "github.com/YaleOpenLab/opensolar/consts"
	mqttlib "github.com/YaleOpenLab/opensolar/teller/mqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sparrc/go-ping"
)

func TrackProject(projIndex int, errs chan error) {
	go func() {
		err := track(projIndex)
		if err != nil {
			errs <- fmt.Errorf("track project error: %s", err.Error())
		}
		close(errs)
	}()
}

func track(projIndex int) error {
	project, err := RetrieveProject(projIndex)
	if err != nil {
		return errors.Wrap(err, "could not retrieve project")
	}

	pinger, err := ping.NewPinger(project.BrokerUrl)
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

	log.Println("tracking project by starting an MQTT client: ", projIndex)
	// start a subscriber connected to a specific topic
	projectString, err := utils.ToString(projIndex)
	if err != nil {
		return errors.Wrap(err, "could not convert to string")
	}

	mqttopts := mqtt.NewClientOptions()
	mqttopts.AddBroker(project.BrokerUrl)
	mqttopts.SetClientID("platformID" + projectString)
	mqttopts.SetUsername("platform" + projectString)

	err = mqttlib.SubscribeMessage(mqttopts, project.TellerPublishTopic, consts.TellerQos, consts.TellerListenNum)
	if err != nil {
		return errors.Wrap(err, "could not subscribe to topic / broker")
	}

	return nil
}
