package twopc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/fileutil"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/db"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/gogo/protobuf/proto"
	"os"
)

var (
	proposalSetPrefix           = []byte("proposalSet:")
	prepareVotesPrefix          = []byte("prepareVotes:")
	confirmVotesPrefix          = []byte("confirmVotes:")
	proposalPeerInfoCachePrefix = []byte("proposalPeerInfoCache:")
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
	if err != nil  {
		savePath, cache, handles = "savePathState", 16, 16
	} else {
		var jsonfile jsonFile
		if err := fileutil.LoadJSON(configFile, &jsonfile); err != nil {
			log.Errorf("Failed to load `--consensus-state-file` on Start twoPC, file: {%s}, err: {%s}", configFile, err)
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

func (w *walDB) GetProposalSetKey(hash common.Hash, partyId string) []byte {
	result := append(proposalSetPrefix, hash.Bytes()...)
	return append(result, []byte(partyId)...)
}

func (w *walDB) GetProposalPeerInfoCacheKey(hash common.Hash) []byte {
	return append(proposalPeerInfoCachePrefix, hash.Bytes()...)
}

func (w *walDB) GetPrepareVotesKey(hash common.Hash, partyId string) []byte {
	result := append(prepareVotesPrefix, hash.Bytes()...)
	return append(result, []byte(partyId)...)
}

func (w *walDB) GetConfirmVotesKey(hash common.Hash, partyId string) []byte {
	result := append(confirmVotesPrefix, hash.Bytes()...)
	return append(result, []byte(partyId)...)
}

func (w *walDB) UpdateOrgProposalState(proposalId common.Hash, sender *apicommonpb.TaskOrganization, orgState *ctypes.OrgProposalState) {
	pbObj := &libtypes.OrgProposalState{
		TaskId:             orgState.TaskId,
		TaskSender:         sender,
		PrePeriodStartTime: orgState.PrePeriodStartTime,
		PeriodStartTime:    orgState.PeriodStartTime,
		DeadlineDuration:   orgState.DeadlineDuration,
		CreateAt:           orgState.CreateAt,
		TaskRole:           orgState.TaskRole,
		TaskOrg:            orgState.TaskOrg,
		PeriodNum:          uint32(orgState.PeriodNum),
	}
	data, err := proto.Marshal(pbObj)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	if err := w.db.Put(w.GetProposalSetKey(proposalId, orgState.TaskOrg.PartyId), data); err != nil {
		log.Warning("UpdateOrgProposalState to db fail,proposalId:", proposalId)
	}
}

func (w *walDB) UpdateConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *twopcpb.ConfirmTaskPeerInfo) {
	data, err := proto.Marshal(peerDesc)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	fmt.Println("proposalId:", proposalId)
	if err := w.db.Put(w.GetProposalPeerInfoCacheKey(proposalId), data); err != nil {
		log.Warning("UpdateConfirmTaskPeerInfo to db fail,proposalId:", proposalId)
	}
}

func (w *walDB) UpdatePrepareVotes(vote *types.PrepareVote) {
	pbObj := &libtypes.PrepareVoteState{
		MsgOption: &libtypes.MsgOption{
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
	}
	data, err := proto.Marshal(pbObj)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	if err := w.db.Put(w.GetPrepareVotesKey(vote.MsgOption.ProposalId, vote.MsgOption.SenderPartyId), data); err != nil {
		log.Warning("UpdatePrepareVotes to db fail,proposalId:", vote.MsgOption.ProposalId)
	}
}

func (w *walDB) UpdateConfirmVotes(vote *types.ConfirmVote) {
	pbObj := &libtypes.ConfirmVoteState{
		MsgOption: &libtypes.MsgOption{
			SenderRole:      vote.MsgOption.SenderRole,
			SenderPartyId:   vote.MsgOption.SenderPartyId,
			ReceiverRole:    vote.MsgOption.ReceiverRole,
			ReceiverPartyId: vote.MsgOption.ReceiverPartyId,
			Owner:           vote.MsgOption.Owner,
		},
		VoteOption: uint32(vote.VoteOption),
		CreateAt:   vote.CreateAt,
		Sign:       vote.Sign,
	}
	data, err := proto.Marshal(pbObj)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	if err := w.db.Put(w.GetConfirmVotesKey(vote.MsgOption.ProposalId, vote.MsgOption.SenderPartyId), data); err != nil {
		log.Warning("UpdateConfirmVotes to db fail,proposalId:", vote.MsgOption.ProposalId)
	}
}

func (w *walDB) DeleteState(key []byte) {
	err := errors.New("")
	prefix := append(bytes.Split(key, []byte(":"))[0], []byte(":")...)
	if bytes.Equal(prefix, proposalPeerInfoCachePrefix) {
		err = w.db.Delete(key)
	} else {
		partyId := string(key[len(prefix)+32:])
		if partyId != "" {
			err = w.db.Delete(key)
		} else {
			iter := w.db.NewIteratorWithPrefixAndStart(prefix, nil)
			iter.Release()
			for iter.Next() {
				err = w.db.Delete(iter.Key())
			}
		}
	}
	if err != nil {
		log.WithError(err).Fatal("Failed to delete from leveldb, key is ", key)
	}
}

func (w *walDB) RecoveryState() *state {
	// region recovery proposalSet
	iter := w.db.NewIteratorWithPrefixAndStart(proposalSetPrefix, nil)
	iter.Release()
	prefixLength := len(proposalSetPrefix)
	proposalSet := make(map[common.Hash]*ctypes.ProposalState, 0)
	stateCache := make(map[string]*ctypes.OrgProposalState, 0)
	for iter.Next() {
		key := iter.Key()
		proposalId, partyId := common.BytesToHash(key[prefixLength:prefixLength+32]), string(key[prefixLength+32:])

		value := iter.Value()
		libOrgProposalState := &libtypes.OrgProposalState{}
		if err := proto.Unmarshal(value, libOrgProposalState); err != nil {
			log.Fatal("marshaling error: ", err)
		}
		if _, ok := proposalSet[proposalId]; !ok {
			stateCache = make(map[string]*ctypes.OrgProposalState, 0)
			proposalSet[proposalId] = ctypes.RecoveryProposalState(proposalId, libOrgProposalState.TaskId, libOrgProposalState.TaskSender, stateCache)
		}
		stateCache[partyId] = &ctypes.OrgProposalState{
			PrePeriodStartTime: libOrgProposalState.PrePeriodStartTime,
			PeriodStartTime:    libOrgProposalState.PeriodStartTime,
			DeadlineDuration:   libOrgProposalState.DeadlineDuration,
			CreateAt:           libOrgProposalState.CreateAt,
			TaskId:             libOrgProposalState.TaskId,
			TaskRole:           libOrgProposalState.TaskRole,
			TaskOrg:            libOrgProposalState.TaskOrg,
			PeriodNum:          ctypes.ProposalStatePeriod(libOrgProposalState.PeriodNum),
		}
	}
	//endregion

	//region recovery prepareVotes
	iter = w.db.NewIteratorWithPrefixAndStart(prepareVotesPrefix, nil)
	prefixLength = len(prepareVotesPrefix)
	prepareVotes := make(map[common.Hash]*prepareVoteState, 0)
	votes := make(map[string]*types.PrepareVote, 0)
	yesVotes := make(map[apicommonpb.TaskRole]uint32, 0)
	voteStatus := make(map[apicommonpb.TaskRole]uint32, 0)
	for iter.Next() {
		key := iter.Key()
		proposalId, partyId := common.BytesToHash(key[prefixLength:prefixLength+33]), string(key[prefixLength+34:])

		value := iter.Value()
		libPrepareVoteState := &libtypes.PrepareVoteState{}
		if err := proto.Unmarshal(value, libPrepareVoteState); err != nil {
			log.Fatal("marshaling error: ", err)
		}
		if _, ok := prepareVotes[proposalId]; !ok {
			votes = make(map[string]*types.PrepareVote, 0)
			yesVotes = make(map[apicommonpb.TaskRole]uint32, 0)
			voteStatus = make(map[apicommonpb.TaskRole]uint32, 0)
			prepareVotes[proposalId] = &prepareVoteState{
				votes:      votes,
				yesVotes:   yesVotes,
				voteStatus: voteStatus,
			}
		}
		votes[partyId] = &types.PrepareVote{
			MsgOption: &types.MsgOption{
				ProposalId:      proposalId,
				SenderRole:      libPrepareVoteState.MsgOption.SenderRole,
				SenderPartyId:   partyId,
				ReceiverRole:    libPrepareVoteState.MsgOption.ReceiverRole,
				ReceiverPartyId: libPrepareVoteState.MsgOption.ReceiverPartyId,
				Owner:           libPrepareVoteState.MsgOption.Owner,
			},
			VoteOption: types.VoteOption(libPrepareVoteState.VoteOption),
			PeerInfo: &types.PrepareVoteResource{
				Id:      libPrepareVoteState.PeerInfo.Id,
				Ip:      libPrepareVoteState.PeerInfo.Ip,
				Port:    libPrepareVoteState.PeerInfo.Port,
				PartyId: libPrepareVoteState.PeerInfo.PartyId,
			},
			CreateAt: libPrepareVoteState.CreateAt,
			Sign:     libPrepareVoteState.Sign,
		}
		if types.VoteOption(libPrepareVoteState.VoteOption) == types.Yes {
			if _, ok := yesVotes[libPrepareVoteState.MsgOption.SenderRole]; !ok {
				yesVotes[libPrepareVoteState.MsgOption.SenderRole] = 1
			} else {
				yesVotes[libPrepareVoteState.MsgOption.SenderRole] += 1
			}
		}

		if _, ok := voteStatus[libPrepareVoteState.MsgOption.SenderRole]; !ok {
			yesVotes[libPrepareVoteState.MsgOption.SenderRole] = 1
		} else {
			yesVotes[libPrepareVoteState.MsgOption.SenderRole] += 1
		}
	}
	//endregion

	// region  recovery confirmVotes
	iter = w.db.NewIteratorWithPrefixAndStart(confirmVotesPrefix, nil)
	prefixLength = len(confirmVotesPrefix)
	confirmVotes := make(map[common.Hash]*confirmVoteState, 0)
	votesConfirm := make(map[string]*types.ConfirmVote, 0)
	yesVotes = make(map[apicommonpb.TaskRole]uint32, 0)
	voteStatus = make(map[apicommonpb.TaskRole]uint32, 0)
	for iter.Next() {
		key := iter.Key()
		proposalId, partyId := common.BytesToHash(key[prefixLength:prefixLength+32]), string(key[prefixLength+32:])
		value := iter.Value()

		libConfirmVoteState := &libtypes.ConfirmVoteState{}
		if err := proto.Unmarshal(value, libConfirmVoteState); err != nil {
			log.Fatal("marshaling error: ", err)
		}
		if _, ok := confirmVotes[proposalId]; !ok {
			votesConfirm = make(map[string]*types.ConfirmVote, 0)
			yesVotes = make(map[apicommonpb.TaskRole]uint32, 0)
			voteStatus = make(map[apicommonpb.TaskRole]uint32, 0)
			confirmVotes[proposalId] = &confirmVoteState{
				votes:      votesConfirm,
				yesVotes:   yesVotes,
				voteStatus: voteStatus,
			}
		}
		votes[partyId] = &types.PrepareVote{
			MsgOption: &types.MsgOption{
				ProposalId:      proposalId,
				SenderRole:      libConfirmVoteState.MsgOption.SenderRole,
				SenderPartyId:   partyId,
				ReceiverRole:    libConfirmVoteState.MsgOption.ReceiverRole,
				ReceiverPartyId: libConfirmVoteState.MsgOption.ReceiverPartyId,
				Owner:           libConfirmVoteState.MsgOption.Owner,
			},
			VoteOption: types.VoteOption(libConfirmVoteState.VoteOption),
			CreateAt:   libConfirmVoteState.CreateAt,
			Sign:       libConfirmVoteState.Sign,
		}
		if types.VoteOption(libConfirmVoteState.VoteOption) == types.Yes {
			if _, ok := yesVotes[libConfirmVoteState.MsgOption.SenderRole]; !ok {
				yesVotes[libConfirmVoteState.MsgOption.SenderRole] = 1
			} else {
				yesVotes[libConfirmVoteState.MsgOption.SenderRole] += 1
			}
		}

		if _, ok := voteStatus[libConfirmVoteState.MsgOption.SenderRole]; !ok {
			yesVotes[libConfirmVoteState.MsgOption.SenderRole] = 1
		} else {
			yesVotes[libConfirmVoteState.MsgOption.SenderRole] += 1
		}
	}
	//endregion

	// region recovery proposalPeerInfoCache
	iter = w.db.NewIteratorWithPrefixAndStart(proposalPeerInfoCachePrefix, nil)
	proposalPeerInfoCache := make(map[common.Hash]*twopcpb.ConfirmTaskPeerInfo, 0)
	prefixLength = len(proposalPeerInfoCachePrefix)
	libProposalPeerInfoCache := &twopcpb.ConfirmTaskPeerInfo{}
	for iter.Next() {
		key := iter.Key()
		proposalId := common.BytesToHash(key[prefixLength:])
		if err := proto.Unmarshal(iter.Value(), libProposalPeerInfoCache); err != nil {
			log.Fatal("marshaling error: ", err)
		}
		proposalPeerInfoCache[proposalId] = libProposalPeerInfoCache
	}
	// endregion

	return recoveryState(proposalSet, prepareVotes, confirmVotes, proposalPeerInfoCache,w)
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
		fmt.Println(proposalId.String(), string(libProposalPeerInfoCache.OwnerPeerInfo.PartyId))
	}
}
