package stringutil

import "strconv"

func StringToUInt64(str string) uint64 {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return uint64(i)
}
