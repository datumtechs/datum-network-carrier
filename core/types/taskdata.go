package types

import libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"

// TaskDataArray is a Transaction slice type for basic sorting.
type TaskDataArray []*libTypes.TaskData

// Len returns the length of s.
func (s TaskDataArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s TaskDataArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s TaskDataArray) GetPb(i int) []byte {
	enc, _ := s[i].Marshal()
	return enc
}

