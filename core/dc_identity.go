package core

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	apipb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

// about identity on datacenter
func (dc *DataCenter) HasIdentity(identity *apipb.Organization) (bool, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	responses, err := dc.client.GetIdentityList(dc.ctx, &api.IdentityListRequest{
		LastUpdated: uint64(timeutils.Now().Second()),
	})
	if err != nil {
		return false, err
	}
	for _, organization := range responses.Identities {
		if strings.EqualFold(organization.IdentityId, identity.IdentityId) {
			return true, nil
		}
	}
	return false, nil
}

func (dc *DataCenter) InsertIdentity(identity *types.Identity) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveIdentity(dc.ctx, types.NewSaveIdentityRequest(identity))
	if err != nil {
		//log.WithError(err).WithField("hash", identity.Hash()).Errorf("InsertIdentity failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert indentity error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) RevokeIdentity(identity *types.Identity) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.RevokeIdentityJoin(dc.ctx, &api.RevokeIdentityJoinRequest{
		Member: &apipb.Organization{
			NodeName:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},
	})
	if err != nil {
		return err
	}
	if response.GetStatus() != 0 {
		return fmt.Errorf("revokeIdeneity err: %s", response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) GetIdentityList() (types.IdentityArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	identityListResponse, err := dc.client.GetIdentityList(dc.ctx, &api.IdentityListRequest{LastUpdated: uint64(timeutils.UnixMsec())})
	return types.NewIdentityArrayFromIdentityListResponse(identityListResponse), err
}

// 存储元数据鉴权申请记录
func (dc *DataCenter) SaveMetadataAuthority(request *types.MetadataAuth) error {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	if request.Data().AuthRecordId == "" {
		request.Data().AuthRecordId = request.Hash().Hex()
	}
	response, err := dc.client.SaveMetadataAuthority(dc.ctx, &api.SaveMetadataAuthorityRequest{
		User:                 request.Data().User,
		UserType:             request.Data().UserType,
		Auth:                 request.Data().DataRecord,
		MetadataAuthId:       request.Data().AuthRecordId,
	})
	if err != nil {
		return err
	}
	if response.GetStatus() == 0 {
		return nil
	}
	return errors.New(response.GetMsg())
}

func (dc *DataCenter) AuditMetadataAuthority(authId string, result apipb.AuditMetadataOption) error {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.AuditMetadataAuthority(dc.ctx, &api.AuditMetadataAuthorityRequest{
		MetadataAuthId:       authId,
		Audit:                result,
	})
	if err != nil {
		return err
	}
	if response.GetStatus() != 0 {
		return errors.New(response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) GetMetadataAuthorityList(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.GetMetadataAuthorityList(dc.ctx, &api.MetadataAuthorityListRequest{
		IdentityId:           identityId,
		LastUpdated:           lastUpdate,
	})
	if err != nil {
		return nil, err
	}
	if response.GetStatus() != 0 {
		return nil, errors.New(response.GetMsg())
	}
	return types.NewMetadataAuthArrayFromResponse(response.GetAuthorities()), nil
}