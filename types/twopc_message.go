package types

import (
	"encoding/json"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	pb "github.com/RosettaFlow/Carrier-Go/lib/consensus/twopc"
	"github.com/ethereum/go-ethereum/rlp"
	"sync/atomic"
)

const (
	SendTaskDir ProposalTaskDir = 0x00
	RecvTaskDir ProposalTaskDir = 0x01
)

type ProposalTaskDir uint8

func (dir ProposalTaskDir) String() string {
	if dir == SendTaskDir {
		return "sendTask"
	} else {
		return "recvTask"
	}
}

// ------------------------------- About PrepareMsg -------------------------------
type PrepareMsgWrap struct {
	*pb.PrepareMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}
func (msg *PrepareMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PrepareMsgWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *PrepareMsgWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, []interface{}{
		msg.ProposalId,
		msg.TaskRole,
		msg.Owner,
		msg.TaskInfo,
		msg.CreateAt,
	})

	hasher.Sum(hash[:0])
	return hash
}
func (msg *PrepareMsgWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *PrepareMsgWrap) Signature() []byte {return msg.Sign}


// ------------------------------- About PrepareVote -------------------------------
type PrepareVoteWrap struct {
	*pb.PrepareVote
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}
func (msg *PrepareVoteWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PrepareVoteWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *PrepareVoteWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, []interface{}{
		msg.ProposalId,
		msg.TaskRole,
		msg.Owner,
		msg.VoteOption,
		msg.PeerInfo,
		msg.CreateAt,
	})

	hasher.Sum(hash[:0])
	return hash
}
func (msg *PrepareVoteWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *PrepareVoteWrap) Signature() []byte {return msg.Sign}



// ------------------------------- About ConfirmMsg -------------------------------
type ConfirmMsgWrap struct {
	*pb.ConfirmMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}
func (msg *ConfirmMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *ConfirmMsgWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *ConfirmMsgWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, []interface{}{
		msg.ProposalId,
		msg.Owner,
		msg.CreateAt,
	})

	hasher.Sum(hash[:0])
	return hash
}
func (msg *ConfirmMsgWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *ConfirmMsgWrap) Signature() []byte {return msg.Sign}




// ------------------------------- About ConfirmVote -------------------------------
type ConfirmVoteWrap struct {
	*pb.ConfirmVote
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}
func (msg *ConfirmVoteWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *ConfirmVoteWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *ConfirmVoteWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, []interface{}{
		msg.ProposalId,
		msg.TaskRole,
		msg.Owner,
		msg.VoteOption,
		msg.CreateAt,
	})

	hasher.Sum(hash[:0])
	return hash
}
func (msg *ConfirmVoteWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *ConfirmVoteWrap) Signature() []byte {return msg.Sign}





// ------------------------------- About CommitMsg -------------------------------
type CommitMsgWrap struct {
	*pb.CommitMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}
func (msg *CommitMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *CommitMsgWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *CommitMsgWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, []interface{}{
		msg.ProposalId,
		msg.Owner,
		msg.CreateAt,
	})

	hasher.Sum(hash[:0])
	return hash
}
func (msg *CommitMsgWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *CommitMsgWrap) Signature() []byte {return msg.Sign}


// ------------------------------- About TaskResultMsg -------------------------------
type TaskResultMsgWrap struct {
	*pb.TaskResultMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}
func (msg *TaskResultMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil{
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *TaskResultMsgWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *TaskResultMsgWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, []interface{}{
		msg.ProposalId,
		msg.TaskRole,
		msg.TaskId,
		msg.TaskEventList,
		msg.CreateAt,
	})

	hasher.Sum(hash[:0])
	return hash
}
func (msg *TaskResultMsgWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *TaskResultMsgWrap) Signature() []byte {return msg.Sign}