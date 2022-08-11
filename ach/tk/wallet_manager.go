package tk

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/datumtechs/datum-network-carrier/ach/tk/kms"
	"github.com/datumtechs/datum-network-carrier/carrierdb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	"sync"
)

var (
	walletManagerOnce sync.Once
)

var walletManager *WalletManager

type WalletManager struct {
	dataCenter    carrierdb.CarrierDB
	kms           kms.KmsService
	priKey        *ecdsa.PrivateKey
	pubKey        *ecdsa.PublicKey
	walletAddress common.Address
}

func WalletManagerInstance() *WalletManager {
	return walletManager
}

func InitWalletManager(db carrierdb.CarrierDB, kmsConfig *kms.Config) {
	walletManagerOnce.Do(func() {
		if kmsConfig != nil {
			kmsService := kms.NewAliKms(kmsConfig)
			walletManager = &WalletManager{dataCenter: db, kms: kmsService}
		} else {
			walletManager = &WalletManager{dataCenter: db, kms: nil}
		}

		walletManager.loadPrivateKey()
	})
}

// GenerateWallet generates a new wallet if there's no wallet for current organization, if there is an organization wallet already, just reuse it.
func (m *WalletManager) GenerateWallet() (common.Address, error) {
	m.loadPrivateKey()
	if m.priKey != nil {
		log.Warnf("organization wallet exists, address: %s", m.walletAddress)
		return m.walletAddress, nil
	}

	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	keyHex := hex.EncodeToString(crypto.FromECDSA(key))
	if m.kms != nil {
		if cipher, err := m.kms.Encrypt(keyHex); err != nil {
			return common.Address{}, err
		} else {
			keyHex = cipher
		}
	}

	// datacenter存储的是加密后的私钥
	if err := m.dataCenter.SaveOrgPriKey(keyHex); err != nil {
		return common.Address{}, err
	}

	m.priKey = key
	m.pubKey = &key.PublicKey
	m.walletAddress = addr

	log.Infof("generate organization wallet successful, address: %s", m.walletAddress)
	return addr, nil
}

// GetAddress returns the organization wallet address
func (m *WalletManager) GetAddress() common.Address {
	return m.walletAddress
}

// GetPrivateKey returns the organization private key
func (m *WalletManager) GetPrivateKey() *ecdsa.PrivateKey {
	return m.priKey
}

func (m *WalletManager) SetPrivateKey(privateKey *ecdsa.PrivateKey) {
	m.priKey = privateKey
	m.pubKey = &privateKey.PublicKey
	m.walletAddress = crypto.PubkeyToAddress(privateKey.PublicKey)
}

// GetPrivateKey returns the organization private key
func (m *WalletManager) GetPublicKey() *ecdsa.PublicKey {
	return m.pubKey
}

// loadPrivateKey loads private key from local DB, and caches the private key
func (m *WalletManager) loadPrivateKey() {
	/*if m.priKey != nil {
		return
	}*/

	// datacenter存储的是加密后的私钥
	priKeyHex, err := m.dataCenter.FindOrgPriKey()
	if nil != err {
		log.WithError(err).Error("failed to load organization wallet. ", err)
		return
	}

	if len(priKeyHex) == 0 {
		log.Warn("failed to load organization wallet. ")
		return
	}
	if m.kms != nil {
		if key, err := m.kms.Decrypt(priKeyHex); err != nil {
			log.Errorf("decrypt organization wallet private key error: %v", err)
			return
		} else {
			priKey, err := crypto.HexToECDSA(key)
			if err != nil {
				log.Errorf("convert organization wallet private key to ECDSA error: %v", err)
				return
			} else {
				m.priKey = priKey
				m.pubKey = &priKey.PublicKey
				m.walletAddress = crypto.PubkeyToAddress(priKey.PublicKey)
				log.Debugf("success to load organization wallet:%s", m.walletAddress)
				return
			}
		}
	} else {
		priKey, err := crypto.HexToECDSA(priKeyHex)
		if err != nil {
			log.Errorf("not a valid private key hex: %v", err)
			return
		}
		m.priKey = priKey
		m.pubKey = &priKey.PublicKey
		m.walletAddress = crypto.PubkeyToAddress(priKey.PublicKey)
		log.Debugf("success to load organization wallet:%s", m.walletAddress)
		return
	}
}
