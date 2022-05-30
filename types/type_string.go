package types

import (
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"strings"
)

func UtilOrgPowerArrString(powers []*carriertypespb.TaskOrganization) string {
	arr := make([]string, len(powers))
	for i, power := range powers {
		arr[i] = power.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func UtilOrgPowerResourceArrString(powers []*carriertypespb.TaskPowerResourceOption) string {
	arr := make([]string, len(powers))
	for i, power := range powers {
		arr[i] = power.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}
func UtilLocalResourceArrString(resources []*LocalResourceTable) string {
	arr := make([]string, len(resources))
	for i, r := range resources {
		arr[i] = r.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func UtilDataResourceArrString(resources []*DataResourceTable) string {
	arr := make([]string, len(resources))
	for i, r := range resources {
		arr[i] = r.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

