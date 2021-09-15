package core

import (
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)


// about identity on local
func (dc *DataCenter) StoreIdentity(identity *apicommonpb.Organization) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.WriteLocalIdentity(dc.db, identity)
	return nil
}

func (dc *DataCenter) RemoveIdentity() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	rawdb.DeleteLocalIdentity(dc.db)
	return nil
}

func (dc *DataCenter) GetIdentityId() (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	identity, err := rawdb.ReadLocalIdentity(dc.db)
	if nil != err {
		return "", err
	}
	return identity.GetIdentityId(), nil
}

func (dc *DataCenter) GetIdentity() (*apicommonpb.Organization, error) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	identity, err := rawdb.ReadLocalIdentity(dc.db)
	if nil != err {
		return nil, err
	}
	return identity, nil
}


// about identity on datacenter
func (dc *DataCenter) HasIdentity(identity *apicommonpb.Organization) (bool, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	responses, err := dc.client.GetIdentityList(dc.ctx, &api.ListIdentityRequest{
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
	response, err := dc.client.RevokeIdentityJoin(dc.ctx, &api.RevokeIdentityRequest{
		/*Member: &apicommonpb.Organization{
			NodeName:       identity.Name(),
			NodeId:     identity.NodeId(),
			IdentityId: identity.IdentityId(),
		},*/
		IdentityId: identity.IdentityId(),
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
	identityListResponse, err := dc.client.GetIdentityList(dc.ctx, &api.ListIdentityRequest{LastUpdated: uint64(timeutils.UnixMsec())})
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

func (dc *DataCenter) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {
	response, err := dc.client.FindMetadataAuthority(dc.ctx, &api.FindMetadataAuthorityRequest{
		MetadataAuthId: metadataAuthId,
	})
	if err != nil {
		return nil, err
	}
	return types.NewMetadataAuthority(response.MetadataAuthority), nil
}

func (dc *DataCenter) GetMetadataAuthorityListByIds (metadataAuthIds []string) (types.MetadataAuthArray, error) {

	return nil, nil
}

func (dc *DataCenter) GetMetadataAuthorityListByIdentityId(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.GetMetadataAuthorityList(dc.ctx, &api.ListMetadataAuthorityRequest{
		IdentityId:  identityId,
		LastUpdated: lastUpdate,
	})
	if err != nil {
		return nil, err
	}
	return types.NewMetadataAuthArrayFromResponse(response.GetMetadataAuthorities()), nil
}

func (dc *DataCenter) GetMetadataAuthorityListByUser (userType apicommonpb.UserType, user string, lastUpdate uint64) (types.MetadataAuthArray, error) {
	return nil, nil
}