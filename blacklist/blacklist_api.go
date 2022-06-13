package blacklist

import (
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/common"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/sirupsen/logrus"
	"sort"
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
	DeleteBlackOrg(key []byte) error
	GetOrgBlacklistCachePrefix() []byte
}

type OrganizationTaskInfo struct {
	TaskId     string
	NodeId     string
	ProposalId string
}

type IdentityBackListCache struct {
	engine BackListEngineAPI
	db     WalDB
	// identityId -> [{taskId1, proposalId1}, {taskId2, proposalId2}, ..., {taskIdN, proposalIdN}]
	// OR identityId -> [{taskId1, proposalId1}, {taskId1, proposalId2}, ..., {taskIdN, proposalIdN}]
	orgBlacklistCache map[string][]*OrganizationTaskInfo
	orgBlacklistLock  sync.RWMutex
}

func NewIdentityBackListCache() *IdentityBackListCache {
	return &IdentityBackListCache{}
}

func (iBlc *IdentityBackListCache) SetEngineAndWal(engine BackListEngineAPI, db WalDB) {
	iBlc.engine = engine
	iBlc.db = db
	iBlc.FindBlackOrgByWalPrefix()
}

func (iBlc *IdentityBackListCache) CheckConsensusResultOfNotExistVote(proposalId common.Hash, task *types.Task) {
	iBlc.orgBlacklistLock.RLock()
	defer iBlc.orgBlacklistLock.RUnlock()
	// [TaskOrganization1, TaskOrganization2, ..., TaskOrganizationN]
	mergeTaskOrg := append(append(task.GetTaskData().GetDataSuppliers(), task.GetTaskData().GetPowerSuppliers()...), task.GetTaskData().GetReceivers()...)
	// Sort by the identityId field of taskOrg
	sort.Slice(mergeTaskOrg, func(i, j int) bool {
		return mergeTaskOrg[i].IdentityId == mergeTaskOrg[j].IdentityId
	})
	// Check and judge each taskOrg, and add the eligible taskOrg to blackList or remove it from blackList
	var (
		identityId string
		// tempCount is used to mark whether all taskOrg in the same organization do not have votes
		tempCount int
		// How many taskOrg are used by sameIdentityIdTaskOrgCount to mark the same identityId
		sameIdentityIdTaskOrgCount int
	)
	jump := true
	taskId := task.GetTaskId()
	n := len(mergeTaskOrg)
	// TaskOrg with the same identityId in mergeTaskOrg are adjacent
	for index, org := range mergeTaskOrg {
		if identityId != org.IdentityId {
			identityId = org.IdentityId
			jump = false
			tempCount, sameIdentityIdTaskOrgCount = 0, 0
		}
		sameIdentityIdTaskOrgCount += 1
		if jump {
			continue
		}
		orgBlacklistCache, ok := iBlc.orgBlacklistCache[identityId]
		orgBlacklistCacheCount := len(orgBlacklistCache)
		if !ok {
			orgBlacklistCache = make([]*OrganizationTaskInfo, 0)
		}
		if !iBlc.HasVoting(proposalId, org) {
			tempCount += 1
		} else {
			jump = true
			if orgBlacklistCacheCount < thresholdCount && orgBlacklistCacheCount > 0 {
				delete(iBlc.orgBlacklistCache, identityId)
				iBlc.RemoveBlackOrgByIdentity(identityId)
			}
		}
		if (index+1 < n && mergeTaskOrg[index+1].IdentityId != identityId) || (index+1 == n) {
			if tempCount == sameIdentityIdTaskOrgCount && orgBlacklistCacheCount < thresholdCount {
				orgBlacklistCache = append(orgBlacklistCache, &OrganizationTaskInfo{
					TaskId:     taskId,
					NodeId:     org.GetNodeId(),
					ProposalId: proposalId.String(),
				})
				iBlc.orgBlacklistCache[identityId] = orgBlacklistCache
				iBlc.db.StoreBlackTaskOrg(identityId, orgBlacklistCache)
			}
		}
	}
}

func (iBlc *IdentityBackListCache) RemoveBlackOrgByIdentity(identityId string) {
	iBlc.orgBlacklistLock.RLock()
	defer iBlc.orgBlacklistLock.RUnlock()
	delete(iBlc.orgBlacklistCache, identityId)
	err := iBlc.db.DeleteBlackOrg(iBlc.db.GetOrgBlacklistCacheKey(identityId))
	if err != nil {
		log.WithError(err).Errorf("RemoveBlackOrgByIdentity fail,identityId is %s", identityId)
	}
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

func (iBlc *IdentityBackListCache) FindBlackOrgByWalPrefix() {
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
	iBlc.orgBlacklistCache = orgBlacklistCache
}

func (iBlc *IdentityBackListCache) GetAllBlackListInfo() map[string][]*OrganizationTaskInfo {
	return iBlc.orgBlacklistCache
}

// QueryBlackListCountByIdentity is reservation method,not called yet
func (iBlc *IdentityBackListCache) QueryBlackListCountByIdentity(identityId string) int {
	result, ok := iBlc.orgBlacklistCache[identityId]
	if !ok {
		return 0
	}
	return len(result)
}

func (iBlc *IdentityBackListCache) HasVoting(proposalId common.Hash, taskOrg *carriertypespb.TaskOrganization) bool {
	return iBlc.engine.HasPrepareVoting(proposalId, taskOrg) && iBlc.engine.HasConfirmVoting(proposalId, taskOrg)
}
