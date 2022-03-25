package core

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"github.com/RosettaFlow/Carrier-Go/types"
)

// on local
func (dc *DataCenter) StoreInternalMetadata(metadata *types.Metadata) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalMetadata(dc.db, metadata)
}

func (dc *DataCenter) IsInternalMetadataByDataId(metadataId string) (bool, error) {
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

func (dc *DataCenter) QueryInternalMetadataByDataId(metadataId string) (*types.Metadata, error) {
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
	if response.Status != 0 {
		return fmt.Errorf("insert metadata error: %s", response.Msg)
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
	if response.Status != 0 {
		return fmt.Errorf("revoke metadata error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) QueryMetadataById(metadataId string) (*types.Metadata, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metadataByIdResponse, err := dc.client.GetMetadataById(dc.ctx, &api.FindMetadataByIdRequest{
		MetadataId: metadataId,
	})
	return types.NewMetadataFromResponse(metadataByIdResponse), err
}

func (dc *DataCenter) QueryMetadataList(lastUpdate, pageSize uint64) (types.MetadataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metaDataListResponse, err := dc.client.GetMetadataList(dc.ctx, &api.ListMetadataRequest{
		LastUpdated: lastUpdate,
		PageSize:    pageSize,
	})
	return types.NewMetadataArrayFromDetailListResponse(metaDataListResponse), err
}

func (dc *DataCenter) QueryMetadataListByIdentity(identityId string, lastUpdate, pageSize uint64) (types.MetadataArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	metaDataListResponse, err := dc.client.GetMetadataListByIdentityId(dc.ctx, &api.ListMetadataByIdentityIdRequest{
		LastUpdated: lastUpdate,
		PageSize:    pageSize,
		IdentityId:  identityId,
	})
	return types.NewMetadataArrayFromDetailListResponse(metaDataListResponse), err
}