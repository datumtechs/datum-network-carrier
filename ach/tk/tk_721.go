package tk

import (
	"context"
	"errors"
	"github.com/datumtechs/datum-network-carrier/ach/tk/contracts"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
)

var (
	NotOwner                = errors.New("task sponsor is not the tk owner")
	NotForExpectedAlgorithm = errors.New("the tk is not for current algorithm")
	TermExpired             = errors.New("the tk is expired")
	TermError               = errors.New("the tk term is error")
	SysError                = errors.New("system internal error")
)

func (m *DatumPayManager) inspectTkErc721ExtInfo(taskSponsorAddress common.Address, tk *carrierapipb.TkItem) error {
	extInfo, err := m.getTkErc721ExtInfo(tk)
	if err != nil {
		log.WithError(err).Errorf("cannot fetch erc721 token ext info: tk:%s, id: %d", tk.TkAddress, tk.Id)
		return SysError
	}

	if extInfo.Owner != taskSponsorAddress {
		log.WithError(err).Errorf("task sponsor is not the tk owner, tk: %s, id: %d", tk.TkAddress, tk.Id)
		return NotOwner
	}

	expiredBlockNo, err := strconv.ParseUint(extInfo.Term, 10, 64)
	if err != nil {
		log.WithError(err).Errorf("the tk term is error, tk: %s, id: %d", tk.TkAddress, tk.Id)
		return TermError
	}

	currentBlockNo, err := m.ethContext.BlockNumber(context.Background())
	if err != nil {
		log.WithError(err).Errorf("get the current block number error, tk: %s, id: %d", tk.TkAddress, tk.Id)
		return SysError
	}

	if expiredBlockNo < currentBlockNo {
		log.WithError(err).Errorf("the tk is Expired, tk: %s, id: %d", tk.TkAddress, tk.Id)
		return TermExpired
	}

	return nil

}

func (m *DatumPayManager) getTkErc721ExtInfo(tk *carrierapipb.TkItem) (*struct {
	Owner         common.Address
	Term          string
	ForEncryptAlg bool
}, error) {

	instance, err := contracts.NewErc721(common.HexToAddress(tk.TkAddress), m.ethContext.GetClient())
	if err != nil {
		return nil, err
	}

	resp, err := instance.GetExtInfo(nil, big.NewInt(0).SetUint64(tk.GetId()))
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
