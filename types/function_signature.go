package types

import (
	"github.com/datumtechs/datum-network-carrier/common"
)

type LatSingFunc func(hash common.Hash, sig []byte) (string, error)
