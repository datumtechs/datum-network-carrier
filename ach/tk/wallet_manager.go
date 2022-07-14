package tk

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/datumtechs/datum-network-carrier/ach/tk/kms"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
)

type WalletManager struct {
	dataCenter carrierdb.CarrierDB
	kms        kms.KmsService
}

func NewWalletManager(db carrierdb.CarrierDB, kmsConfig *kms.Config) *WalletManager {
	if kmsConfig != nil {
		kmsService := kms.NewAliKms(kmsConfig)
		return &WalletManager{db, kmsService}
	} else {
		return &WalletManager{db, nil}
	}
}

// FindOrgPriKey returns the organization wallet address
func (m *WalletManager) GetOrgWalletAddress() (*common.Address, error) {
	// datacenter存储的是加密后的私钥
	priKey := m.LoadPrivateKey()
	if priKey == nil {
		return nil, errors.New("failed to query organization private key")
	}
	addr := crypto.PubkeyToAddress(priKey.PublicKey)
	return &addr, nil
}

// GenerateOrgWallet generates a new wallet if there's no wallet for current organization, if there is an organization wallet already, just reuse it.
func (m *WalletManager) GenerateOrgWallet() {
	priKey := m.LoadPrivateKey()
	if priKey != nil {
		log.Warnf("organization wallet exists, address: %s", crypto.PubkeyToAddress(priKey.PublicKey))
		return
	}

	key, _ := crypto.GenerateKey()
	keyHex := hex.EncodeToString(crypto.FromECDSA(key))
	if m.kms != nil {
		if cipher, err := m.kms.Encrypt(keyHex); err != nil {
			log.WithError(err).Error("cannot encrypt organization wallet private key")
			return
		} else {
			keyHex = cipher
		}
	}

	// datacenter存储的是加密后的私钥
	if err := m.dataCenter.SaveOrgPriKey(keyHex); err != nil {
		log.WithError(err).Error("Failed to store organization wallet")
		return
	}
	log.Warnf("generate organization wallet successful, address: %s", crypto.PubkeyToAddress(key.PublicKey))
}

func (m *WalletManager) LoadPrivateKey() *ecdsa.PrivateKey {
	// datacenter存储的是加密后的私钥
	priKeyHex, err := m.dataCenter.FindOrgPriKey()
	if nil != err {
		log.WithError(err).Error("failed to query organization wallet. ", err)
		return nil
	}

	if m.kms != nil {
		if key, err := m.kms.Decrypt(priKeyHex); err != nil {
			log.Errorf("decrypt organization wallet private key error: %v", err)
			return nil
		} else {
			priKey, err := crypto.HexToECDSA(key)
			if err != nil {
				log.Errorf("convert organization wallet private key to ECDSA error: %v", err)
				return nil
			} else {
				return priKey
			}
		}
	} else {
		priKey, err := crypto.HexToECDSA(priKeyHex)
		if err != nil {
			log.Errorf("not a valid private key hex: %v", err)
			return nil
		}
		return priKey
	}
}
