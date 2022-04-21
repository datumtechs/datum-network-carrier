package core

import (
	"github.com/Metisnetwork/Metis-Carrier/core/iface"
	"github.com/Metisnetwork/Metis-Carrier/types"
)

type CarrierDB interface {
	iface.LocalStoreCarrierDB
	iface.MetadataCarrierDB
	iface.MetadataAuthorityCarrierDB
	iface.ResourceCarrierDB
	iface.IdentityCarrierDB
	iface.TaskCarrierDB
	InsertData(blocks types.Blocks) (int, error)
	Stop()
}
