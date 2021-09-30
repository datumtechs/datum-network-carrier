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

func (am *AuthorityManager) GetLocalMetadataAuthorityList () (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetLocalMetadataAuthorityList()
}

func (am *AuthorityManager) GetGlobalMetadataAuthorityList () (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetGlobalMetadataAuthorityList()
}

func (am *AuthorityManager) GetMetadataAuthorityListByIds (metadataAuthIds  []string) (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetMetadataAuthorityListByIds(metadataAuthIds)
}

func (am *AuthorityManager) HasValidLastMetadataAuth (userType apicommonpb.UserType, user, metadataId string) (bool, error) {
	return am.metadataAuth.HasValidLastMetadataAuth(userType, user, metadataId)
}

func (am *AuthorityManager) VerifyMetadataAuth (userType apicommonpb.UserType, user, metadataId string) bool {
	return am.metadataAuth.VerifyMetadataAuth(userType, user, metadataId)
}

func (am *AuthorityManager) StoreUserMetadataAuthUsed (userType apicommonpb.UserType, user, metadataAuthId string)  error {
	return am.metadataAuth.StoreUserMetadataAuthUsed(userType, user, metadataAuthId)
}

func (am *AuthorityManager) StoreUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId, metadataAuthId string) error {
	return am.metadataAuth.StoreUserMetadataAuthIdByMetadataId(userType, user, metadataId, metadataAuthId)
}

func  (am *AuthorityManager) QueryMetadataAuthIdByMetadataId(userType apicommonpb.UserType, user, metadataId string) (string, error) {
	return am.metadataAuth.QueryMetadataAuthIdByMetadataId(userType, user, metadataId)
}

func  (am *AuthorityManager) RemoveUserMetadataAuthIdByMetadataId (userType apicommonpb.UserType, user, metadataId string) error {
	return am.metadataAuth.RemoveUserMetadataAuthIdByMetadataId(userType, user, metadataId)
}