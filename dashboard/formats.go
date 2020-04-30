package main

import "time"

type linkFormat struct {
	Link string
	Text string
}

type pingFormat struct {
	Link string
	Text string
	URL  string
}

type personFormat struct {
	Name     string
	Username string
	Email    string
}

type adminFormat struct {
	Username   string
	Password   string
	AdminToken string
	RecpToken  string
}

type content struct {
	Title            string
	Name             string
	OpensStatus      pingFormat
	OpenxStatus      pingFormat
	BuildsStatus     pingFormat
	WebStatus        pingFormat
	Validate         linkFormat
	NextInterval     linkFormat
	TellerEnergy     linkFormat
	DateLastPaid     linkFormat
	DateLastStart    linkFormat
	DeviceID         linkFormat
	DABalance        linkFormat
	PBBalance        linkFormat
	AccountBalance1  linkFormat
	AccountBalance2  linkFormat
	EscrowBalance    linkFormat
	Recipient        personFormat
	Investor         personFormat
	Developer        personFormat
	PastEnergyValues []uint32
	DeviceLocation   string
	StateHashes      []string
	PaybackPeriod    time.Duration
	BalanceLeft      float64
	OwnershipShift   float64
	DateInitiated    string
	Stage            int
	DateFunded       string
	InvAssetCode     string
	ProjCount        linkFormat
	UserCount        linkFormat
	InvCount         linkFormat
	RecpCount        linkFormat
	Admin            adminFormat
	Date             string
}

type length struct {
	Length int
}
