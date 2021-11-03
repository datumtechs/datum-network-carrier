package registration

import (
	"github.com/RosettaFlow/Carrier-Go/common/fileutil"
	"github.com/RosettaFlow/Carrier-Go/common/flags"
	"github.com/RosettaFlow/Carrier-Go/common/sliceutil"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	datadirStaticNodes     = "static-nodes.json"  // Path within the datadir to the static node list
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

func readStaticNodesFromJSON(fileName string) ([]string, error) {
	var nodeList []string
	if err := fileutil.LoadJSON(fileName, &nodeList); err != nil {
		//log.WithError(err).Errorf("Can't load sttaic node file from JSON %s: %v", fileName, err)
		return nil, err
	}
	return nodeList, nil
}

func P2PStaticNodeAddrs(cliCtx *cli.Context, dataDir string) ([]string, error) {
	// Static ENR may be a filepath to a YAML file
	staticNodeAddrs := make([]string, 0)
	cmdStatics := sliceutil.SplitCommaSeparated(cliCtx.StringSlice(flags.StaticPeers.Name))
	if len(cmdStatics) != 0 {
		staticNodeAddrs = append(staticNodeAddrs, cmdStatics...)
	}
	path := ResolvePath(dataDir, datadirStaticNodes)
	log.Debugf("resolve path for static nodes, path: %s", path)
	if path == "" {
		return staticNodeAddrs, nil
	}
	if _, err := os.Stat(path); err != nil {
		return staticNodeAddrs, nil
	}
	if filepath.Ext(path) == ".yaml" {
		fileNodes, err := readbootNodes(path)
		if err != nil {
			log.WithError(err).Infof("Can't load static node file from YML file %s: %v", path, err)
			return staticNodeAddrs, nil
		}
		staticNodeAddrs = append(staticNodeAddrs, fileNodes...)
	}
	if filepath.Ext(path) == ".json" {
		fileNodes, err := readStaticNodesFromJSON(path)
		if err != nil {
			log.WithError(err).Errorf("Can't load static node file from JSON %s: %v", path, err)
			return staticNodeAddrs, nil
		}
		staticNodeAddrs = append(staticNodeAddrs, fileNodes...)
	}
	return staticNodeAddrs, nil
}

func ResolvePath(dataDir, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if dataDir == "" {
		return ""
	}
	return filepath.Join(dataDir, path)
}