package main

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/cmd"
	dbcommand "github.com/RosettaFlow/Carrier-Go/cmd/carrier/db"
	"github.com/RosettaFlow/Carrier-Go/cmd/common"
	"github.com/RosettaFlow/Carrier-Go/common/debug"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/common/logutil"
	"github.com/RosettaFlow/Carrier-Go/node"
	gethlog "github.com/ethereum/go-ethereum/log"
	golog "github.com/ipfs/go-log/v2"
	joonix "github.com/joonix/log"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"runtime"
	runtimeDebug "runtime/debug"

	"github.com/urfave/cli/v2"
)

var (
	appFlags = []cli.Flag{
		flags.SetGCPercent,
		flags.DeveloperFlag,
		flags.TestnetFlag,
	}

	nodeFlags = []cli.Flag{
		flags.DataDirFlag,
		flags.ClearDB,
		flags.VerbosityFlag,
		flags.RestoreSourceFileFlag,
		flags.RestoreTargetDirFlag,
		flags.ConfigFileFlag,
		flags.LogFormat,
		flags.LogFileName,
	}

	rpcFlags = []cli.Flag{
		flags.RPCHost,
		flags.RPCPort,
		flags.CertFlag,
		flags.KeyFlag,
		flags.DisableGRPCGateway,
		flags.GPRCGatewayCorsDomain,
		flags.GRPCGatewayHost,
		flags.GRPCGatewayPort,
		flags.GRPCDataCenterHost,
		flags.GRPCDataCenterPort,
		flags.GrpcMaxCallRecvMsgSizeFlag,
		flags.GrpcMaxCallSendMsgSizeFlag,
	}

	p2pFlags = []cli.Flag{
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
	}

	debugFlags = []cli.Flag{
		debug.DebugFlag,
		flags.EnableDebugRPCEndpoints,
		debug.PProfFlag,
		debug.PProfPortFlag,
		debug.PProfAddrFlag,
		debug.MemProfileRateFlag,
		debug.MutexProfileFractionFlag,
		debug.BlockProfileRateFlag,
		debug.CPUProfileFlag,
		debug.TraceFlag,
		flags.EnableTracingFlag,
		flags.TracingProcessNameFlag,
		flags.TracingEndpointFlag,
		flags.TraceSampleFractionFlag,
	}

	mockFlags = []cli.Flag{
		flags.MockIdentityIdFile,
	}

	consensusFlags = []cli.Flag{
		flags.ConsensusStateWalDir,
	}
)

func init() {
	appFlags = cmd.WrapFlags(appFlags)
	nodeFlags = cmd.WrapFlags(nodeFlags)
	rpcFlags = cmd.WrapFlags(rpcFlags)
	p2pFlags = cmd.WrapFlags(p2pFlags)
	debugFlags = cmd.WrapFlags(debugFlags)
	mockFlags = cmd.WrapFlags(mockFlags)
	consensusFlags = cmd.WrapFlags(consensusFlags)
}

func main() {
	app := cli.App{}
	app.Name = "carrier"
	app.Usage = "this is a carrier network implementation for Carrier Node"
	// set action func.
	app.Action = startNode
	app.Version = common.Version()
	app.Commands = []*cli.Command{
		dbcommand.Commands,
	}
	app.Flags = append(app.Flags, appFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, p2pFlags...)
	app.Flags = append(app.Flags, debugFlags...)
	app.Flags = append(app.Flags, mockFlags...)
	app.Flags = append(app.Flags, consensusFlags...)

	app.Before = func(ctx *cli.Context) error {
		// Load flags from config file, if specified.
		if err := flags.LoadFlagsFromConfig(ctx, app.Flags); err != nil {
			return err
		}

		format := ctx.String(flags.LogFormat.Name)
		switch format {
		case "text":
			formatter := new(prefixed.TextFormatter)
			formatter.TimestampFormat = "2006-01-02 15:04:05.000"
			formatter.FullTimestamp = true
			// If persistent log files are written - we disable the log messages coloring because
			// the colors are ANSI codes and seen as gibberish in the log files.
			formatter.DisableColors = ctx.String(flags.LogFileName.Name) != ""
			logrus.SetFormatter(formatter)
		case "fluentd":
			f := joonix.NewFormatter()
			if err := joonix.DisableTimestampFormat(f); err != nil {
				panic(err)
			}
			logrus.SetFormatter(f)
		case "json":
			logrus.SetFormatter(&logrus.JSONFormatter{})
		default:
			return fmt.Errorf("unknown log format %s", format)
		}

		logFileName := ctx.String(flags.LogFileName.Name)
		if logFileName != "" {
			if err := logutil.ConfigurePersistentLogging(logFileName); err != nil {
				log.WithError(err).Error("Failed to configuring logging to disk.")
			}
		}
		if ctx.IsSet(debug.DebugFlag.Name) {
			filenameHook := filename.NewHook()
			filenameHook.Field = "zline"
			logrus.AddHook(filenameHook)
		}
		if ctx.IsSet(flags.SetGCPercent.Name) {
			runtimeDebug.SetGCPercent(ctx.Int(flags.SetGCPercent.Name))
		}
		runtime.GOMAXPROCS(runtime.NumCPU())
		return debug.Setup(ctx)
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
	if args := ctx.Args(); args.Len() > 0 {
		return fmt.Errorf("invalid command: %q", args.Get(0))
	}

	// setting log level.
	verbosity := ctx.String(flags.VerbosityFlag.Name)
	level, err := logrus.ParseLevel(verbosity)
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	if level == logrus.TraceLevel {
		// libp2p specific logging.（special）
		golog.SetAllLoggers(golog.LevelDebug)
		// Geth specific logging.
		glogger := gethlog.NewGlogHandler(gethlog.StreamHandler(os.Stderr, gethlog.TerminalFormat(true)))
		glogger.Verbosity(gethlog.LvlTrace)
		gethlog.Root().SetHandler(glogger)
	}

	// initial no and start.
	node, err := node.New(ctx)
	if err != nil {
		return err
	}
	node.Start()
	return nil
}
