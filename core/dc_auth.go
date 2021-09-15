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

// InsertMetadataAuthority saves the metadataAuthority to the storage.
func (dc *DataCenter) InsertMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	if metadataAuth.Data().MetadataAuthId == "" {
		metadataAuth.Data().MetadataAuthId = metadataAuth.Hash().Hex()
	}
	response, err := dc.client.SaveMetadataAuthority(dc.ctx, &api.MetadataAuthorityRequest{
		MetadataAuthority: metadataAuth.Data(),
	})
	if err != nil {
		return err
	}
	if response.GetStatus() == 0 {
		return nil
	}
	return errors.New(response.GetMsg())
}

// Functions that do not need to be implemented.
func (dc *DataCenter) RevokeMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	// todo  delete
	return nil
}

func (dc *DataCenter) AuditMetadataAuthority(metadataAuthId, suggestion string, option apicommonpb.AuditMetadataOption) error {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()

	metadataAuthority, err := dc.GetMetadataAuthority(metadataAuthId)
	if err != nil {
		return err
	}
	if metadataAuthority.Data() == nil {
		return errors.New("authority record not found")
	}
	metadataAuthority.Data().AuditSuggestion = suggestion
	metadataAuthority.Data().AuditAt = uint64(timeutils.UnixMsec())
	metadataAuthority.Data().AuditOption = option
	response, err := dc.client.AuditMetadataAuthority(dc.ctx, &api.MetadataAuthorityRequest{
		MetadataAuthority: metadataAuthority.Data(),
	})
	if err != nil {
		return err
	}
	if response.GetStatus() != 0 {
		return errors.New(response.GetMsg())
	}
	return nil
}

// GetMetadataAuthority retrieves the data with the specified metadataAuthId from the data store.
func (dc *DataCenter) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()

	response, err := dc.client.GetMetadataAuthorityList(dc.ctx, &api.ListMetadataAuthorityRequest{
		MetadataAuthId: metadataAuthId,
	})
	if err != nil {
		return nil, err
	}
	if response.GetStatus() != 0 {
		return nil, errors.New(response.GetMsg())
	}
	return types.NewMetadataAuthority(response.GetMetadataAuthorities()[0]), nil
}

func (dc *DataCenter) GetMetadataAuthorityListByIds (metadataAuthIds []string) (types.MetadataAuthArray, error) {

	return nil, nil
}

func (dc *DataCenter) GetMetadataAuthorityListByIdentityId(identityId string, lastUpdate uint64) (types.MetadataAuthArray, error) {
	dc.serviceMu.RLock()
	defer dc.serviceMu.RUnlock()
	response, err := dc.client.GetMetadataAuthorityList(dc.ctx, &api.ListMetadataAuthorityRequest{
		IdentityId:           identityId,
		LastUpdated:           lastUpdate,
	})
	if err != nil {
		return nil, err
	}
	if response.GetStatus() != 0 {
		return nil, errors.New(response.GetMsg())
	}
	return types.NewMetadataAuthArrayFromResponse(response.GetMetadataAuthorities()), nil
}

func (dc *DataCenter) GetMetadataAuthorityListByUser (userType apicommonpb.UserType, user string, lastUpdate uint64) (types.MetadataAuthArray, error) {
	return nil, nil
}