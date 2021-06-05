package types

import libTypes "github.com/RosettaFlow/Carrier-Go/lib/types"

// MetadataArray is a Transaction slice type for basic sorting.
type MetadataArray []*libTypes.MetaData

// Len returns the length of s.
func (s MetadataArray) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s MetadataArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s MetadataArray) GetPb(i int) []byte {
	enc, _ := s[i].Marshal()
	return enc
}

func (s MetadataArray) Copy() []libTypes.MetaData {
	result := make([]libTypes.MetaData, 0, 1)
	for _, v := range s {
		cpy := *v
		result = append(result, cpy)
	}
	return result
}

func (s MetadataArray) BuildFrom(metadata []libTypes.MetaData) []*libTypes.MetaData {
	result := make([]*libTypes.MetaData, 0)
	for _, meta := range metadata {
		result = append(result, &meta)
	}
	copy(s, result)
	return result
}

