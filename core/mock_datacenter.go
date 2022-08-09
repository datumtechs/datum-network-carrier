package core

import (
	"errors"
	"fmt"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/types"
)

type MockDataCenter struct {
}

func (mc MockDataCenter) InsertData(blocks types.Blocks) (int, error)      { return 0, nil }
func (mc MockDataCenter) Stop()                                            {}
func (mc MockDataCenter) SetConfig(config *types.CarrierChainConfig) error { return nil }

// about carrier
func (mc MockDataCenter) QueryYarnName() (string, error)                       { return "", nil }
func (mc MockDataCenter) SetSeedNode(seed *carrierapipb.SeedPeer) error        { return nil }
func (mc MockDataCenter) RemoveSeedNode(addr string) error                     { return nil }
func (mc MockDataCenter) QuerySeedNodeList() ([]*carrierapipb.SeedPeer, error) { return nil, nil }
func (mc MockDataCenter) SetRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) error {
	return nil
}
func (mc MockDataCenter) DeleteRegisterNode(typ carrierapipb.RegisteredNodeType, id string) error {
	return nil
}
func (mc MockDataCenter) QueryRegisterNode(typ carrierapipb.RegisteredNodeType, id string) (*carrierapipb.YarnRegisteredPeerDetail, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryRegisterNodeList(typ carrierapipb.RegisteredNodeType) ([]*carrierapipb.YarnRegisteredPeerDetail, error) {
	return nil, nil
}

// about powerId -> jobNodeId
func (mc MockDataCenter) StoreJobNodeIdIdByPowerId(powerId, jobNodeId string) error { return nil }
func (mc MockDataCenter) RemoveJobNodeIdByPowerId(powerId string) error             { return nil }
func (mc MockDataCenter) QueryJobNodeIdByPowerId(powerId string) (string, error)    { return "", nil }

// about jobRerource   (jobNodeId -> {jobNodeId, powerId, resource, slotTotal, slotUsed})
func (mc MockDataCenter) StoreLocalResourceTable(resource *types.LocalResourceTable) error {
	return nil
}
func (mc MockDataCenter) RemoveLocalResourceTable(resourceId string) error { return nil }
func (mc MockDataCenter) StoreLocalResourceTables(resources []*types.LocalResourceTable) error {
	return nil
}
func (mc MockDataCenter) QueryLocalResourceTable(resourceId string) (*types.LocalResourceTable, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalResourceTables() ([]*types.LocalResourceTable, error) {
	return nil, nil
}

// about DataResourceTable (dataNodeId -> {dataNodeId, totalDisk, usedDisk})
func (mc MockDataCenter) StoreDataResourceTable(StoreDataResourceTables *types.DataResourceTable) error {
	return nil
}
func (mc MockDataCenter) StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error {
	return nil
}
func (mc MockDataCenter) RemoveDataResourceTable(nodeId string) error { return nil }
func (mc MockDataCenter) QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryDataResourceTables() ([]*types.DataResourceTable, error) {
	return nil, nil
}

// about DataResourceDataUpload (originId -> {originId, dataNodeId, metaDataId, filePath})
func (mc MockDataCenter) StoreDataResourceDataUpload(dataResourceDataUpload *types.DataResourceDataUpload) error {
	return nil
}
func (mc MockDataCenter) StoreDataResourceDataUploads(dataResourceDataUploads []*types.DataResourceDataUpload) error {
	return nil
}
func (mc MockDataCenter) RemoveDataResourceDataUpload(originId string) error { return nil }
func (mc MockDataCenter) QueryDataResourceDataUpload(originId string) (*types.DataResourceDataUpload, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryDataResourceDataUploads() ([]*types.DataResourceDataUpload, error) {
	return nil, nil
}

// about DataResourceDiskUsed (metaDataId -> {metaDataId, dataNodeId, diskUsed})
func (mc MockDataCenter) StoreDataResourceDiskUsed(dataResourceDiskUsed *types.DataResourceDiskUsed) error {
	return nil
}
func (mc MockDataCenter) RemoveDataResourceDiskUsed(metaDataId string) error { return nil }
func (mc MockDataCenter) QueryDataResourceDiskUsed(metaDataId string) (*types.DataResourceDiskUsed, error) {
	return nil, nil
}

// v 0.5.0  about user metadataAuthStatus by metadataId
//(userType + user + metadataId -> metadataAuthTotalCount
// | userType + user + metadataId -> metadataAuthValidCount
// | userType + user + metadataId + metadataAuthId -> status)
func (mc MockDataCenter) StoreValidUserMetadataAuthStatusByMetadataId(userType commonconstantpb.UserType, user, metadataId, metadataAuthId string, status uint16) error {
	return nil
}
func (mc MockDataCenter) QueryValidUserMetadataAuthIdsByMetadataId(userType commonconstantpb.UserType, user, metadataId string) ([]string, error) {
	return nil, nil
}
func (mc MockDataCenter) HasValidUserMetadataAuthStatusByMetadataId(userType commonconstantpb.UserType, user, metadataId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) HasUserMetadataAuthIdByMetadataId(userType commonconstantpb.UserType, user, metadataId, metadataAuthId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) RemoveUserMetadataAuthStatusByMetadataId(userType commonconstantpb.UserType, user, metadataId, metadataAuthId string) error {
	return nil
}

// v 0.2.0 about metadata used taskId    (metadataId -> [taskId, taskId, ..., taskId])
func (mc MockDataCenter) StoreMetadataHistoryTaskId(metadataId, taskId string) error { return nil }
func (mc MockDataCenter) HasMetadataHistoryTaskId(metadataId, taskId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) QueryMetadataHistoryTaskIdCount(metadataId string) (uint32, error) {
	return 0, nil
}
func (mc MockDataCenter) QueryMetadataHistoryTaskIds(metadataId string) ([]string, error) {
	return nil, nil
}

// v 0.2.0 about Message Cache
func (mc MockDataCenter) StoreMessageCache(value interface{}) error          { return nil }
func (mc MockDataCenter) RemovePowerMsg(powerId string) error                { return nil }
func (mc MockDataCenter) RemoveAllPowerMsg() error                           { return nil }
func (mc MockDataCenter) RemoveMetadataMsg(metadataId string) error          { return nil }
func (mc MockDataCenter) RemoveMetadataUpdateMsg(metadataId string) error    { return nil }
func (mc MockDataCenter) RemoveAllMetadataMsg() error                        { return nil }
func (mc MockDataCenter) RemoveMetadataAuthMsg(metadataAuthId string) error  { return nil }
func (mc MockDataCenter) RemoveAllMetadataAuthMsg() error                    { return nil }
func (mc MockDataCenter) RemoveTaskMsg(taskId string) error                  { return nil }
func (mc MockDataCenter) RemoveAllTaskMsg() error                            { return nil }
func (mc MockDataCenter) QueryPowerMsgArr() (types.PowerMsgArr, error)       { return nil, nil }
func (mc MockDataCenter) QueryMetadataMsgArr() (types.MetadataMsgArr, error) { return nil, nil }
func (mc MockDataCenter) QueryMetadataUpdateMsgArr() (types.MetadataUpdateMsgArr, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryMetadataAuthorityMsgArr() (types.MetadataAuthorityMsgArr, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskMsgArr() (types.TaskMsgArr, error) { return nil, nil }
func (mc MockDataCenter) SaveOrgPriKey(priKey string) error          { return nil }

// FindOrgPriKey does not return ErrNotFound if the organization wallet not found.
func (mc MockDataCenter) FindOrgPriKey() (string, error) { return "", nil }

func (mc MockDataCenter) InsertMetadata(metadata *types.Metadata) error { return nil }
func (mc MockDataCenter) RevokeMetadata(metadata *types.Metadata) error { return nil }
func (mc MockDataCenter) QueryMetadataById(metadataId string) (*types.Metadata, error) {
	return generateTestMetadata(metadataId)
}
func (mc MockDataCenter) QueryMetadataByIds(metadataIds []string) ([]*types.Metadata, error) {
	return nil, nil
} // add by v 0.4.0
func (mc MockDataCenter) QueryMetadataList(lastUpdate, pageSize uint64) (types.MetadataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryMetadataListByIdentity(identityId string, lastUpdate, pageSize uint64) (types.MetadataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) UpdateGlobalMetadata(metadata *types.Metadata) error { return nil } // add by v 0.4.0
// v 0.3.0 about internal metadata (generate from a task result file)
func (mc MockDataCenter) StoreInternalMetadata(metadata *types.Metadata) error   { return nil }
func (mc MockDataCenter) IsInternalMetadataById(metadataId string) (bool, error) { return false, nil }
func (mc MockDataCenter) QueryInternalMetadataById(metadataId string) (*types.Metadata, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryInternalMetadataList() (types.MetadataArray, error) { return nil, nil }

// about global power resource
func (mc MockDataCenter) InsertResource(resource *types.Resource) error { return nil }
func (mc MockDataCenter) RevokeResource(resource *types.Resource) error { return nil }
func (mc MockDataCenter) QueryGlobalResourceSummaryList() (types.ResourceArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryGlobalResourceDetailList(lastUpdate, pageSize uint64) (types.ResourceArray, error) {
	return nil, nil
}
func (mc MockDataCenter) SyncPowerUsed(resource *types.LocalResource) error { return nil }

// about local resource (local jobNode resource)
func (mc MockDataCenter) StoreLocalResource(resource *types.LocalResource) error { return nil }
func (mc MockDataCenter) RemoveLocalResource(jobNodeId string) error             { return nil }
func (mc MockDataCenter) QueryLocalResource(jobNodeId string) (*types.LocalResource, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalResourceList() (types.LocalResourceArray, error) { return nil, nil }

// about global identity
func (mc MockDataCenter) InsertIdentity(identity *types.Identity) error { return nil }
func (mc MockDataCenter) RevokeIdentity(identity *types.Identity) error { return nil }
func (mc MockDataCenter) QueryIdentityById(identityId string) (*types.Identity, error) {
	return nil, nil
} // add by v0.5.0
func (mc MockDataCenter) QueryIdentityList(lastUpdate, pageSize uint64) (types.IdentityArray, error) {
	return nil, nil
}

//QueryIdentityListByIds(identityIds []string) (types.IdentityArray, error)
// about local identity
func (mc MockDataCenter) HasIdentity(identity *carriertypespb.Organization) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) StoreIdentity(identity *carriertypespb.Organization) error { return nil }
func (mc MockDataCenter) RemoveIdentity() error                                     { return nil }
func (mc MockDataCenter) QueryIdentityId() (string, error)                          { return "", nil }
func (mc MockDataCenter) QueryIdentity() (*carriertypespb.Organization, error)      { return nil, nil }

// v2.0
func (mc MockDataCenter) InsertMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	return nil
}
func (mc MockDataCenter) UpdateMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	return nil
}

//RevokeMetadataAuthority(metadataAuth *types.MetadataAuthority) error
func (mc MockDataCenter) QueryMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryMetadataAuthorityListByIdentityId(identityId string, lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	return nil, nil
}
func (mc MockDataCenter) UpdateIdentityCredential(identityId, credential string) error { return nil }

func (mc MockDataCenter) StoreLocalTask(task *types.Task) error             { return nil }
func (mc MockDataCenter) RemoveLocalTask(taskId string) error               { return nil }
func (mc MockDataCenter) HasLocalTask(taskId string) (bool, error)          { return false, nil }
func (mc MockDataCenter) QueryLocalTask(taskId string) (*types.Task, error) { return nil, nil }
func (mc MockDataCenter) QueryLocalTaskListByIds(taskIds []string) (types.TaskDataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalTaskList() (types.TaskDataArray, error) { return nil, nil }
func (mc MockDataCenter) QueryLocalTaskAndEvents(taskId string) (*types.Task, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalTaskAndEventsListByIds(taskIds []string) (types.TaskDataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalTaskAndEventsList() (types.TaskDataArray, error) { return nil, nil }

// about local task event
func (mc MockDataCenter) StoreTaskEvent(event *carriertypespb.TaskEvent) error { return nil }
func (mc MockDataCenter) QueryTaskEventList(taskId string) ([]*carriertypespb.TaskEvent, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskEventListByPartyId(taskId, partyId string) ([]*carriertypespb.TaskEvent, error) {
	return nil, nil
}
func (mc MockDataCenter) RemoveTaskEventList(taskId string) error                   { return nil }
func (mc MockDataCenter) RemoveTaskEventListByPartyId(taskId, partyId string) error { return nil }

// about global task on datacenter
func (mc MockDataCenter) InsertTask(task *types.Task) error { return nil }
func (mc MockDataCenter) QueryGlobalTaskList(lastUpdate, pageSize uint64) (types.TaskDataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskListByIdentityId(identityId string, lastUpdate, pageSize uint64) (types.TaskDataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskListByTaskIds(taskIds []string) (types.TaskDataArray, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskEventListByTaskId(taskId string) ([]*carriertypespb.TaskEvent, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskEventListByTaskIds(taskIds []string) ([]*carriertypespb.TaskEvent, error) {
	return nil, nil
}

// v 1.0 about TaskPowerUsed  (prefix + taskId + partyId -> {taskId, partId, jobNodeId, slotCount})
func (mc MockDataCenter) StoreLocalTaskPowerUsed(taskPowerUsed *types.LocalTaskPowerUsed) error {
	return nil
}
func (mc MockDataCenter) StoreLocalTaskPowerUseds(taskPowerUseds []*types.LocalTaskPowerUsed) error {
	return nil
}
func (mc MockDataCenter) HasLocalTaskPowerUsed(taskId, partyId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) RemoveLocalTaskPowerUsed(taskId, partyId string) error { return nil }
func (mc MockDataCenter) RemoveLocalTaskPowerUsedByTaskId(taskId string) error  { return nil }
func (mc MockDataCenter) QueryLocalTaskPowerUsed(taskId, partyId string) (*types.LocalTaskPowerUsed, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalTaskPowerUsedsByTaskId(taskId string) ([]*types.LocalTaskPowerUsed, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryLocalTaskPowerUseds() ([]*types.LocalTaskPowerUsed, error) {
	return nil, nil
}

// about JobNodeTaskPartyId (prefix + jobNodeId + taskId -> [partyId, ..., partyId])
func (mc MockDataCenter) StoreJobNodeTaskPartyId(jobNodeId, taskId, partyId string) error {
	return nil
}
func (mc MockDataCenter) RemoveJobNodeTaskPartyId(jobNodeId, taskId, partyId string) error {
	return nil
}
func (mc MockDataCenter) RemoveJobNodeTaskIdAllPartyIds(jobNodeId, taskId string) error { return nil }
func (mc MockDataCenter) QueryJobNodeRunningTaskIdList(jobNodeId string) ([]string, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryJobNodeRunningTaskCount(jobNodeId string) (uint32, error) {
	return 0, nil
}
func (mc MockDataCenter) QueryJobNodeRunningTaskIdsAndPartyIdsPairs(jobNodeId string) (map[string][]string, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryJobNodeRunningTaskAllPartyIdList(jobNodeId, taskId string) ([]string, error) {
	return nil, nil
}
func (mc MockDataCenter) HasJobNodeRunningTaskId(jobNodeId, taskId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) HasJobNodeTaskPartyId(jobNodeId, taskId, partyId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) QueryJobNodeTaskPartyIdCount(jobNodeId, taskId string) (uint32, error) {
	return 0, nil
}

// v 2.0 about jobNode history task count (prefix + jobNodeId -> history task count AND prefix + jobNodeId + taskId -> index)
func (mc MockDataCenter) StoreJobNodeHistoryTaskId(jobNodeId, taskId string) error { return nil }
func (mc MockDataCenter) HasJobNodeHistoryTaskId(jobNodeId, taskId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) QueryJobNodeHistoryTaskCount(jobNodeId string) (uint32, error) {
	return 0, nil
}

// v 2.0  about TaskResultData  (taskId -> {taskId, originId, metadataId})
func (mc MockDataCenter) StoreTaskUpResultData(turf *types.TaskUpResultData) error { return nil }
func (mc MockDataCenter) QueryTaskUpResulData(taskId string) (*types.TaskUpResultData, error) {
	return nil, nil
}
func (mc MockDataCenter) QueryTaskUpResultDataList() ([]*types.TaskUpResultData, error) {
	return nil, nil
}
func (mc MockDataCenter) RemoveTaskUpResultData(taskId string) error { return nil }

// v 2.0 about task partyIds of all partners (prefix + taskId -> [partyId, ..., partyId]  for task sender)
func (mc MockDataCenter) StoreTaskPartnerPartyIds(taskId string, partyIds []string) error {
	return nil
}
func (mc MockDataCenter) HasTaskPartnerPartyIds(taskId string) (bool, error)       { return false, nil }
func (mc MockDataCenter) QueryTaskPartnerPartyIds(taskId string) ([]string, error) { return nil, nil }
func (mc MockDataCenter) RemoveTaskPartnerPartyId(taskId, partyId string) error    { return nil }
func (mc MockDataCenter) RemoveTaskPartnerPartyIds(taskId string) error            { return nil }

// v 1.0 -> v 2.0 about task exec status (prefix + taskId + partyId -> "consensusing"|"running|terminating")
func (mc MockDataCenter) StoreLocalTaskExecuteStatusValConsensusByPartyId(taskId, partyId string) error {
	return nil
}
func (mc MockDataCenter) StoreLocalTaskExecuteStatusValRunningByPartyId(taskId, partyId string) error {
	return nil
}
func (mc MockDataCenter) StoreLocalTaskExecuteStatusValTerminateByPartyId(taskId, partyId string) error {
	return nil
} // add by v 0.3.0
func (mc MockDataCenter) RemoveLocalTaskExecuteStatusByPartyId(taskId, partyId string) error {
	return nil
}
func (mc MockDataCenter) HasLocalTaskExecuteStatusParty(taskId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) HasLocalTaskExecuteStatusByPartyId(taskId, partyId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) HasLocalTaskExecuteStatusConsensusByPartyId(taskId, partyId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) HasLocalTaskExecuteStatusRunningByPartyId(taskId, partyId string) (bool, error) {
	return false, nil
}
func (mc MockDataCenter) HasLocalTaskExecuteStatusTerminateByPartyId(taskId, partyId string) (bool, error) {
	return false, nil
} // add by v 0.3.0
// v 2.0 about NeedExecuteTask
func (mc MockDataCenter) StoreNeedExecuteTask(task *types.NeedExecuteTask) error      { return nil }
func (mc MockDataCenter) RemoveNeedExecuteTaskByPartyId(taskId, partyId string) error { return nil }
func (mc MockDataCenter) RemoveNeedExecuteTask(taskId string) error                   { return nil }
func (mc MockDataCenter) ForEachNeedExecuteTaskWithPrefix(prifix []byte, f func(key, value []byte) error) error {
	return nil
}
func (mc MockDataCenter) ForEachNeedExecuteTask(f func(key, value []byte) error) error { return nil }

// v 2.0 about taskbullet
func (mc MockDataCenter) StoreTaskBullet(bullet *types.TaskBullet) error           { return nil }
func (mc MockDataCenter) RemoveTaskBullet(taskId string) error                     { return nil }
func (mc MockDataCenter) ForEachTaskBullets(f func(key, value []byte) error) error { return nil }

func generateTestMetadata(metadataId string) (*types.Metadata, error) {
	testMetadata := make(map[string]*types.Metadata, 0)
	id := "MetadataId001"
	metadataOption := `{"originId": "originId001", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 2}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x67e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 1000000, \"plainAlgoConsumeUnit\": 7}]", "[\"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF11\", \"0x79e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
	testMetadata[id] = types.NewMetadata(&carriertypespb.MetadataPB{
		MetadataId:     id,
		MetadataOption: metadataOption,
	})
	id = "MetadataId002"
	metadataOption = `{"originId": "originId002", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"cryptoAlgoConsumeUnit\": 2000000, \"plainAlgoConsumeUnit\": 3}, {\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 2000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 2000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x77e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 2000000, \"plainAlgoConsumeUnit\": 12}]", "[\"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF11\", \"0x89e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
	testMetadata[id] = types.NewMetadata(&carriertypespb.MetadataPB{
		MetadataId:     id,
		MetadataOption: metadataOption,
	})
	id = "MetadataId003"
	metadataOption = `{"originId": "originId003", "dataPath": "dataPath001", "rows": 12, "columns": 7, "size": 56, "hasTitle": true, "metadataColumns": [], "consumeTypes": [1, 2, 3], "consumeOptions": ["[]", "[{\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF03\", \"cryptoAlgoConsumeUnit\": 3000000, \"plainAlgoConsumeUnit\": 5}, {\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF04\", \"cryptoAlgoConsumeUnit\": 3000000, \"plainAlgoConsumeUnit\": 4}, {\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF06\", \"cryptoAlgoConsumeUnit\": 3000000, \"plainAlgoConsumeUnit\": 11}, {\"contract\": \"0x87e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"cryptoAlgoConsumeUnit\": 3000000, \"plainAlgoConsumeUnit\": 7}]", "[\"0x97e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\", \"0x97e4b947F015f3f7C06E5173C2CfF41F2DDBAF11\", \"0x99e4b947F015f3f7C06E5173C2CfF41F2DDBAF07\"]"]}`
	testMetadata[id] = types.NewMetadata(&carriertypespb.MetadataPB{
		MetadataId:     id,
		MetadataOption: metadataOption,
	})
	if result, ok := testMetadata[metadataId]; !ok {
		return nil, errors.New(fmt.Sprintf("metadataId %s not exits", metadataId))
	} else {
		return result, nil
	}
}
