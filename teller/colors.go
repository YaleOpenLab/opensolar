package main

import (
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

// ColorOutput prints the string in the passed color
func ColorOutput(msg string, gColor color.Attribute) {
	x := color.New(gColor)
	x.Fprintf(color.Output, "%s\n", msg)
}
