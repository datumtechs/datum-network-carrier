package twopc

import (
	"fmt"
	ctypes "github.com/datumtechs/datum-network-carrier/consensus/twopc/types"
	"github.com/datumtechs/datum-network-carrier/p2p"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
)

func (t *Twopc) validateTaskOfPrepareMsg(pid peer.ID, msg *types.PrepareMsg) error {

	// todo validate user and task sign

	partyIdCache := make(map[string]struct{}, 0)

	// validate the sender
	if types.IsNotSameTaskOrg(msg.GetMsgOption().GetOwner(), msg.GetTask().GetTaskSender()) {
		return fmt.Errorf("task sender of prepareMsg and prepareMsg sender must be same one")
	}
	if err := t.validateOranization(msg.GetMsgOption().GetOwner()); nil != err {
		return fmt.Errorf("verify task sender of prepareMsg %s", err)
	}
	partyIdCache[msg.GetTask().GetTaskSender().GetPartyId()] = struct {}{}

	// validate the algo supplier
	if err := t.validateOranization(msg.GetTask().GetTaskData().GetAlgoSupplier()); nil != err {
		return fmt.Errorf("verify task algoSupplier of prepareMsg %s", err)
	}
	if _, ok := partyIdCache[msg.GetTask().GetTaskData().GetAlgoSupplier().GetPartyId()]; ok {
		return fmt.Errorf("task algoSupplier partyId of prepareMsg has already exists")
	}
	partyIdCache[msg.GetTask().GetTaskData().GetAlgoSupplier().GetPartyId()] = struct {}{}

	// validate data suppliers
	if len(msg.GetTask().GetTaskData().GetDataSuppliers()) == 0 {
		return fmt.Errorf("the task dataSuppliers of prepareMsg is empty")
	}
	for _, supplier := range msg.GetTask().GetTaskData().GetDataSuppliers() {
		if err := t.validateOranization(supplier); nil != err {
			return fmt.Errorf("verify task dataSupplier of prepareMsg %s", err)
		}
		if _, ok := partyIdCache[supplier.GetPartyId()]; ok {
			return fmt.Errorf("task dataSupplier partyId of prepareMsg has already exists")
		}
		partyIdCache[supplier.GetPartyId()] = struct {}{}
	}

	// validate power suppliers
	if len(msg.GetTask().GetTaskData().GetPowerSuppliers()) == 0 {
		return fmt.Errorf("the task powerSuppliers of prepareMsg is empty")
	}
	for _, supplier := range msg.GetTask().GetTaskData().GetPowerSuppliers() {
		if err := t.validateOranization(supplier); nil != err {
			return fmt.Errorf("verify task powerSupplier of prepareMsg %s", err)
		}
		if _, ok := partyIdCache[supplier.GetPartyId()]; ok {
			return fmt.Errorf("task powerSupplier partyId of prepareMsg has already exists")
		}
		partyIdCache[supplier.GetPartyId()] = struct {}{}
	}

	// validate receicers
	if len(msg.GetTask().GetTaskData().GetReceivers()) == 0 {
		return fmt.Errorf("the task receivers of prepareMsg is empty")
	}
	for _, supplier := range msg.GetTask().GetTaskData().GetReceivers() {
		if err := t.validateOranization(supplier); nil != err {
			return fmt.Errorf("verify task receivers of prepareMsg %s", err)
		}
		if _, ok := partyIdCache[supplier.GetPartyId()]; ok {
			return fmt.Errorf("task receiver partyId of prepareMsg has already exists")
		}
		partyIdCache[supplier.GetPartyId()] = struct {}{}
	}

	// validate contractCode
	if "" == strings.Trim(msg.GetTask().GetTaskData().GetAlgorithmCode(), "") && "" == strings.Trim(msg.GetTask().GetTaskData().GetMetaAlgorithmId(), "") {
		return ctypes.ErrProposalTaskAlgorithmEmpty
	}

	// validate task create time
	if msg.GetTask().GetTaskData().GetCreateAt() >= msg.GetCreateAt() {
		return ctypes.ErrProposalParamsInvalid
	}

	return nil
}

func (t *Twopc) validateOranization(identity *carriertypespb.TaskOrganization) error {
	if "" == identity.GetNodeName() {
		return ctypes.ErrOrganizationIdentity
	}
	_, err := p2p.HexID(identity.GetNodeId())
	if nil != err {
		return ctypes.ErrOrganizationIdentity
	}
	has, err := t.resourceMng.GetDB().HasIdentity(&carriertypespb.Organization{
		NodeName:   identity.GetNodeName(),
		NodeId:     identity.GetNodeId(),
		IdentityId: identity.GetIdentityId(),
	})
	if nil != err {
		return err
	}
	if !has {
		return ctypes.ErrOrganizationIdentity
	}

	return nil
}