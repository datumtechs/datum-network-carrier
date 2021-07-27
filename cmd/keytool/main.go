package main

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/cmd/utils"
	"github.com/urfave/cli/v2"
	"os"
)

const (
	defaultKeyfileName = "keyfile.json"
)

// Git SHA1 commit hash of the release (set via linker flags)
var gitCommit = ""

var app *cli.App

func init() {
	app = utils.NewApp(gitCommit, "an Carrier key manager")
	app.Commands = []*cli.Command{
		commandGenkeypair,
	}
}

// Commonly used command line flags.
var (
	passphraseFlag = &cli.StringFlag{
		Name:  "passwordfile",
		Usage: "the file that contains the passphrase for the keyfile",
	}
	jsonFlag = &cli.BoolFlag{
		Name:  "json",
		Usage: "output JSON instead of human-readable format",
	}
)

func main() {
	cli.CommandHelpTemplate = utils.OriginCommandHelpTemplate
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
