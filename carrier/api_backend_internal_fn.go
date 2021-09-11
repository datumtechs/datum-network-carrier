package carrier

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/types"
	"strings"
)

func utilLocalTaskPowerUsedArrString(used []*types.LocalTaskPowerUsed) string {
	arr := make([]string, len(used))
	for i, u := range used {
		arr[i] = u.String()
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

func utilLocalTaskPowerUsedMapString(taskPowerUsedMap map[string][]*types.LocalTaskPowerUsed) string {
	arr := make([]string, 0)
	for jobNodeId, useds := range taskPowerUsedMap {
		arr = append(arr, fmt.Sprintf(`{"%s": %s}`, jobNodeId, utilLocalTaskPowerUsedArrString(useds)))
	}
	if len(arr) != 0 {
		return "[" + strings.Join(arr, ",") + "]"
	}
	return "[]"
}

