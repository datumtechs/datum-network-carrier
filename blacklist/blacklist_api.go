package blacklist

import (
	"encoding/json"
	"fmt"
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

	dataSuppliersIndex := len(task.GetTaskData().GetDataSuppliers())
	powerSuppliersIndex := dataSuppliersIndex + len(task.GetTaskData().GetPowerSuppliers())
	mergeTaskOrgsSize := powerSuppliersIndex + len(task.GetTaskData().GetReceivers())

	// [TaskOrganization1, TaskOrganization2, ..., TaskOrganizationN]
	mergeTaskOrgs := make([]*carriertypespb.TaskOrganization, mergeTaskOrgsSize)
	copy(mergeTaskOrgs[:dataSuppliersIndex], task.GetTaskData().GetDataSuppliers())
	copy(mergeTaskOrgs[dataSuppliersIndex:powerSuppliersIndex], task.GetTaskData().GetPowerSuppliers())
	copy(mergeTaskOrgs[powerSuppliersIndex:], task.GetTaskData().GetReceivers())

	// Sort by the identityId field of taskOrg
	sort.Slice(mergeTaskOrgs, func(i, j int) bool {
		return mergeTaskOrgs[i].GetIdentityId() == mergeTaskOrgs[j].GetIdentityId()
	})
	// Check and judge each taskOrg, and add the eligible taskOrg to blackList or remove it from blackList
	var (
		identityId string
		// identityHasNotVoteCount is used to mark whether all taskOrg in the same organization do not have votes
		identityHasNotVoteCount int
		// How many taskOrg are used by sameIdentityIdTaskOrgCount to mark the same identityId
		sameIdentityIdTaskOrgCount int
		skip                       = true
	)

	// TaskOrg with the same identityId in mergeTaskOrg are adjacent
	for index, org := range mergeTaskOrgs {
		if identityId != org.GetIdentityId() {
			identityId = org.GetIdentityId()
			skip = false
			identityHasNotVoteCount, sameIdentityIdTaskOrgCount = 0, 0
		}
		sameIdentityIdTaskOrgCount += 1
		if skip {
			continue
		}
		consensusProposalTicks, ok := iBlc.orgConsensusProposalTickInfosCache[identityId]
		if !ok {
			consensusProposalTicks = make([]*ConsensusProposalTickInfo, 0)
		}
		consensusProposalTicksCount := len(consensusProposalTicks)

		// #### NOTE ####
		// Check whether to remove a 'identityid' related information from the blacklist.
		// ##############
		//
		// Whether to remove from the blacklist depends
		// on whether the 'Taskorganization' of the current partyId of the current organization
		// passes the "proposed" ticket.
		if iBlc.hasNotVoting(proposalId, org) {
			identityHasNotVoteCount += 1
		} else {

			// As long as a 'Taskorganization' of the organization has voted,
			// it can skip the subsequent 'Taskorganization' check of the organization.
			skip = true

			// If the identity has voted (any partner in the local task),
			// it should be removed from the blacklist directly.
			//
			// NOTE:
			//
			// condition: consensusProposalTicksCount > 0,
			//			  the above conditions can only be met
			//	 		  when we first process the 'concensusproposalticks' of this organization.
			if consensusProposalTicksCount < thresholdCount && consensusProposalTicksCount > 0 {
				delete(iBlc.orgConsensusProposalTickInfosCache, identityId)
				iBlc.RemoveBlackOrgByIdentity(identityId)
			}
		}

		// #### NOTE ####
		// Check whether need to add a 'identityid' related information to the blacklist.
		// ##############
		//
		// If all the 'mergetaskorgs' or all the' taskorganizations' of a 'identityid' are processed,
		// it can determine whether to add them to the blacklist.
		//
		// (NOTE: If the '(index+1 = = mergetaskorgssize)' condition exists,
		//  the condition '(mergetaskorgs[index+1] GetIdentityId() != Identityid)` out of bounds)
		if (index+1 == mergeTaskOrgsSize) || (mergeTaskOrgs[index+1].GetIdentityId() != identityId) {

			// When all the 'partyids' of the 'identityid' in a single task do not vote,
			// and the 'Concensusproposaltick' of the 'identityid' participating in the [not voting],
			// that is, the count of proposalids is less than the blacklist threshold,
			// the count of 'concensusproposaltick' will continue to be added.
			//
			// (NOTE: When the number of 'consumusproposalstick' reaches the blacklist threshold,
			//        no more 'consumusproposalstick' will be added)
			if identityHasNotVoteCount == sameIdentityIdTaskOrgCount && consensusProposalTicksCount < thresholdCount {
				consensusProposalTicks = append(consensusProposalTicks, &ConsensusProposalTickInfo{
					TaskId:     task.GetTaskId(),
					NodeId:     org.GetNodeId(),
					ProposalId: proposalId.String(),
				})
				iBlc.orgConsensusProposalTickInfosCache[identityId] = consensusProposalTicks
				iBlc.db.StoreBlackTaskOrg(identityId, consensusProposalTicks)
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
		cache[ticks[0].NodeId] = identityId
	}
	return cache
}

// QueryConsensusProposalTickInfoCountByIdentity is reservation method,not called yet
func (iBlc *IdentityBackListCache) QueryConsensusProposalTickInfoCountByIdentity(identityId string) int {

	iBlc.orgConsensusProposalTickInfosCacheLock.RLock()
	defer iBlc.orgConsensusProposalTickInfosCacheLock.RUnlock()

	result, ok := iBlc.orgConsensusProposalTickInfosCache[identityId]
	if !ok {
		return 0
	}
	return len(result)
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
	return iBlc.engine.HasPrepareVoting(proposalId, taskOrg) || iBlc.engine.HasConfirmVoting(proposalId, taskOrg)
}

func (iBlc *IdentityBackListCache) hasNotVoting(proposalId common.Hash, taskOrg *carriertypespb.TaskOrganization) bool {
	return !iBlc.hasVoting(proposalId, taskOrg)
}
