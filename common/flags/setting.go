package flags

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/carrier"
	"github.com/RosettaFlow/Carrier-Go/node"
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
)

// SetNodeConfig applies node-related command line flags to the config.
func SetNodeConfig(ctx *cli.Context, cfg *node.Config) {
	setHTTP(ctx, cfg)
	switch {
	case ctx.IsSet(DataDirFlag.Name):
		cfg.DataDir = ctx.String(DataDirFlag.Name)
	case ctx.IsSet(DeveloperFlag.Name):
		cfg.DataDir = "" // unless explicitly requested, use memory databases
	case ctx.IsSet(TestnetFlag.Name):
		cfg.DataDir = filepath.Join(node.DefaultDataDir(), "testnet")
	}
	// todo: more setting....
}

// setHTTP creates the HTTP RPC listener interface string from the set
// command line flags, returning empty if the HTTP endpoint is disabled.
func setHTTP(ctx *cli.Context, cfg *node.Config) {
	if ctx.IsSet(RPCEnabledFlag.Name) && cfg.HTTPHost == "" {
		cfg.HTTPHost = "127.0.0.1"
		if ctx.IsSet(RPCListenAddrFlag.Name) {
			cfg.HTTPHost = ctx.String(RPCListenAddrFlag.Name)
		}
	}
	if ctx.IsSet(RPCPortFlag.Name) {
		cfg.HTTPPort = ctx.Uint64(RPCPortFlag.Name)
	}
	if ctx.IsSet(RPCCORSDomainFlag.Name) {
		cfg.HTTPCors = splitAndTrim(ctx.String(RPCCORSDomainFlag.Name))
	}
	if ctx.IsSet(RPCApiFlag.Name) {
		cfg.HTTPModules = splitAndTrim(ctx.String(RPCApiFlag.Name))
	}
}

// SetCarrierConfig applies eth-related command line flags to the config.
func SetCarrierConfig(ctx *cli.Context, cfg *carrier.Config) {
	// Avoid conflicting network flags
	// todo:
	//checkExclusive(ctx, DeveloperFlag, TestnetFlag)
	cfg.DatabaseHandles = makeDatabaseHandles()

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
		Fatalf("Flags %v can't be used at the same time", strings.Join(set, ", "))
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
		Fatalf("Failed to retrieve file descriptor allowance: %v", err)
	}
	raised, err := fdlimit.Raise(uint64(limit))
	if err != nil {
		Fatalf("Failed to raise file descriptor allowance: %v", err)
	}
	return int(raised / 2) // Leave half for networking and other stuff
}