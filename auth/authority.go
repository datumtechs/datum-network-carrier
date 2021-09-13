package auth

import (
	"github.com/RosettaFlow/Carrier-Go/auth/metadata"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
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



func (am *AuthorityManager) AuditMetadataAuthority (audit *types.AuditMetadataAuth) error {
	return nil
}

//func (am *AuthorityManager) QueryMetadataAuthorityAudit (metadataAuthId string) (apicommonpb.MetadataAuthorityState, error) {
//	return 0, nil
//}


func (am *AuthorityManager) GetMetadataAuthority (metadataAuthId string) (*pb.GetMetadataAuthority, error) {

	return nil, nil
}

func (am *AuthorityManager) GetMetadataAuthorityList (metadataAuthIds  []string) (*pb.GetMetadataAuthority, error) {

	return nil, nil
}

func (am *AuthorityManager) GetMetadataAuthorityListByUser (userType apicommonpb.UserType) (*pb.GetMetadataAuthority, error) {

	return nil, nil
}

