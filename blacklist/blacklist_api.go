package blacklist

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	carrierrpcdebugpbv1 "github.com/datumtechs/datum-network-carrier/pb/carrier/rpc/debug/v1"
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
	StoreBlackTaskOrg(identityId string, info []*ConsensusProposalTickInfo)
	RemoveBlackTaskOrg(identityId string) error
	GetOrgBlacklistCachePrefix() []byte
}

type ConsensusProposalTickInfo struct {
	TaskId     string
	NodeId     string
	ProposalId string
}

type IdentityBackListCache struct {
	engine BackListEngineAPI
	db     WalDB
	// identityId -> [{taskId1, proposalId1}, {taskId2, proposalId2}, ..., {taskIdN, proposalIdN}]
	// OR identityId -> [{taskId1, proposalId1}, {taskId1, proposalId2}, ..., {taskIdN, proposalIdN}]
	orgConsensusProposalTickInfosCache     map[string][]*ConsensusProposalTickInfo
	orgConsensusProposalTickInfosCacheLock sync.RWMutex
}

func NewIdentityBackListCache() *IdentityBackListCache {
	return &IdentityBackListCache{
		orgConsensusProposalTickInfosCache: make(map[string][]*ConsensusProposalTickInfo, 0),
	}
}

func (iBlc *IdentityBackListCache) SetEngineAndWal(engine BackListEngineAPI, db WalDB) {
	iBlc.engine = engine
	iBlc.db = db
	iBlc.recoveryBlackOrg()
}

func (iBlc *IdentityBackListCache) CheckConsensusResultOfNotExistVote(proposalId common.Hash, task *types.Task) {

	iBlc.orgConsensusProposalTickInfosCacheLock.Lock()
	defer iBlc.orgConsensusProposalTickInfosCacheLock.Unlock()

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
		orgBlacklistCache, ok := iBlc.orgConsensusProposalTickInfosCache[identityId]
		orgBlacklistCacheCount := len(orgBlacklistCache)
		if !ok {
			orgBlacklistCache = make([]*ConsensusProposalTickInfo, 0)
		}
		if !iBlc.hasVoting(proposalId, org) {
			tempCount += 1
		} else {
			jump = true
			if orgBlacklistCacheCount < thresholdCount && orgBlacklistCacheCount > 0 {
				delete(iBlc.orgConsensusProposalTickInfosCache, identityId)
				iBlc.RemoveBlackOrgByIdentity(identityId)
			}
		}
		if (index+1 == n) || (mergeTaskOrg[index+1].IdentityId != identityId) {
			if tempCount == sameIdentityIdTaskOrgCount && orgBlacklistCacheCount < thresholdCount {
				orgBlacklistCache = append(orgBlacklistCache, &ConsensusProposalTickInfo{
					TaskId:     taskId,
					NodeId:     org.GetNodeId(),
					ProposalId: proposalId.String(),
				})
				iBlc.orgConsensusProposalTickInfosCache[identityId] = orgBlacklistCache
				iBlc.db.StoreBlackTaskOrg(identityId, orgBlacklistCache)
			}
		}
	}
}

func (iBlc *IdentityBackListCache) RemoveBlackOrgByIdentity(identityId string) {

	iBlc.orgConsensusProposalTickInfosCacheLock.Lock()
	delete(iBlc.orgConsensusProposalTickInfosCache, identityId)
	iBlc.orgConsensusProposalTickInfosCacheLock.Unlock()

	if err := iBlc.db.RemoveBlackTaskOrg(identityId); nil != err {
		log.WithError(err).Errorf("RemoveBlackOrgByIdentity failed, identityId: %s", identityId)
	}
}

func (iBlc *IdentityBackListCache) QueryBlackListIdentityIds() []string {

	blackListOrgArr := make([]string, 0)

	iBlc.orgConsensusProposalTickInfosCacheLock.RLock()
	defer iBlc.orgConsensusProposalTickInfosCacheLock.RUnlock()

	for identityId, ticks := range iBlc.orgConsensusProposalTickInfosCache {
		if len(ticks) == thresholdCount {
			blackListOrgArr = append(blackListOrgArr, identityId)
		}
	}
	return blackListOrgArr
}

func (iBlc *IdentityBackListCache) GetBlackListOrgSymbolCache() map[string]string {

	cache := make(map[string]string, 0)

	iBlc.orgConsensusProposalTickInfosCacheLock.RLock()
	defer iBlc.orgConsensusProposalTickInfosCacheLock.RUnlock()

	for identityId, ticks := range iBlc.orgConsensusProposalTickInfosCache {
		if len(ticks) == thresholdCount {
			cache[ticks[0].NodeId] = identityId
		}
	}
	return cache
}

// QueryConsensusProposalTickInfoCountByIdentity is reservation method,not called yet
func (iBlc *IdentityBackListCache) QueryConsensusProposalTickInfoCountByIdentity(identityId string) int {
	result, ok := iBlc.orgConsensusProposalTickInfosCache[identityId]
	if !ok {
		return 0
	}
	return len(result)
}

func (iBlc *IdentityBackListCache) GetAllBlackOrg() (*carrierrpcdebugpbv1.GetConsensusBlackOrgResponse, error) {
	result := make([]*carrierrpcdebugpbv1.GetConsensusBlackOrgResponse_ConsensusProposalList, 0)
	for identityId, taskOrgArr := range iBlc.orgConsensusProposalTickInfosCache {
		if len(taskOrgArr) == thresholdCount {
			savePbOrgArr := make([]*carrierrpcdebugpbv1.ConsensusProposalTickInfo, 0)
			for _, org := range taskOrgArr {
				savePbOrgArr = append(savePbOrgArr, &carrierrpcdebugpbv1.ConsensusProposalTickInfo{
					TaskId:     org.TaskId,
					NodeId:     org.NodeId,
					ProposalId: org.ProposalId,
				})
			}
			result = append(result, &carrierrpcdebugpbv1.GetConsensusBlackOrgResponse_ConsensusProposalList{
				IdentityId:       identityId,
				ProposalInfoList: savePbOrgArr,
			})
		}
	}
	return &carrierrpcdebugpbv1.GetConsensusBlackOrgResponse{
		AllBlackOrg: result,
	}, nil
}

// internal methods ...

func (iBlc *IdentityBackListCache) recoveryBlackOrg() {

	prefix := iBlc.db.GetOrgBlacklistCachePrefix()
	prefixLength := len(prefix)

	if err := iBlc.db.ForEachKVWithPrefix(prefix, func(key, value []byte) error {
		identityId := string(key[prefixLength:])
		var proposalTicks []*ConsensusProposalTickInfo
		if err := json.Unmarshal(value, &proposalTicks); nil != err {
			return fmt.Errorf("cannot json unmarshal proposalTicks of blacklist, identityId: %s", identityId)
		}
		iBlc.orgConsensusProposalTickInfosCache[identityId] = proposalTicks
		return nil
	}); err != nil {
		log.WithError(err).Warnf("recoveryBlackOrg failed")
	}
}

func (iBlc *IdentityBackListCache) hasVoting(proposalId common.Hash, taskOrg *carriertypespb.TaskOrganization) bool {
	return iBlc.engine.HasPrepareVoting(proposalId, taskOrg) && iBlc.engine.HasConfirmVoting(proposalId, taskOrg)
}