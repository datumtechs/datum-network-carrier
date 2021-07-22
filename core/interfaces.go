package core

import (
	"github.com/RosettaFlow/Carrier-Go/core/iface"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type CarrierDB interface {
	iface.LocalStoreCarrierDB
	iface.MetadataCarrierDB
	iface.ResourceCarrierDB
	iface.IdentityCarrierDB
	iface.TaskCarrierDB
	InsertData(blocks types.Blocks) (int, error)
	Stop()
}
