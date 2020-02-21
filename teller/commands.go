package main

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

var (
	// color scheme follows
	// cyan - non error logs
	// red - error logs
	// yellow - banners
	// green - successful calls
	// magenta - blockchain / hash messages
	// white - others

	// WhiteColor is a pretty handler for the default colors defined
	// WhiteColor = color.FgHiWhite

	// GreenColor is a pretty handler for the default colors defined
	GreenColor = color.FgHiGreen
	// RedColor is a pretty handler for the default colors defined
	RedColor = color.FgHiRed

	// CyanColor is a pretty handler for the default colors defined
	CyanColor = color.FgHiCyan
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
func colorOutput(gColor color.Attribute, msg ...interface{}) {
	x := color.New(gColor)
	x.Println(msg...)
}

var commands = []string{"qq", "help", "ping", "receive", "display", "info", "update", "hh"}

// ParseInput parses user input
func ParseInput(input []string) {
	if len(input) == 0 {
		fmt.Println("List of commands: ", commands)
		return
	}

	command := input[0]
	switch command {
	case commands[0]:
		// handler to quit and test the teller without hashing the state and committing two transactions
		// each time we start the teller
		log.Fatal("qq emergency exit")
	case commands[1]:
		fmt.Println("List of commands: ", commands)
	case commands[2]:
		err := ping()
		if err != nil {
			log.Println(err)
		}
	case commands[3]:
		if len(input) != 2 {
			fmt.Println("USAGE: receive xlm")
			return
		}
	case commands[4]:
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
			colorOutput(MagentaColor, "Displaying balance in "+subsubcommand+" for user: ")

			switch subsubcommand {
			case "xlm":
				balance, err = getNativeBalance()
			default:
				balance, err = getAssetBalance(input[3])
			}

			if err != nil {
				log.Println(err)
				return
			}
			colorOutput(MagentaColor, balance)
		default:
			// handle defaults here
			log.Println("Invalid command or need more parameters")
		}
	case commands[5]:
		fmt.Println("          PROJECT INDEX: ", LocalProject.Index)
		fmt.Println("          Money Raised: ", LocalProject.MoneyRaised)
		fmt.Println("          Metadata: ", LocalProject.Metadata)
		fmt.Println("          Debt Asset Code: ", LocalProject.DebtAssetCode)
		fmt.Println("          Payback Asset Code: ", LocalProject.PaybackAssetCode)
		fmt.Println("          Balance Left: ", LocalProject.BalLeft)
		fmt.Println("          Date Initiated: ", 0)
		fmt.Println("          Date Last Paid: ", time.Unix(LocalProject.DateLastPaid, 0))
	// end of display
	case commands[6]:
		if len(input) != 1 {
			fmt.Println("USAGE: update")
			return
		}
		updateState(true)
	case commands[7]:
		// hh = hashchain header
		if len(input) != 1 {
			fmt.Println("USAGE: hh")
			return
		}
		log.Println("HASHCHAIN HEADER: ", HashChainHeader)
	}
}
