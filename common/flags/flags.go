package flags

import (
	"github.com/datumtechs/datum-network-carrier/common/fileutil"
	"github.com/datumtechs/datum-network-carrier/params"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"path/filepath"
	"runtime"
)

var (
	// CertFlag defines a flag for the node's TLS certificate.
	CertFlag = &cli.StringFlag{
		Name:  "tls-cert",
		Usage: "Certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
	// KeyFlag defines a flag for the node's TLS key.
	KeyFlag = &cli.StringFlag{
		Name:  "tls-key",
		Usage: "Key for secure gRPC. Pass this and the tls-cert flag in order to use gRPC securely.",
	}
	// DisableGRPCGateway for JSON-HTTP requests to the beacon node.
	DisableGRPCGateway = &cli.BoolFlag{
		Name:  "disable-grpc-gateway",
		Usage: "Disable the gRPC gateway for JSON-HTTP requests",
	}
	// GPRCGatewayCorsDomain serves preflight requests when serving gRPC JSON gateway.
	GPRCGatewayCorsDomain = &cli.StringFlag{
		Name: "grpc-gateway-corsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests " +
			"(browser enforced). This flag has no effect if not used with --grpc-gateway-port.",
		Value: "http://localhost:4200,http://localhost:7500,http://127.0.0.1:4200,http://127.0.0.1:7500,http://0.0.0.0:4200,http://0.0.0.0:7500",
	}
	// RPCHost defines the host on which the RPC server should listen.
	RPCHost = &cli.StringFlag{
		Name:  "rpc-host",
		Usage: "Host on which the RPC server should listen",
		Value: "127.0.0.1",
	}
	// RPCPort defines a beacon node RPC port to open.
	RPCPort = &cli.IntFlag{
		Name:  "rpc-port",
		Usage: "RPC port exposed by a carrier node",
		Value: 4000,
	}
	// GRPCGatewayHost specifies a gRPC gateway host for Carrier.
	GRPCGatewayHost = &cli.StringFlag{
		Name:  "grpc-gateway-host",
		Usage: "The host on which the gateway server runs on",
		Value: "127.0.0.1",
	}
	// GRPCGatewayPort enables a gRPC gateway to be exposed for Carrier.
	GRPCGatewayPort = &cli.IntFlag{
		Name:  "grpc-gateway-port",
		Usage: "Enable gRPC gateway for JSON requests",
		Value: 3500,
	}
	// GrpcMaxCallRecvMsgSizeFlag defines the max call message size for GRPC
	GrpcMaxCallRecvMsgSizeFlag = &cli.IntFlag{
		Name:  "grpc-max-msg-size",
		Usage: "Integer to define max receive message call size (default: /*4194304*/ 8388608 (for 8MB))",
		Value: 1 << 23,
	}
	// GrpcMaxCallSendMsgSizeFlag defines the max call message size for GRPC
	GrpcMaxCallSendMsgSizeFlag = &cli.IntFlag{
		Name:  "grpc-max-send-msg-size",
		Usage: "Integer to define max send message call size (default: /*4194304*/ 8388608 (for 8MB))",
		Value: 1 << 23,
	}
	// EnableDebugRPCEndpoints
	EnableDebugRPCEndpoints = &cli.BoolFlag{
		Name:  "enable-debug-rpc-endpoints",
		Usage: "Enables the debug rpc service",
	}
	// DataDirFlag defines a path on disk.
	DataDirFlag = &cli.StringFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: DefaultDataDir(),
	}
	// ClearDB prompts user to see if they want to remove any previously stored data at the data directory.
	ClearDB = &cli.BoolFlag{
		Name:  "clear-db",
		Usage: "Prompt for clearing any previously stored data at the data directory",
	}
	// VerbosityFlag defines the logrus configuration.
	VerbosityFlag = &cli.StringFlag{
		Name:  "verbosity",
		Usage: "Logging verbosity (trace, debug, info=default, warn, error, fatal, panic)",
		Value: "info",
	}
	// RestoreSourceFileFlag specifies the filepath to the backed-up database file
	// which will be used to restore the database.
	RestoreSourceFileFlag = &cli.StringFlag{
		Name:  "restore-source-file",
		Usage: "Filepath to the backed-up database file which will be used to restore the database",
	}
	// RestoreTargetDirFlag specifies the target directory of the restored database.
	RestoreTargetDirFlag = &cli.StringFlag{
		Name:  "restore-target-dir",
		Usage: "Target directory of the restored database",
		Value: DefaultDataDir(),
	}
	// ConfigFileFlag specifies the filepath to load flag values.
	ConfigFileFlag = &cli.StringFlag{
		Name:  "config-file",
		Usage: "The filepath to a yaml file with flag values",
	}
	// LogFormat specifies the log output format.
	LogFormat = &cli.StringFlag{
		Name:  "log-format",
		Usage: "Specify log formatting. Supports: text, json, fluentd, journald.",
		Value: "text",
	}
	// LogFileName specifies the log output file name.
	LogFileName = &cli.StringFlag{
		Name:  "log-file",
		Usage: "Specify log file name, relative or absolute",
	}
	DeveloperFlag = &cli.BoolFlag{
		Name:  "dev",
		Usage: "Ephemeral proof-of-authority network with a pre-funded developer account, mining enabled",
	}
	TestnetFlag = &cli.BoolFlag{
		Name:  "testnet",
		Usage: "Ropsten network: pre-configured proof-of-work test network",
	}

	// ****************************** P2P module ************************************
	EnableFakeNetwork = &cli.BoolFlag{
		Name:  "enable-fake-network",
		Usage: "Enable only local network p2p and connect to nodes by special.",
	}
	// NoDiscovery specifies whether we are running a local network and have no need for connecting
	// to the bootstrap nodes in the cloud
	NoDiscovery = &cli.BoolFlag{
		Name:  "no-discovery",
		Usage: "Enable only local network p2p and do not connect to cloud bootstrap nodes.",
	}
	// StaticPeers specifies a set of peers to connect to explicitly.
	StaticPeers = &cli.StringSliceFlag{
		Name:  "peer",
		Usage: "Connect with this peer. This flag may be used multiple times.",
	}
	// BootstrapNode tells the beacon node which bootstrap node to connect to
	BootstrapNode = &cli.StringSliceFlag{
		Name:  "bootstrap-node",
		Usage: "The address of bootstrap node. Beacon node will connect for peer discovery via DHT.  Multiple nodes can be passed by using the flag multiple times but not comma-separated. You can also pass YAML files containing multiple nodes.",
		Value: cli.NewStringSlice(params.CarrierNetworkConfig().BootstrapNodes...),
	}
	// RelayNode tells the beacon node which relay node to connect to.
	RelayNode = &cli.StringFlag{
		Name: "relay-node",
		Usage: "The address of relay node. The beacon node will connect to the " +
			"relay node and advertise their address via the relay node to other peers",
		Value: "",
	}
	// P2PUDPPort defines the port to be used by discv5.
	P2PUDPPort = &cli.IntFlag{
		Name:  "p2p-udp-port",
		Usage: "The port used by discv5.",
		Value: 12000,
	}
	// P2PTCPPort defines the port to be used by libp2p.
	P2PTCPPort = &cli.IntFlag{
		Name:  "p2p-tcp-port",
		Usage: "The port used by libp2p.",
		Value: 13000,
	}
	// P2PIP defines the local IP to be used by libp2p.
	P2PIP = &cli.StringFlag{
		Name:  "p2p-local-ip",
		Usage: "The local ip address to listen for incoming data.",
		Value: "",
	}
	// P2PHost defines the host IP to be used by libp2p.
	P2PHost = &cli.StringFlag{
		Name:  "p2p-host-ip",
		Usage: "The IP address advertised by libp2p. This may be used to advertise an external IP.",
		Value: "",
	}
	// P2PHostDNS defines the host DNS to be used by libp2p.
	P2PHostDNS = &cli.StringFlag{
		Name:  "p2p-host-dns",
		Usage: "The DNS address advertised by libp2p. This may be used to advertise an external DNS.",
		Value: "",
	}
	// P2PPrivKey defines a flag to specify the location of the private key file for libp2p.
	P2PPrivKey = &cli.StringFlag{
		Name:  "p2p-priv-key",
		Usage: "The file containing the private key to use in communications with other peers.",
		Value: "",
	}
	// P2PMetadata defines a flag to specify the location of the peer metadata file.
	P2PMetadata = &cli.StringFlag{
		Name:  "p2p-metadata",
		Usage: "The file containing the metadata to communicate with other peers.",
		Value: "",
	}
	// P2PMaxPeers defines a flag to specify the max number of peers in libp2p.
	P2PMaxPeers = &cli.IntFlag{
		Name:  "p2p-max-peers",
		Usage: "The max number of p2p peers to maintain.",
		Value: 45,
	}
	// P2PAllowList defines a CIDR subnet to exclusively allow connections.
	P2PAllowList = &cli.StringFlag{
		Name: "p2p-allowlist",
		Usage: "The CIDR subnet for allowing only certain peer connections. Example: " +
			"192.168.0.0/16 would permit connections to peers on your local network only. The " +
			"default is to accept all connections.",
	}
	// P2PDenyList defines a list of CIDR subnets to disallow connections from them.
	P2PDenyList = &cli.StringSliceFlag{
		Name: "p2p-denylist",
		Usage: "The CIDR subnets for denying certainy peer connections. Example: " +
			"192.168.0.0/16 would deny connections from peers on your local network only. The " +
			"default is to accept all connections.",
	}
	// EnableUPnPFlag specifies if UPnP should be enabled or not. The default value is false.
	EnableUPnPFlag = &cli.BoolFlag{
		Name:  "enable-upnp",
		Usage: "Enable the service (Carrier chain or Validator) to use UPnP when possible.",
	}
	// DisableDiscv5 disables running discv5.
	DisableDiscv5 = &cli.BoolFlag{
		Name:  "disable-discv5",
		Usage: "Does not run the discoveryV5 dht.",
	}
	// SetGCPercent is the percentage of current live allocations at which the garbage collector is to run.
	SetGCPercent = &cli.IntFlag{
		Name:  "gc-percent",
		Usage: "The percentage of freshly allocated data to live data on which the gc will be run again.",
		Value: 100,
	}

	// ================================= Tracing Flags ===========================================
	// EnableTracingFlag defines a flag to enable p2p message tracing.
	EnableTracingFlag = &cli.BoolFlag{
		Name:  "enable-tracing",
		Usage: "Enable request tracing.",
	}
	// TracingProcessNameFlag defines a flag to specify a process name.
	TracingProcessNameFlag = &cli.StringFlag{
		Name:  "tracing-process-name",
		Usage: "The name to apply to tracing tag \"process_name\"",
	}
	// TracingEndpointFlag flag defines the http endpoint for serving traces to Jaeger.
	TracingEndpointFlag = &cli.StringFlag{
		Name:  "tracing-endpoint",
		Usage: "Tracing endpoint defines where network traces are exposed to Jaeger.",
		Value: "http://127.0.0.1:14268/api/traces",
	}
	// TraceSampleFractionFlag defines a flag to indicate what fraction of p2p
	// messages are sampled for tracing.
	TraceSampleFractionFlag = &cli.Float64Flag{
		Name:  "trace-sample-fraction",
		Usage: "Indicate what fraction of p2p messages are sampled for tracing.",
		Value: 0.20,
	}

	// ================================= Service Discovery Center Flags ===========================================
	DiscoveryServerIP = &cli.StringFlag{
		Name:  "discovery-server-ip",
		Usage: "The ip of discovery center",
		Value: "",
	}

	DiscoveryServerPort = &cli.IntFlag{
		Name:  "discovery-server-port",
		Usage: "The port of discovery center",
		Value: 0,
	}

	DiscoveryServiceId = &cli.StringFlag{
		Name:  "discovery-service-id",
		Usage: "The id of the current service registered with the discovery center",
		Value: "",
	}

	DiscoveryServiceName = &cli.StringFlag{
		Name:  "discovery-service-name",
		Usage: "The name of the current service registered with the discovery center",
		Value: params.CarrierConfig().DiscoveryServiceName,
	}

	DiscoveryServerTags = &cli.StringSliceFlag{
		Name:  "discovery-service-tags",
		Usage: "Tags registered with the service discovery center",
		Value: cli.NewStringSlice(params.CarrierConfig().DiscoveryServiceTags...),
	}

	DiscoveryServiceHealthCheckInterval = &cli.IntFlag{
		Name:  "discovery-service-health-check-interval",
		Usage: "Health check interval between service discovery center and this service (unit: ms, default 3s)",
		Value: 3000,
	}

	DiscoveryServiceHealthCheckDeregister = &cli.IntFlag{
		Name:  "discovery-service-health-check-deregister",
		Usage: "When the service leaves, the service discovery center removes the service information (unit: ms, default 10s)",
		Value: 10000,
	}

	// +++++++++++++++++++++++++++++++++++++++++ Mock Flags +++++++++++++++++++++++++++++++++++++++++
	MockIdentityIdFile = &cli.StringFlag{
		Name:  "mock-identity-file",
		Usage: "Specifies the file path of the identityid information required by mock.",
		Value: "",
	}

	// ================================= consensus Flags ===========================================

	// consensus state
	ConsensusStateWalDir = &cli.StringFlag{
		Name:  "consensus-state-wal-dir",
		Usage: "Configuration dir required to persist the consensus state.",
		Value: "",
	}

	// ================================= Blockchain Flags ===========================================

	// add by v0.4.0
	BlockChain = &cli.StringFlag{
		Name:  "blockchain.url",
		Usage: "blockchain node endpoint.",
		Value: "",
	}

	// ================================= KMS config Flags ===========================================

	// add by v0.4.0
	KMSKeyId = &cli.StringFlag{
		Name:  "kms.key-id",
		Usage: "KMS KeyId.",
		Value: "",
	}

	KMSRegionId = &cli.StringFlag{
		Name:  "kms.region-id",
		Usage: "KMS RegionId.",
		Value: "",
	}

	KMSAccessKeyId = &cli.StringFlag{
		Name:  "kms.access-key-id",
		Usage: "KMS AccessKeyId.",
		Value: "",
	}

	KMSAccessKeySecret = &cli.StringFlag{
		Name:  "kms.access-key-secret",
		Usage: "KMS AccessKeySecret.",
		Value: "",
	}

	// ================================= TaskManager Flags ===========================================

	// add by v.04.0
	TaskReplayScheduleChanSize = &cli.IntFlag{
		Name:  "task.replay-schedule-chan-size",
		Usage: "The buffer size for task participants to receive task information from the task sender [re-calculation operations are required] (default 600)",
		Value: 600,
	}

	TaskNeedExecuteChanSize = &cli.IntFlag{
		Name:  "task.need-execute-chan-size",
		Usage: "The buffer size for the task manager to receive tasks that are about to perform subsequent operations [that have failed or succeeded in consensus] (default 600)",
		Value: 600,
	}

	TaskConsResultChanSize = &cli.IntFlag{
		Name:  "task.consensus-result-chan-size",
		Usage: "The buffer size for task sender will handling has finished one epoch consensus [failed or successful] (default 600)",
		Value: 600,
	}

	TaskMetadataConsumeOption = &cli.IntFlag{
		Name:  "task.metadata-consume-option",
		Usage: "Metadata consumption options for tasks (default 0, 0 mean nothing, 1 mean metadataAuth, 2 mean datatoken)",
		Value: 0,
	}
)

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := fileutil.HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Carrier")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Local", "Carrier")
		} else {
			return filepath.Join(home, ".carrier")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

// LoadFlagsFromConfig sets flags values from config file if ConfigFileFlag is set.
func LoadFlagsFromConfig(cliCtx *cli.Context, flags []cli.Flag) error {
	if cliCtx.IsSet(ConfigFileFlag.Name) {
		if err := altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(ConfigFileFlag.Name))(cliCtx); err != nil {
			return err
		}
	}
	return nil
}
