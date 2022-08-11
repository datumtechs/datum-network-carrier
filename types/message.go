package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/common/rlputil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"sync/atomic"
)

const (
	PREFIX_POWER_ID         = "power:"
	PREFIX_METADATA_ID      = "metadata:"
	PREFIX_TASK_ID          = "task:"
	PREFIX_METADATA_AUTH_ID = "metadataAuth:"

	MSG_IDENTITY                      = "identityMsg"
	MSG_UPDATE_IDENTITY_CREDENTIALMSG = "updateIdentityCredentialMsg"
	MSG_IDENTITY_REVOKE               = "identityRevokeMsg"
	MSG_POWER                         = "powerMsg"
	MSG_POWER_REVOKE                  = "powerRevokeMsg"
	MSG_METADATA                      = "metaDataMsg"
	MSG_UPDATEMETADATA                = "updateMetadataMsg"
	MSG_METADATA_REVOKE               = "metaDataRevokeMsg"
	MSG_TASK                          = "taskMsg"
	MSG_TASK_TERMINATE                = "taskTerminateMsg"

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
	organization *carriertypespb.Organization
	CreateAt     uint64 `json:"createAt"`
}

func NewIdentityMessageFromRequest(req *carrierapipb.ApplyIdentityJoinRequest) *IdentityMsg {
	return &IdentityMsg{
		organization: &carriertypespb.Organization{
			NodeName:   req.GetInformation().GetNodeName(),
			NodeId:     req.GetInformation().GetNodeId(),
			IdentityId: req.GetInformation().GetIdentityId(),
			ImageUrl:   req.GetInformation().GetImageUrl(),
			Details:    req.GetInformation().GetDetails(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}
}

func (msg *IdentityMsg) ToDataCenter() *Identity {
	return NewIdentity(&carriertypespb.IdentityPB{
		NodeName:   msg.GetOrganization().GetNodeName(),
		NodeId:     msg.GetOrganization().GetNodeId(),
		IdentityId: msg.GetOrganization().GetIdentityId(),
		ImageUrl:   msg.GetOrganization().GetImageUrl(),
		Details:    msg.GetOrganization().GetDetails(),
		DataId:     "",
		DataStatus: commonconstantpb.DataStatus_DataStatus_Valid,
		Status:     commonconstantpb.CommonStatus_CommonStatus_Valid,
		Credential: "",
		Nonce:      msg.GetOrganization().GetNonce(),
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
func (msg *IdentityMsg) MsgType() string                               { return MSG_IDENTITY }
func (msg *IdentityMsg) GetOrganization() *carriertypespb.Organization { return msg.organization }
func (msg *IdentityMsg) GetOwnerName() string                          { return msg.GetOrganization().NodeName }
func (msg *IdentityMsg) GetOwnerNodeId() string                        { return msg.GetOrganization().NodeId }
func (msg *IdentityMsg) GetOwnerIdentityId() string                    { return msg.GetOrganization().IdentityId }
func (msg *IdentityMsg) GetCreateAt() uint64                           { return msg.CreateAt }
func (msg *IdentityMsg) SetOwnerNodeId(nodeId string)                  { msg.GetOrganization().NodeId = nodeId }

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

func NewPowerMessageFromRequest(req *carrierapipb.PublishPowerRequest) *PowerMsg {
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

	// caches
	hash atomic.Value
}

func NewPowerRevokeMessageFromRequest(req *carrierapipb.RevokePowerRequest) *PowerRevokeMsg {
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
func (msg *PowerRevokeMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	var buf bytes.Buffer
	buf.Write([]byte(msg.GetPowerId()))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

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
	MetadataSummary *carriertypespb.MetadataSummary `json:"metadataSummary"`
	CreateAt        uint64                          `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewMetadataMessageFromRequest(req *carrierapipb.PublishMetadataRequest) *MetadataMsg {
	metadataMsg := &MetadataMsg{
		MetadataSummary: &carriertypespb.MetadataSummary{
			/**
			MetadataId     string
			MetadataName   string
			MetadataType   constant.MetadataType
			DataHash       string
			Desc           string
			LocationType   constant.DataLocationTyp
			DataType       constant.OrigindataType
			Industry       string
			State          constant.MetadataState
			PublishAt      uint64
			UpdateAt       uint64
			Nonce          uint64
			MetadataOption string
			// add by v0.5.0
			User                 string
			UserType             constant.UserType
			Sign                 []byte
			*/
			MetadataId:     req.GetInformation().GetMetadataId(),
			MetadataName:   req.GetInformation().GetMetadataName(),
			MetadataType:   req.GetInformation().GetMetadataType(),
			DataHash:       req.GetInformation().GetDataHash(),
			Desc:           req.GetInformation().GetDesc(),
			LocationType:   req.GetInformation().GetLocationType(),
			DataType:       req.GetInformation().GetDataType(),
			Industry:       req.GetInformation().GetIndustry(),
			State:          req.GetInformation().GetState(),
			PublishAt:      req.GetInformation().GetPublishAt(),
			UpdateAt:       req.GetInformation().GetUpdateAt(),
			Nonce:          req.GetInformation().GetNonce(),
			MetadataOption: req.GetInformation().GetMetadataOption(),
			UserType:       req.GetInformation().GetUserType(),
			User:           req.GetInformation().GetUser(),
			Sign:           req.GetInformation().GetSign(),
		},
		CreateAt: timeutils.UnixMsecUint64(),
	}

	return metadataMsg
}

func (msg *MetadataMsg) ToDataCenter(identity *carriertypespb.Organization) *Metadata {
	return NewMetadata(&carriertypespb.MetadataPB{
		/**
		MetadataId     string
		Owner          *Organization
		DataId         string
		DataStatus     constant.DataStatus
		MetadataName   string
		MetadataType   constant.MetadataType
		DataHash       string
		Desc           string
		LocationType   constant.DataLocationTy
		DataType       constant.OrigindataType
		Industry       string
		State          constant.MetadataState
		PublishAt      uint64
		UpdateAt       uint64
		Nonce          uint64
		MetadataOption string
		// add by v0.5.0
		User                 string
		UserType             constant.UserType
		Sign                 []byte
		*/
		MetadataId:     msg.GetMetadataId(),
		Owner:          identity,
		DataId:         msg.GetMetadataId(),
		DataStatus:     commonconstantpb.DataStatus_DataStatus_Valid,
		MetadataName:   msg.GetMetadataName(),
		MetadataType:   msg.GetMetadataType(),
		DataHash:       msg.GetDataHash(),
		Desc:           msg.GetDesc(),
		LocationType:   msg.GetLocationType(),
		DataType:       msg.GetDataType(),
		Industry:       msg.GetIndustry(),
		State:          commonconstantpb.MetadataState_MetadataState_Released, // metaData status, eg: create/release/revoke
		PublishAt:      timeutils.UnixMsecUint64(),
		UpdateAt:       timeutils.UnixMsecUint64(),
		Nonce:          msg.GetNonce(),
		MetadataOption: msg.GetMetadataOption(),
		UserType:       msg.GetUserType(),
		User:           msg.GetUser(),
		Sign:           msg.GetSign(),
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
func (msg *MetadataMsg) GetMetadataSummary() *carriertypespb.MetadataSummary {
	return msg.MetadataSummary
}

func (msg *MetadataMsg) GetMetadataName() string { return msg.GetMetadataSummary().MetadataName }
func (msg *MetadataMsg) GetMetadataType() commonconstantpb.MetadataType {
	return msg.GetMetadataSummary().MetadataType
}
func (msg *MetadataMsg) GetDataHash() string { return msg.GetMetadataSummary().DataHash }
func (msg *MetadataMsg) GetDesc() string     { return msg.GetMetadataSummary().Desc }
func (msg *MetadataMsg) GetLocationType() commonconstantpb.DataLocationType {
	return msg.GetMetadataSummary().LocationType
}
func (msg *MetadataMsg) GetDataType() commonconstantpb.OrigindataType {
	return msg.GetMetadataSummary().DataType
}
func (msg *MetadataMsg) GetNonce() uint64 { return msg.GetMetadataSummary().Nonce }
func (msg *MetadataMsg) GetState() commonconstantpb.MetadataState {
	return msg.GetMetadataSummary().State
}
func (msg *MetadataMsg) GetIndustry() string       { return msg.GetMetadataSummary().Industry }
func (msg *MetadataMsg) GetMetadataOption() string { return msg.GetMetadataSummary().MetadataOption }
func (msg *MetadataMsg) GetCreateAt() uint64       { return msg.CreateAt }
func (msg *MetadataMsg) GetMetadataId() string     { return msg.GetMetadataSummary().MetadataId }
func (msg *MetadataMsg) GetUser() string           { return msg.GetMetadataSummary().User }
func (msg *MetadataMsg) GetUserType() commonconstantpb.UserType {
	return msg.GetMetadataSummary().UserType
}
func (msg *MetadataMsg) GetSign() []byte { return msg.GetMetadataSummary().Sign }
func (msg *MetadataMsg) GenMetadataId() string {
	if "" != msg.GetMetadataId() {
		return msg.GetMetadataId()
	}
	msg.GetMetadataSummary().MetadataId = PREFIX_METADATA_ID + msg.HashByCreateTime().Hex()
	return msg.GetMetadataId()
}
func (msg *MetadataMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	/**
	MetadataName         string
	MetadataType         MetadataType
	DataHash             string
	Desc                 string
	LocationType         DataLocationType
	DataType             OrigindataType
	Industry             string
	State                MetadataState
	MetadataOption       string
	AllowExpose          bool

	*/
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetMetadataName()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetMetadataType())))
	buf.Write([]byte(msg.GetDataHash()))
	buf.Write([]byte(msg.GetDesc()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetLocationType())))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetDataType())))
	buf.Write([]byte(msg.GetIndustry()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetState())))
	//buf.Write(bytesutil.Uint64ToBytes(msg.GetNonce()))
	buf.Write([]byte(msg.GetMetadataOption()))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *MetadataMsg) HashByCreateTime() common.Hash {

	var buf bytes.Buffer
	buf.Write([]byte(msg.GetMetadataName()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetMetadataType())))
	buf.Write([]byte(msg.GetDataHash()))
	buf.Write([]byte(msg.GetDesc()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetLocationType())))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetDataType())))
	buf.Write([]byte(msg.GetIndustry()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetState())))
	//buf.Write(bytesutil.Uint64ToBytes(msg.GetNonce()))
	buf.Write([]byte(msg.GetMetadataOption()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))

	return rlputil.RlpHash(buf.Bytes())
}

type MetadataUpdateMsg struct {
	MetadataSummary *carriertypespb.MetadataSummary `json:"metadataSummary"`
	CreateAt        uint64                          `json:"createAt"`
	// caches
	hash atomic.Value
}

func (msg *MetadataUpdateMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *MetadataUpdateMsg) Unmarshal(b []byte) error { return nil }
func (msg *MetadataUpdateMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *MetadataUpdateMsg) MsgType() string { return MSG_UPDATEMETADATA }
func (msg *MetadataUpdateMsg) GetMetadataSummary() *carriertypespb.MetadataSummary {
	return msg.MetadataSummary
}
func (msg *MetadataUpdateMsg) GetMetadataId() string   { return msg.GetMetadataSummary().GetMetadataId() }
func (msg *MetadataUpdateMsg) GetCreateAt() uint64     { return msg.CreateAt }
func (msg *MetadataUpdateMsg) GetMetadataName() string { return msg.GetMetadataSummary().MetadataName }
func (msg *MetadataUpdateMsg) GetMetadataType() commonconstantpb.MetadataType {
	return msg.GetMetadataSummary().MetadataType
}
func (msg *MetadataUpdateMsg) GetDataHash() string { return msg.GetMetadataSummary().DataHash }
func (msg *MetadataUpdateMsg) GetDesc() string     { return msg.GetMetadataSummary().Desc }
func (msg *MetadataUpdateMsg) GetLocationType() commonconstantpb.DataLocationType {
	return msg.GetMetadataSummary().LocationType
}
func (msg *MetadataUpdateMsg) GetDataType() commonconstantpb.OrigindataType {
	return msg.GetMetadataSummary().DataType
}
func (msg *MetadataUpdateMsg) GetNonce() uint64 { return msg.GetMetadataSummary().Nonce }
func (msg *MetadataUpdateMsg) GetState() commonconstantpb.MetadataState {
	return msg.GetMetadataSummary().State
}
func (msg *MetadataUpdateMsg) GetIndustry() string { return msg.GetMetadataSummary().Industry }
func (msg *MetadataUpdateMsg) GetMetadataOption() string {
	return msg.GetMetadataSummary().MetadataOption
}
func (msg *MetadataUpdateMsg) GetUserType() commonconstantpb.UserType {
	return msg.GetMetadataSummary().UserType
}
func (msg *MetadataUpdateMsg) GetUser() string {
	return msg.GetMetadataSummary().User
}
func (msg *MetadataUpdateMsg) GetPublishAt() uint64 { return msg.GetMetadataSummary().PublishAt }
func (msg *MetadataUpdateMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetMetadataName()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetMetadataType())))
	buf.Write([]byte(msg.GetDataHash()))
	buf.Write([]byte(msg.GetDesc()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetLocationType())))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetDataType())))
	buf.Write([]byte(msg.GetIndustry()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetState())))
	//buf.Write(bytesutil.Uint64ToBytes(msg.GetNonce()))
	buf.Write([]byte(msg.GetMetadataOption()))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

type UpdateIdentityCredentialMsg struct {
	IdentityId string `json:"identity_id"`
	Credential string `json:"credential"`
	CreateAt   uint64 `json:"createAt"`
	// caches
	hash atomic.Value
}

func (msg *UpdateIdentityCredentialMsg) String() string {
	result, err := json.Marshal(msg)
	if err != nil {
		return "Failed to generate string"
	}
	return string(result)
}
func (msg *UpdateIdentityCredentialMsg) Unmarshal(b []byte) error { return nil }
func (msg *UpdateIdentityCredentialMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *UpdateIdentityCredentialMsg) MsgType() string          { return MSG_UPDATE_IDENTITY_CREDENTIALMSG }

type MetadataRevokeMsg struct {
	MetadataId string `json:"metadataId"`
	CreateAt   uint64 `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewMetadataRevokeMessageFromRequest(req *carrierapipb.RevokeMetadataRequest) *MetadataRevokeMsg {
	return &MetadataRevokeMsg{
		MetadataId: req.GetMetadataId(),
		CreateAt:   timeutils.UnixMsecUint64(),
	}
}

func (msg *MetadataRevokeMsg) GetMetadataId() string { return msg.MetadataId }
func (msg *MetadataRevokeMsg) GetCreateAt() uint64   { return msg.CreateAt }
func (msg *MetadataRevokeMsg) ToDataCenter(identity *carriertypespb.Organization) *Metadata {
	return NewMetadata(&carriertypespb.MetadataPB{
		MetadataId: msg.GetMetadataId(),
		Owner:      identity,
		DataId:     msg.GetMetadataId(),
		DataStatus: commonconstantpb.DataStatus_DataStatus_Invalid,
		// metaData status, eg: create/release/revoke
		State:    commonconstantpb.MetadataState_MetadataState_Revoked,
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
func (msg *MetadataRevokeMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	/**
	MetadataId         string
	*/
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetMetadataId()))
	v := rlputil.RlpHash(buf.Bytes())

	msg.hash.Store(v)
	return v
}

type MetadataMsgArr []*MetadataMsg
type MetadataRevokeMsgArr []*MetadataRevokeMsg
type MetadataUpdateMsgArr []*MetadataUpdateMsg

// Len returns the length of s.
func (s MetadataMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// Len returns the length of s.
func (s MetadataUpdateMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataUpdateMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataUpdateMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// Len returns the length of s.
func (s MetadataRevokeMsgArr) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataRevokeMsgArr) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s MetadataRevokeMsgArr) Less(i, j int) bool { return s[i].GetCreateAt() < s[j].GetCreateAt() }

// ------------------- metadata authority apply -------------------

type MetadataAuthorityMsg struct {
	MetadataAuthId string                            `json:"metaDataAuthId"`
	User           string                            `json:"user"`
	UserType       commonconstantpb.UserType         `json:"userType"`
	Auth           *carriertypespb.MetadataAuthority `json:"auth"`
	Sign           []byte                            `json:"sign"`
	CreateAt       uint64                            `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewMetadataAuthorityMessageFromRequest(req *carrierapipb.ApplyMetadataAuthorityRequest) *MetadataAuthorityMsg {
	return &MetadataAuthorityMsg{
		User:     req.GetUser(),
		UserType: req.GetUserType(),
		Auth:     req.GetAuth(),
		Sign:     req.GetSign(),
		CreateAt: timeutils.UnixMsecUint64(),
	}
}

func (msg *MetadataAuthorityMsg) GetMetadataAuthId() string              { return msg.MetadataAuthId }
func (msg *MetadataAuthorityMsg) GetUser() string                        { return msg.User }
func (msg *MetadataAuthorityMsg) GetUserType() commonconstantpb.UserType { return msg.UserType }
func (msg *MetadataAuthorityMsg) GetMetadataAuthority() *carriertypespb.MetadataAuthority {
	return msg.Auth
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwner() *carriertypespb.Organization {
	return msg.Auth.GetOwner()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerIdentityId() string {
	return msg.Auth.GetOwner().GetIdentityId()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerNodeId() string {
	return msg.Auth.GetOwner().GetNodeId()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerNodeName() string {
	return msg.Auth.GetOwner().GetNodeName()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerImageUrl() string {
	return msg.Auth.GetOwner().GetImageUrl()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityOwnerDetails() string {
	return msg.Auth.GetOwner().GetDetails()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityMetadataId() string {
	return msg.Auth.GetMetadataId()
}
func (msg *MetadataAuthorityMsg) GetMetadataAuthorityUsageRule() *carriertypespb.MetadataUsageRule {
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
	/**
	MetadataId
	User
	UserType
	NodeName
	NodeId
	IdentityId
	ImageUrl
	Details
	UsageType
	StartAt     `
	EndAt          `
	Times
	*/
	var buf bytes.Buffer

	buf.Write([]byte(msg.GetMetadataAuthorityMetadataId()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerNodeName()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerNodeId()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerIdentityId()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerImageUrl()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerDetails()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetMetadataAuthorityUsageRule().GetUsageType())))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMetadataAuthorityUsageRule().GetStartAt()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMetadataAuthorityUsageRule().GetEndAt()))
	buf.Write(bytesutil.Uint32ToBytes(msg.GetMetadataAuthorityUsageRule().GetTimes()))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *MetadataAuthorityMsg) HashByCreateTime() common.Hash {

	var buf bytes.Buffer

	buf.Write([]byte(msg.GetMetadataAuthorityMetadataId()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerNodeName()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerNodeId()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerIdentityId()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerImageUrl()))
	buf.Write([]byte(msg.GetMetadataAuthorityOwnerDetails()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetMetadataAuthorityUsageRule().GetUsageType())))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMetadataAuthorityUsageRule().GetStartAt()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetMetadataAuthorityUsageRule().GetEndAt()))
	buf.Write(bytesutil.Uint32ToBytes(msg.GetMetadataAuthorityUsageRule().GetTimes()))
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
	MetadataAuthId string
	User           string
	UserType       commonconstantpb.UserType
	Sign           []byte
	CreateAt       uint64
	// caches
	hash atomic.Value
}

func NewMetadataAuthorityRevokeMessageFromRequest(req *carrierapipb.RevokeMetadataAuthorityRequest) *MetadataAuthorityRevokeMsg {
	return &MetadataAuthorityRevokeMsg{
		MetadataAuthId: req.GetMetadataAuthId(),
		User:           req.GetUser(),
		UserType:       req.GetUserType(),
		Sign:           req.GetSign(),
		CreateAt:       timeutils.UnixMsecUint64(),
	}
}

func (msg *MetadataAuthorityRevokeMsg) GetMetadataAuthId() string              { return msg.MetadataAuthId }
func (msg *MetadataAuthorityRevokeMsg) GetUser() string                        { return msg.User }
func (msg *MetadataAuthorityRevokeMsg) GetUserType() commonconstantpb.UserType { return msg.UserType }
func (msg *MetadataAuthorityRevokeMsg) GetSign() []byte                        { return msg.Sign }
func (msg *MetadataAuthorityRevokeMsg) GetCreateAt() uint64                    { return msg.CreateAt }

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

func (msg *MetadataAuthorityRevokeMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	/**
	MetadataAuthId
	User
	UserType
	*/
	var buf bytes.Buffer

	buf.Write([]byte(msg.GetMetadataAuthId()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

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

type TaskMsg struct {
	Data *Task
	// caches
	hash atomic.Value
}

func NewTaskMessageFromRequest(req *carrierapipb.PublishTaskDeclareRequest) *TaskMsg {
	return &TaskMsg{
		Data: NewTask(&carriertypespb.TaskPB{

			TaskId:        "",
			DataId:        "",
			DataStatus:    commonconstantpb.DataStatus_DataStatus_Valid,
			User:          req.GetUser(),
			UserType:      req.GetUserType(),
			TaskName:      req.GetTaskName(),
			Sender:        req.GetSender(),
			AlgoSupplier:  req.GetAlgoSupplier(),
			DataSuppliers: req.GetDataSuppliers(),
			// PowerSuppliers: ,
			Receivers:                req.GetReceivers(),
			DataPolicyTypes:          req.GetDataPolicyTypes(),
			DataPolicyOptions:        req.GetDataPolicyOptions(),
			PowerPolicyTypes:         req.GetPowerPolicyTypes(),
			PowerPolicyOptions:       req.GetPowerPolicyOptions(),
			ReceiverPolicyTypes:      req.GetReceiverPolicyTypes(),
			ReceiverPolicyOptions:    req.GetReceiverPolicyOptions(),
			DataFlowPolicyTypes:      req.GetDataFlowPolicyTypes(),
			DataFlowPolicyOptions:    req.GetDataFlowPolicyOptions(),
			OperationCost:            req.GetOperationCost(),
			AlgorithmCode:            req.GetAlgorithmCode(),
			MetaAlgorithmId:          req.GetMetaAlgorithmId(),
			AlgorithmCodeExtraParams: req.GetAlgorithmCodeExtraParams(),
			// PowerResourceOptions:
			State:    commonconstantpb.TaskState_TaskState_Pending,
			Reason:   "",
			Desc:     req.GetDesc(),
			CreateAt: timeutils.UnixMsecUint64(),
			EndAt:    0,
			StartAt:  0,
			// TaskEvents:
			Sign: req.GetSign(),
		}),
	}
}

func (msg *TaskMsg) Marshal() ([]byte, error) { return nil, nil }
func (msg *TaskMsg) Unmarshal(b []byte) error { return nil }
func (msg *TaskMsg) String() string {
	return fmt.Sprintf(`{"taskId": %s, "task": %s}`, msg.GetTask().GetTaskId(), msg.GetTask().GetTaskData().String())
}
func (msg *TaskMsg) MsgType() string { return MSG_TASK }

func (msg *TaskMsg) GetTask() *Task                      { return msg.Data }
func (msg *TaskMsg) GetTaskData() *carriertypespb.TaskPB { return msg.GetTask().GetTaskData() }
func (msg *TaskMsg) GetTaskId() string                   { return msg.GetTask().GetTaskData().GetTaskId() }
func (msg *TaskMsg) GetUser() string                     { return msg.GetTask().GetTaskData().GetUser() }
func (msg *TaskMsg) GetUserType() commonconstantpb.UserType {
	return msg.GetTask().GetTaskData().GetUserType()
}
func (msg *TaskMsg) GetTaskName() string { return msg.GetTask().GetTaskData().GetTaskName() }
func (msg *TaskMsg) GetSender() *carriertypespb.TaskOrganization {
	return msg.GetTask().GetTaskSender()
}
func (msg *TaskMsg) GetSenderName() string   { return msg.GetTask().GetTaskSender().GetNodeName() }
func (msg *TaskMsg) GetSenderNodeId() string { return msg.GetTask().GetTaskSender().GetNodeId() }
func (msg *TaskMsg) GetSenderIdentityId() string {
	return msg.GetTask().GetTaskSender().GetIdentityId()
}
func (msg *TaskMsg) GetSenderPartyId() string { return msg.GetTask().GetTaskSender().GetPartyId() }
func (msg *TaskMsg) GetAlgoSupplier() *carriertypespb.TaskOrganization {
	return &carriertypespb.TaskOrganization{
		PartyId:    msg.GetTask().GetTaskData().GetAlgoSupplier().GetPartyId(),
		NodeName:   msg.GetTask().GetTaskData().GetAlgoSupplier().GetNodeName(),
		NodeId:     msg.GetTask().GetTaskData().GetAlgoSupplier().GetNodeId(),
		IdentityId: msg.GetTask().GetTaskData().GetAlgoSupplier().GetIdentityId(),
	}
}
func (msg *TaskMsg) GetDataSuppliers() []*carriertypespb.TaskOrganization {
	return msg.GetTask().GetTaskData().GetDataSuppliers()
}
func (msg *TaskMsg) GetPowerSuppliers() []*carriertypespb.TaskOrganization {
	return msg.Data.GetTaskData().GetPowerSuppliers()
}
func (msg *TaskMsg) GetReceivers() []*carriertypespb.TaskOrganization {
	return msg.Data.GetTaskData().GetReceivers()
}
func (msg *TaskMsg) GetDataPolicyTypes() []uint32 { return msg.Data.GetTaskData().GetDataPolicyTypes() }
func (msg *TaskMsg) GetDataPolicyOptions() []string {
	return msg.Data.GetTaskData().GetDataPolicyOptions()
}
func (msg *TaskMsg) GetPowerPolicyTypes() []uint32 {
	return msg.Data.GetTaskData().GetPowerPolicyTypes()
}
func (msg *TaskMsg) GetPowerPolicyOptions() []string {
	return msg.Data.GetTaskData().GetPowerPolicyOptions()
}
func (msg *TaskMsg) GetReceiverPolicyTypes() []uint32 {
	return msg.Data.GetTaskData().GetReceiverPolicyTypes()
}
func (msg *TaskMsg) GetReceiverPolicyOptions() []string {
	return msg.Data.GetTaskData().GetReceiverPolicyOptions()
}
func (msg *TaskMsg) GetDataFlowPolicyTypes() []uint32 {
	return msg.Data.GetTaskData().GetDataFlowPolicyTypes()
}
func (msg *TaskMsg) GetDataFlowPolicyOptions() []string {
	return msg.Data.GetTaskData().GetDataFlowPolicyOptions()
}
func (msg *TaskMsg) GetOperationCost() *carriertypespb.TaskResourceCostDeclare {
	return msg.Data.GetTaskData().GetOperationCost()
}
func (msg *TaskMsg) GetAlgorithmCode() string   { return msg.Data.GetTaskData().GetAlgorithmCode() }
func (msg *TaskMsg) GetMetaAlgorithmId() string { return msg.Data.GetTaskData().GetMetaAlgorithmId() }
func (msg *TaskMsg) GetAlgorithmCodeExtraParams() string {
	return msg.Data.GetTaskData().GetAlgorithmCodeExtraParams()
}
func (msg *TaskMsg) GetPowerResourceOptions() []*carriertypespb.TaskPowerResourceOption {
	return msg.Data.GetTaskData().GetPowerResourceOptions()
}
func (msg *TaskMsg) GetState() commonconstantpb.TaskState { return msg.Data.GetTaskData().GetState() }
func (msg *TaskMsg) GetReason() string                    { return msg.Data.GetTaskData().GetReason() }
func (msg *TaskMsg) GetDesc() string                      { return msg.Data.GetTaskData().GetDesc() }
func (msg *TaskMsg) GetCreateAt() uint64                  { return msg.GetTask().GetTaskData().GetCreateAt() }
func (msg *TaskMsg) GetEndAt() uint64                     { return msg.GetTask().GetTaskData().GetEndAt() }
func (msg *TaskMsg) GetStartAt() uint64                   { return msg.GetTask().GetTaskData().GetStartAt() }
func (msg *TaskMsg) GetSign() []byte                      { return msg.Data.GetTaskData().GetSign() }

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

	/**
	User                     string
	UserType                 UserType
	TaskName                 string
	Sender.NodeName
	Sender.NodeId
	Sender.IdentityId
	Sender.PartyId
	AlgoSupplier.NodeName
	AlgoSupplier.NodeId
	AlgoSupplier.IdentityId
	AlgoSupplier.PartyId
	len DataSuppliers            []*TaskOrganization
	len PowerSuppliers           []*TaskOrganization
	len Receivers                []*TaskOrganization
	len DataPolicyTypes          uint32
	len DataPolicyOptions        string
	len PowerPolicyTypes         uint32
	len PowerPolicyOptions       string
	len ReceiverPolicyTypes      uint32
	len ReceiverPolicyOptions    string
	len DataFlowPolicyTypes      uint32
	len DataFlowPolicyOptions    string
	OperationCost.Processor
	OperationCost.Bandwidth
	OperationCost.Memory
	OperationCost.Duration
	AlgorithmCode            string
	MetaAlgorithmId          string
	AlgorithmCodeExtraParams string
	len PowerResourceOptions     []*TaskPowerResourceOption
	State                    TaskState
	Desc                     string
	*/

	var buf bytes.Buffer

	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))
	buf.Write([]byte(msg.GetTaskName()))
	buf.Write([]byte(msg.GetSender().GetNodeName()))
	buf.Write([]byte(msg.GetSender().GetNodeId()))
	buf.Write([]byte(msg.GetSender().GetIdentityId()))
	buf.Write([]byte(msg.GetSender().GetPartyId()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetNodeName()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetNodeId()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetIdentityId()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetPartyId()))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetDataSuppliers()))))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetPowerSuppliers()))))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetReceivers()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetPowerPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetPowerPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetReceiverPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetReceiverPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataFlowPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataFlowPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(msg.GetOperationCost().GetProcessor()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetOperationCost().GetBandwidth()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetOperationCost().GetMemory()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetOperationCost().GetDuration()))
	buf.Write([]byte(msg.GetAlgorithmCode()))
	buf.Write([]byte(msg.GetMetaAlgorithmId()))
	buf.Write([]byte(msg.GetAlgorithmCodeExtraParams()))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetPowerResourceOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetState())))
	buf.Write([]byte(msg.GetDesc()))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *TaskMsg) HashByCreateTime() common.Hash {

	var buf bytes.Buffer

	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))
	buf.Write([]byte(msg.GetTaskName()))
	buf.Write([]byte(msg.GetSender().GetNodeName()))
	buf.Write([]byte(msg.GetSender().GetNodeId()))
	buf.Write([]byte(msg.GetSender().GetIdentityId()))
	buf.Write([]byte(msg.GetSender().GetPartyId()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetNodeName()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetNodeId()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetIdentityId()))
	buf.Write([]byte(msg.GetAlgoSupplier().GetPartyId()))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetDataSuppliers()))))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetPowerSuppliers()))))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetReceivers()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetPowerPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetPowerPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetReceiverPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetReceiverPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataFlowPolicyTypes()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(len(msg.GetDataFlowPolicyOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(msg.GetOperationCost().GetProcessor()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetOperationCost().GetBandwidth()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetOperationCost().GetMemory()))
	buf.Write(bytesutil.Uint64ToBytes(msg.GetOperationCost().GetDuration()))
	buf.Write([]byte(msg.GetAlgorithmCode()))
	buf.Write([]byte(msg.GetMetaAlgorithmId()))
	buf.Write([]byte(msg.GetAlgorithmCodeExtraParams()))
	buf.Write(bytesutil.Uint16ToBytes(uint16(len(msg.GetPowerResourceOptions()))))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetState())))
	buf.Write([]byte(msg.GetDesc()))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))

	return rlputil.RlpHash(buf.Bytes())
}

type TaskTerminateMsg struct {
	TaskId   string                    `json:"taskId"`
	User     string                    `json:"user"`
	UserType commonconstantpb.UserType `json:"userType"`
	Sign     []byte                    `json:"sign"`
	CreateAt uint64                    `json:"createAt"`
	// caches
	hash atomic.Value
}

func NewTaskTerminateMsg(userType commonconstantpb.UserType, user, taskId string, sign []byte) *TaskTerminateMsg {
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
func (msg *TaskTerminateMsg) MsgType() string                        { return MSG_TASK_TERMINATE }
func (msg *TaskTerminateMsg) GetUserType() commonconstantpb.UserType { return msg.UserType }
func (msg *TaskTerminateMsg) GetUser() string                        { return msg.User }
func (msg *TaskTerminateMsg) GetTaskId() string                      { return msg.TaskId }
func (msg *TaskTerminateMsg) GetSign() []byte                        { return msg.Sign }
func (msg *TaskTerminateMsg) GetCreateAt() uint64                    { return msg.CreateAt }
func (msg *TaskTerminateMsg) Hash() common.Hash {
	if hash := msg.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	/**
	TaskId,
	User,
	UserType,
	*/

	var buf bytes.Buffer

	buf.Write([]byte(msg.GetTaskId()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))

	v := rlputil.RlpHash(buf.Bytes())
	msg.hash.Store(v)
	return v
}

func (msg *TaskTerminateMsg) HashByCreateTime() common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(msg.GetTaskId()))
	buf.Write([]byte(msg.GetUser()))
	buf.Write(bytesutil.Uint32ToBytes(uint32(msg.GetUserType())))
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

func (b *TaskBullet) IsNewBullet() bool {

	var flag int

	if 0 != b.GetTerm() {
		flag |= 1
	}

	if 0 != b.GetResched() {
		flag |= 1
	}

	if 0 != flag {
		return false
	}

	return true
}
func (b *TaskBullet) IsNotNewBullet() bool { return !b.IsNewBullet() }

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

/**
Example:
{
  "self_cfg_params": {
      "party_id": "data1",    # 本方party_id
      "input_data": [
        {
            "input_type": 1,     # 输入数据的类型. 0: unknown, 1: origin_data, 2: psi_output 3: model
            "data_type": 1,      # 数据的格式. 0:unknown, 1:csv, 2:folder, 3:xls, 4:txt, 5:json, 6:mysql, 7:bin
            "data_path": "path/to/data",  # 数据所在的本地路径
            "key_column": "col1",  # ID列名
            "selected_columns": ["col2", "col3"]  # 自变量(特征列名)
        },
        {
            "input_type": 2,     # 输入数据的类型. 0: unknown, 1: origin_data, 2: psi_output 3: model
            "data_type": 1,
            "data_path": "path/to/data1/psi_result.csv",
            "key_column": "",
            "selected_columns": []
        }
      ]
  },
  "algorithm_dynamic_params": {
      "use_psi": true,           # 是否使用psi
      "label_owner": "data1",       # 标签所在方的party_id
      "label_column": "Y",       # 因变量(标签)
      "hyperparams": {           # 逻辑回归的超参数
          "epochs": 10,            # 训练轮次，大于0的整数
          "batch_size": 256,       # 批量大小，大于0的整数
          "learning_rate": 0.1,    # 学习率，，大于0的数
          "use_validation_set": true,  # 是否使用验证集，true-用，false-不用
          "validation_set_rate": 0.2,  # 验证集占输入数据集的比例，值域(0,1)
          "predict_threshold": 0.5     # 验证集预测结果的分类阈值，值域[0,1]
      }
  }
}
*/

type SelfCfgParams struct {
	PartyId   string      `json:"party_id"`
	InputData interface{} `json:"input_data""`
}

/**
{
    "input_type": 3,  # 输入数据的类型，(算法用标识数据使用方式). 0:unknown, 1:origin_data, 2:psi_output, 3:model
    "access_type": 1, # 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local, 2:url
    "data_type": 0,   # 数据的格式，(算法用标识数据格式). 0:unknown, 1:csv, 2:dir, 3:binary, 4:xls, 5:xlsx, 6:txt, 7:json
    "data_path": "/task_result/task:0xdeefff3434..556/"  # 数据所在的本地路径
}
*/
type InputDataDIR struct {
	InputType  uint32 `json:"input_type"`
	AccessType uint32 `json:"access_type"` // 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local <default>, 2:http, 3:https, 4:ftp
	DataType   uint32 `json:"data_type"`
	DataPath   string `json:"data_path"`
}

/**
{
    "input_type": 3,  # 输入数据的类型，(算法用标识数据使用方式). 0:unknown, 1:origin_data, 2:psi_output, 3:model
    "access_type": 1, # 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local, 2:url
    "data_type": 0,   # 数据的格式，(算法用标识数据格式). 0:unknown, 1:csv, 2:dir, 3:binary, 4:xls, 5:xlsx, 6:txt, 7:json
    "data_path": "/task_result/task:0xdeefff3434..556/"  # 数据所在的本地路径
}
*/
type InputDataBINARY struct {
	InputType  uint32 `json:"input_type"`
	AccessType uint32 `json:"access_type"` // 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local <default>, 2:http, 3:https, 4:ftp
	DataType   uint32 `json:"data_type"`
	DataPath   string `json:"data_path"`
}

/**
{
   "input_type": 1,  # 输入数据的类型，(算法用标识数据使用方式). 0:unknown, 1:origin_data, 2:psi_output, 3:model
   "access_type": 1, # 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local, 2:url
   "data_type": 0,   # 数据的格式，(算法用标识数据格式). 0:unknown, 1:csv, 2:binary, 3:dir, 4:xls, 5:xlsx, 6:txt, 7:json
   "data_path": "/metadata/20220427_预测银行数据.csv",  # 数据所在的本地路径
   "key_column": "col1",  # ID列名
   "selected_columns": ["col2", "col3"]  # 自变量(特征列名)
}
*/
type InputDataCSV struct {
	InputType       uint32   `json:"input_type"`
	AccessType      uint32   `json:"access_type"` // 访问数据的方式，(fighter用决定是否预先加载数据). 0:unknown, 1:local <default>, 2:http, 3:https, 4:ftp
	DataType        uint32   `json:"data_type"`
	DataPath        string   `json:"data_path"`
	KeyColumn       string   `json:"key_column"`
	SelectedColumns []string `json:"selected_columns"`
}
