package metadata

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/types"
)

type MetadataAuthority struct {
	dataCenter core.CarrierDB
}

func NewMetadataAuthority(dataCenter core.CarrierDB) *MetadataAuthority {
	return &MetadataAuthority{
		dataCenter: dataCenter,
	}
}

func (ma *MetadataAuthority) ApplyMetadataAuthority(metadataAuth *types.MetadataAuthority) error {
	return ma.dataCenter.InsertMetadataAuthority(metadataAuth)
}

func (ma *MetadataAuthority) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (apicommonpb.AuditMetadataOption, error) {

	// verify
	metadataAuth, err := ma.GetMetadataAuthority(audit.GetMetadataAuthId())
	if nil != err {
		log.Errorf("Failed to query old metadataAuth on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			audit.GetMetadataAuthId(), err)
		return apicommonpb.AuditMetadataOption_Audit_Pending, err
	}

	if metadataAuth.GetData().State != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth state is not release")
	}

	if metadataAuth.GetData().GetAuditOption() != apicommonpb.AuditMetadataOption_Audit_Pending {
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}",
			audit.GetMetadataAuthId())
		return metadataAuth.GetData().GetAuditOption(), fmt.Errorf("the old metadataAuth has already audited")
	}

	metadataAuth.GetData().AuditOption = audit.GetAuditOption()
	metadataAuth.GetData().AuditSuggestion = audit.GetAuditSuggestion()

	err = ma.dataCenter.InsertMetadataAuthority(metadataAuth)
	if nil != err {
		log.Errorf("Failed to store metadataAuth after audit on MetadataAuthority.AuditMetadataAuthority(), metadataAuthId: {%s}, audit option:{%s}, err: {%s}",
			audit.GetMetadataAuthId(), audit.GetAuditOption().String(), err)
	}

	return metadataAuth.GetData().GetAuditOption(), err
}

func (ma *MetadataAuthority) ConsumeMetadataAuthority(metadataAuthId string) error {

	// verify
	metadataAuth, err := ma.GetMetadataAuthority(metadataAuthId)
	if nil != err {
		log.Errorf("Failed to query old metadataAuth on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			metadataAuthId, err)
		return err
	}

	if metadataAuth.GetData().State != apicommonpb.MetadataAuthorityState_MAState_Released {
		log.Errorf("the old metadataAuth state is not release on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("the old metadataAuth state is not release")
	}

	if metadataAuth.GetData().GetAuditOption() == apicommonpb.AuditMetadataOption_Audit_Refused {
		log.Errorf("the old metadataAuth has already audited on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("the old metadataAuth was refused")
	}

	usageRule := metadataAuth.GetData().GetAuth().GetUsageRule()
	usedQuo := metadataAuth.GetData().GetUsedQuo()

	switch usageRule.UsageType {
	case apicommonpb.MetadataUsageType_Usage_Period:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Period
		if uint64(timeutils.UnixMsec()) >= usageRule.GetEndAt() {
			usedQuo.Expire = true
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		} else {
			usedQuo.Expire = false
		}
	case apicommonpb.MetadataUsageType_Usage_Times:
		usedQuo.UsageType = apicommonpb.MetadataUsageType_Usage_Times
		if usedQuo.GetUsedTimes() < usageRule.GetTimes() {
			usedQuo.UsedTimes += 1
		} else {
			metadataAuth.GetData().State = apicommonpb.MetadataAuthorityState_MAState_Invalid
		}
	default:
		log.Errorf("usageRule state of the old metadataAuth is invalid on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}",
			metadataAuthId)
		return fmt.Errorf("usageRule state of the old metadataAuth is invalid")
	}


	metadataAuth.GetData().UsedQuo = usedQuo

	err = ma.dataCenter.InsertMetadataAuthority(metadataAuth)
	if nil != err {
		log.Errorf("Failed to update metadataAuth after consume on MetadataAuthority.ConsumeMetadataAuthority(), metadataAuthId: {%s}, err: {%s}",
			metadataAuthId, err)
		return err
	}
	return nil
}

func (ma *MetadataAuthority) GetMetadataAuthority(metadataAuthId string) (*types.MetadataAuthority, error) {
	return ma.dataCenter.GetMetadataAuthority(metadataAuthId)
}

func (ma *MetadataAuthority) GetMetadataAuthorityList() (types.MetadataAuthArray, error) {
	identityId, err := ma.dataCenter.GetIdentityId()
	if nil != err {
		return nil, err
	}
	return ma.dataCenter.GetMetadataAuthorityListByIdentityId(identityId, uint64(timeutils.UnixMsec()))
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByIds(metadataAuthIds []string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.GetMetadataAuthorityListByIds(metadataAuthIds)
}

func (ma *MetadataAuthority) GetMetadataAuthorityListByUser(userType apicommonpb.UserType, user string) (types.MetadataAuthArray, error) {
	return ma.dataCenter.GetMetadataAuthorityListByUser(userType, user, uint64(timeutils.UnixMsec()))
}

func (ma *MetadataAuthority) VerifyMetadataAuth (user, metadataId string, userType apicommonpb.UserType) bool {

	// TODO 检查 user  和 metadataId 的关联关系

	// 获取 metadata
	metadata, err := ma.dataCenter.GetMetadataByDataId(metadataId)
	if nil != err {

	}

	// verify
	metadataAuth, err := ma.GetMetadataAuthority(metadataId)
	if nil != err {
		log.Errorf("Failed to GetMetadataAuthority on VerifyMetadataAuth, metadataAuthId: {%s}, userType: {%s}, user: {%s}, err: {%s}",
			metadataId, userType.String(), user, err)
		return false
	}

	metadataAuth.

}