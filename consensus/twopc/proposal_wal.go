package twopc

import (
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/blacklist"
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/fileutil"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	"github.com/datumtechs/datum-network-carrier/db"
	carriertwopcpb "github.com/datumtechs/datum-network-carrier/pb/carrier/netmsg/consensus/twopc"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/gogo/protobuf/proto"
	"os"
	"path/filepath"
)

var (
	proposalTaskCachePrefix     = []byte("proposalTaskCache:")     //	taskId -> partyId -> proposalTask
	proposalSetPrefix           = []byte("proposalSet:")           // 	proposalId -> partyId -> orgState
	prepareVotesPrefix          = []byte("prepareVotes:")          //  proposalId -> partyId -> prepareVote
	confirmVotesPrefix          = []byte("confirmVotes:")          //  proposalId -> partyId -> confirmVote
	proposalPeerInfoCachePrefix = []byte("proposalPeerInfoCache:") //  proposalId -> ConfirmTaskPeerInfo
	orgBlacklistCachePrefix     = []byte("orgBlacklistCache:")     // identityId  -> map[string][]*organizationTaskInfo
)

type jsonFile struct {
	SavePath string `json:"savePath"`
	Cache    int    `json:"cache"`
	Handles  int    `json:"handles"`
}

type walDB struct {
	db *db.LDBDatabase
}

func initLDB(conf *Config) (*db.LDBDatabase, error) {
	configFile := conf.ConsensusStateFile
	var (
		savePath string
		cache    int
		handles  int
	)
	_, err := os.Stat(configFile)
	if err != nil {
		savePath, cache, handles = filepath.Join(conf.DefaultConsensusWal, "consensuswal"), conf.DatabaseCache, conf.DatabaseHandles
	} else {
		var jsonfile jsonFile
		if err := fileutil.LoadJSON(configFile, &jsonfile); err != nil {
			log.WithError(err).Errorf("Failed to load `--consensus-state-file` on Start twopc, file: {%s}", configFile)
			return nil, err
		} else {
			savePath, cache, handles = jsonfile.SavePath, jsonfile.Cache, jsonfile.Handles
		}
	}
	return db.NewLDBDatabase(savePath, cache, handles)
}

func ldbObj(conf *Config) *db.LDBDatabase {
	ldb, err := initLDB(conf)
	if err != nil {
		return nil
	} else {
		return ldb
	}
}

func newWal(conf *Config) *walDB {
	return &walDB{
		db: ldbObj(conf),
	}
}

func (w *walDB) GetProposalTaskCacheKey(taskId, partyId string) []byte {
	return append(append(proposalTaskCachePrefix, []byte(taskId)...), []byte(partyId)...)
}

func (w *walDB) GetProposalSetKey(proposalId common.Hash, partyId string) []byte {
	return append(append(proposalSetPrefix, proposalId.Bytes()...), []byte(partyId)...)
}

func (w *walDB) GetPrepareVotesKey(proposalId common.Hash, partyId string) []byte {
	return append(append(prepareVotesPrefix, proposalId.Bytes()...), []byte(partyId)...)
}

func (w *walDB) GetConfirmVotesKey(proposalId common.Hash, partyId string) []byte {
	return append(append(confirmVotesPrefix, proposalId.Bytes()...), []byte(partyId)...)
}

func (w *walDB) GetProposalPeerInfoCacheKey(proposalId common.Hash) []byte {
	return append(proposalPeerInfoCachePrefix, proposalId.Bytes()...)
}

func (w *walDB) GetOrgBlacklistCacheKey(identityId string) []byte {
	return append(orgBlacklistCachePrefix, []byte(identityId)...)
}

func (w *walDB) StoreProposalTask(partyId string, task *ctypes.ProposalTask) {
	data, err := proto.Marshal(&carriertypespb.ProposalTask{
		ProposalId: task.GetProposalId().String(),
		TaskId:     task.GetTaskId(),
		CreateAt:   task.GetCreateAt(),
	})
	if err != nil {
		log.WithError(err).Fatalf("marshal proposalTask failed, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			task.GetProposalId().String(), task.GetTaskId(), partyId)
	}
	if err := w.db.Put(w.GetProposalTaskCacheKey(task.GetTaskId(), partyId), data); err != nil {
		log.WithError(err).Fatalf("store proposalTask failed, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			task.GetProposalId().String(), task.GetTaskId(), partyId)
	}
}

func (w *walDB) StoreOrgProposalState(orgState *ctypes.OrgProposalState) {
	data, err := proto.Marshal(&carriertypespb.OrgProposalState{
		TaskId:           orgState.GetTaskId(),
		TaskSender:       orgState.GetTaskSender(),
		StartAt:          orgState.GetStartAt(),
		DeadlineDuration: orgState.GetDeadlineDuration(),
		CreateAt:         orgState.GetCreateAt(),
		TaskRole:         orgState.GetTaskRole(),
		TaskOrg:          orgState.GetTaskOrg(),
		PeriodNum:        uint32(orgState.GetPeriodNum()),
	})
	if err != nil {
		log.WithError(err).Fatalf("marshal org proposalState failed, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())
	}
	if err := w.db.Put(w.GetProposalSetKey(orgState.GetProposalId(), orgState.GetTaskOrg().PartyId), data); err != nil {
		log.WithError(err).Fatalf("store org proposalState failed, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			orgState.GetProposalId().String(), orgState.GetTaskId(), orgState.GetTaskOrg().GetPartyId())
	}
}

func (w *walDB) StorePrepareVote(vote *types.PrepareVote) {
	data, err := proto.Marshal(&carriertypespb.PrepareVote{
		MsgOption: &carriertypespb.MsgOption{
			ProposalId:      vote.MsgOption.ProposalId.String(),
			SenderRole:      vote.MsgOption.SenderRole,
			SenderPartyId:   vote.MsgOption.SenderPartyId,
			ReceiverRole:    vote.MsgOption.ReceiverRole,
			ReceiverPartyId: vote.MsgOption.ReceiverPartyId,
			Owner:           vote.MsgOption.Owner,
		},
		VoteOption: uint32(vote.VoteOption),
		PeerInfo: &carriertypespb.PrepareVoteResource{
			Id:      vote.PeerInfo.Id,
			Ip:      vote.PeerInfo.Ip,
			Port:    vote.PeerInfo.Port,
			PartyId: vote.PeerInfo.PartyId,
		},
		CreateAt: vote.CreateAt,
		Sign:     vote.Sign,
	})
	if err != nil {
		log.WithError(err).Fatalf("marshal org prepareVote failed, proposalId: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.SenderPartyId)
	}
	if err := w.db.Put(w.GetPrepareVotesKey(vote.MsgOption.ProposalId, vote.MsgOption.SenderPartyId), data); err != nil {
		log.WithError(err).Fatalf("store org prepareVote failed, proposalId: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.SenderPartyId)
	}
}

func (w *walDB) StoreConfirmVote(vote *types.ConfirmVote) {
	data, err := proto.Marshal(&carriertypespb.ConfirmVote{
		MsgOption: &carriertypespb.MsgOption{
			ProposalId:      vote.MsgOption.ProposalId.String(),
			SenderRole:      vote.MsgOption.SenderRole,
			SenderPartyId:   vote.MsgOption.SenderPartyId,
			ReceiverRole:    vote.MsgOption.ReceiverRole,
			ReceiverPartyId: vote.MsgOption.ReceiverPartyId,
			Owner:           vote.MsgOption.Owner,
		},
		VoteOption: uint32(vote.VoteOption),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	})
	if err != nil {
		log.WithError(err).Fatalf("marshal org confirmVote failed, proposalId: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.SenderPartyId)
	}
	if err := w.db.Put(w.GetConfirmVotesKey(vote.MsgOption.ProposalId, vote.MsgOption.SenderPartyId), data); err != nil {
		log.WithError(err).Fatalf("store org confirmVote failed, proposalId: {%s}, partyId: {%s}",
			vote.MsgOption.ProposalId.String(), vote.MsgOption.SenderPartyId)
	}
}

func (w *walDB) StoreConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *carriertwopcpb.ConfirmTaskPeerInfo) {
	data, err := proto.Marshal(peerDesc)
	if err != nil {
		log.WithError(err).Fatalf("marshal confirmTaskPeerInfo failed, proposalId: {%s}",
			proposalId.String())
	}
	if err := w.db.Put(w.GetProposalPeerInfoCacheKey(proposalId), data); err != nil {
		log.WithError(err).Fatalf("store confirmTaskPeerInfo failed, proposalId: {%s}",
			proposalId.String())
	}
}

func (w *walDB) DeleteState(key []byte) error {
	log.Debugf("Delete state,key is:%s", string(key))
	return w.DeleteByKey(key)
}

func (w *walDB) RemoveConsensusProposalTicks(identityId string) error {
	return w.DeleteByKey(w.GetOrgBlacklistCacheKey(identityId))
}

func (w *walDB) DeleteByKey(key []byte) error {
	has, err := w.db.Has(key)
	switch {
	case rawdb.IsNoDBNotFoundErr(err):
		return err
	case rawdb.IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return w.db.Delete(key)
}

func (w *walDB) ForEachKV(f func(key, value []byte) error) error {
	it := w.db.NewIterator()
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

func (w *walDB) ForEachKVWithPrefix(prefix []byte, f func(key, value []byte) error) error {
	it := w.db.NewIteratorWithPrefix(prefix)
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

func (w *walDB) UnmarshalTest() {
	it := w.db.NewIteratorWithPrefixAndStart(proposalPeerInfoCachePrefix, nil)
	defer it.Release()
	proposalPeerInfoCache := make(map[common.Hash]*carriertwopcpb.ConfirmTaskPeerInfo, 0)
	prefixLength := len(proposalPeerInfoCachePrefix)
	libProposalPeerInfoCache := &carriertwopcpb.ConfirmTaskPeerInfo{}
	for it.Next() {
		key := it.Key()
		proposalId := common.BytesToHash(key[prefixLength:])
		if err := proto.Unmarshal(it.Value(), libProposalPeerInfoCache); err != nil {
			log.WithError(err).Fatalf("marshaling confirmTaskPeerInfo failed, proposalId: {%s}", proposalId.String())
		}
		proposalPeerInfoCache[proposalId] = libProposalPeerInfoCache
	}
}

func (w *walDB) StoreConsensusProposalTicks(identityId string, arr []*blacklist.ConsensusProposalTickInfo) {
	key := w.GetOrgBlacklistCacheKey(identityId)
	_, err := w.db.Get(key)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("can not query old `consensusProposalTicks` from db, identityId: {%s}", identityId)
		return
	}
	value, err := json.Marshal(arr)
	if nil != err {
		log.WithError(err).Errorf("can not json marshal new `consensusProposalTicks`, identityId: {%s}", identityId)
		return
	}
	if err := w.db.Put(key, value); err != nil {
		log.WithError(err).Errorf("can not store new `consensusProposalTicks` into db, identityId: {%s}", identityId)
	}
}

func (w *walDB) GetOrgBlacklistCachePrefix() []byte {
	return orgBlacklistCachePrefix
}
