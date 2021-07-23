package main

import (
	"github.com/RosettaFlow/Carrier-Go/common/debug"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/urfave/cli/v2"
	"io"
	"sort"
)

var appHelpTemplate = `NAME:
   {{.App.Name}} - {{.App.Usage}}
USAGE:
   {{.App.HelpName}} [options]{{if .App.Commands}} command [command options]{{end}} {{if .App.ArgsUsage}}{{.App.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if .App.Version}}
AUTHOR:
   {{range .App.Authors}}{{ . }}{{end}}
   {{end}}{{if .App.Commands}}
GLOBAL OPTIONS:
   {{range .App.Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .FlagGroups}}
{{range .FlagGroups}}{{.Name}} OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
{{end}}{{end}}{{if .App.Copyright }}
COPYRIGHT:
   {{.App.Copyright}}
VERSION:
   {{.App.Version}}
   {{end}}{{if len .App.Authors}}
   {{end}}
`


type flagGroup struct {
	Name  string
	Flags []cli.Flag
}

var appHelpFlagGroups = []flagGroup{
	{
		Name: "cmd",
		Flags: []cli.Flag{
			flags.RPCHost,
			flags.RPCPort,
			flags.DisableGRPCGateway,
			flags.GPRCGatewayCorsDomain,
			flags.GRPCGatewayHost,
			flags.GRPCGatewayPort,
			flags.CertFlag,
			flags.KeyFlag,
			flags.GrpcMaxCallRecvMsgSizeFlag,
		},
	},
	{
		Name: "debug",
		Flags: []cli.Flag{
			debug.PProfFlag,
			debug.PProfAddrFlag,
			debug.PProfPortFlag,
			debug.MemProfileRateFlag,
			debug.CPUProfileFlag,
			debug.TraceFlag,
			debug.BlockProfileRateFlag,
			debug.MutexProfileFractionFlag,
			debug.DebugFlag,
			flags.EnableDebugRPCEndpoints,
		},
	},
	{
		Name: "carrier",
		Flags: []cli.Flag{
			flags.DeveloperFlag,
			flags.TestnetFlag,
			flags.DataDirFlag,
			flags.ClearDB,
			flags.VerbosityFlag,
			flags.RestoreSourceFileFlag,
			flags.RestoreTargetDirFlag,
			flags.ConfigFileFlag,
		},
	},
	{
		Name: "p2p",
		Flags: []cli.Flag{
			flags.P2PIP,
			flags.P2PHost,
			flags.P2PAllowList,
			flags.P2PHostDNS,
			flags.P2PDenyList,
			flags.P2PMaxPeers,
			flags.P2PMetadata,
			flags.P2PPrivKey,
			flags.P2PTCPPort,
			flags.P2PUDPPort,
			flags.NoDiscovery,
			flags.BootstrapNode,
			flags.EnableUPnPFlag,
			flags.DisableDiscv5,
			flags.StaticPeers,
			flags.RelayNode,
		},
	},
	{
		Name: "log",
		Flags: []cli.Flag{
			flags.LogFormat,
			flags.LogFileName,
		},
	},
}

func init() {
	cli.AppHelpTemplate = appHelpTemplate

	type helpData struct {
		App        interface{}
		FlagGroups []flagGroup
	}

	originalHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		if tmpl == appHelpTemplate {
			for _, group := range appHelpFlagGroups {
				sort.Sort(cli.FlagsByName(group.Flags))
			}
			originalHelpPrinter(w, tmpl, helpData{data, appHelpFlagGroups})
		} else {
			originalHelpPrinter(w, tmpl, data)
		}
	}
}


