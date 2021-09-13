package auth

import (
	"github.com/RosettaFlow/Carrier-Go/auth/metadata"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type AuthorityManager struct {
	metadataAuth     *metadata.MetadataAuthority
}


func NewAuthorityManager() *AuthorityManager {
	return &AuthorityManager{
		metadataAuth: metadata.NewMetadataAuthority(),
	}
}

func (am *AuthorityManager) ApplyMetadataAuthority (metadataAuth *types.MetadataAuthApply) error {
	return nil
}

func (am *AuthorityManager) AuditMetadataAuthority (audit *types.MetadataAuthAudit) (apicommonpb.MetadataAuthorityState, error) {

	return 0, nil
}

func (am *AuthorityManager) ConsumeMetadataAuthority (metadataAuthId string) error {

	return nil
}

func (am *AuthorityManager) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {

	return nil, nil
}

func (am *AuthorityManager) GetMetadataAuthorityList () (types.MetadataAuthArray, error) {

	return nil, nil
}

func (am *AuthorityManager) GetMetadataAuthorityListByIds (metadataAuthIds  []string) (types.MetadataAuthArray, error) {

	return nil, nil
}

func (am *AuthorityManager) GetMetadataAuthorityListByUser (userType apicommonpb.UserType, user string) (types.MetadataAuthArray, error) {

	return nil, nil
}

