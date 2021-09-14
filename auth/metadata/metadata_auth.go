package metadata

import (
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type MetadataAuthority struct {
	resourceMng *resource.Manager
}


func NewMetadataAuthority (resourceMng *resource.Manager) *MetadataAuthority {
	return &MetadataAuthority{

	}
}



func (ma *MetadataAuthority) ApplyMetadataAuthority (metadataAuth *types.MetadataAuthApply) error {
	return nil
}

func (ma *MetadataAuthority) AuditMetadataAuthority (audit *types.MetadataAuthAudit) (apicommonpb.MetadataAuthorityState, error) {

	return 0, nil
}

func (ma *MetadataAuthority) ConsumeMetadataAuthority (metadataAuthId string) error {

	return nil
}

func (ma *MetadataAuthority) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {

	return nil, nil
}

func (ma *MetadataAuthority) GetMetadataAuthorityList () (types.MetadataAuthArray, error) {

	return nil, nil
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByIds (metadataAuthIds  []string) (types.MetadataAuthArray, error) {

	return nil, nil
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByUser (userType apicommonpb.UserType, user string) (types.MetadataAuthArray, error) {

	return nil, nil
}
