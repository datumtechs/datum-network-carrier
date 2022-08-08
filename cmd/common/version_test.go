package common

import (
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"math"
	"testing"
)

func TestVersion_DefaultValue(t *testing.T) {
	v := Version()
	t.Logf("version:%s", v)
	assert.Equal(t, "Moments ago", v)
}

func TestVersion_SetValue(t *testing.T) {
	buildDate = "{DATE}"
	v := Version()
	t.Logf("version:%s", v)
	require.NotEqual(t, "Moments ago", v)
}

func TestG(t *testing.T) {
	t.Log(0xFFFF)
	t.Log(math.MaxUint16)
	t.Log(0xFFFF << 8)
	t.Log(0xFF << 8)
	t.Log(0xFF<<8 + 0xFF)

	oldStatus := uint16(commonconstantpb.MetadataAuthorityState_MAState_Invalid)
	oldStatus |= uint16(commonconstantpb.AuditMetadataOption_Audit_Passed << 8)

	t.Log(oldStatus &^ (0xFF << 8))

	t.Log((oldStatus &^ 0xFF) >> 8)

	t.Log(uint16(commonconstantpb.MetadataAuthorityState_MAState_Invalid)&(oldStatus&^(0xFF<<8)) == uint16(commonconstantpb.MetadataAuthorityState_MAState_Invalid))
	t.Log(uint16(commonconstantpb.AuditMetadataOption_Audit_Passed)&((oldStatus&^0xFF)>>8) == uint16(commonconstantpb.AuditMetadataOption_Audit_Passed))

	t.Log(oldStatus&^(0xFF<<8) == uint16(commonconstantpb.MetadataAuthorityState_MAState_Invalid))
	t.Log((oldStatus&^0xFF)>>8 == uint16(commonconstantpb.AuditMetadataOption_Audit_Passed))

}
