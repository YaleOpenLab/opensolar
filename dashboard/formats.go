package main

import "time"

type LinkFormat struct {
	Link string
	Text string
}

type PingFormat struct {
	Link string
	Text string
	URL  string
}

type PersonFormat struct {
	Name     string
	Username string
	Email    string
}

type AdminFormat struct {
	Username   string
	Password   string
	AdminToken string
	RecpToken  string
}

type Content struct {
	Title            string
	Name             string
	OpensStatus      PingFormat
	OpenxStatus      PingFormat
	BuildsStatus     PingFormat
	WebStatus        PingFormat
	Validate         LinkFormat
	NextInterval     LinkFormat
	TellerEnergy     LinkFormat
	DateLastPaid     LinkFormat
	DateLastStart    LinkFormat
	DeviceID         LinkFormat
	DABalance        LinkFormat
	PBBalance        LinkFormat
	AccountBalance1  LinkFormat
	AccountBalance2  LinkFormat
	EscrowBalance    LinkFormat
	Recipient        PersonFormat
	Investor         PersonFormat
	Developer        PersonFormat
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
	ProjCount        LinkFormat
	UserCount        LinkFormat
	InvCount         LinkFormat
	RecpCount        LinkFormat
	Admin            AdminFormat
	Date             string
}

type length struct {
	Length int
}
