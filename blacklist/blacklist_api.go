package blacklist

import (
	"github.com/Metisnetwork/Metis-Carrier/common"
	libtypes "github.com/Metisnetwork/Metis-Carrier/lib/types"
	"github.com/Metisnetwork/Metis-Carrier/types"
	"sync"
)

const thresholdCount = 10

type BackListEngineAPI interface {
	// HasPrepareVoting HasConfirmVoting consensus 获取 state 的 方法 (这里的方法是给  blackList 用的)
	HasPrepareVoting(proposalId common.Hash, org *libtypes.TaskOrganization) bool
	HasConfirmVoting(proposalId common.Hash, org *libtypes.TaskOrganization) bool
}

type organizationTaskInfo struct {
	taskId     string
	nodeId     string
	proposalId string
	partyId    string
}

type IdentityBackListCache struct {
	engine            BackListEngineAPI
	orgBlacklistCache map[string][]*organizationTaskInfo
	orgBlacklistLock  sync.RWMutex
}

func NewIdentityBackListCache() *IdentityBackListCache {
	return &IdentityBackListCache{}
}

func (iBlc *IdentityBackListCache) SetEngine(engine BackListEngineAPI) {
	iBlc.engine = engine
}

// CheckConsensusResultOfNoVote 这些方法是  Black List 自有的 (consensus  和 p2p 直接用)
func (iBlc *IdentityBackListCache) CheckConsensusResultOfNoVote(proposalId common.Hash, task *types.Task) {
	iBlc.orgBlacklistLock.RLock()
	defer iBlc.orgBlacklistLock.RUnlock()

	mergeTaskOrgByIdentityId := make(map[string][]*libtypes.TaskOrganization, 0)
	mergeTaskByOrg := func(org *libtypes.TaskOrganization) {
		taskOrgInfo, ok := mergeTaskOrgByIdentityId[org.GetIdentityId()]
		if !ok {
			taskOrgInfo = make([]*libtypes.TaskOrganization, 0)
		}
		taskOrgInfo = append(taskOrgInfo, org)
		mergeTaskOrgByIdentityId[org.GetIdentityId()] = taskOrgInfo
	}
	for _, org := range task.GetTaskData().GetDataSuppliers() {
		mergeTaskByOrg(org)
	}
	for _, org := range task.GetTaskData().GetPowerSuppliers() {
		mergeTaskByOrg(org)
	}
	for _, org := range task.GetTaskData().GetReceivers() {
		mergeTaskByOrg(org)
	}

	taskId := task.GetTaskId()
	for identityId, taskOrgArr := range mergeTaskOrgByIdentityId {
		tempCount := 0
		for _, taskOrg := range taskOrgArr {
			orgBlacklistCache, ok := iBlc.orgBlacklistCache[identityId]
			if !ok {
				orgBlacklistCache = make([]*organizationTaskInfo, 0)
			}

			if !(iBlc.engine.HasPrepareVoting(proposalId, taskOrg) && iBlc.engine.HasConfirmVoting(proposalId, taskOrg)) {
				tempCount += 1
			} else {
				if len(orgBlacklistCache) < thresholdCount && len(orgBlacklistCache) > 0 {
					delete(iBlc.orgBlacklistCache, identityId)
				}
			}
			if len(taskOrgArr) == tempCount && tempCount != 0 {
				if len(orgBlacklistCache) < thresholdCount {
					orgBlacklistCache = append(orgBlacklistCache, &organizationTaskInfo{
						taskId:     taskId,
						nodeId:     taskOrg.GetNodeId(),
						proposalId: proposalId.String(),
						partyId:    taskOrg.GetPartyId(),
					})
					iBlc.orgBlacklistCache[identityId] = orgBlacklistCache
				}
			}
		}
	}
}

func (iBlc *IdentityBackListCache) QueryBlackListByIdentity(identityId string) int {
	result, ok := iBlc.orgBlacklistCache[identityId]
	if !ok {
		return 0
	}
	return len(result)
}

func (iBlc *IdentityBackListCache) RemoveBlackOrgByIdentity(identityId string) {
	iBlc.orgBlacklistLock.RLock()
	defer iBlc.orgBlacklistLock.RUnlock()
	delete(iBlc.orgBlacklistCache, identityId)
}

func (iBlc *IdentityBackListCache) FighterEqualThresholdCountOrg() []string {
	blackOrg := make([]string, 0)
	for identityId, value := range iBlc.orgBlacklistCache {
		if len(value) == thresholdCount {
			blackOrg = append(blackOrg, identityId)
		}
	}
	return blackOrg
}
