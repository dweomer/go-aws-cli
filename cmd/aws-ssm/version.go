package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/urfave/cli"
)

var (
	// Version is the urfave/cli App.Version
	Version = "unknown"
)

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%v version %v (%v=%v)\n", c.App.HelpName, c.App.Version, aws.SDKName, aws.SDKVersion)
	}

	app.Commands = append(app.Commands, cli.Command{
		Name:   "version",
		Usage:  "Print the version",
		Action: cli.ShowVersion,
	})
}
