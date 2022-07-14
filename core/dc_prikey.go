package core

import (
	"github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
)

func (dc *DataCenter) SaveOrgPriKey(priKey string) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.SaveOrgPriKey(dc.db, priKey)
}

// FindOrgPriKey does not return ErrNotFound if the organization private key not found.
func (dc *DataCenter) FindOrgPriKey() (string, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.FindOrgPriKey(dc.db)
}
