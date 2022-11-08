package node

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrier"
	"github.com/datumtechs/datum-network-carrier/common/flags"
	"github.com/datumtechs/datum-network-carrier/params"
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

// SetNodeConfig applies node-related command line flags to the config.
func SetNodeConfig(ctx *cli.Context, cfg *Config) {
	switch {
	case ctx.IsSet(flags.DataDirFlag.Name):
		cfg.DataDir = ctx.String(flags.DataDirFlag.Name)
	case ctx.IsSet(flags.DeveloperFlag.Name):
		cfg.DataDir = "" // unless explicitly requested, use memory databases
	case ctx.IsSet(flags.TestnetFlag.Name):
		cfg.DataDir = filepath.Join(DefaultDataDir(), "testnet")
	}
	// todo: more setting....
}

// SetCarrierConfig applies eth-related command line flags to the config.
func SetCarrierConfig(ctx *cli.Context, cfg *carrier.Config) {
	// Avoid conflicting network flags
	//checkExclusive(ctx, DeveloperFlag, TestnetFlag)
	cfg.DatabaseHandles = makeDatabaseHandles()
	cfg.DiscoverServiceConfig = &params.DiscoverServiceConfig{
		DiscoveryServerIP:                     ctx.String(flags.DiscoveryServerIP.Name),
		DiscoveryServerPort:                   ctx.Int(flags.DiscoveryServerPort.Name),
		DiscoveryServiceId:                    ctx.String(flags.DiscoveryServiceId.Name),
		DiscoveryServiceName:                  ctx.String(flags.DiscoveryServiceName.Name),
		DiscoveryServerTags:                   ctx.StringSlice(flags.DiscoveryServerTags.Name),
		DiscoveryServiceHealthCheckInterval:   ctx.Int(flags.DiscoveryServiceHealthCheckInterval.Name),
		DiscoveryServiceHealthCheckDeregister: ctx.Int(flags.DiscoveryServiceHealthCheckDeregister.Name),
	}
	cfg.TaskManagerConfig = &params.TaskManagerConfig{
		//MetadataConsumeOption:          ctx.Int(flags.TaskMetadataConsumeOption.Name),
		NeedReplayScheduleTaskChanSize: ctx.Int(flags.TaskReplayScheduleChanSize.Name),
		NeedExecuteTaskChanSize:        ctx.Int(flags.TaskNeedExecuteChanSize.Name),
		//TaskConsResultChanSize:         ctx.Int(flags.TaskConsResultChanSize.Name),
	}

	// override any default configs.
	switch {
	case ctx.IsSet(flags.TestnetFlag.Name):
		log.Warn("Running on Testnet")
		params.UseTestnetConfig()
		params.UseTestnetNetworkConfig()
	}

	// Override any default configs for hard coded networks.
	// todo: more setting...
}

// checkExclusive verifies that only a single instance of the provided flags was
// set by the user. Each flag might optionally be followed by a string type to
// specialize it further.
func checkExclusive(ctx *cli.Context, args ...interface{}) {
	set := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		// Make sure the next argument is a flag and skip if not set
		flag, ok := args[i].(cli.Flag)
		if !ok {
			panic(fmt.Sprintf("invalid argument, not cli.Flag type: %T", args[i]))
		}
		// Check if next arg extends current and expand its name if so
		name := flag.String()

		if i+1 < len(args) {
			switch option := args[i+1].(type) {
			case string:
				// Extended flag check, make sure value set doesn't conflict with passed in option
				if ctx.String(flag.String()) == option {
					name += "=" + option
					set = append(set, "--"+name)
				}
				// shift arguments and continue
				i++
				continue

			case cli.Flag:
			default:
				panic(fmt.Sprintf("invalid argument, not cli.Flag or string extension: %T", args[i+1]))
			}
		}
		// Mark the flag if it's set
		if ctx.IsSet(flag.String()) {
			set = append(set, "--"+name)
		}
	}
	if len(set) > 1 {
		flags.Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
	}
}

// splitAndTrim splits input separated by a comma
// and trims excessive white space from the substrings.
func splitAndTrim(input string) []string {
	result := strings.Split(input, ",")
	for i, r := range result {
		result[i] = strings.TrimSpace(r)
	}
	return result
}

// makeDatabaseHandles raises out the number of allowed file handles per process
// for Geth and returns half of the allowance to assign to the database.
func makeDatabaseHandles() int {
	limit, err := fdlimit.Maximum()
	if err != nil {
		flags.Fatalf("Failed to retrieve file descriptor allowance: %v", err)
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		flags.Fatalf("Failed to raise file descriptor allowance: %v", err)
	}
	return int(raised / 2) // Leave half for networking and other stuff
}
