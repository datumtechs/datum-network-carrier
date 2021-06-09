package params

import "fmt"

// RosettaConfig is the core config which determines the rosettaNet settings.
type RosettaConfig struct {

}

// String implements the fmt.Stringer interface.
func (c *RosettaConfig) String() string {
	return fmt.Sprintf("{%v}", "config")
}

// DataCenterConfig is the datacenter service config.
type DataCenterConfig struct {
	GrpcUrl string
	Port uint64
}
