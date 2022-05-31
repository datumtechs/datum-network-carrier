package main

import (
	"github.com/datumtechs/datum-network-carrier/common/debug"
	"github.com/datumtechs/datum-network-carrier/common/flags"
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
			flags.GrpcMaxCallSendMsgSizeFlag,
			flags.SetGCPercent,
			flags.DiscoveryServerIP,
			flags.DiscoveryServerPort,
			flags.DiscoveryServiceId,
			flags.DiscoveryServiceName,
			flags.DiscoveryServerTags,
			flags.DiscoveryServiceHealthCheckDeregister,
			flags.DiscoveryServiceHealthCheckInterval,
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
			flags.EnableTracingFlag,
			flags.TracingProcessNameFlag,
			flags.TracingEndpointFlag,
			flags.TraceSampleFractionFlag,
			flags.MockIdentityIdFile,
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
			flags.ConsensusStateWalDir,
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
			flags.EnableFakeNetwork,
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
	{
		Name: "chain",
		Flags: []cli.Flag{
			flags.BlockChain,
		},
	},
	// add by v0.4.0
	{
		Name: "kms",
		Flags: []cli.Flag{
			flags.KMSKeyId,
			flags.KMSRegionId,
			flags.KMSAccessKeyId,
			flags.KMSAccessKeySecret,
		},
	},
	// add by v0.4.0
	{
		Name: "task",
		Flags: []cli.Flag{
			flags.TaskReplayScheduleChanSize,
			flags.TaskNeedExecuteChanSize,
			flags.TaskConsResultChanSize,
			flags.TaskMetadataConsumeOption,
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
