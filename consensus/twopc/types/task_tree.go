package types

import "github.com/RosettaFlow/Carrier-Go/types"

type TaskTree struct {

}


type taskExt struct {
	*types.ConsensusResult
	CreateAt  uint64
}