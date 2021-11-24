package mathutil


func Min3number(a, b, c uint64) uint64 {
	var min uint64
	if a > b {
		min = b
	} else {
		min = a
	}

	if min > c {
		min = c
	}
	return min
}

func Max3number(a, b, c uint64) uint64 {
	var max uint64
	if a > b {
		max = a
	} else {
		max = b
	}

	if max < c {
		max = c
	}
	return max
}


func DivCeilUint64(a, b uint64) uint64 {
	div := a / b
	mod := a % b

	if mod > 0 {
		div += 1
	}
	return div
}

func DivCeilUint32(a, b uint32) uint32 {
	div := a / b
	mod := a % b

	if mod > 0 {
		div += 1
	}
	return div
}