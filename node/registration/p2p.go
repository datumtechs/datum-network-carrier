package registration

import (
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// P2PPreregistration prepares data for p2p.Service's registration.
func P2PPreregistration(cliCtx *cli.Context) (bootstrapNodeAddrs []string, dataDir string, err error) {
	// Bootnode ENR may be a filepath to a YAML file
	bootnodesTemp := params.CarrierNetworkConfig().BootstrapNodes // actual CLI values
	bootstrapNodeAddrs = make([]string, 0)                       // dest of final list of nodes
	for _, addr := range bootnodesTemp {
		if filepath.Ext(addr) == ".yaml" {
			fileNodes, err := readbootNodes(addr)
			if err != nil {
				return nil, "", err
			}
			bootstrapNodeAddrs = append(bootstrapNodeAddrs, fileNodes...)
		} else {
			bootstrapNodeAddrs = append(bootstrapNodeAddrs, addr)
		}
	}

	dataDir = cliCtx.String(flags.DataDirFlag.Name)
	if dataDir == "" {
		dataDir = flags.DefaultDataDir()
		if dataDir == "" {
			log.Fatal(
				"Could not determine your system's HOME path, please specify a --datadir you wish " +
					"to use for your chain data",
			)
		}
	}

	return
}

func readbootNodes(fileName string) ([]string, error) {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	listNodes := make([]string, 0)
	err = yaml.Unmarshal(fileContent, &listNodes)
	if err != nil {
		return nil, err
	}
	return listNodes, nil
}
