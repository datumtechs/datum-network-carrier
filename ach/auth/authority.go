package auth

import (
	"fmt"
	metadata2 "github.com/datumtechs/datum-network-carrier/ach/auth/metadata"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/datumtechs/datum-network-carrier/rpc/backend"
	"github.com/datumtechs/datum-network-carrier/types"
	"time"
)

type AuthorityManager struct {
	metadataAuth     *metadata2.MetadataAuthority
	quit  chan struct{}
}

func NewAuthorityManager(dataCenter carrierdb.CarrierDB) *AuthorityManager {
	return &AuthorityManager{
		metadataAuth: metadata2.NewMetadataAuthority(dataCenter),
		quit:         make(chan struct{}),
	}
}

func (am *AuthorityManager) Start() error {
	//go am.loop()
	log.Info("Started authorityManager ...")
	return nil
}

func (am *AuthorityManager) Stop() error {
	close(am.quit)
	return nil
}

func (am *AuthorityManager) loop () {
	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			am.refreshMetadataAuthority()
		case <-am.quit:
			log.Info("Stopped AuthorityManager ...")
			ticker.Stop()
			return
		}
	}
}

func (am *AuthorityManager) refreshMetadataAuthority () {
	list, err := am.metadataAuth.GetLocalMetadataAuthorityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		return
	}

	//log.Debugf("Started call AuthorityManager.refreshMetadataAuthority()")

	for _, metadataAuth := range list {

		// Regularly check the validity of metadata auth information in 'pending' status,
		// and decide whether to automatically issue 'refused' audit suggestions.
		if metadataAuth.GetData().GetAuditOption() == commonconstantpb.AuditMetadataOption_Audit_Pending {
			var invalid bool

			switch metadataAuth.GetData().GetAuth().GetUsageRule().GetUsageType() {
			case commonconstantpb.MetadataUsageType_Usage_Period:
				if timeutils.UnixMsecUint64() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetEndAt() {
					metadataAuth.GetData().GetUsedQuo().Expire = true
					metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
					// refuse it for audit suggestion.
					// update audit things.
					metadataAuth.GetData().AuditOption = commonconstantpb.AuditMetadataOption_Audit_Refused
					metadataAuth.GetData().AuditSuggestion = "metadataAuth has expired, refused it"
					metadataAuth.GetData().AuditAt = timeutils.UnixMsecUint64()
					invalid = true
				}
			case commonconstantpb.MetadataUsageType_Usage_Times:
				if metadataAuth.GetData().GetUsedQuo().GetUsedTimes() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetTimes() {
					metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid

					// refuse it for audit suggestion.
					// update audit things.
					metadataAuth.GetData().AuditOption = commonconstantpb.AuditMetadataOption_Audit_Refused
					metadataAuth.GetData().AuditSuggestion = "metadataAuth has no enough remain times, refused it"
					metadataAuth.GetData().AuditAt = timeutils.UnixMsecUint64()
					invalid = true
				}
			default:
				log.Errorf("unknown usageType of the old metadataAuth on AuthorityManager.refreshMetadataAuthority(), metadataAuthId: {%s}", metadataAuth.GetData().GetMetadataAuthId())
				continue
			}

			if invalid {

				// update the metadataAuth when it was refused audit.
				if err := am.metadataAuth.UpdateMetadataAuthority(metadataAuth); nil != err {
					log.WithError(err).Errorf("Failed to update metadataAuth after audit on MetadataAuthority.refreshMetadataAuthority(), metadataAuthId: {%s}, audit option:{%s}",
						metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuditOption().String())
				}
				// remove the invaid metadataAuthId from local db
				if err := am.metadataAuth.RemoveUserMetadataAuthIdByMetadataId(metadataAuth.GetUserType(), metadataAuth.GetUser(), metadataAuth.GetData().GetAuth().GetMetadataId()); nil != err {
					log.WithError(err).Errorf("Failed to remove metadataId and metadataAuthId mapping while metadataAuth has invalid on MetadataAuthority.refreshMetadataAuthority(), metadataAuthId: {%s}, metadataId: {%s}, userType: {%s}, user:{%s}",
						metadataAuth.GetData().GetMetadataAuthId(), metadataAuth.GetData().GetAuth().GetMetadataId(), metadataAuth.GetUserType(), metadataAuth.GetUser())
				}
			}
		}
	}
}

func (am *AuthorityManager) ApplyMetadataAuthority (metadataAuth *types.MetadataAuthority) error {
	return am.metadataAuth.ApplyMetadataAuthority(metadataAuth)
}

func (am *AuthorityManager) AuditMetadataAuthority (audit *types.MetadataAuthAudit) (commonconstantpb.AuditMetadataOption, error) {
	return am.metadataAuth.AuditMetadataAuthority(audit)
}

func (am *AuthorityManager) ConsumeMetadataAuthority (metadataAuthId string) error {
	return am.metadataAuth.ConsumeMetadataAuthority(metadataAuthId)
}

func filterMetadataAuth (list types.MetadataAuthArray) (types.MetadataAuthArray, error) {
	for i, metadataAuth := range list {
		switch metadataAuth.GetData().GetAuth().GetUsageRule().GetUsageType() {
		case commonconstantpb.MetadataUsageType_Usage_Period:
			if timeutils.UnixMsecUint64() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetEndAt() {
				metadataAuth.GetData().GetUsedQuo().Expire = true
				metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
			}
		case commonconstantpb.MetadataUsageType_Usage_Times:
			if metadataAuth.GetData().GetUsedQuo().GetUsedTimes() >= metadataAuth.GetData().GetAuth().GetUsageRule().GetTimes() {
				metadataAuth.GetData().State = commonconstantpb.MetadataAuthorityState_MAState_Invalid
			}
		default:
			log.Errorf("unknown usageType of the old metadataAuth on AuthorityManager.filterMetadataAuth(), metadataAuthId: {%s}", metadataAuth.GetData().GetMetadataAuthId())
			return nil, fmt.Errorf("unknown usageType of the old metadataAuth")
		}

		list[i] = metadataAuth
	}
	return list, nil
}

func (am *AuthorityManager) GetMetadataAuthority (metadataAuthId string) (*types.MetadataAuthority, error) {
	//metadataAuth, err := am.metadataAuth.GetMetadataAuthority(metadataAuthId)
	//if nil != err {
	//	return nil, err
	//}
	//list , err := filterMetadataAuth(types.MetadataAuthArray{metadataAuth})
	//if nil != err {
	//	return nil, err
	//}
	//return list[0], nil

	return am.metadataAuth.GetMetadataAuthority(metadataAuthId)
}

func (am *AuthorityManager) GetLocalMetadataAuthorityList (lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	//list, err := am.metadataAuth.GetLocalMetadataAuthorityList()
	//if nil != err {
	//	return nil, err
	//}
	//return filterMetadataAuth(list)

	return am.metadataAuth.GetLocalMetadataAuthorityList(lastUpdate, pageSize)
}

func (am *AuthorityManager) GetGlobalMetadataAuthorityList (lastUpdate uint64, pageSize uint64) (types.MetadataAuthArray, error) {
	//list, err := am.metadataAuth.GetGlobalMetadataAuthorityList()
	//if nil != err {
	//	return nil, err
	//}
	//return filterMetadataAuth(list)
	return am.metadataAuth.GetGlobalMetadataAuthorityList(lastUpdate, pageSize)
}

func (am *AuthorityManager) GetMetadataAuthorityListByIds (metadataAuthIds  []string) (types.MetadataAuthArray, error) {
	return am.metadataAuth.GetMetadataAuthorityListByIds(metadataAuthIds)
}

func (am *AuthorityManager)  HasValidMetadataAuth(userType commonconstantpb.UserType, user, identityId, metadataId string) (bool, error) {
	return am.metadataAuth.HasValidMetadataAuth(userType, user, identityId, metadataId)
}

func (am *AuthorityManager) VerifyMetadataAuth (userType commonconstantpb.UserType, user, metadataId string) error {
	return am.metadataAuth.VerifyMetadataAuth(userType, user, metadataId)
}

func  (am *AuthorityManager) QueryMetadataAuthIdByMetadataId(userType commonconstantpb.UserType, user, metadataId string) (string, error) {
	return am.metadataAuth.QueryMetadataAuthIdByMetadataId(userType, user, metadataId)
}

