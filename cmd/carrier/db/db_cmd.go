package db

import (
	"github.com/RosettaFlow/Carrier-Go/cmd"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var log = logrus.WithField("prefix", "db")

// Commands for interacting with a carrier node database.
var Commands = &cli.Command{
	Name:     "db",
	Category: "db",
	Usage:    "defines commands for interacting with carrier node database",
	Subcommands: []*cli.Command{
		{
			Name:        "restore",
			Description: `restores a database from a backup file`,
			Flags: cmd.WrapFlags([]cli.Flag{
				flags.RestoreSourceFileFlag,
				flags.RestoreTargetDirFlag,
			}),
			Before: func(context *cli.Context) error {
				return nil
			},
			Action: func(cliCtx *cli.Context) error {

				return nil
			},
		},
	},
}
