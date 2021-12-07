package types

import (
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"strings"
)

func UtilOrgPowerArrString(powers []*libtypes.TaskPowerSupplier) string {
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

