package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"strings"
	"sync/atomic"
)

const (
	PREFIX_POWER_ID         = "power:"
	PREFIX_METADATA_ID      = "metadata:"
	PREFIX_TASK_ID          = "task:"
	PREFIX_METADATA_AUTH_ID = "metadataAuth:"

	MSG_IDENTITY        = "identityMsg"
	MSG_IDENTITY_REVOKE = "identityRevokeMsg"
	MSG_POWER           = "powerMsg"
	MSG_POWER_REVOKE    = "powerRevokeMsg"
	MSG_METADATA        = "metaDataMsg"
	MSG_METADATA_REVOKE = "metaDataRevokeMsg"
	MSG_TASK            = "taskMsg"
	MSG_TASK_TERMINATE  = "taskTerminateMsg"

	MSG_METADATAAUTHORITY       = "MetadataAuthorityMsg"
	MSG_METADATAAUTHORITYREVOKE = "MetadataAuthorityRevokeMsg"

	MaxNodeNameLen   = 30 // 30 char
	MaxNodeIdLen     = 140
	MaxIdentityIdLen = 140
	MaxImageUrlLen   = 140
	MaxDetailsLen    = 280
)

type MessageType string

type Msg interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) error
	String() string
	MsgType() string
}

// ------------------- identity -------------------

type IdentityMsg struct {
	organization *apicommonpb.Organization
	CreateAt     uint64 `json:"createAt"`
}

func NewIdentityMessageFromRequest(req *pb.ApplyIdentityJoinRequest) *IdentityMsg {
	return &IdentityMsg{
		organization: &apicommonpb.Organization{
			NodeName:   req.GetMember().GetNodeName(),
			NodeId:     req.GetMember().GetNodeId(),
			IdentityId: req.GetMember().GetIdentityId(),
			ImageUrl:   req.GetMember().GetImageUrl(),
			Details:    req.GetMember().GetDetails(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}
}

func (msg *IdentityMsg) ToDataCenter() *Identity {
	return NewIdentity(&libtypes.IdentityPB{
		NodeName:   msg.GetOrganization().GetNodeName(),
		NodeId:     msg.GetOrganization().GetNodeId(),
		IdentityId: msg.GetOrganization().GetIdentityId(),
		ImageUrl:   msg.GetOrganization().GetImageUrl(),
		Details:    msg.GetOrganization().GetDetails(),
		DataId:     "",
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		Status:     apicommonpb.CommonStatus_CommonStatus_Normal,
		Credential: "",
	})
}
func (msg *IdentityMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *IdentityMsg) Unmarshal(b []byte) error { return nil }
func (msg *IdentityMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *IdentityMsg) MsgType() string                            { return MSG_IDENTITY }
func (msg *IdentityMsg) GetOrganization() *apicommonpb.Organization { return msg.organization }
func (msg *IdentityMsg) GetOwnerName() string                       { return msg.GetOrganization().NodeName }
func (msg *IdentityMsg) GetOwnerNodeId() string                     { return msg.GetOrganization().NodeId }
func (msg *IdentityMsg) GetOwnerIdentityId() string                 { return msg.GetOrganization().IdentityId }
func (msg *IdentityMsg) GetCreateAt() uint64                        { return msg.CreateAt }
func (msg *IdentityMsg) SetOwnerNodeId(nodeId string)               { msg.GetOrganization().NodeId = nodeId }

func (msg *IdentityMsg) CheckLength() error {

	if len(msg.GetOrganization().GetNodeName()) > MaxNodeNameLen {
		return fmt.Errorf("NodeName overflow, got len is: %d, max len is: %d", len(msg.GetOrganization().GetNodeName()), MaxNodeNameLen)
	}
	if len(msg.GetOrganization().GetNodeId()) > MaxNodeIdLen {
		return fmt.Errorf("NodeId overflow, got len is: %d, max len is: %d", len(msg.GetOrganization().GetNodeId()), MaxNodeIdLen)
	}
	if len(msg.GetOrganization().GetIdentityId()) > MaxIdentityIdLen {
		return fmt.Errorf("IdentityId overflow, got len is: %d, max len is: %d", len(msg.GetOrganization().GetIdentityId()), MaxIdentityIdLen)
	}
	if len(msg.GetOrganization().GetImageUrl()) > MaxImageUrlLen {
		return fmt.Errorf("ImageUrl overflow, got len is: %d, max len is: %d", len(msg.GetOrganization().GetImageUrl()), MaxImageUrlLen)
	}
	if len(msg.GetOrganization().GetDetails()) > MaxDetailsLen {
		return fmt.Errorf("Details overflow, got len is: %d, max len is: %d", len(msg.GetOrganization().GetDetails()), MaxDetailsLen)
	}
	return nil
}

type IdentityRevokeMsg struct {
	CreateAt uint64 `json:"createAt"`
}

func NewIdentityRevokeMessage() *IdentityRevokeMsg {
	return &IdentityRevokeMsg{
		CreateAt: timeutils.UnixMsecUint64(),
	}
}

func (msg *IdentityRevokeMsg) GetCreateAt() uint64      { return msg.CreateAt }
func (msg *IdentityRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *IdentityRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *IdentityRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *IdentityRevokeMsg) MsgType() string { return MSG_IDENTITY_REVOKE }

type IdentityMsgArr []*IdentityMsg
type IdentityRevokeMsgArr []*IdentityRevokeMsg

// Len returns the length of s.
func (s IdentityMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s IdentityMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IdentityMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// Len returns the length of s.
func (s IdentityRevokeMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s IdentityRevokeMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IdentityRevokeMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// ------------------- power -------------------

type PowerMsg struct {
	// This is only used when marshaling to JSON.
	PowerId   string `json:"powerId"`
	JobNodeId string `json:"jobNodeId"`
	CreateAt  uint64 `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewPowerMessageFromRequest(req *pb.PublishPowerRequest) *PowerMsg {
	msg := &PowerMsg{
		JobNodeId: req.GetJobNodeId(),
		CreateAt:  timeutils.UnixMsecUint64(),
	}
	msg.GenPowerId()
	return msg
}

func (msg *PowerMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *PowerMsg) Unmarshal(b []byte) error { return nil }
func (msg *PowerMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PowerMsg) MsgType() string { return MSG_POWER }

func (msg *PowerMsg) GetJobNodeId() string { return msg.JobNodeId }
func (msg *PowerMsg) GetPowerId() string   { return msg.PowerId }
func (msg *PowerMsg) GetCreateAt() uint64  { return msg.CreateAt }
func (msg *PowerMsg) GenPowerId() string {
	if "" != msg.GetPowerId() {
		return msg.GetPowerId()
	}
	msg.PowerId = PREFIX_POWER_ID + msg.HashByCreateTime().Hex()
	return msg.GetPowerId()
}

func (msg *PowerMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	var buf bytes.Buffer
	buf.Write([]byte(msg.GetJobNodeId()))
	buf.Write([]byte(msg.GetPowerId()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *PowerMsg) HashByCreateTime() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetJobNodeId()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	return rlputil.RlpHash(buf.Bytes())
}

type PowerRevokeMsg struct {
	PowerId  string `json:"powerId"`
	CreateAt uint64 `json:"createAt"`
}

func NewPowerRevokeMessageFromRequest(req *pb.RevokePowerRequest) *PowerRevokeMsg {
	return &PowerRevokeMsg{
		PowerId:  req.GetPowerId(),
		CreateAt: timeutils.UnixMsecUint64(),
	}
}

func (msg *PowerRevokeMsg) GetPowerId() string       { return msg.PowerId }
func (msg *PowerRevokeMsg) GetCreateAt() uint64      { return msg.CreateAt }
func (msg *PowerRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *PowerRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *PowerRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *PowerRevokeMsg) MsgType() string { return MSG_POWER_REVOKE }

type PowerMsgArr []*PowerMsg
type PowerRevokeMsgArr []*PowerRevokeMsg

// Len returns the length of s.
func (s PowerMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s PowerMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PowerMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// Len returns the length of s.
func (s PowerRevokeMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s PowerRevokeMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PowerRevokeMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// ------------------- metaData -------------------

type MetadataMsg struct {
	MetadataId      string                    `json:"metadataId"`
	MetadataSummary *libtypes.MetadataSummary `json:"metadataSummary"`
	CreateAt        uint64                    `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewMetadataMessageFromRequest(req *pb.PublishMetadataRequest) *MetadataMsg {
	metadataMsg := &MetadataMsg{
		MetadataSummary: &libtypes.MetadataSummary{
			MetadataId:     req.GetInformation().GetMetadataSummary().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataSummary().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataSummary().GetMetadataType(),
			FileHash:       req.GetInformation().GetMetadataSummary().GetFileHash(),
			Desc:           req.GetInformation().GetMetadataSummary().GetDesc(),
			FileType:       req.GetInformation().GetMetadataSummary().GetFileType(),
			Industry:       req.GetInformation().GetMetadataSummary().GetIndustry(),
			State:          req.GetInformation().GetMetadataSummary().GetState(),
			PublishAt:      req.GetInformation().GetMetadataSummary().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetMetadataSummary().GetUpdateAt(),
			Nonce:          req.GetInformation().GetMetadataSummary().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataSummary().GetMetadataOption(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}

	metadataMsg.GenMetadataId()
	return metadataMsg
}

func (msg *MetadataMsg) ToDataCenter(identity *apicommonpb.Organization) *Metadata {
	return NewMetadata(&libtypes.MetadataPB{

		MetadataId:   msg.GetMetadataId(),
		Owner:        identity,
		DataId:       msg.GetMetadataId(),
		DataStatus:   apicommonpb.DataStatus_DataStatus_Valid,
		MetadataName: msg.GetMetadataName(),
		MetadataType: msg.GetMetadataType(),
		FileHash:     msg.GetFileHash(),
		Desc:         msg.GetDesc(),
		FileType:     msg.GetFileType(),
		Industry:     msg.GetIndustry(),
		// metaData status, eg: create/release/revoke
		State:          apicommonpb.MetadataState_MetadataState_Released,
		PublishAt:      timeutils.UnixMsecUint64(),
		UpdateAt:       timeutils.UnixMsecUint64(),
		Nonce:          msg.GetNonce(),
		MetadataOption: msg.GetMetadataOption(),
	})
}
func (msg *MetadataMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataMsg) MsgType() string { return MSG_METADATA }
func (msg *MetadataMsg) GetMetadataSummary() *libtypes.MetadataSummary {
	return msg.MetadataSummary
}

func (msg *MetadataMsg) GetMetadataName() string { return msg.GetMetadataSummary().MetadataName }
func (msg *MetadataMsg) GetMetadataType() uint32 { return msg.GetMetadataSummary().MetadataType }
func (msg *MetadataMsg) GetFileHash() string     { return msg.GetMetadataSummary().FileHash }
func (msg *MetadataMsg) GetDesc() string         { return msg.GetMetadataSummary().Desc }
func (msg *MetadataMsg) GetFileType() apicommonpb.OriginFileType {
	return msg.GetMetadataSummary().FileType
}
func (msg *MetadataMsg) GetNonce() uint64 { return msg.GetMetadataSummary().Nonce }
func (msg *MetadataMsg) GetState() apicommonpb.MetadataState {
	return msg.GetMetadataSummary().State
}
func (msg *MetadataMsg) GetIndustry() string       { return msg.GetMetadataSummary().Industry }
func (msg *MetadataMsg) GetMetadataOption() string { return msg.GetMetadataSummary().MetadataOption }
func (msg *MetadataMsg) GetCreateAt() uint64       { return msg.CreateAt }
func (msg *MetadataMsg) GetMetadataId() string     { return msg.MetadataId }

func (msg *MetadataMsg) GenMetadataId() string {
	if "" != msg.GetMetadataId() {
		return msg.GetMetadataId()
	}
	msg.MetadataId = PREFIX_METADATA_ID + msg.HashByCreateTime().Hex()
	return msg.GetMetadataId()
}
func (msg *MetadataMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	var buf bytes.Buffer
	buf.Write([]byte(msg.GetMetadataId()))
	buf.Write([]byte(msg.GetMetadataName()))
	buf.Write(bytesutil.Uint32ToBytes(msg.GetMetadataType()))
	buf.Write([]byte(msg.GetFileHash()))
	buf.Write([]byte(msg.GetDesc()))
	buf.Write([]byte(msg.GetFileType().String()))
	buf.Write([]byte(msg.GetIndustry()))
	buf.Write([]byte(msg.GetState().String()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetNonce()))
	buf.Write([]byte(msg.GetMetadataOption()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *MetadataMsg) HashByCreateTime() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetFileHash()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	return rlputil.RlpHash(buf.Bytes())
}

type MetadataRevokeMsg struct {
	MetadataId string `json:"metadataId"`
	CreateAt   uint64 `json:"createAt"`
}

func NewMetadataRevokeMessageFromRequest(req *pb.RevokeMetadataRequest) *MetadataRevokeMsg {
	return &MetadataRevokeMsg{
		MetadataId: req.GetMetadataId(),
		CreateAt:   timeutils.UnixMsecUint64(),
	}
}

func (msg *MetadataRevokeMsg) GetMetadataId() string { return msg.MetadataId }
func (msg *MetadataRevokeMsg) GetCreateAt() uint64   { return msg.CreateAt }
func (msg *MetadataRevokeMsg) ToDataCenter(identity *apicommonpb.Organization) *Metadata {
	return NewMetadata(&libtypes.MetadataPB{
		MetadataId: msg.GetMetadataId(),
		Owner:      identity,
		DataId:     msg.GetMetadataId(),
		DataStatus: apicommonpb.DataStatus_DataStatus_Invalid,
		// metaData status, eg: create/release/revoke
		State:    apicommonpb.MetadataState_MetadataState_Revoked,
		UpdateAt: timeutils.UnixMsecUint64(),
	})
}
func (msg *MetadataRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataRevokeMsg) MsgType() string { return MSG_METADATA_REVOKE }

type MetadataMsgArr []*MetadataMsg
type MetadataRevokeMsgArr []*MetadataRevokeMsg

// Len returns the length of s.
func (s MetadataMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// Len returns the length of s.
func (s MetadataRevokeMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataRevokeMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataRevokeMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// ------------------- metadata authority apply -------------------

type MetadataAuthorityMsg struct {
	MetadataAuthId string                   `json:"metaDataAuthId"`
	User           string                   `json:"user"`
	UserType       apicommonpb.UserType     `json:"userType"`
	Auth           *types.MetadataAuthority `json:"auth"`
	Sign           []byte                   `json:"sign"`
	CreateAt       uint64                   `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewMetadataAuthorityMessageFromRequest(req *pb.ApplyMetadataAuthorityRequest) *MetadataAuthorityMsg {
	metadataAuthorityMsg := &MetadataAuthorityMsg{
		User:     req.GetUser(),
		UserType: req.GetUserType(),
		Auth:     req.GetAuth(),
		Sign:     req.GetSign(),
		CreateAt: timeutils.UnixMsecUint64(),
	}

	metadataAuthorityMsg.GenMetadataAuthId()

	return metadataAuthorityMsg
}

func (msg *MetadataAuthorityMsg) GetMetadataAuthId() string                      { return msg.MetadataAuthId }
func (msg *MetadataAuthorityMsg) GetUser() string                                { return msg.User }
func (msg *MetadataAuthorityMsg) GetUserType() apicommonpb.UserType              { return msg.UserType }
func (msg *MetadataAuthorityMsg) GetMetadataAuthority() *types.MetadataAuthority { return msg.Auth }
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwner() *apicommonpb.Organization {
	return msg.Auth.GetOwner()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerIdentity() string {
	return msg.Auth.GetOwner().GetIdentityId()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerNodeId() string {
	return msg.Auth.GetOwner().GetNodeId()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerNodeName() string {
	return msg.Auth.GetOwner().GetNodeName()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityMetadataId() string {
	return msg.Auth.GetMetadataId()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityUsageRule() *types.MetadataUsageRule {
	return msg.Auth.GetUsageRule()
}
func (msg *MetadataAuthorityMsg) GetSign() []byte     { return msg.Sign }
func (msg *MetadataAuthorityMsg) GetCreateAt() uint64 { return msg.CreateAt }

func (msg *MetadataAuthorityMsg) GenMetadataAuthId() string {
	if "" != msg.GetMetadataAuthId() {
		return msg.GetMetadataAuthId()
	}
	msg.MetadataAuthId = PREFIX_METADATA_AUTH_ID + msg.HashByCreateTime().Hex()
	return msg.GetMetadataAuthId()
}
func (msg *MetadataAuthorityMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	b, _ := msg.GetMetadataAuthority().Marshal()
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetUser()))
	buf.Write([]byte(msg.GetUserType().String()))
	buf.Write([]byte(msg.GetMetadataAuthorityMetadataId()))
	buf.Write(b)
	buf.Write(msg.GetSign())
	buf.Write(bytesutil.Uint64ToBytes(msg.CreateAt))
	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *MetadataAuthorityMsg) HashByCreateTime() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetUser()))
	buf.Write([]byte(msg.GetUserType().String()))
	buf.Write([]byte(msg.GetMetadataAuthorityMetadataId()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	return rlputil.RlpHash(buf.Bytes())
}

func (msg *MetadataAuthorityMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataAuthorityMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataAuthorityMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataAuthorityMsg) MsgType() string { return MSG_METADATAAUTHORITY }

type MetadataAuthorityRevokeMsg struct {
	User           string
	UserType       apicommonpb.UserType
	MetadataAuthId string
	Sign           []byte
	CreateAt       uint64
}

func NewMetadataAuthorityRevokeMessageFromRequest(req *pb.RevokeMetadataAuthorityRequest) *MetadataAuthorityRevokeMsg {
	return &MetadataAuthorityRevokeMsg{
		User:           req.GetUser(),
		UserType:       req.GetUserType(),
		MetadataAuthId: req.GetMetadataAuthId(),
		Sign:           req.GetSign(),
		CreateAt:       timeutils.UnixMsecUint64(),
	}
}

func (msg *MetadataAuthorityRevokeMsg) GetMetadataAuthId() string         { return msg.MetadataAuthId }
func (msg *MetadataAuthorityRevokeMsg) GetUser() string                   { return msg.User }
func (msg *MetadataAuthorityRevokeMsg) GetUserType() apicommonpb.UserType { return msg.UserType }
func (msg *MetadataAuthorityRevokeMsg) GetSign() []byte                   { return msg.Sign }
func (msg *MetadataAuthorityRevokeMsg) GetCreateAt() uint64               { return msg.CreateAt }

func (msg *MetadataAuthorityRevokeMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataAuthorityRevokeMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataAuthorityRevokeMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataAuthorityRevokeMsg) MsgType() string { return MSG_METADATAAUTHORITYREVOKE }

type MetadataAuthorityMsgArr []*MetadataAuthorityMsg
type MetadataAuthorityRevokeMsgArr []*MetadataAuthorityRevokeMsg

// Len returns the length of s.
func (s MetadataAuthorityMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataAuthorityMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataAuthorityMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// Len returns the length of s.
func (s MetadataAuthorityRevokeMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataAuthorityRevokeMsgArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s MetadataAuthorityRevokeMsgArr) Less(i, j int) bool {
	return s[i].GetCreateAt() < s[j].GetCreateAt()
}

// ------------------- task -------------------

type TaskBullet struct {
	TaskId      string
	Starve      bool
	Term        uint32
	Resched     uint32
	InQueueFlag bool // true: in queue/starveQueue,  false: was pop ?
}

func NewTaskBullet(taskId string) *TaskBullet {
	return &TaskBullet{
		TaskId: taskId,
	}
}

func (b *TaskBullet) GetTaskId() string    { return b.TaskId }
func (b *TaskBullet) IsStarve() bool       { return b.Starve }
func (b *TaskBullet) GetTerm() uint32      { return b.Term }
func (b *TaskBullet) GetResched() uint32   { return b.Resched }
func (b *TaskBullet) GetInQueueFlag() bool { return b.InQueueFlag }

func (b *TaskBullet) IncreaseResched() { b.Resched++ }
func (b *TaskBullet) DecreaseResched() {
	if b.GetResched() > 0 {
		b.Resched--
	}
}

func (b *TaskBullet) IncreaseTerm() { b.Term++ }
func (b *TaskBullet) DecreaseTerm() {
	if b.GetTerm() > 0 {
		b.Term--
	}
}
func (b *TaskBullet) IsOverlowReschedThreshold(reschedMaxCount uint32) bool {
	if b.GetResched() >= reschedMaxCount {
		return true
	}
	return false
}

type TaskBullets []*TaskBullet

func (h TaskBullets) Len() int           { return len(h) }
func (h TaskBullets) Less(i, j int) bool { return h[i].GetTerm() > h[j].GetTerm() } // term:  a.3 > c.2 > b.1,  So order is: a c b
func (h TaskBullets) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TaskBullets) Push(x interface{}) {
	*h = append(*h, x.(*TaskBullet))
}

func (h *TaskBullets) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func (h *TaskBullets) IncreaseTerm() {
	for i := range *h {
		(*h)[i].IncreaseTerm()
	}
}
func (h *TaskBullets) DecreaseTerm() {
	for i := range *h {
		(*h)[i].DecreaseTerm()
	}
}

func (h *TaskBullets) IncreaseTermByCallbackFn(f func(b *TaskBullet)) {
	for i := range *h {
		b := (*h)[i]
		b.IncreaseTerm()
		f(b)
	}
}
func (h *TaskBullets) DecreaseTermByCallbackFn(f func(b *TaskBullet)) {
	for i := range *h {
		b := (*h)[i]
		b.DecreaseTerm()
		f(b)
	}
}

type TaskMsg struct {
	Data          *Task
	// caches
	hash atomic.Value
}

func NewTaskMessageFromRequest(req *pb.PublishTaskDeclareRequest) *TaskMsg {
	return &TaskMsg{
		Data: NewTask(&libtypes.TaskPB{
			TaskId:       "",
			TaskName:     req.GetTaskName(),
			UserType:     req.GetUserType(),
			User:         req.GetUser(),
			Sender:       req.GetSender(),
			DataId:       "",
			DataStatus:   apicommonpb.DataStatus_DataStatus_Valid,
			State:        apicommonpb.TaskState_TaskState_Pending,
			Reason:       "",
			Desc:         req.GetDesc(),
			CreateAt:     timeutils.UnixMsecUint64(),
			EndAt:        0,
			StartAt:      0,
			AlgoSupplier: req.GetAlgoSupplier(),
			DataSuppliers: req.GetDataSuppliers(),
			//PowerSuppliers: ,
			Receivers:             req.GetReceivers(),
			OperationCost:         req.GetOperationCost(),
			CalculateContractCode: req.GetCalculateContractCode(),
			DataSplitContractCode: req.GetDataSplitContractCode(),
			ContractExtraParams:   req.GetContractExtraParams(),
			Sign:                  req.GetSign(),
		}),
	}
}

func (msg *TaskMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskMsg) String() string {
	return fmt.Sprintf(`{"taskId": %s, "powerPartyIds": %s, "task": %s}`,
		msg.GetTask().GetTaskId(), "["+strings.Join(msg.PowerPartyIds, ",")+"]", msg.GetTask().GetTaskData().String())
}
func (msg *TaskMsg) MsgType() string { return MSG_TASK }

func (msg *TaskMsg) GetTask() *Task             { return msg.Data }
func (msg *TaskMsg) GetTaskData() *types.TaskPB { return msg.GetTask().GetTaskData() }
func (msg *TaskMsg) GetUserType() apicommonpb.UserType {
	return msg.GetTask().GetTaskData().GetUserType()
}
func (msg *TaskMsg) GetUser() string { return msg.GetTask().GetTaskData().GetUser() }
func (msg *TaskMsg) GetSender() *apicommonpb.TaskOrganization {
	return &apicommonpb.TaskOrganization{
		PartyId:    msg.GetTask().GetTaskData().GetPartyId(),
		NodeName:   msg.GetTask().GetTaskData().GetNodeName(),
		NodeId:     msg.GetTask().GetTaskData().GetNodeId(),
		IdentityId: msg.GetTask().GetTaskData().GetIdentityId(),
	}
}
func (msg *TaskMsg) GetSenderName() string       { return msg.GetTask().GetTaskData().GetNodeName() }
func (msg *TaskMsg) GetSenderNodeId() string     { return msg.GetTask().GetTaskData().GetNodeId() }
func (msg *TaskMsg) GetSenderIdentityId() string { return msg.GetTask().GetTaskData().GetIdentityId() }
func (msg *TaskMsg) GetSenderPartyId() string    { return msg.GetTask().GetTaskData().GetPartyId() }
func (msg *TaskMsg) GetTaskId() string           { return msg.GetTask().GetTaskData().GetTaskId() }
func (msg *TaskMsg) GetTaskName() string         { return msg.GetTask().GetTaskData().GetTaskName() }
func (msg *TaskMsg) GetAlgoSupplier() *apicommonpb.TaskOrganization {
	return &apicommonpb.TaskOrganization{
		PartyId:    msg.GetTask().GetTaskData().GetAlgoSupplier().GetPartyId(),
		NodeName:   msg.GetTask().GetTaskData().GetAlgoSupplier().GetNodeName(),
		NodeId:     msg.GetTask().GetTaskData().GetAlgoSupplier().GetNodeId(),
		IdentityId: msg.GetTask().GetTaskData().GetAlgoSupplier().GetIdentityId(),
	}
}
func (msg *TaskMsg) GetTaskMetadataSuppliers() []*apicommonpb.TaskOrganization {

	partners := make([]*apicommonpb.TaskOrganization, len(msg.GetTask().GetTaskData().GetDataSuppliers()))
	for i, v := range msg.Data.GetTaskData().GetDataSuppliers() {

		partners[i] = &apicommonpb.TaskOrganization{
			PartyId:    v.GetOrganization().GetPartyId(),
			NodeName:   v.GetOrganization().GetNodeName(),
			NodeId:     v.GetOrganization().GetNodeId(),
			IdentityId: v.GetOrganization().GetIdentityId(),
		}
	}
	return partners
}
func (msg *TaskMsg) GetTaskMetadataSupplierDatas() []*libtypes.TaskDataSupplier {
	return msg.GetTask().GetTaskData().GetDataSuppliers()
}

func (msg *TaskMsg) GetTaskResourceSuppliers() []*apicommonpb.TaskOrganization {
	powers := make([]*apicommonpb.TaskOrganization, len(msg.GetTask().GetTaskData().GetPowerSuppliers()))
	for i, v := range msg.GetTask().GetTaskData().GetPowerSuppliers() {
		powers[i] = &apicommonpb.TaskOrganization{
			PartyId:    v.GetOrganization().GetPartyId(),
			NodeName:   v.GetOrganization().GetNodeName(),
			NodeId:     v.GetOrganization().GetNodeId(),
			IdentityId: v.GetOrganization().GetIdentityId(),
		}
	}
	return powers
}
func (msg *TaskMsg) GetTaskResourceSupplierDatas() []*libtypes.TaskPowerSupplier {
	return msg.Data.GetTaskData().GetPowerSuppliers()
}
func (msg *TaskMsg) GetPowerPartyIds() []string { return msg.PowerPartyIds }
func (msg *TaskMsg) GetReceivers() []*apicommonpb.TaskOrganization {
	return msg.Data.GetTaskData().GetReceivers()
}

func (msg *TaskMsg) GetCalculateContractCode() string {
	return msg.Data.GetTaskData().GetCalculateContractCode()
}
func (msg *TaskMsg) GetDataSplitContractCode() string {
	return msg.Data.GetTaskData().GetDataSplitContractCode()
}
func (msg *TaskMsg) GetContractExtraParams() string {
	return msg.Data.GetTaskData().GetContractExtraParams()
}
func (msg *TaskMsg) GetOperationCost() *apicommonpb.TaskResourceCostDeclare {
	return msg.Data.GetTaskData().GetOperationCost()
}
func (msg *TaskMsg) GetCreateAt() uint64 { return msg.GetTask().GetTaskData().GetCreateAt() }
func (msg *TaskMsg) GenTaskId() string {
	if "" != msg.GetTask().GetTaskId() {
		return msg.GetTask().GetTaskId()
	}
	msg.GetTask().GetTaskData().TaskId = PREFIX_TASK_ID + msg.HashByCreateTime().Hex()
	return msg.GetTask().GetTaskId()
}
func (msg *TaskMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlputil.RlpHash(msg.GetTask().GetTaskData())
	msg.hash.Store(v)
	return v
}

func (msg *TaskMsg) HashByCreateTime() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetTask().GetTaskData().GetIdentityId()))
	buf.Write([]byte(msg.GetTask().GetTaskData().GetUserType().String()))
	buf.Write([]byte(msg.GetTask().GetTaskData().GetUser()))
	buf.Write([]byte(msg.GetTask().GetTaskData().GetPartyId()))
	buf.Write([]byte(msg.GetTask().GetTaskData().GetTaskName()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	return rlputil.RlpHash(buf.Bytes())
}

type TaskTerminateMsg struct {
	UserType apicommonpb.UserType `json:"userType"`
	User     string               `json:"user"`
	TaskId   string               `json:"taskId"`
	Sign     []byte               `json:"sign"`
	CreateAt uint64               `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewTaskTerminateMsg(userType apicommonpb.UserType, user, taskId string, sign []byte) *TaskTerminateMsg {
	return &TaskTerminateMsg{
		UserType: userType,
		User:     user,
		TaskId:   taskId,
		Sign:     sign,
		CreateAt: timeutils.UnixMsecUint64(),
	}
}

func (msg *TaskTerminateMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskTerminateMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskTerminateMsg) String() string {
	return fmt.Sprintf(`{"userType": %s, "user": %s, "taskId": %s, "sign": %v, "createAt": %d}`,
		msg.UserType.String(), msg.User, msg.TaskId, msg.Sign, msg.CreateAt)
}
func (msg *TaskTerminateMsg) MsgType() string                   { return MSG_TASK_TERMINATE }
func (msg *TaskTerminateMsg) GetUserType() apicommonpb.UserType { return msg.UserType }
func (msg *TaskTerminateMsg) GetUser() string                   { return msg.User }
func (msg *TaskTerminateMsg) GetTaskId() string                 { return msg.TaskId }
func (msg *TaskTerminateMsg) GetSign() []byte                   { return msg.Sign }
func (msg *TaskTerminateMsg) GetCreateAt() uint64               { return msg.CreateAt }
func (msg *TaskTerminateMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetUserType().String()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write([]byte(msg.GetTaskId()))
	buf.Write(msg.GetSign())
	buf.Write(bytesutil.Uint64ToBytes(msg.GetCreateAt()))
	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *TaskTerminateMsg) HashByCreateTime() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetUserType().String()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write([]byte(msg.GetTaskId()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	return rlputil.RlpHash(buf.Bytes())
}

type BadTaskMsg struct {
	msg *TaskMsg
	s   string // the error reason
}

func NewBadTaskMsg(msg *TaskMsg, s string) *BadTaskMsg {
	return &BadTaskMsg{
		msg: msg,
		s:   s,
	}
}
func (bmsg *BadTaskMsg) GetTaskMsg() *TaskMsg { return bmsg.msg }
func (bmsg *BadTaskMsg) GetErr() error        { return fmt.Errorf(bmsg.s) }
func (bmsg *BadTaskMsg) GetErrStr() string    { return bmsg.s }

type TaskMsgArr []*TaskMsg
type TaskTerminateMsgArr []*TaskTerminateMsg
type BadTaskMsgArr []*BadTaskMsg

// Len returns the length of s.
func (s TaskMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskMsgArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TaskMsgArr) Less(i, j int) bool {
	return s[i].GetTask().GetTaskData().GetCreateAt() < s[j].GetTask().GetTaskData().GetCreateAt()
}

// Len returns the length of s.
func (s TaskTerminateMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskTerminateMsgArr) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TaskTerminateMsgArr) Less(i, j int) bool {
	return s[i].GetCreateAt() < s[j].GetCreateAt()
}

/**
Example:
{
  "party_id": "p0",
  "data_party": {
      "input_file": "../data/bank_predict_data.csv",
       "key_column": "CLIENT_ID",
       "selected_columns": ["col1", "col2"]
    },
  "dynamic_parameter": {
    "model_restore_party": "p0",
    "train_task_id": "task_id"
  }
}

or:

{
  "party_id": "p0",
  "data_party": {
    "input_file": "../data/bank_train_data.csv",
    "key_column": "CLIENT_ID",
    "selected_columns": ["col1", "col2"]
  },
  "dynamic_parameter": {
    "label_owner": "p0",
    "label_column_name": "Y",
    "algorithm_parameter": {
      "epochs": 10,
      "batch_size": 256,
      "learning_rate": 0.1
    }
  }
}
*/
type FighterTaskReadyGoReqContractCfg struct {
	PartyId   string `json:"party_id"`
	DataParty struct {
		InputFile       string   `json:"input_file"`
		KeyColumn       string   `json:"key_column"`
		SelectedColumns []string `json:"selected_columns"`
	} `json:"data_party"`
	DynamicParameter map[string]interface{} `json:"dynamic_parameter"`
}
