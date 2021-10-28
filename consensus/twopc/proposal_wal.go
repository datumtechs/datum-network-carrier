package twopc

import (
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/fileutil"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/db"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/gogo/protobuf/proto"
	"os"
	"path/filepath"
)

var (
	proposalTaskCachePrefix     = []byte("proposalTaskCache:")  	//	taskId -> partyId -> proposalTask
	proposalSetPrefix           = []byte("proposalSet:")			// 	proposalId -> partyId -> orgState
	prepareVotesPrefix          = []byte("prepareVotes:")			//  proposalId -> partyId -> prepareVote
	confirmVotesPrefix          = []byte("confirmVotes:")			//  proposalId -> partyId -> confirmVote
	proposalPeerInfoCachePrefix = []byte("proposalPeerInfoCache:")	//  proposalId -> ConfirmTaskPeerInfo
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
			log.Errorf("Failed to load `--consensus-state-file` on Start twopc, file: {%s}, err: {%s}", configFile, err)
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

func (w *walDB) GetProposalSetKey(hash common.Hash, partyId string) []byte {
	return append(append(proposalSetPrefix, hash.Bytes()...), []byte(partyId)...)
}

func (w *walDB) GetPrepareVotesKey(hash common.Hash, partyId string) []byte {
	return append(append(prepareVotesPrefix, hash.Bytes()...), []byte(partyId)...)
}

func (w *walDB) GetConfirmVotesKey(hash common.Hash, partyId string) []byte {
	return append(append(confirmVotesPrefix, hash.Bytes()...), []byte(partyId)...)
}

func (w *walDB) GetProposalPeerInfoCacheKey(hash common.Hash) []byte {
	return append(proposalPeerInfoCachePrefix, hash.Bytes()...)
}

func (w *walDB) StoreProposalTask(partyId string, task *types.ProposalTask) {
	data, err := proto.Marshal(&libtypes.ProposalTask{
		ProposalId: task.GetProposalId().String(),
		TaskId: task.GetTaskId(),
		CreateAt: task.GetCreateAt(),
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

func (w *walDB) StoreOrgProposalState(proposalId common.Hash, sender *apicommonpb.TaskOrganization, orgState *ctypes.OrgProposalState) {
	data, err := proto.Marshal(&libtypes.OrgProposalState{
		TaskId:             orgState.GetTaskId(),
		TaskSender:         sender,
		PrePeriodStartTime: orgState.PrePeriodStartTime,
		PeriodStartTime:    orgState.PeriodStartTime,
		DeadlineDuration:   orgState.DeadlineDuration,
		CreateAt:           orgState.CreateAt,
		TaskRole:           orgState.TaskRole,
		TaskOrg:            orgState.TaskOrg,
		PeriodNum:          uint32(orgState.PeriodNum),
	})
	if err != nil {
		log.WithError(err).Fatalf("marshal org proposalState failed, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			 proposalId.String(), orgState.GetTaskId(), orgState.TaskOrg.GetPartyId())
	}
	if err := w.db.Put(w.GetProposalSetKey(proposalId, orgState.TaskOrg.PartyId), data); err != nil {
		log.WithError(err).Fatalf("store org proposalState failed, proposalId: {%s}, taskId: {%s}, partyId: {%s}",
			 proposalId.String(), orgState.GetTaskId(), orgState.TaskOrg.GetPartyId())
	}
}

func (w *walDB) StorePrepareVote(vote *types.PrepareVote) {
	data, err := proto.Marshal(&libtypes.PrepareVote{
		MsgOption: &libtypes.MsgOption{
			ProposalId:      vote.MsgOption.ProposalId.String(),
			SenderRole:      vote.MsgOption.SenderRole,
			SenderPartyId:   vote.MsgOption.SenderPartyId,
			ReceiverRole:    vote.MsgOption.ReceiverRole,
			ReceiverPartyId: vote.MsgOption.ReceiverPartyId,
			Owner:           vote.MsgOption.Owner,
		},
		VoteOption: uint32(vote.VoteOption),
		PeerInfo: &libtypes.PrepareVoteResource{
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
	data, err := proto.Marshal(&libtypes.ConfirmVote{
		MsgOption: &libtypes.MsgOption{
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

func (w *walDB) StoreConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *twopcpb.ConfirmTaskPeerInfo) {
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
	has, err := w.db.Has(key)
	switch {
	case rawdb.IsNoDBNotFoundErr(err):
		return err
	case rawdb.IsDBNotFoundErr(err), nil == err && !has:
		return nil
	}
	return w.db.Delete(key)
}

func (w *walDB)  ForEachKV (f func(key, value []byte) error) error {
	it := w.db.NewIterator()
	defer it.Release()
	for it.Next() {
		if err := f(it.Key(), it.Value()); nil != err {
			return err
		}
	}
	return nil
}

func (w *walDB)  ForEachKVWithPrefix (prefix []byte, f func(key, value []byte) error) error {
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
	iter := w.db.NewIteratorWithPrefixAndStart(proposalPeerInfoCachePrefix, nil)
	proposalPeerInfoCache := make(map[common.Hash]*twopcpb.ConfirmTaskPeerInfo, 0)
	prefixLength := len(proposalPeerInfoCachePrefix)
	libProposalPeerInfoCache := &twopcpb.ConfirmTaskPeerInfo{}
	for iter.Next() {
		key := iter.Key()
		proposalId := common.BytesToHash(key[prefixLength:])
		if err := proto.Unmarshal(iter.Value(), libProposalPeerInfoCache); err != nil {
			log.Fatal("marshaling error: ", err)
		}
		proposalPeerInfoCache[proposalId] = libProposalPeerInfoCache
	}
}
