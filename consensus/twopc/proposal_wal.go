package twopc

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	ctypes "github.com/RosettaFlow/Carrier-Go/consensus/twopc/types"
	"github.com/RosettaFlow/Carrier-Go/db"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	twopcpb "github.com/RosettaFlow/Carrier-Go/lib/netmsg/consensus/twopc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/gogo/protobuf/proto"
)

var (
	proposalSetPrefix           = []byte("proposalSet:")
	prepareVotesPrefix          = []byte("prepareVotes:")
	confirmVotesPrefix          = []byte("confirmVotes:")
	proposalPeerInfoCachePrefix = []byte("proposalPeerInfoCache:")
	databasePath                = "./consensus_state_save"
)

func GetProposalSetKey(hash common.Hash, partyId string) []byte {
	result := append(proposalSetPrefix, hash.Bytes()...)
	return append(result, []byte(partyId)...)
}

func GetProposalPeerInfoCacheKey(hash common.Hash) []byte {
	return append(proposalPeerInfoCachePrefix, hash.Bytes()...)
}

func GetPrepareVotesKey(hash common.Hash, partyId string) []byte {
	result := append(prepareVotesPrefix, hash.Bytes()...)
	return append(result, []byte(partyId)...)
}

func GetConfirmVotesKey(hash common.Hash, partyId string) []byte {
	result := append(confirmVotesPrefix, hash.Bytes()...)
	return append(result, []byte(partyId)...)
}

func OpenDatabase(dbpath string, cache int, handles int) (db.Database, error) {
	db, err := db.NewLDBDatabase(dbpath, cache, handles)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func UpdateOrgProposalState(proposalId common.Hash, sender *apicommonpb.TaskOrganization, orgState *ctypes.OrgProposalState) {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()
	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}

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
	data, er := proto.Marshal(pbObj)
	if er != nil {
		log.Fatal("marshaling error: ", err)
	}
	if err := db.Put(GetProposalSetKey(proposalId, orgState.TaskOrg.PartyId), data); err != nil {
		log.Warning("UpdateOrgProposalState to db fail,proposalId:", proposalId)
	}
}

func UpdateConfirmTaskPeerInfo(proposalId common.Hash, peerDesc *twopcpb.ConfirmTaskPeerInfo) {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()
	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}
	data, er := proto.Marshal(peerDesc)
	if er != nil {
		log.Fatal("marshaling error: ", err)
	}
	fmt.Println("proposalId:", proposalId)
	if err := db.Put(GetProposalPeerInfoCacheKey(proposalId), data); err != nil {
		log.Warning("UpdateConfirmTaskPeerInfo to db fail,proposalId:", proposalId)
	}
}

func UpdatePrepareVotes(vote *types.PrepareVote) {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()
	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}

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
	data, er := proto.Marshal(pbObj)
	if er != nil {
		log.Fatal("marshaling error: ", err)
	}
	if err := db.Put(GetPrepareVotesKey(vote.MsgOption.ProposalId, vote.MsgOption.SenderPartyId), data); err != nil {
		log.Warning("UpdatePrepareVotes to db fail,proposalId:", vote.MsgOption.ProposalId)
	}
}

func UpdateConfirmVotes(vote *types.ConfirmVote) {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()
	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}

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
	data, er := proto.Marshal(pbObj)
	if er != nil {
		log.Fatal("marshaling error: ", err)
	}
	if err := db.Put(GetConfirmVotesKey(vote.MsgOption.ProposalId, vote.MsgOption.SenderPartyId), data); err != nil {
		log.Warning("UpdateConfirmVotes to db fail,proposalId:", vote.MsgOption.ProposalId)
	}
}

func DeleteState(key []byte) {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()

	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}
	er := db.Delete(key)
	if err != nil {
		log.WithError(er).Fatal("Failed to delete from leveldb, key is ", key)
	}
}

func RecoveryState() *state {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()

	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}

	// region recovery proposalSet
	iter := db.NewIteratorWithPrefixAndStart(proposalSetPrefix, nil)
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
	iter = db.NewIteratorWithPrefixAndStart(prepareVotesPrefix, nil)
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
	iter = db.NewIteratorWithPrefixAndStart(confirmVotesPrefix, nil)
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
	iter = db.NewIteratorWithPrefixAndStart(proposalPeerInfoCachePrefix, nil)
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

	return recoveryState(proposalSet, prepareVotes, confirmVotes, proposalPeerInfoCache)
}

func UnmarshalTest() {
	db, err := OpenDatabase(databasePath, 0, 0)
	defer db.Close()

	if err != nil {
		log.Warning("open leveldb fail!,err:", err)
	}
	iter := db.NewIteratorWithPrefixAndStart(proposalPeerInfoCachePrefix, nil)
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
