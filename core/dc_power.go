package core

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/lib/center/api"
	"github.com/RosettaFlow/Carrier-Go/types"
)

// about power on local
func (dc *DataCenter) InsertLocalResource(resource *types.LocalResource) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreLocalResource(dc.db, resource)
}

func (dc *DataCenter) RemoveLocalResource(jobNodeId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveLocalResource(dc.db, jobNodeId)
}

func (dc *DataCenter) QueryLocalResource(jobNodeId string) (*types.LocalResource, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryLocalResource(dc.db, jobNodeId)
}

func (dc *DataCenter) QueryLocalResourceList() (types.LocalResourceArray, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryAllLocalResource(dc.db)
}

func (dc *DataCenter) StoreJobNodeIdIdByPowerId(powerId, jobNodeId string) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.StoreJobNodeIdByPowerId(dc.db, powerId, jobNodeId)
}

func (dc *DataCenter) RemoveJobNodeIdByPowerId(powerId string) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.RemoveJobNodeIdByPowerId(dc.db, powerId)
}

func (dc *DataCenter) QueryJobNodeIdByPowerId(powerId string) (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryJobNodeIdByPowerId(dc.db, powerId)
}

// about power on datacenter
func (dc *DataCenter) InsertResource(resource *types.Resource) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SaveResource(dc.ctx, types.NewPublishPowerRequest(resource))
	if err != nil {
		log.WithError(err).WithField("hash", resource.Hash()).Errorf("InsertResource failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("insert resource error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) RevokeResource(resource *types.Resource) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.RevokeResource(dc.ctx, types.RevokePowerRequest(resource))
	if err != nil {
		log.WithError(err).WithField("hash", resource.Hash()).Errorf("RevokeResource failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("revoke resource error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) SyncPowerUsed(resource *types.LocalResource) error {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	response, err := dc.client.SyncPower(dc.ctx, types.NewSyncPowerRequest(resource))
	if err != nil {
		log.WithError(err).WithField("hash", resource.Hash()).Errorf("SyncPowerUsed failed")
		return err
	}
	if response.Status != 0 {
		return fmt.Errorf("sync resource used error: %s", response.Msg)
	}
	return nil
}

func (dc *DataCenter) GetResourceListByIdentityId(identityId string) (types.ResourceArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	powerTotalSummaryResponse, err := dc.client.GetPowerSummaryByIdentityId(dc.ctx, &api.GetPowerSummaryByIdentityRequest{
		IdentityId: identityId,
	})
	return types.NewResourceFromResponse(powerTotalSummaryResponse), err
}

func (dc *DataCenter) QueryGlobalResourceSummaryList() (types.ResourceArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	powerListRequest, err := dc.client.GetPowerGlobalSummaryList(dc.ctx)
	return types.NewResourceArrayFromPowerTotalSummaryListResponse(powerListRequest), err
}

func (dc *DataCenter) QueryGlobalResourceDetailList() (types.ResourceArray, error) {
	dc.serviceMu.Lock()
	defer dc.serviceMu.Unlock()
	powerListRequest, err := dc.client.GetPowerList(dc.ctx, &api.ListPowerRequest{LastUpdated: timeutils.BeforeYearUnixMsecUint64()})
	return types.NewResourceArrayFromPowerDetailListResponse(powerListRequest), err
}

// For ResourceManager
// about jobRerource
func (dc *DataCenter) StoreLocalResourceTable(resource *types.LocalResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreNodeResource(dc.db, resource)
}

func (dc *DataCenter) StoreLocalResourceTables(resources []*types.LocalResourceTable) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.StoreNodeResources(dc.db, resources)
}

func (dc *DataCenter) RemoveLocalResourceTable(resourceId string) error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	return rawdb.RemoveNodeResource(dc.db, resourceId)
}

func (dc *DataCenter) QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryNodeResource(dc.db, resourceId)
}

func (dc *DataCenter) QueryLocalResourceTables() ([]*types.LocalResourceTable, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryNodeResources(dc.db)
}
