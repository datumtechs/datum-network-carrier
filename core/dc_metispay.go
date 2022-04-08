package core

import (
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/types"
)

func (dc *DataCenter) StoreOrgWallet(sysWallet *types.OrgWallet) error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.StoreOrgWallet(dc.db, sysWallet)
}

func (dc *DataCenter) QueryOrgWallet() (*types.OrgWallet, error) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	return rawdb.QueryOrgWallet(dc.db)
}
