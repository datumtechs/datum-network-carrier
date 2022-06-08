package signutil

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	"github.com/datumtechs/datum-network-carrier/p2p"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func SignMsg(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// signature the msg
	sign, err := crypto.Sign(hash, privateKey)
	if nil != err {
		return nil, err
	}
	return sign, nil
}

func VerifyMsgSign(nodeId string, m []byte, sig []byte) (bool, error) {
	// Verify the signature
	if len(sig) < 0 {
		return false, ctypes.ErrMsgSignEmpty
	}
	if len(sig) < types.MsgSignLength {
		return false, ctypes.ErrMsgSignLenShort
	}
	if len(sig) > types.MsgSignLength {
		return false, ctypes.ErrMsgSignLenLong
	}

	recPubKey, err := crypto.Ecrecover(m, sig)
	if err != nil {
		return false, err
	}

	ownerNodeId, err := p2p.HexID(nodeId)
	if nil != err {
		return false, ctypes.ErrMsgOwnerNodeIdInvalid
	}
	ownerPubKey, err := ownerNodeId.Pubkey()
	if nil != err {
		return false, ctypes.ErrMsgOwnerNodeIdInvalid
	}
	pbytes := elliptic.Marshal(ownerPubKey.Curve, ownerPubKey.X, ownerPubKey.Y)
	if !bytes.Equal(pbytes, recPubKey) {
		return false, ctypes.ErrMsgSignInvalid
	}
	return true, nil
}

