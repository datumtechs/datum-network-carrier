package task

import (
	"github.com/datumtechs/datum-network-carrier/types"
)

type taskConsumeOption struct {
	PartyId           string
	MetadataId        string
	DataConsumeOption types.DataConsumeOption
}

func (option *taskConsumeOption) GetPartyId() string    { return option.PartyId }
func (option *taskConsumeOption) GetMetadataId() string { return option.MetadataId }
func (option *taskConsumeOption) GetDataConsumeOption() types.DataConsumeOption {
	return option.DataConsumeOption
}

//type metadataConsumeOption struct {
//	PartyId           string
//	MetadataId        string
//	DataConsumeOption types.MetadataConsumeOptionTK
//}
//
//func (option *metadataConsumeOption) GetPartyId() string    { return option.PartyId }
//func (option *metadataConsumeOption) GetMetadataId() string { return option.MetadataId }
//func (option *metadataConsumeOption) GetDataConsumeOption() types.MetadataConsumeOptionTK {
//	return option.DataConsumeOption
//}
