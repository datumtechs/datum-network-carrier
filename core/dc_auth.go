package core

import (
	"errors"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/core/rawdb"
	"github.com/datumtechs/datum-network-carrier/pb/datacenter/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
	"strings"
)

// about identity on local
func (dc *DataCenter) StoreIdentity(identity *carriertypespb.Organization) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalIdentity(dc.db, identity)
}

func (dc *DataCenter) RemoveIdentity() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalIdentity(dc.db)
}

func (dc *DataCenter) QueryIdentityId() (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	identity, err := rawdb.QueryLocalIdentity(dc.db)
	if nil != err {
		return "", err
	}
	return identity.GetIdentityId(), nil
}

func (dc *DataCenter) QueryIdentity() (*carriertypespb.Organization, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.QueryLocalIdentity(dc.db)
}

// about identity on datacenter
func (dc *DataCenter) HasIdentity(identity *carriertypespb.Organization) (bool, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	responses, err := dc.client.GetIdentityList(dc.ctx, &api.ListIdentityRequest{
		LastUpdated: timeutils.BeforeYearUnixMsecUint64(),
		PageSize:    backend.DefaultMaxPageSize,
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
	response, err := dc.client.RevokeIdentityJoin(dc.ctx, &api.RevokeIdentityRequest{
		IdentityId: identity.GetIdentityId(),
	})
	if err != nil {
		return err
	}
	if response.GetStatus() != 0 {
		return fmt.Errorf("revokeIdeneity err: %s", response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) QueryIdentityList(lastUpdate uint64, pageSize uint64) (types.IdentityArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	identityListResponse, err := dc.client.GetIdentityList(dc.ctx, &api.ListIdentityRequest{LastUpdated: lastUpdate, PageSize: pageSize})
	return types.NewIdentityArrayFromIdentityListResponse(identityListResponse), err
}

func (dc *DataCenter) InsertMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	_, err := dc.client.SaveMetadataAuthority(dc.ctx, &api.MetadataAuthorityRequest{
		MetadataAuthority: metadataAuth.GetData(),
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateMetadataAuthority updates the data already stored.
func (dc *DataCenter) UpdateMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.UpdateMetadataAuthority(dc.ctx, &api.MetadataAuthorityRequest{
		MetadataAuthority: metadataAuth.GetData(),
	})
	if err != nil {
		return err
	}
	if response.GetStatus() != 0 {
		return errors.New(response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) QueryMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	response, err := dc.client.FindMetadataAuthority(dc.ctx, &api.FindMetadataAuthorityRequest{
		MetadataAuthId: metadataAuthId,
	})
	if err != nil {
		return nil, err
	}
	return types.NewMetadataAuthority(response.MetadataAuthority), nil
}

func (dc *DataCenter) QueryMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	//todo: Discussion: Do we need to achieve?
	panic("implements me")
	return nil, nil
}

func (dc *DataCenter) QueryMetadataAuthorityListByIdentityId(identityId string, lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.GetMetadataAuthorityList(dc.ctx, &api.ListMetadataAuthorityRequest{
		IdentityId:  identityId,
		LastUpdated: lastUpdate,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, err
	}
	return types.NewMetadataAuthArrayFromResponse(response.GetMetadataAuthorities()), nil
}

func (dc *DataCenter) QueryMetadataAuthorityList(lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.GetMetadataAuthorityList(dc.ctx, &api.ListMetadataAuthorityRequest{
		LastUpdated: lastUpdate,
		PageSize:    pageSize,
	})
	if err != nil {
		return nil, err
	}
	return types.NewMetadataAuthArrayFromResponse(response.GetMetadataAuthorities()), nil
}
