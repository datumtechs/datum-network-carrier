// Copyright (C) 2021 The RosettaNet Authors.

package rawdb

import "github.com/RosettaFlow/Carrier-Go/common"

// ReadMetaDataId retrieves the dataId assigned to a canonical number and index.
func ReadMetaDataId(db DatabaseReader, number uint64, index uint64) string {
	data, _ := db.Get(metadataIdKey(number, index))
	if len(data) == 0 {
		return ""
	}
	return common.Bytes2Hex(data)
}

// WriteMetadataId stores the dataId assigned to a canonical number and index.
func WriteMetadataId(db DatabaseWriter, dataId string, number uint64, index uint64)  {
	if err := db.Put(metadataIdKey(number, index), common.Hex2Bytes(dataId)); err != nil {
		log.WithError(err).Error("Failed to store number-index to dataId mapping")
	}
}

// DeleteMetadataId removes the number-index to dataId canonical mapping.
func DeleteMetadataId(db DatabaseDeleter, number uint64, index uint64) {
	if err := db.Delete(metadataIdKey(number, index)); err != nil {
		log.WithError(err).Error("Failed to delete number-index to dataId mapping")
	}
}