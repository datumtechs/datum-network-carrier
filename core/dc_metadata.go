package core

import (
	"fmt"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	datacenterapipb "github.com/datumtechs/datum-network-carrier/pb/datacenter/api"
	"github.com/datumtechs/datum-network-carrier/types"
)

// on local
func (dc *DataCenter) StoreInternalMetadata(metadata *types.Metadata) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalMetadata(dc.db, metadata)
}

func (dc *DataCenter) IsInternalMetadataById(metadataId string) (bool, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	_, err := rawdb.QueryLocalMetadata(dc.db, metadataId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return false, err
	}
	if rawdb.IsDBNotFoundErr(err) {
		return false, nil
	}
	return true, nil
}

func (dc *DataCenter) QueryInternalMetadataById(metadataId string) (*types.Metadata, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalMetadata(dc.db, metadataId)
}

func (dc *DataCenter) QueryInternalMetadataList() (types.MetadataArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryAllLocalMetadata(dc.db)
}

// on datecenter
func (dc *DataCenter) InsertMetadata(metadata *types.Metadata) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveMetadata(dc.ctx, types.NewMetadataSaveRequest(metadata))
	if err != nil {
		log.WithError(err).WithField("hash", metadata.Hash()).Errorf("InsertMetadata failed")
		return err
	}
	if response.GetStatus() != 0 {
		return fmt.Errorf("insert metadata error: %s", response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) RevokeMetadata(metadata *types.Metadata) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.RevokeMetadata(dc.ctx, types.NewMetadataRevokeRequest(metadata))
	if err != nil {
		log.WithError(err).WithField("hash", metadata.Hash()).Errorf("RevokeMetadata failed")
		return err
	}
	if response.GetStatus() != 0 {
		return fmt.Errorf("revoke metadata error: %s", response.GetMsg())
	}
	return nil
}

func (dc *DataCenter) QueryMetadataById(metadataId string) (*types.Metadata, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metadataByIdResponse, err := dc.client.GetMetadataById(dc.ctx, &datacenterapipb.FindMetadataByIdRequest{
		MetadataId: metadataId,
	})
	if metadataByIdResponse.GetStatus() != 0 {
		return nil, fmt.Errorf("call datacenter GetMetadataById error, status %d,error msg %s", metadataByIdResponse.GetStatus(), metadataByIdResponse.GetMsg())
	}
	return types.NewMetadataFromResponse(metadataByIdResponse), err
}

func (dc *DataCenter) QueryMetadataByIds(metadataIds []string) ([]*types.Metadata, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metaDataListResponse, err := dc.client.GetMetadataByIds(dc.ctx, &datacenterapipb.FindMetadataByIdsRequest{
		MetadataIds: metadataIds,
	})
	if metaDataListResponse.GetStatus() != 0 {
		return nil, fmt.Errorf("call datacenter GetMetadataByIds error, status %d,error msg %s", metaDataListResponse.GetStatus(), metaDataListResponse.GetMsg())
	}
	return types.NewMetadataArrayFromDetailListResponse(metaDataListResponse), err
}

func (dc *DataCenter) QueryMetadataList(lastUpdate, pageSize uint64) (types.MetadataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metaDataListResponse, err := dc.client.GetMetadataList(dc.ctx, &datacenterapipb.ListMetadataRequest{
		LastUpdated: lastUpdate,
		PageSize:    pageSize,
	})
	if metaDataListResponse.GetStatus() != 0 {
		return nil, fmt.Errorf("call datacenter GetMetadataList error, status %d,error msg %s", metaDataListResponse.GetStatus(), metaDataListResponse.GetMsg())
	}
	return types.NewMetadataArrayFromDetailListResponse(metaDataListResponse), err
}

func (dc *DataCenter) QueryMetadataListByIdentity(identityId string, lastUpdate, pageSize uint64) (types.MetadataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metaDataListResponse, err := dc.client.GetMetadataListByIdentityId(dc.ctx, &datacenterapipb.ListMetadataByIdentityIdRequest{
		LastUpdated: lastUpdate,
		PageSize:    pageSize,
		IdentityId:  identityId,
	})
	if metaDataListResponse.GetStatus() != 0 {
		return nil, fmt.Errorf("call datacenter GetMetadataListByIdentityId error, status %d,error msg %s", metaDataListResponse.GetStatus(), metaDataListResponse.GetMsg())
	}
	return types.NewMetadataArrayFromDetailListResponse(metaDataListResponse), err
}

func (dc *DataCenter) UpdateGlobalMetadata(metadata *types.Metadata) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.UpdateMetadata(dc.ctx, types.NewMetadataUpdateRequest(metadata))
	if err != nil {
		log.WithError(err).WithField("hash", metadata.Hash()).Errorf("UpdateMetadata failed")
		return err
	}
	if response.GetStatus() != 0 {
		return fmt.Errorf("update metadata failed: resp.status: %d, resp.msg: %s", response.GetStatus(), response.GetMsg())
	}
	return nil
}
