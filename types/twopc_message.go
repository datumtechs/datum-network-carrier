package types

import (
	"bytes"
	"encoding/json"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/common/rlputil"
	"github.com/datumtechs/datum-network-carrier/crypto/sha3"
	commonpb "github.com/datumtechs/datum-network-carrier/lib/netmsg/common"
	twopcpb "github.com/datumtechs/datum-network-carrier/lib/netmsg/consensus/twopc"
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
	*twopcpb.PrepareMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}

func (msg *PrepareMsgWrap) GetData() *twopcpb.PrepareMsg { return msg.PrepareMsg }

func (msg *PrepareMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
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
	var buf bytes.Buffer
	buf.Write(msg.GetMsgOption().GetProposalId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetSenderRole()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetReceiverRole()))
	buf.Write(msg.GetMsgOption().GetSenderPartyId())
	buf.Write(msg.GetMsgOption().GetReceiverPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetName())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetNodeId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetIdentityId())
	buf.Write(msg.GetTaskInfo())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	rlp.Encode(hasher, buf.Bytes())
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
func (msg *PrepareMsgWrap) Signature() []byte { return msg.Sign }

// ------------------------------- About PrepareVote -------------------------------
type PrepareVoteWrap struct {
	*twopcpb.PrepareVote
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}

func (msg *PrepareVoteWrap) GetData () *twopcpb.PrepareVote { return msg.PrepareVote }

func (msg *PrepareVoteWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
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
	var buf bytes.Buffer
	buf.Write(msg.GetMsgOption().GetProposalId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetSenderRole()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetReceiverRole()))
	buf.Write(msg.GetMsgOption().GetSenderPartyId())
	buf.Write(msg.GetMsgOption().GetReceiverPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetName())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetNodeId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetIdentityId())
	buf.Write(msg.GetVoteOption())
	buf.Write(msg.GetPeerInfo().GetPartyId())
	buf.Write(msg.GetPeerInfo().GetIp())
	buf.Write(msg.GetPeerInfo().GetPort())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	rlp.Encode(hasher, buf.Bytes())
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
func (msg *PrepareVoteWrap) Signature() []byte { return msg.Sign }

// ------------------------------- About ConfirmMsg -------------------------------
type ConfirmMsgWrap struct {
	*twopcpb.ConfirmMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}

func (msg *ConfirmMsgWrap) GetData () *twopcpb.ConfirmMsg { return msg.ConfirmMsg }

func (msg *ConfirmMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
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
	var buf bytes.Buffer
	buf.Write(msg.GetMsgOption().GetProposalId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetSenderRole()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetReceiverRole()))
	buf.Write(msg.GetMsgOption().GetSenderPartyId())
	buf.Write(msg.GetMsgOption().GetReceiverPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetName())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetNodeId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetIdentityId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	rlp.Encode(hasher, buf.Bytes())
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
func (msg *ConfirmMsgWrap) Signature() []byte { return msg.Sign }

// ------------------------------- About ConfirmVote -------------------------------
type ConfirmVoteWrap struct {
	*twopcpb.ConfirmVote
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}

func (msg *ConfirmVoteWrap) GetData () *twopcpb.ConfirmVote { return msg.ConfirmVote }

func (msg *ConfirmVoteWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
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
	var buf bytes.Buffer
	buf.Write(msg.GetMsgOption().GetProposalId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetSenderRole()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetReceiverRole()))
	buf.Write(msg.GetMsgOption().GetSenderPartyId())
	buf.Write(msg.GetMsgOption().GetReceiverPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetName())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetNodeId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetIdentityId())
	buf.Write(msg.GetVoteOption())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	rlp.Encode(hasher, buf.Bytes())
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
func (msg *ConfirmVoteWrap) Signature() []byte { return msg.Sign }

// ------------------------------- About CommitMsg -------------------------------
type CommitMsgWrap struct {
	*twopcpb.CommitMsg
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}

func (msg *CommitMsgWrap) GetData() *twopcpb.CommitMsg { return msg.CommitMsg }

func (msg *CommitMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
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
	var buf bytes.Buffer
	buf.Write(msg.GetMsgOption().GetProposalId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetSenderRole()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMsgOption().GetReceiverRole()))
	buf.Write(msg.GetMsgOption().GetSenderPartyId())
	buf.Write(msg.GetMsgOption().GetReceiverPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetName())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetNodeId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetIdentityId())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	rlp.Encode(hasher, buf.Bytes())
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
func (msg *CommitMsgWrap) Signature() []byte { return msg.Sign }

// ------------------------------- About InterruptMsg -------------------------------
type TerminateConsensusMsgWrap struct {
	MsgOption *commonpb.MsgOption
	TaskId    []byte
	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
	hash     atomic.Value `json:"-" rlp:"-"`
}

func NewInterruptMsgWrap(taskId string, msgOption *commonpb.MsgOption) *TerminateConsensusMsgWrap {
	return &TerminateConsensusMsgWrap{
		TaskId: []byte(taskId),
		MsgOption: msgOption,
	}
}
func (msg *TerminateConsensusMsgWrap) GetTaskId() string                 { return string(msg.TaskId) }
func (msg *TerminateConsensusMsgWrap) GetMsgOption() *commonpb.MsgOption { return msg.MsgOption }
func (msg *TerminateConsensusMsgWrap) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *TerminateConsensusMsgWrap) SealHash() common.Hash {
	if sealHash := msg.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := msg._sealHash()
	msg.sealHash.Store(v)
	return v
}
func (msg *TerminateConsensusMsgWrap) _sealHash() (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetTaskId()))
	buf.Write(msg.GetMsgOption().GetSenderPartyId())
	buf.Write(msg.GetMsgOption().GetReceiverPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetPartyId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetName())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetNodeId())
	buf.Write(msg.GetMsgOption().GetMsgOwner().GetIdentityId())
	rlp.Encode(hasher, buf.Bytes())
	hasher.Sum(hash[:0])
	return hash
}
func (msg *TerminateConsensusMsgWrap) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg)
	msg.hash.Store(v)
	return v
}
func (msg *TerminateConsensusMsgWrap) Signature() []byte { return nil }