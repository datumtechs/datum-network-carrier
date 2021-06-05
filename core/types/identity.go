package types

import libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"

// IdentityArray is a Transaction slice type for basic sorting.
type IdentityArray []*libTypes.IdentityData

// Len returns the length of s.
func (s IdentityArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s IdentityArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s IdentityArray) GetPb(i int) []byte {
	enc, _ := s[i].Marshal()
	return enc
}

