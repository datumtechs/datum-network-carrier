package params

import "fmt"

// CarrierConfig is the core config which determines the rosettaNet settings.
type CarrierConfig struct {

}

// String implements the fmt.Stringer interface.
func (c *CarrierConfig) String() string {
	return fmt.Sprintf("{%v}", "config")
}

// DataCenterConfig is the datacenter service config.
type DataCenterConfig struct {
	GrpcUrl string
	Port uint64
}
