package main

import (
	"github.com/RosettaFlow/Carrier-Go/cmd/common"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/node"
	"os"
	runtimeDebug "runtime/debug"

	"github.com/urfave/cli/v2"
)

var (
	appFlags = []cli.Flag{
		flags.RPCHost,
		flags.RPCPort,
		flags.GRPCGatewayHost,
		flags.GRPCGatewayPort,
		// todo: more flags could be define here.
	}
)

func init() {
	appFlags = WrapFlags(appFlags)
}

func main() {
	app := cli.App{}
	app.Name = "carrier"
	app.Usage = "this is a carrier network implementation for RosettaNet"
	// set action func.
	app.Action = startNode
	app.Version = common.Version()
	app.Commands = []*cli.Command {
		// todo: some command cloud define here
	}

	app.Flags = appFlags

	app.Before = func(ctx *cli.Context) error {
		// todo:
		return nil
	}

	defer func() {
		if x := recover(); x != nil {
			log.Errorf("Runtime panic: %v\n%v", x, string(runtimeDebug.Stack()))
			panic(x)
		}
	}()

	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
	}
}

func startNode(ctx *cli.Context) error {
	// todo: some logic could be added here
	node, err := node.New(ctx)
	if err != nil {
		return err
	}
	node.Start()
	return nil
}
