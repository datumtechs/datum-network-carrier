package types

import libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"

// ResourceArray is a Transaction slice type for basic sorting.
type ResourceArray []*libTypes.ResourceData

// Len returns the length of s.
func (s ResourceArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s ResourceArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ResourceArray) GetPb(i int) []byte {
	enc, _ := s[i].Marshal()
	return enc
}

