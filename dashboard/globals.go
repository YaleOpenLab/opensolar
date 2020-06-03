package main

import (
	core "github.com/YaleOpenLab/opensolar/core"
)

var platformURL = "https://api2.openx.solar"

// AdminToken is a global used to track the admin token
var AdminToken string

// Token is a global used to track the recipient's token
var Token string

// Pnb is a global used to store the primary native balance across threads
var Pnb string

// Pub is a global used to store the primary usd balance across threads
var Pub string

// Snb is a global used to store the secondary native balance across threads
var Snb string

// Sub is a global used to store the secondary usd balance across threads
var Sub string

// XlmUSD is a global used to track the XLM:USD ticker
var XlmUSD float64

// Project is a global used to track the retrieve and use the project across threads
var Project core.Project

// Recipient is a global used to track the project's Recipient
var Recipient core.Recipient

// Return is a global used to store the dashboard's content
var Return content

// Developer is a global used to track the project's Developer
var Developer core.Entity
