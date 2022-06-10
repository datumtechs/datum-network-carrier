package blacklist

import (
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/common"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/sirupsen/logrus"
	"sync"
)

const thresholdCount = 10

// Global log object, used by the current package.
var log = logrus.WithField("prefix", "blacklist")

type BackListEngineAPI interface {
	HasPrepareVoting(proposalId common.Hash, org *carriertypespb.TaskOrganization) bool
	HasConfirmVoting(proposalId common.Hash, org *carriertypespb.TaskOrganization) bool
}

type WalDB interface {
	ForEachKVWithPrefix(prefix []byte, f func(key, value []byte) error) error
	GetOrgBlacklistCacheKey(identityId string) []byte
	StoreBlackTaskOrg(identityId string, info []*OrganizationTaskInfo)
	DeleteState(key []byte) error
	GetOrgBlacklistCachePrefix() []byte
}

type OrganizationTaskInfo struct {
	TaskId     string
	NodeId     string
	ProposalId string
}

type IdentityBackListCache struct {
	engine            BackListEngineAPI
	db                WalDB
	orgBlacklistCache map[string][]*OrganizationTaskInfo
	orgBlacklistLock  sync.RWMutex
}

func NewIdentityBackListCache() *IdentityBackListCache {
	return &IdentityBackListCache{}
}

func (iBlc *IdentityBackListCache) SetEngineAndWal(engine BackListEngineAPI, db WalDB) {
	iBlc.engine = engine
	iBlc.db = db
	iBlc.orgBlacklistCache = iBlc.FindBlackOrgByWalPrefix()
}

func (iBlc *IdentityBackListCache) CheckConsensusResultOfNoVote(proposalId common.Hash, task *types.Task) {
	iBlc.orgBlacklistLock.RLock()
	defer iBlc.orgBlacklistLock.RUnlock()

	// identityId -> [TaskOrganization1, TaskOrganization2, ..., TaskOrganizationN]
	mergeTaskOrgByIdentityId := make(map[string][]*carriertypespb.TaskOrganization, 0)
	mergeTaskByOrg := func(org *carriertypespb.TaskOrganization) {
		taskOrgInfo, ok := mergeTaskOrgByIdentityId[org.GetIdentityId()]
		if !ok {
			taskOrgInfo = make([]*carriertypespb.TaskOrganization, 0)
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
				orgBlacklistCache = make([]*OrganizationTaskInfo, 0)
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
					orgBlacklistCache = append(orgBlacklistCache, &OrganizationTaskInfo{
						TaskId:     taskId,
						NodeId:     taskOrg.GetNodeId(),
						ProposalId: proposalId.String(),
					})
					iBlc.orgBlacklistCache[identityId] = orgBlacklistCache
					iBlc.db.StoreBlackTaskOrg(identityId, orgBlacklistCache)
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
	iBlc.db.DeleteState(iBlc.db.GetOrgBlacklistCacheKey(identityId))
}

func (iBlc *IdentityBackListCache) FilterEqualThresholdCountOrg() map[string]struct{} {
	blackOrg := make(map[string]struct{}, 0)
	for identityId, value := range iBlc.orgBlacklistCache {
		if len(value) == thresholdCount {
			blackOrg[identityId] = struct{}{}
		}
	}
	return blackOrg
}

func (iBlc *IdentityBackListCache) FindBlackOrgByWalPrefix() map[string][]*OrganizationTaskInfo {
	prefix := iBlc.db.GetOrgBlacklistCachePrefix()
	prefixLength := len(prefix)

	orgBlacklistCache := make(map[string][]*OrganizationTaskInfo, 0)
	if err := iBlc.db.ForEachKVWithPrefix(prefix, func(key, value []byte) error {
		identityId := string(key[prefixLength:])
		orgBlacklist := make([]*OrganizationTaskInfo, 0)
		err := json.Unmarshal(value, &orgBlacklist)
		orgBlacklistCache[identityId] = orgBlacklist
		return err
	}); err != nil {
		log.WithError(err).Errorf("FindBlackOrgByWalPrefix ->ForEachKVWithPrefix fail")
	}
	return orgBlacklistCache
}

func (iBlc *IdentityBackListCache) HasVoting(proposalId common.Hash, taskOrg *carriertypespb.TaskOrganization) bool {
	return iBlc.engine.HasPrepareVoting(proposalId, taskOrg) && iBlc.engine.HasConfirmVoting(proposalId, taskOrg)
}
