package policy

import (
	"github.com/datumtechs/datum-network-carrier/carrierdb"
)

//type ResourcEngineDB interface {
//	GetDB() core.CarrierDB
//}

type PolicyEngine struct {
	db carrierdb.CarrierDB
}

func NewPolicyEngine(db carrierdb.CarrierDB) *PolicyEngine {
	return &PolicyEngine{
		db: db,
	}
}


