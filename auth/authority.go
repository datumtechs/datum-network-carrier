package auth

import (
	"github.com/RosettaFlow/Carrier-Go/auth/metadata"
	"github.com/RosettaFlow/Carrier-Go/core"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type AuthorityManager struct {
	metadataAuth     *metadata.MetadataAuthority
}


func NewAuthorityManager(dataCenter  core.CarrierDB) *AuthorityManager {
	return &AuthorityManager{
		metadataAuth: metadata.NewMetadataAuthority(dataCenter),
	}
}

func (am *AuthorityManager) ApplyMetadataAuthority (metadataAuth *types.MetadataAuthority) error {
	return am.metadataAuth.ApplyMetadataAuthority(metadataAuth)
}

func (am *AuthorityManager) AuditMetadataAuthority (audit *types.MetadataAuthAudit) (apicommonpb.AuditMetadataOption, error) {
	return am.metadataAuth.AuditMetadataAuthority(audit)
}

func (am *AuthorityManager) ConsumeMetadataAuthority (metadataAuthId string) error {
	return am.metadataAuth.ConsumeMetadataAuthority(metadataAuthId)
}

func (am *AuthorityManager) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {
	return am.metadataAuth.GetMetadataAuthority(metadataAuthId)
}

func (am *AuthorityManager) GetMetadataAuthorityList () (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetMetadataAuthorityList()
}

func (am *AuthorityManager) GetMetadataAuthorityListByIds (metadataAuthIds  []string) (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetMetadataAuthorityList()
}

func (am *AuthorityManager) GetMetadataAuthorityListByUser (userType apicommonpb.UserType, user string) (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetMetadataAuthorityListByUser(userType, user)
}


func (am *AuthorityManager) VerifyMetadataAuth (user, metadataId string, userType apicommonpb.UserType) bool {
	return am.metadataAuth.VerifyMetadataAuth(user, metadataId, userType)
}
