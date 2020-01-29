package main

import (
	"fmt"
	"log"

	utils "github.com/Varunram/essentials/utils"
	"github.com/fatih/color"
)

var (
	// WhiteColor is a pretty handler for the default colors defined
	WhiteColor = color.FgHiWhite
	// GreenColor is a pretty handler for the default colors defined
	GreenColor = color.FgHiGreen
	// RedColor is a pretty handler for the default colors defined
	RedColor = color.FgHiRed

	// CyanColor is a pretty handler for the default colors defined
	// CyanColor = color.FgHiCyan
	// HiYellowColor is a pretty handler for the default colors defined
	// HiYellowColor = color.FgHiYellow

	// YellowColor is a pretty handler for the default colors defined
	YellowColor = color.FgYellow
	// MagentaColor is a pretty handler for the default colors defined
	MagentaColor = color.FgMagenta

	// HiWhiteColor is a pretty handler for the default colors defined
	// HiWhiteColor = color.FgHiWhite
	// FaintColor is a pretty handler for the default colors defined
	// FaintColor = color.Faint
)

// colorOutput prints the string in the passed color
func colorOutput(msg string, gColor color.Attribute) {
	x := color.New(gColor)
	x.Fprintf(color.Output, "%s\n", msg)
}

// ParseInput parses user input
func ParseInput(input []string) {
	if len(input) == 0 {
		fmt.Println("List of commands: ping, receive, display, update")
		return
	}

	command := input[0]
	switch command {
	case "qq":
		// handler to quit and test the teller without hashing the state and committing two transactions
		// each time we start the teller
		log.Fatal("qq emergency exit")
	case "help":
		fmt.Println("List of commands: ping, receive, display, info, update")
	case "ping":
		err := ping()
		if err != nil {
			log.Println(err)
		}
	case "receive":
		if len(input) != 2 {
			fmt.Println("USAGE: receive xlm")
			return
		}
		err := askXLM() // the rpc allows people to only ask for coins to their publickey, so we should be okay here
		if err != nil {
			log.Println(err)
		}
	case "display":
		if len(input) < 2 {
			fmt.Println("USAGE: display <balance, info>")
			return
		}
		subcommand := input[1]
		switch subcommand {
		case "balance":
			if len(input) < 3 {
				fmt.Println("USAGE: display balance <xlm, asset>")
				return
			}

			subsubcommand := input[2]
			var balance float64
			var err error
			colorOutput("Displaying balance in "+subsubcommand+" for user: ", WhiteColor)

			switch subsubcommand {
			case "xlm":
				balance, err = getNativeBalance()
			default:
				balance, err = getAssetBalance(subsubcommand)
			}

			if err != nil {
				log.Println(err)
				return
			}

			balanceS, err := utils.ToString(balance)
			if err != nil {
				log.Println(err)
				return
			}
			colorOutput(balanceS, MagentaColor)
		case "info":
			fmt.Println("          PROJECT INDEX: ", LocalProject.Index)
			fmt.Println("          Panel Size: ", LocalProject.PanelSize)
			fmt.Println("          Total Value: ", LocalProject.TotalValue)
			fmt.Println("          Location: ", LocalProject.State)
			fmt.Println("          Money Raised: ", LocalProject.MoneyRaised)
			fmt.Println("          Metadata: ", LocalProject.Metadata)
			fmt.Println("          Years: ", LocalProject.EstimatedAcquisition)
			fmt.Println("          Auction Type: ", LocalProject.AuctionType)
			fmt.Println("          Debt Asset Code: ", LocalProject.DebtAssetCode)
			fmt.Println("          Payback Asset Code: ", LocalProject.PaybackAssetCode)
			fmt.Println("          Balance Left: ", LocalProject.BalLeft)
			fmt.Println("          Date Initiated: ", LocalProject.DateInitiated)
			fmt.Println("          Date Last Paid: ", LocalProject.DateLastPaid)
		default:
			// handle defaults here
			log.Println("Invalid command or need more parameters")
		} // end of display
	case "update":
		if len(input) != 1 {
			fmt.Println("USAGE: update")
			return
		}
		updateState(true)
	case "hh":
		// hh = hashchain header
		if len(input) != 1 {
			fmt.Println("USAGE: hh")
			return
		}
		log.Println("HASHCHAIN HEADER: ", HashChainHeader)
	}
}
