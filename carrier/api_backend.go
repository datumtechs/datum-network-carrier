package carrier

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	rawdb "github.com/datumtechs/datum-network-carrier/carrierdb/rawdb"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/common/rlputil"
	"github.com/datumtechs/datum-network-carrier/common/timeutils"
	"github.com/datumtechs/datum-network-carrier/core/policy"
	"github.com/datumtechs/datum-network-carrier/grpclient"
	"github.com/datumtechs/datum-network-carrier/params"
	carrierapipb "github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	commonconstantpb "github.com/datumtechs/datum-network-carrier/pb/common/constant"
	fighterapicomputepb "github.com/datumtechs/datum-network-carrier/pb/fighter/api/compute"
	"github.com/datumtechs/datum-network-carrier/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p-core/peer"
	"math/big"
	"strings"
)

// CarrierAPIBackend implements rpc.Backend for Carrier
type CarrierAPIBackend struct {
	carrier *Service
}

func NewCarrierAPIBackend(carrier *Service) *CarrierAPIBackend {
	return &CarrierAPIBackend{carrier: carrier}
}

// add by v0.4.0
func (s *CarrierAPIBackend) GetCarrierChainConfig() *types.CarrierChainConfig {
	return params.CarrierConfig()
}

// add by v0.4.0
func (s *CarrierAPIBackend) GetPolicyEngine() *policy.PolicyEngine {
	return s.carrier.policyEngine
}

func (s *CarrierAPIBackend) SendMsg(msg types.Msg) error {
	return s.carrier.mempool.Add(msg)
}

// system (the yarn node self info)
func (s *CarrierAPIBackend) GetNodeInfo() (*carrierapipb.YarnNodeInfo, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		log.WithError(err).Warnf("Warn to get identity, on GetNodeInfo()")
	}
	var identityId string
	var nodeId string
	var nodeName string
	if nil != identity {
		identityId = identity.IdentityId
		nodeId = identity.NodeId
		nodeName = identity.NodeName
	}

	// local bootstrap node
	enr := s.carrier.config.P2P.ENR()
	enc, _ := rlp.EncodeToBytes(&enr)
	b64 := base64.RawURLEncoding.EncodeToString(enc)
	enrStr := "enr:" + b64

	nodeInfo := &carrierapipb.YarnNodeInfo{
		NodeType:           carrierapipb.NodeType_NodeType_YarnNode,
		NodeId:             nodeId,
		IdentityType:       types.IDENTITY_TYPE_DID, // default: DID
		IdentityId:         identityId,
		Name:               nodeName,
		State:              carrierapipb.YarnNodeState_State_Active,
		RelatePeers:        uint32(len(s.carrier.config.P2P.Peers().Active())),
		LocalBootstrapNode: enrStr,
	}

	for _, _peer := range s.carrier.config.P2P.Peers().Active() {
		log.Debugf("show active peer id:{%s}", _peer.String())
	}
	multiAddr := s.carrier.config.P2P.Host().Addrs()
	if len(multiAddr) != 0 {
		// /ip4/192.168.35.1/tcp/16788
		multiAddrParts := strings.Split(multiAddr[0].String(), "/")
		nodeInfo.ExternalIp = multiAddrParts[2]
		nodeInfo.ExternalPort = multiAddrParts[4]
	}

	selfMultiAddrs, _ := s.carrier.config.P2P.DiscoveryAddresses()
	if len(selfMultiAddrs) != 0 {
		nodeInfo.LocalMultiAddr = selfMultiAddrs[0].String()
	}

	if addr, err := s.carrier.datumPayManager.WalletManager.GetOrgWalletAddress(); err == nil {
		nodeInfo.ObserverProxyWalletAddress = addr.Hex()
	} else {
		log.WithError(err).Errorf("cannot load organization wallet of node info: %v", err)
		return nil, err
	}

	return nodeInfo, nil
}

func (s *CarrierAPIBackend) SetSeedNode(seed *carrierapipb.SeedPeer) (carrierapipb.ConnState, error) {
	// format: enr:-xxxxxx
	if err := s.carrier.config.P2P.AddPeer(seed.GetAddr()); nil != err {
		log.WithError(err).Errorf("Failed to call p2p.AddPeer() with seed.Addr on SetSeedNode(), addr: {%s}", seed.GetAddr())
		return carrierapipb.ConnState_ConnState_UnConnected, err
	}
	if err := s.carrier.carrierDB.SetSeedNode(seed); nil != err {
		log.WithError(err).Errorf("Failed to call SetSeedNode() to store seedNode on SetSeedNode(), seed: {%s}", seed.String())
		return carrierapipb.ConnState_ConnState_UnConnected, err
	}
	addrs, err := s.carrier.config.P2P.PeerFromAddress([]string{seed.GetAddr()})
	if err != nil || len(addrs) == 0 {
		log.WithError(err).Errorf("Failed to parse addr")
		return carrierapipb.ConnState_ConnState_UnConnected, err
	}
	addrInfo, err := peer.AddrInfoFromP2pAddr(addrs[0])
	if nil != err {
		log.WithError(err).Errorf("Failed to call peer.AddrInfoFromP2pAddr() with multiAddr on SetSeedNode(), addr: {%s}", addrs[0].String())
		return carrierapipb.ConnState_ConnState_UnConnected, err
	}
	for _, active := range s.carrier.config.P2P.Peers().Active() {
		if active.String() == addrInfo.ID.String() {
			return carrierapipb.ConnState_ConnState_Connected, nil
		}
	}
	return carrierapipb.ConnState_ConnState_UnConnected, nil
}

func (s *CarrierAPIBackend) DeleteSeedNode(addrStr string) error {
	addrs, err := s.carrier.config.P2P.PeerFromAddress([]string{addrStr})
	if nil != err || len(addrs) == 0 {
		log.WithError(err).Errorf("Failed to convert multiAddr from string on DeleteSeedNode(), addr: {%s}", addrStr)
		return err
	}
	addrInfo, err := peer.AddrInfoFromP2pAddr(addrs[0])
	if nil != err {
		log.WithError(err).Errorf("Failed to call peer.AddrInfoFromP2pAddr() with multiAddr on DeleteSeedNode(), addr: {%s}", addrs[0].String())
		return err
	}

	if err := s.carrier.config.P2P.Disconnect(addrInfo.ID); nil != err {
		log.WithError(err).Errorf("Failed to call p2p.Disconnect() with peerId on DeleteSeedNode(), peerId: {%s}", addrInfo.ID.String())
		return err
	}
	return s.carrier.carrierDB.RemoveSeedNode(addrStr)
}

func (s *CarrierAPIBackend) GetSeedNodeList() ([]*carrierapipb.SeedPeer, error) {
	// load seed node from default bootstrap
	bootstrapNodeStrs, err := s.carrier.config.P2P.BootstrapAddresses()
	if nil != err {
		log.WithError(err).Errorf("Failed to call p2p.BootstrapAddresses() on GetSeedNodeList()")
		return nil, err
	}

	bootstrapNodes := make([]*carrierapipb.SeedPeer, len(bootstrapNodeStrs))
	for i, addr := range bootstrapNodeStrs {
		bootstrapNodes[i] = &carrierapipb.SeedPeer{
			Addr:      addr,
			IsDefault: true,
			ConnState: carrierapipb.ConnState_ConnState_UnConnected,
		}
	}

	// query seed node arr from db
	seeds, err := s.carrier.carrierDB.QuerySeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed to call QuerySeedNodeList() on GetSeedNodeList()")
		return nil, err
	}

	seeds = append(bootstrapNodes, seeds...)

	peerIds := s.carrier.config.P2P.Peers().Active()
	tmp := make(map[peer.ID]struct{}, len(peerIds))
	for _, peerId := range peerIds {
		tmp[peerId] = struct{}{}
	}

	for i, seed := range seeds {
		addrs, err := s.carrier.config.P2P.PeerFromAddress([]string{seed.GetAddr()})
		if nil != err || len(addrs) == 0 {
			log.WithError(err).Errorf("Failed to convert multiAddr from string on GetSeedNodeList, addr: {%s}", seed.GetAddr())
			continue
		}
		addrInfo, err := peer.AddrInfoFromP2pAddr(addrs[0])
		if nil != err {
			log.WithError(err).Errorf("Failed to call peer.AddrInfoFromP2pAddr() with multiAddr on GetSeedNodeList(), addr: {%s}", addrs[0].String())
			continue
		}

		if _, ok := tmp[addrInfo.ID]; ok {
			seeds[i].ConnState = carrierapipb.ConnState_ConnState_Connected
		}
	}
	return seeds, nil
}

func (s *CarrierAPIBackend) storeLocalResource(identity *carriertypespb.Organization, jobNodeId string, jobNodeStatus *fighterapicomputepb.GetStatusReply) error {

	// store into local db
	if err := s.carrier.carrierDB.StoreLocalResource(types.NewLocalResource(&carriertypespb.LocalResourcePB{
		Owner:      identity,
		JobNodeId:  jobNodeId,
		DataId:     "", // can not own powerId now, because power have not publish
		DataStatus: commonconstantpb.DataStatus_DataStatus_Valid,
		// resource status, eg: create/release/revoke
		State: commonconstantpb.PowerState_PowerState_Created,
		// unit: byte
		TotalMem: jobNodeStatus.GetTotalMemory(),
		UsedMem:  0,
		// number of cpu cores.
		TotalProcessor: jobNodeStatus.GetTotalCpu(),
		UsedProcessor:  0,
		// unit: byte
		TotalBandwidth: jobNodeStatus.GetTotalBandwidth(),
		UsedBandwidth:  0,
		TotalDisk:      jobNodeStatus.GetTotalDisk(),
		UsedDisk:       0,
	})); nil != err {
		log.WithError(err).Errorf("Failed to store power to local on MessageHandler with broadcast, jobNodeId: {%s}",
			jobNodeId)
		return err
	}
	return nil
}

func (s *CarrierAPIBackend) SetRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) (carrierapipb.ConnState, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case carrierapipb.PrefixTypeDataNode, carrierapipb.PrefixTypeJobNode:
	default:
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("invalid nodeType")
	}
	if typ == carrierapipb.PrefixTypeJobNode {
		client, err := grpclient.NewJobNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreJobNodeClient(node.GetId(), client)

		jobNodeStatus, err := client.GetStatus()
		if err != nil {
			client.Close()
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect jobNode query status failed, %s", err)
		}
		// add resource usage first, but not own power now (mem, proccessor, bandwidth)
		if err = s.storeLocalResource(identity, node.GetId(), jobNodeStatus); nil != err {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("store jobNode local resource failed, %s", err)
		}
	}
	if typ == carrierapipb.PrefixTypeDataNode {
		client, err := grpclient.NewDataNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreDataNodeClient(node.GetId(), client)

		dataNodeStatus, err := client.GetStatus()
		if err != nil {
			client.Close()
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect dataNode query status failed, %s", err)
		}
		// add data resource  (disk)
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk(), true))
		//err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
		log.Debugf("Add dataNode resource table succeed, nodeId: {%s}, totalDisk: {%d}, usedDisk: {%d}", node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk())
	}
	node.ConnState = carrierapipb.ConnState_ConnState_Connected
	if err = s.carrier.carrierDB.SetRegisterNode(typ, node); err != nil {
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("Store registerNode to db failed, %s", err)
	}
	return carrierapipb.ConnState_ConnState_Connected, nil
}

func (s *CarrierAPIBackend) UpdateRegisterNode(typ carrierapipb.RegisteredNodeType, node *carrierapipb.YarnRegisteredPeerDetail) (carrierapipb.ConnState, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case carrierapipb.PrefixTypeDataNode, carrierapipb.PrefixTypeJobNode:
	default:
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("invalid nodeType")
	}
	if typ == carrierapipb.PrefixTypeJobNode {

		// The published jobNode cannot be updated directly
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(node.GetId())
		if rawdb.IsNoDBNotFoundErr(err) {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on UpdateRegisterNode, %s", resourceTable.String())
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("still have the published computing power information on old jobNode failed, input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				node.Id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// First check whether there is a task being executed on jobNode or not.
		runningTaskCount, err := s.carrier.carrierDB.QueryJobNodeRunningTaskCount(node.GetId())
		if rawdb.IsNoDBNotFoundErr(err) {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceManager.QueryJobNodeClient(node.GetId()); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceManager.RemoveJobNodeClient(node.GetId())
		}

		// generate new client
		client, err := grpclient.NewJobNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreJobNodeClient(node.GetId(), client)

		jobNodeStatus, err := client.GetStatus()
		if err != nil {
			client.Close()
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect jobNode query status failed, %s", err)
		}
		// add resource usage first, but not own power now (mem, proccessor, bandwidth)
		if err = s.storeLocalResource(identity, node.GetId(), jobNodeStatus); nil != err {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("store jobNode local resource failed, %s", err)
		}

	}

	if typ == carrierapipb.PrefixTypeDataNode {

		if client, ok := s.carrier.resourceManager.QueryDataNodeClient(node.GetId()); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceManager.RemoveDataNodeClient(node.GetId())
		}

		// remove old data resource  (disk)
		if err := s.carrier.carrierDB.RemoveDataResourceTable(node.GetId()); rawdb.IsNoDBNotFoundErr(err) {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("remove disk summary of old dataNode, %s", err)
		}

		client, err := grpclient.NewDataNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreDataNodeClient(node.GetId(), client)

		dataNodeStatus, err := client.GetStatus()
		if err != nil {
			client.Close()
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("connect dataNode query status failed, %s", err)
		}
		// add new data resource  (disk)
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk(), true))
		//err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
		log.Debugf("Update dataNode resource table succeed, nodeId: {%s}, totalDisk: {%d}, usedDisk: {%d}", node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk())
	}

	// remove  old jobNode from db
	if err := s.carrier.carrierDB.DeleteRegisterNode(typ, node.GetId()); nil != err {
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("remove old registerNode from db failed, %s", err)
	}

	// add new node to db
	node.ConnState = carrierapipb.ConnState_ConnState_Connected
	if err = s.carrier.carrierDB.SetRegisterNode(typ, node); err != nil {
		return carrierapipb.ConnState_ConnState_UnConnected, fmt.Errorf("update registerNode to db failed, %s", err)
	}
	return carrierapipb.ConnState_ConnState_Connected, nil
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ carrierapipb.RegisteredNodeType, id string) error {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case carrierapipb.PrefixTypeDataNode, carrierapipb.PrefixTypeJobNode:
	default:
		return fmt.Errorf("invalid nodeType")
	}
	if typ == carrierapipb.PrefixTypeJobNode {

		// release resource of jobNode by jobNodeId
		if err = s.carrier.resourceManager.UnLockLocalResourceWithJobNodeId(id); nil != err {
			log.WithError(err).Errorf("Failed to unlock local resource with jobNodeId on RemoveRegisterNode(), jobNodeId: {%s}",
				id)
			return fmt.Errorf("unlock jobnode local resource failed, %d", err)
		}

		// The published jobNode cannot be updated directly
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			//log.Debugf("still have the published computing power information on old jobNode on RemoveRegisterNode, %s", resourceTable.String())
			//return fmt.Errorf("still have the published computing power information on old jobNode failed,input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
			//	id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
			log.Warnf("still have the published computing power information by the jobNode, that need revoke power short circuit on RemoveRegisterNode, %s",
				resourceTable.String())
			// ##############################
			// A. remove power about jobNode
			// ##############################

			// 1. remove jobNodeId and powerId mapping
			if err := s.carrier.carrierDB.RemoveJobNodeIdByPowerId(resourceTable.GetPowerId()); nil != err {
				log.WithError(err).Errorf("Failed to call RemoveJobNodeIdByPowerId() on RemoveRegisterNode with revoke power, powerId: {%s}, jobNodeId: {%s}",
					resourceTable.GetPowerId(), id)
				return err
			}

			// 2. remove local jobNode resource table
			if err := s.carrier.carrierDB.RemoveLocalResourceTable(id); nil != err {
				log.WithError(err).Errorf("Failed to RemoveLocalResourceTable on RemoveRegisterNode with revoke power, powerId: {%s}, jobNodeId: {%s}",
					resourceTable.GetPowerId(), id)
				return err
			}

			// 3. revoke power about jobNode from global
			if err := s.carrier.carrierDB.RevokeResource(types.NewResource(&carriertypespb.ResourcePB{
				Owner:  identity,
				DataId: resourceTable.GetPowerId(),
				// the status of data, N means normal, D means deleted.
				DataStatus: commonconstantpb.DataStatus_DataStatus_Invalid,
				// resource status, eg: create/release/revoke
				State:    commonconstantpb.PowerState_PowerState_Revoked,
				UpdateAt: timeutils.UnixMsecUint64(),
			})); nil != err {
				log.WithError(err).Errorf("Failed to remove dataCenter resource on RemoveRegisterNode with revoke power, powerId: {%s}, jobNodeId: {%s}",
					resourceTable.GetPowerId(), id)
				return err
			}
		}
		// ##############################
		// B. remove except power things about jobNode
		// ##############################

		// 1. remove all running task
		taskIdList, _ := s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(id)
		for _, taskId := range taskIdList {
			s.carrier.carrierDB.RemoveJobNodeTaskIdAllPartyIds(id, taskId)
		}
		// 2. remove local jobNode reource
		// remove jobNode local resource
		if err = s.carrier.carrierDB.RemoveLocalResource(id); nil != err {
			log.WithError(err).Errorf("Failed to remove jobNode local resource on RemoveRegisterNode, jobNodeId: {%s}",
				id)
			return err
		}
		// 3. remove rpc client
		if client, ok := s.carrier.resourceManager.QueryJobNodeClient(id); ok {
			client.Close()
			s.carrier.resourceManager.RemoveJobNodeClient(id)
		}
		// 4. goto `Finally`
	}

	if typ == carrierapipb.PrefixTypeDataNode {

		// 1. remove data resource  (disk)
		if err := s.carrier.carrierDB.RemoveDataResourceTable(id); rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to remove disk summary of old dataNode on RemoveRegisterNode")
			return err
		}
		// 2. remove rpc client
		if client, ok := s.carrier.resourceManager.QueryDataNodeClient(id); ok {
			client.Close()
			s.carrier.resourceManager.RemoveDataNodeClient(id)
		}
		// 3. goto `Finally`
	}
	// `Finally`: remove local jobNode info
	return s.carrier.carrierDB.DeleteRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNode(typ carrierapipb.RegisteredNodeType, id string) (*carrierapipb.YarnRegisteredPeerDetail, error) {
	node, err := s.carrier.carrierDB.QueryRegisterNode(typ, id)
	if nil != err {
		return nil, err
	}

	if typ == carrierapipb.PrefixTypeJobNode {

		client, ok := s.carrier.resourceManager.QueryJobNodeClient(id)
		if !ok {
			node.ConnState = carrierapipb.ConnState_ConnState_UnConnected
		} else {
			if !client.IsConnected() {
				node.ConnState = carrierapipb.ConnState_ConnState_UnConnected
			} else {
				node.ConnState = carrierapipb.ConnState_ConnState_Connected
			}
		}

	} else {

		client, ok := s.carrier.resourceManager.QueryDataNodeClient(id)
		if !ok {
			node.ConnState = carrierapipb.ConnState_ConnState_UnConnected
		} else {
			if !client.IsConnected() {
				node.ConnState = carrierapipb.ConnState_ConnState_UnConnected
			} else {
				node.ConnState = carrierapipb.ConnState_ConnState_Connected
			}
		}

	}
	return node, nil
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ carrierapipb.RegisteredNodeType) ([]*carrierapipb.YarnRegisteredPeerDetail, error) {

	if typ != carrierapipb.PrefixTypeJobNode &&
		typ != carrierapipb.PrefixTypeDataNode {

		return nil, fmt.Errorf("Invalid nodeType")
	}

	nodeList, err := s.carrier.carrierDB.QueryRegisterNodeList(typ)
	if nil != err {
		return nil, err
	}

	for i, n := range nodeList {

		var connState carrierapipb.ConnState
		var duration uint64
		var taskCount uint32
		var taskIdList []string
		var fileCount uint32
		var fileTotalSize uint32

		if typ == carrierapipb.PrefixTypeJobNode {

			client, ok := s.carrier.resourceManager.QueryJobNodeClient(n.GetId())
			if !ok {
				connState = carrierapipb.ConnState_ConnState_UnConnected
			} else {
				duration = uint64(client.RunningDuration())
				if !client.IsConnected() {
					connState = carrierapipb.ConnState_ConnState_UnConnected
				} else {
					connState = carrierapipb.ConnState_ConnState_Connected
				}
			}

			taskCount, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskCount(n.GetId())
			taskIdList, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(n.GetId())
			fileCount = 0     // todo need implament this logic
			fileTotalSize = 0 // todo need implament this logic
		} else {

			client, ok := s.carrier.resourceManager.QueryDataNodeClient(n.GetId())
			if !ok {
				connState = carrierapipb.ConnState_ConnState_UnConnected
			} else {
				duration = uint64(client.RunningDuration())
				if !client.IsConnected() {
					connState = carrierapipb.ConnState_ConnState_UnConnected
				} else {
					connState = carrierapipb.ConnState_ConnState_Connected
				}
			}

			taskCount = 0
			taskIdList = nil
			fileCount = 0     // todo need implament this logic
			fileTotalSize = 0 // todo need implament this logic
		}
		n.Duration = duration
		n.ConnState = connState
		n.TaskCount = taskCount
		n.TaskIdList = taskIdList
		n.FileCount = fileCount
		n.FileTotalSize = fileTotalSize

		nodeList[i] = n
	}

	return nodeList, nil
}

func (s *CarrierAPIBackend) SendTaskEvent(event *carriertypespb.TaskEvent) error {
	return s.carrier.TaskManager.SendTaskEvent(event)
}

func (s *CarrierAPIBackend) ReportTaskResourceUsage(nodeType carrierapipb.NodeType, ip, port string, usage *types.TaskResuorceUsage) error {

	if err := s.carrier.TaskManager.HandleReportResourceUsage(usage); nil != err {
		log.WithError(err).Errorf("Failed to call HandleReportResourceUsage() on CarrierAPIBackend.ReportTaskResourceUsage(), taskId: {%s}, partyId: {%s}, nodeType: {%s}, ip:{%s}, port: {%s}",
			usage.GetTaskId(), usage.GetPartyId(), nodeType.String(), ip, port)
		return err
	}
	return nil
}

func (s *CarrierAPIBackend) GenerateObServerProxyWalletAddress() (string, error) {
	if s.carrier.datumPayManager != nil {
		if addr, err := s.carrier.datumPayManager.WalletManager.GetOrgWalletAddress(); nil != err {
			log.WithError(err).Error("Failed to call GenerateOrgWallet() on CarrierAPIBackend.GenerateObServerProxyWalletAddress()")
			return "", err
		} else {
			log.Debugf("Success to generate organization wallet %s", addr)
			return addr.Hex(), nil
		}
	} else {
		return "", errors.New("token20Pay manager not initialized properly")
	}
}

// metadata api
func (s *CarrierAPIBackend) IsInternalMetadata(metadataId string) (bool, error) {
	return s.carrier.carrierDB.IsInternalMetadataById(metadataId)
}

func (s *CarrierAPIBackend) GetInternalMetadataDetail(metadataId string) (*types.Metadata, error) {
	// find internal metadata
	metadata, err := s.carrier.carrierDB.QueryInternalMetadataById(metadataId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, fmt.Errorf("found internal metadata by metadataId failed, %s", err)
	}
	return metadata, nil
}

func (s *CarrierAPIBackend) GetMetadataDetail(metadataId string) (*types.Metadata, error) {
	// find global metadata
	metadata, err := s.carrier.carrierDB.QueryMetadataById(metadataId)
	if nil != err {
		return nil, fmt.Errorf("global metadata by metadataId failed, %s", err)
	}
	return metadata, nil
}

// GetMetadataDetailList returns a list of all metadata details in the network.
func (s *CarrierAPIBackend) GetGlobalMetadataDetailList(lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalMetadataDetail, error) {
	log.Debug("Invoke: GetGlobalMetadataDetailList executing...")
	var (
		arr []*carrierapipb.GetGlobalMetadataDetail
		err error
	)

	publishMetadataArr, err := s.carrier.carrierDB.QueryMetadataList(lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, fmt.Errorf("found global metadata arr failed, %s", err)
	}
	if len(publishMetadataArr) != 0 {
		arr = append(arr, policy.NewGlobalMetadataInfoArrayFromMetadataArray(publishMetadataArr)...)
	}
	//// set metadata used taskCount
	//for i, metadata := range arr {
	//	count, err := s.carrier.carrierDB.QueryMetadataUsedTaskIdCount(metadata.GetInformation().GetMetadataSummary().GetMetadataId())
	//	if nil != err {
	//		log.WithError(err).Warnf("Warn, query metadata used taskIdCount failed on CarrierAPIBackend.GetGlobalMetadataDetailList()")
	//		continue
	//	}
	//	metadata.Information.TotalTaskCount = count
	//	arr[i] = metadata
	//}
	return arr, err
}

// GetGlobalMetadataDetailListByIdentityId returns a list of all metadata details in the network by identityId.
func (s *CarrierAPIBackend) GetGlobalMetadataDetailListByIdentityId(identityId string, lastUpdate, pageSize uint64) ([]*carrierapipb.GetGlobalMetadataDetail, error) {
	log.Debug("Invoke: GetGlobalMetadataDetailListByIdentityId executing...")
	var (
		arr []*carrierapipb.GetGlobalMetadataDetail
		err error
	)

	publishMetadataArr, err := s.carrier.carrierDB.QueryMetadataListByIdentity(identityId, lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, fmt.Errorf("found global metadata arr by identityId failed, %s, identityId: {%s}", err, identityId)
	}
	if len(publishMetadataArr) != 0 {
		arr = append(arr, policy.NewGlobalMetadataInfoArrayFromMetadataArray(publishMetadataArr)...)
	}
	//// set metadata used taskCount
	//for i, metadata := range arr {
	//	count, err := s.carrier.carrierDB.QueryMetadataUsedTaskIdCount(metadata.GetInformation().GetMetadataSummary().GetMetadataId())
	//	if nil != err {
	//		log.WithError(err).Warnf("Warn, query metadata used taskIdCount failed on CarrierAPIBackend.GetGlobalMetadataDetailList()")
	//		continue
	//	}
	//	metadata.Information.TotalTaskCount = count
	//	arr[i] = metadata
	//}
	return arr, err
}

func (s *CarrierAPIBackend) GetLocalMetadataDetailList(lastUpdate uint64, pageSize uint64) ([]*carrierapipb.GetLocalMetadataDetail, error) {
	log.Debug("Invoke: GetLocalMetadataDetailList executing...")

	var (
		arr []*carrierapipb.GetLocalMetadataDetail
		err error
	)

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("found local identity failed, %s", err)
	}

	publishMetadataArr, err := s.carrier.carrierDB.QueryMetadataListByIdentity(identity.GetIdentityId(), lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, fmt.Errorf("found global metadata arr failed on query local metadata arr, %s", err)
	}

	arr = append(arr, policy.NewLocalMetadataInfoArrayFromMetadataArray(nil, publishMetadataArr)...)

	// set metadata used taskCount
	for i, metadata := range arr {
		count, err := s.carrier.carrierDB.QueryMetadataHistoryTaskIdCount(metadata.GetInformation().GetMetadataSummary().GetMetadataId())
		if nil != err {
			log.WithError(err).Warnf("Warn, query global metadata history used taskIdCount failed on CarrierAPIBackend.GetLocalMetadataDetailList()")
			continue
		}
		metadata.Information.TotalTaskCount = count
		arr[i] = metadata
	}

	return arr, nil
}

func (s *CarrierAPIBackend) GetLocalInternalMetadataDetailList() ([]*carrierapipb.GetLocalMetadataDetail, error) {
	log.Debug("Invoke: GetLocalInternalMetadataDetailList executing...")

	var (
		arr []*carrierapipb.GetLocalMetadataDetail
		err error
	)

	_, err = s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("found local identity failed, %s", err)
	}

	internalMetadataArr, err := s.carrier.carrierDB.QueryInternalMetadataList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, fmt.Errorf("found local internal metadata arr failed, %s", err)
	}

	arr = append(arr, policy.NewLocalMetadataInfoArrayFromMetadataArray(internalMetadataArr, nil)...)

	// set metadata used taskCount
	for i, metadata := range arr {
		count, err := s.carrier.carrierDB.QueryMetadataHistoryTaskIdCount(metadata.GetInformation().GetMetadataSummary().GetMetadataId())
		if nil != err {
			log.WithError(err).Warnf("Warn, query internal metadata history used taskIdCount failed on CarrierAPIBackend.GetLocalInternalMetadataDetailList()")
			continue
		}
		metadata.Information.TotalTaskCount = count
		arr[i] = metadata
	}

	return arr, nil
}

func (s *CarrierAPIBackend) GetMetadataUsedTaskIdList(identityId, metadataId string) ([]string, error) {
	taskIds, err := s.carrier.carrierDB.QueryMetadataHistoryTaskIds(metadataId)
	if nil != err {
		return nil, err
	}
	return taskIds, nil
}

func (s *CarrierAPIBackend) UpdateGlobalMetadata(metadata *types.Metadata) error {
	return s.carrier.carrierDB.UpdateGlobalMetadata(metadata)
}

// power api

func (s *CarrierAPIBackend) GetGlobalPowerSummaryList() ([]*carrierapipb.GetGlobalPowerSummary, error) {
	log.Debug("Invoke: GetGlobalPowerSummaryList executing...")
	resourceList, err := s.carrier.carrierDB.QueryGlobalResourceSummaryList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power summary list, len: {%d}", len(resourceList))
	powerList := make([]*carrierapipb.GetGlobalPowerSummary, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &carrierapipb.GetGlobalPowerSummary{
			Owner: resource.GetOwner(),
			Information: &carriertypespb.PowerUsageDetail{
				TotalTaskCount:   0,
				CurrentTaskCount: 0,
				Tasks:            make([]*carriertypespb.PowerTask, 0),
				Information: &carriertypespb.ResourceUsageOverview{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: resource.GetTotalProcessor(),
					UsedProcessor:  resource.GetUsedProcessor(),
					TotalBandwidth: resource.GetTotalBandwidth(),
					UsedBandwidth:  resource.GetUsedBandwidth(),
				},
				State: resource.GetState(),
				// todo Summary is aggregate information and does not require paging, so there are no `publishat` and `updateat`
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetGlobalPowerDetailList(lastUpdate uint64, pageSize uint64) ([]*carrierapipb.GetGlobalPowerDetail, error) {
	log.Debug("Invoke: GetGlobalPowerDetailList executing...")
	resourceList, err := s.carrier.carrierDB.QueryGlobalResourceDetailList(lastUpdate, pageSize)
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power detail list, len: {%d}", len(resourceList))
	powerList := make([]*carrierapipb.GetGlobalPowerDetail, 0, resourceList.Len())
	for _, resource := range resourceList.To() {

		jobNodeId, _ := s.carrier.carrierDB.QueryJobNodeIdByPowerId(resource.GetDataId())
		var (
			totalTaskCount   uint32
			currentTaskCount uint32
		)
		if "" != jobNodeId {
			totalTaskCount, _ = s.carrier.carrierDB.QueryJobNodeHistoryTaskCount(jobNodeId)
			currentTaskCount, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskCount(jobNodeId)
		}

		powerList = append(powerList, &carrierapipb.GetGlobalPowerDetail{
			Owner:   resource.GetOwner(),
			PowerId: resource.GetDataId(),
			Information: &carriertypespb.PowerUsageDetail{
				TotalTaskCount:   totalTaskCount,
				CurrentTaskCount: currentTaskCount,
				Tasks:            make([]*carriertypespb.PowerTask, 0),
				Information: &carriertypespb.ResourceUsageOverview{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: resource.GetTotalProcessor(),
					UsedProcessor:  resource.GetUsedProcessor(),
					TotalBandwidth: resource.GetTotalBandwidth(),
					UsedBandwidth:  resource.GetUsedBandwidth(),
					TotalDisk:      resource.GetTotalDisk(),
					UsedDisk:       resource.GetUsedDisk(),
				},
				State:     resource.GetState(),
				PublishAt: resource.GetPublishAt(),
				UpdateAt:  resource.GetUpdateAt(),
				Nonce:     resource.GetNonce(),
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetLocalPowerDetailList() ([]*carrierapipb.GetLocalPowerDetail, error) {
	log.Debug("Invoke:GetLocalPowerDetailList executing...")
	// query local resource list from db.
	machineList, err := s.carrier.carrierDB.QueryLocalResourceList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetLocalPowerDetailList, call QueryLocalResourceList, machineList: %s", machineList.String())

	buildPowerTaskList := func(jobNodeId string) []*carriertypespb.PowerTask {

		powerTaskList := make([]*carriertypespb.PowerTask, 0)

		taskIds, err := s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(jobNodeId)
		if rawdb.IsNoDBNotFoundErr(err) {
			log.WithError(err).Errorf("Failed to query jobNode runningTaskIds on GetLocalPowerDetailList, jobNodeId: {%s}", jobNodeId)
			return powerTaskList
		}

		for _, taskId := range taskIds {
			// query local task information by taskId
			task, err := s.carrier.carrierDB.QueryLocalTask(taskId)
			if err != nil {
				log.WithError(err).Warnf("Warning query local task on GetLocalPowerDetailList, taskId: {%s}", taskId)
				continue
			}

			partyIds, err := s.carrier.carrierDB.QueryJobNodeRunningTaskAllPartyIdList(jobNodeId, taskId)
			if rawdb.IsNoDBNotFoundErr(err) {
				log.WithError(err).Errorf("Failed to query jobNode runningTask all partyIds on GetLocalPowerDetailList, jobNodeId: {%s}, taskId: {%s}", jobNodeId, taskId)
				continue
			}
			// build powerTask info
			powerTask := &carriertypespb.PowerTask{
				TaskId:   taskId,
				TaskName: task.GetTaskData().TaskName,
				Owner: &carriertypespb.Organization{
					NodeName:   task.GetTaskSender().GetNodeName(),
					NodeId:     task.GetTaskSender().GetNodeId(),
					IdentityId: task.GetTaskSender().GetIdentityId(),
				},
				Receivers: make([]*carriertypespb.Organization, 0),
				OperationCost: &carriertypespb.TaskResourceCostDeclare{
					Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
					Memory:    task.GetTaskData().GetOperationCost().GetMemory(),
					Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
					Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
				},
				OperationSpend: nil, // will be culculating after ...
				CreateAt:       task.GetTaskData().GetCreateAt(),
			}
			// build dataSuppliers of task info
			for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
				powerTask.Partners = append(powerTask.GetPartners(), &carriertypespb.Organization{
					NodeName:   dataSupplier.GetNodeName(),
					NodeId:     dataSupplier.GetNodeId(),
					IdentityId: dataSupplier.GetIdentityId(),
				})
			}
			// build receivers of task info
			for _, receiver := range task.GetTaskData().GetReceivers() {
				powerTask.Receivers = append(powerTask.GetReceivers(), &carriertypespb.Organization{
					NodeName:   receiver.GetNodeName(),
					NodeId:     receiver.GetNodeId(),
					IdentityId: receiver.GetIdentityId(),
				})
			}

			var (
				processor uint32
				memory    uint64
				bandwidth uint64
				duration  uint64
			)
			// culculate power resource used on task
			partyIdTmp := make(map[string]struct{}, 0)
			for _, partyId := range partyIds {
				partyIdTmp[partyId] = struct{}{}
			}

			for _, option := range task.GetTaskData().PowerResourceOptions {
				if _, ok := partyIdTmp[option.GetPartyId()]; ok {
					processor += option.GetResourceUsedOverview().GetUsedProcessor()
					memory += option.GetResourceUsedOverview().GetUsedMem()
					bandwidth += option.GetResourceUsedOverview().GetUsedBandwidth()
				}
			}

			if task.GetTaskData().GetStartAt() != 0 {
				duration = timeutils.UnixMsecUint64() - task.GetTaskData().GetStartAt()
			}
			powerTask.OperationSpend = &carriertypespb.TaskResourceCostDeclare{
				Processor: processor,
				Memory:    memory,
				Bandwidth: bandwidth,
				Duration:  duration,
			}
			powerTaskList = append(powerTaskList, powerTask)
		}
		return powerTaskList
	}

	// culculate task history total count on jobNode
	taskTotalCount := func(jobNodeId string) uint32 {
		count, err := s.carrier.carrierDB.QueryJobNodeHistoryTaskCount(jobNodeId)
		if err != nil {
			log.WithError(err).Errorf("Failed to query task totalCount with jobNodeId on GetLocalPowerDetailList, jobNodeId: {%s}", jobNodeId)
			return 0
		}
		return count
	}
	// culculate task current running count on jobNode
	taskRunningCount := func(jobNodeId string) uint32 {
		count, err := s.carrier.carrierDB.QueryJobNodeRunningTaskCount(jobNodeId)
		if err != nil {
			log.WithError(err).Errorf("Failed to query task runningCount with jobNodeId on GetLocalPowerDetailList, jobNodeId: {%s}", jobNodeId)
			return 0
		}
		return count
	}

	resourceList := machineList.To()
	// handle  resource one by one
	result := make([]*carrierapipb.GetLocalPowerDetail, len(resourceList))
	for i, resource := range resourceList {

		nodePowerDetail := &carrierapipb.GetLocalPowerDetail{
			JobNodeId: resource.GetJobNodeId(),
			PowerId:   resource.GetDataId(),
			Owner:     resource.GetOwner(),
			Power: &carriertypespb.PowerUsageDetail{
				TotalTaskCount:   taskTotalCount(resource.GetJobNodeId()),
				CurrentTaskCount: taskRunningCount(resource.GetJobNodeId()),
				Tasks:            make([]*carriertypespb.PowerTask, 0),
				Information: &carriertypespb.ResourceUsageOverview{
					TotalMem:       resource.GetTotalMem(),
					UsedMem:        resource.GetUsedMem(),
					TotalProcessor: resource.GetTotalProcessor(),
					UsedProcessor:  resource.GetUsedProcessor(),
					TotalBandwidth: resource.GetTotalBandwidth(),
					UsedBandwidth:  resource.GetUsedBandwidth(),
				},
				State: resource.GetState(),
				// local resource power need not they (publishAt and updateAt).
				//PublishAt: ,
				//UpdateAt: ,
				Nonce: resource.GetNonce(),
			},
		}
		nodePowerDetail.GetPower().Tasks = buildPowerTaskList(resource.GetJobNodeId())
		result[i] = nodePowerDetail
	}
	return result, nil
}

// identity api

func (s *CarrierAPIBackend) GetNodeIdentity() (*types.Identity, error) {
	nodeAlias, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, err
	}
	return types.NewIdentity(&carriertypespb.IdentityPB{
		IdentityId: nodeAlias.GetIdentityId(),
		NodeId:     nodeAlias.GetNodeId(),
		NodeName:   nodeAlias.GetNodeName(),
		ImageUrl:   nodeAlias.GetImageUrl(),
		Details:    nodeAlias.GetDetails(),
	}), err
}

func (s *CarrierAPIBackend) GetIdentityList(lastUpdate uint64, pageSize uint64) ([]*types.Identity, error) {
	return s.carrier.carrierDB.QueryIdentityList(lastUpdate, pageSize)
}

// for metadataAuthority

func (s *CarrierAPIBackend) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (commonconstantpb.AuditMetadataOption, error) {
	return s.carrier.authManager.AuditMetadataAuthority(audit)
}

func (s *CarrierAPIBackend) GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	return s.carrier.authManager.GetLocalMetadataAuthorityList(lastUpdate, pageSize)
}

func (s *CarrierAPIBackend) GetGlobalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	return s.carrier.authManager.GetGlobalMetadataAuthorityList(lastUpdate, pageSize)
}

func (s *CarrierAPIBackend) HasValidMetadataAuth(userType commonconstantpb.UserType, user, identityId, metadataId string) (bool, error) {
	return s.carrier.authManager.HasValidMetadataAuth(userType, user, identityId, metadataId)
}

// task api
func (s *CarrierAPIBackend) GetLocalTask(taskId string) (*carriertypespb.TaskDetail, error) {
	// the task is executing.
	localTask, err := s.carrier.carrierDB.QueryLocalTask(taskId)
	if nil != err {
		log.WithError(err).Errorf("Failed to query local task on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}", taskId)
		return nil, err
	}

	if nil == localTask {
		log.Errorf("Not found local task on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}", taskId)
		return nil, fmt.Errorf("not found local task")
	}

	return &carriertypespb.TaskDetail{
		Information: &carriertypespb.TaskDetailSummary{
			/**
			TaskId                   string
			TaskName                 string
			User                     string
			UserType                 constant.UserType
			Sender                   *TaskOrganization
			AlgoSupplier             *TaskOrganization
			DataSuppliers            []*TaskOrganization
			PowerSuppliers           []*TaskOrganization
			Receivers                []*TaskOrganization
			DataPolicyTypes          []uint32
			DataPolicyOptions        []string
			PowerPolicyTypes         []uint32
			PowerPolicyOptions       []string
			ReceiverPolicyTypes      []uint32
			ReceiverPolicyOptions    []string
			DataFlowPolicyTypes      []uint32
			DataFlowPolicyOptions    []string
			OperationCost            *TaskResourceCostDeclare
			AlgorithmCode            string
			MetaAlgorithmId          string
			AlgorithmCodeExtraParams string
			PowerResourceOptions     []*TaskPowerResourceOption
			State                    constant.TaskState
			Reason                   string
			Desc                     string
			CreateAt                 uint64
			StartAt                  uint64
			EndAt                    uint64
			Sign                     []byte
			Nonce                    uint64
			UpdateAt                 uint64
			*/
			TaskId:                   localTask.GetTaskData().GetTaskId(),
			TaskName:                 localTask.GetTaskData().GetTaskName(),
			UserType:                 localTask.GetTaskData().GetUserType(),
			User:                     localTask.GetTaskData().GetUser(),
			Sender:                   localTask.GetTaskSender(),
			AlgoSupplier:             localTask.GetTaskData().GetAlgoSupplier(),
			DataSuppliers:            localTask.GetTaskData().GetDataSuppliers(),
			PowerSuppliers:           localTask.GetTaskData().GetPowerSuppliers(),
			Receivers:                localTask.GetTaskData().GetReceivers(),
			DataPolicyTypes:          localTask.GetTaskData().GetDataPolicyTypes(),
			DataPolicyOptions:        localTask.GetTaskData().GetDataPolicyOptions(),
			PowerPolicyTypes:         localTask.GetTaskData().GetPowerPolicyTypes(),
			PowerPolicyOptions:       localTask.GetTaskData().GetPowerPolicyOptions(),
			ReceiverPolicyTypes:      localTask.GetTaskData().GetReceiverPolicyTypes(),
			ReceiverPolicyOptions:    localTask.GetTaskData().GetReceiverPolicyOptions(),
			DataFlowPolicyTypes:      localTask.GetTaskData().GetDataFlowPolicyTypes(),
			DataFlowPolicyOptions:    localTask.GetTaskData().GetDataFlowPolicyOptions(),
			OperationCost:            localTask.GetTaskData().GetOperationCost(),
			AlgorithmCode:            localTask.GetTaskData().GetAlgorithmCode(),
			MetaAlgorithmId:          localTask.GetTaskData().GetMetaAlgorithmId(),
			AlgorithmCodeExtraParams: localTask.GetTaskData().GetAlgorithmCodeExtraParams(),
			PowerResourceOptions:     localTask.GetTaskData().GetPowerResourceOptions(),
			State:                    localTask.GetTaskData().GetState(),
			Reason:                   localTask.GetTaskData().GetReason(),
			Desc:                     localTask.GetTaskData().GetDesc(),
			CreateAt:                 localTask.GetTaskData().GetCreateAt(),
			StartAt:                  localTask.GetTaskData().GetStartAt(),
			EndAt:                    localTask.GetTaskData().GetEndAt(),
			Sign:                     localTask.GetTaskData().GetSign(),
			Nonce:                    localTask.GetTaskData().GetNonce(),
			UpdateAt:                 localTask.GetTaskData().GetEndAt(), // The endAt of the task is the updateAt in the data center database
		},
	}, nil
}

func (s *CarrierAPIBackend) GetLocalTaskDetailList(lastUpdate, pageSize uint64) ([]*carriertypespb.TaskDetail, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("query local identity failed, %s", err)
	}
	// the task is executing.
	localTaskList, err := s.carrier.carrierDB.QueryLocalTaskList()

	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.QueryTaskListByIdentityId(identity.GetIdentityId(), lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	return s.mergeLocalAndNetworkTasks(identity, localTaskList, networkTaskList), err
}

// v0.4.0
func (s *CarrierAPIBackend) GetGlobalTaskDetailList(lastUpdate, pageSize uint64) ([]*carriertypespb.TaskDetail, error) {

	_, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("query local identity failed, %s", err)
	}

	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.QueryGlobalTaskList(lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	makeTaskViewFn := func(task *types.Task) *carriertypespb.TaskDetail {
		return policy.NewTaskDetailShowFromTaskData(task)
	}
	result := make([]*carriertypespb.TaskDetail, 0)

	for _, networkTask := range networkTaskList {
		if taskView := makeTaskViewFn(networkTask); nil != taskView {
			result = append(result, taskView)
		}
	}
	return result, err
}

// v0.3.0
func (s *CarrierAPIBackend) GetTaskDetailListByTaskIds(taskIds []string) ([]*carriertypespb.TaskDetail, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("query local identity failed, %s", err)
	}

	// the task status is pending and running.
	localTaskList, err := s.carrier.carrierDB.QueryLocalTaskList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	// the task status was finished (failed or succeed).
	networkTaskList, err := s.carrier.carrierDB.QueryTaskListByTaskIds(taskIds)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	return s.mergeLocalAndNetworkTasks(identity, localTaskList, networkTaskList), err
}

func (s *CarrierAPIBackend) mergeLocalAndNetworkTasks(identity *carriertypespb.Organization, localTaskList, networkTaskList []*types.Task) []*carriertypespb.TaskDetail {

	makeTaskViewFn := func(task *types.Task) *carriertypespb.TaskDetail {
		return policy.NewTaskDetailShowFromTaskData(task)
	}

	networkTasks := make([]*carriertypespb.TaskDetail, 0)

	filterTaskIdCache := make(map[string]struct{}, 0)

	for _, networkTask := range networkTaskList {
		if taskView := makeTaskViewFn(networkTask); nil != taskView {
			networkTasks = append(networkTasks, taskView)
			filterTaskIdCache[taskView.GetInformation().GetTaskId()] = struct{}{}
		}
	}

	localTasks := make([]*carriertypespb.TaskDetail, 0)

next:
	for _, task := range localTaskList {

		// 1、If we are not task sender
		//
		// Filter out the local tasks belonging to the computing power provider that have not been started
		// (Note: the tasks under consensus are also tasks that have not been started)
		if identity.GetIdentityId() != task.GetTaskSender().GetIdentityId() {
			for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
				if identity.GetIdentityId() == powerSupplier.GetIdentityId() {
					running, err := s.carrier.carrierDB.HasLocalTaskExecuteStatusRunningByPartyId(task.GetTaskId(), powerSupplier.GetPartyId())
					if nil != err || !running {
						continue next // goto next task, if running one party of current identity is  still not executing this task.
					}
				}
			}
		}

		// 2、If we are task sender
		//
		// For the initiator's local task, when the task has not started execution
		// (i.e. the status is still: pending), the 'powersuppliers' of the task should not be returned.
		if identity.GetIdentityId() == task.GetTaskSender().GetIdentityId() && task.GetTaskData().GetState() == commonconstantpb.TaskState_TaskState_Pending {
			task.RemovePowerSuppliers() // clean powerSupplier when before return.
			task.RemovePowerResources()
		}

		// If the taskId already appears in the finished task,
		// the localtask to which the taskId belongs should be filtered out
		// (sometimes the localtask may not be cleaned up when the task is finished)
		if _, ok := filterTaskIdCache[task.GetTaskId()]; ok {
			continue
		}

		if taskView := makeTaskViewFn(task); nil != taskView {
			localTasks = append(localTasks, taskView)
		}
	}

	return append(localTasks, networkTasks...)
}

func (s *CarrierAPIBackend) GetTaskEventList(taskId string) ([]*carriertypespb.TaskEvent, error) {

	// 1、 If it is a local task, first find out the local eventList of the task
	eventList, err := s.carrier.carrierDB.QueryTaskEventList(taskId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}
	// 2、Then find out the eventList of the task in the data center
	taskEvents, err := s.carrier.carrierDB.QueryTaskEventListByTaskId(taskId)
	if nil != err {
		return nil, err
	}
	eventList = append(eventList, taskEvents...)
	return eventList, nil
}

func (s *CarrierAPIBackend) GetTaskEventListByTaskIds(taskIds []string) ([]*carriertypespb.TaskEvent, error) {

	evenList := make([]*carriertypespb.TaskEvent, 0)

	// 1、 If it is a local task, first find out the local eventList of the task
	for _, taskId := range taskIds {
		localEventList, err := s.carrier.carrierDB.QueryTaskEventList(taskId)
		if rawdb.IsNoDBNotFoundErr(err) {
			return nil, err
		}
		if rawdb.IsDBNotFoundErr(err) {
			continue
		}
		evenList = append(evenList, localEventList...)
	}
	// 2、Then find out the eventList of the task in the data center
	taskEvents, err := s.carrier.carrierDB.QueryTaskEventListByTaskIds(taskIds)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, taskEvents...)
	return evenList, nil
}

func (s *CarrierAPIBackend) HasLocalTask() (bool, error) {
	localTasks, err := s.carrier.carrierDB.QueryLocalTaskList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return false, err
	}
	if rawdb.IsDBNotFoundErr(err) {
		return false, nil
	}
	if len(localTasks) == 0 {
		return false, nil
	}
	return true, nil
}

// about jobResource
func (s *CarrierAPIBackend) QueryPowerRunningTaskList(powerId string) ([]string, error) {

	jobNodeId, err := s.carrier.carrierDB.QueryJobNodeIdByPowerId(powerId)
	if nil != err {
		log.WithError(err).Errorf("Failed query jobNodeId by powerId on CarrierAPIBackend.QueryPowerRunningTaskList(), powerId: {%s}", powerId)
		return nil, err
	}
	return s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(jobNodeId)
}

// about DataResourceTable
func (s *CarrierAPIBackend) StoreDataResourceTable(dataResourceTable *types.DataResourceTable) error {
	return s.carrier.carrierDB.StoreDataResourceTable(dataResourceTable)
}

func (s *CarrierAPIBackend) StoreDataResourceTables(dataResourceTables []*types.DataResourceTable) error {
	return s.carrier.carrierDB.StoreDataResourceTables(dataResourceTables)
}

func (s *CarrierAPIBackend) RemoveDataResourceTable(nodeId string) error {
	return s.carrier.carrierDB.RemoveDataResourceTable(nodeId)
}

func (s *CarrierAPIBackend) QueryDataResourceTable(nodeId string) (*types.DataResourceTable, error) {
	return s.carrier.carrierDB.QueryDataResourceTable(nodeId)
}

func (s *CarrierAPIBackend) QueryDataResourceTables() ([]*types.DataResourceTable, error) {
	return s.carrier.carrierDB.QueryDataResourceTables()
}

// about DataResourceDataUpload
func (s *CarrierAPIBackend) StoreDataResourceDataUpload(dataResourceDataUsed *types.DataResourceDataUpload) error {
	return s.carrier.carrierDB.StoreDataResourceDataUpload(dataResourceDataUsed)
}

func (s *CarrierAPIBackend) StoreDataResourceDataUploads(dataResourceDataUseds []*types.DataResourceDataUpload) error {
	return s.carrier.carrierDB.StoreDataResourceDataUploads(dataResourceDataUseds)
}

func (s *CarrierAPIBackend) RemoveDataResourceDataUpload(originId string) error {
	return s.carrier.carrierDB.RemoveDataResourceDataUpload(originId)
}

func (s *CarrierAPIBackend) QueryDataResourceDataUpload(originId string) (*types.DataResourceDataUpload, error) {
	return s.carrier.carrierDB.QueryDataResourceDataUpload(originId)
}

func (s *CarrierAPIBackend) QueryDataResourceDataUploads() ([]*types.DataResourceDataUpload, error) {
	return s.carrier.carrierDB.QueryDataResourceDataUploads()
}

func (s *CarrierAPIBackend) StoreTaskResultDataSummary(taskId, originId, dataHash, metadataOption, dataNodeId, extra string, dataType commonconstantpb.OrigindataType) error {
	// generate metadataId
	var buf bytes.Buffer
	buf.Write([]byte(originId))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	originIdHash := rlputil.RlpHash(buf.Bytes())

	metadataId := types.PREFIX_METADATA_ID + originIdHash.Hex()

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed query local identity on CarrierAPIBackend.StoreTaskResultDataSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, dataType: {%s}, metadataOption: %s",
			taskId, dataNodeId, originId, metadataId, dataType.String(), metadataOption)
		return err
	}

	// store local metadata (about task result file)
	metadata := types.NewMetadata(&carriertypespb.MetadataPB{
		/**
		MetadataId           string
		Owner                *Organization
		DataId               string
		DataStatus           DataStatus
		MetadataName         string
		MetadataType         MetadataType
		DataHash             string
		Desc                 string
		LocationType         DataLocationType
		DataType             OrigindataType
		Industry             string
		State                MetadataState
		PublishAt            uint64
		UpdateAt             uint64
		Nonce                uint64
		MetadataOption       string
		AllowExpose          bool
		TokenAddress         string
		*/
		MetadataId:     metadataId,
		Owner:          identity,
		DataId:         metadataId,
		DataStatus:     commonconstantpb.DataStatus_DataStatus_Valid,
		MetadataName:   fmt.Sprintf("task `%s` result file", taskId),
		MetadataType:   commonconstantpb.MetadataType_MetadataType_Unknown, // It means this is a module or psi result ??? so we don't known it.
		DataHash:       dataHash,
		Desc:           fmt.Sprintf("the task `%s` result file after executed", taskId),
		LocationType:   commonconstantpb.DataLocationType_DataLocationType_Local,
		DataType:       dataType,
		Industry:       "Unknown",
		State:          commonconstantpb.MetadataState_MetadataState_Created, // metaData status, eg: create/release/revoke
		PublishAt:      0,                                                    // have not publish
		UpdateAt:       timeutils.UnixMsecUint64(),
		Nonce:          0,
		MetadataOption: metadataOption,
		AllowExpose:    false,
		TokenAddress:   "",
	})
	s.carrier.carrierDB.StoreInternalMetadata(metadata)

	// todo whether need to store a dataResourceDiskUsed (metadataId. dataNodeId, diskUsed) ??? 后面需要上传 磁盘使用空间在弄吧

	// store dataResourceDataUpload (about task result file)
	err = s.carrier.carrierDB.StoreDataResourceDataUpload(types.NewDataResourceDataUpload(uint32(metadata.GetData().GetDataType()), dataNodeId, originId, metadataId, metadataOption, dataHash))
	if nil != err {
		log.WithError(err).Errorf("Failed store dataResourceDataUpload about task result file on CarrierAPIBackend.StoreTaskResultDataSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, dataType: {%s}, metadataOption: %s",
			taskId, dataNodeId, originId, metadataId, metadata.GetData().GetDataType(), metadataOption)
		return err
	}
	// 记录原始数据占用资源大小   StoreDataResourceTable  todo 后续考虑是否加上, 目前不加 因为对于系统生成的元数据暂时不需要记录 disk 使用实况 ??
	// 单独记录 metaData 的 GetSize 和所在 dataNodeId   StoreDataResourceDiskUsed  todo 后续考虑是否加上, 目前不加 因为对于系统生成的元数据暂时不需要记录 disk 使用实况 ??

	// store taskId -> TaskUpResultData (about task result file)
	err = s.carrier.carrierDB.StoreTaskUpResultData(types.NewTaskUpResultData(taskId, originId, metadataId, extra))
	if nil != err {
		log.WithError(err).Errorf("Failed store taskUpResultData on CarrierAPIBackend.StoreTaskResultDataSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, dataType: {%s}, metadataOption: %s",
			taskId, dataNodeId, originId, metadataId, metadata.GetData().GetDataType(), metadataOption)
		return err
	}
	return nil
}

func (s *CarrierAPIBackend) QueryTaskResultDataSummary(taskId string) (*types.TaskResultDataSummary, error) {
	summarry, err := s.carrier.carrierDB.QueryTaskUpResulData(taskId)
	if nil != err {
		log.WithError(err).Errorf("Failed query taskUpResultData on CarrierAPIBackend.QueryTaskResultDataSummary(), taskId: {%s}", taskId)
		return nil, err
	}
	dataResourceDataUpload, err := s.carrier.carrierDB.QueryDataResourceDataUpload(summarry.GetOriginId())
	if nil != err {
		log.WithError(err).Errorf("Failed query dataResourceDataUpload on CarrierAPIBackend.QueryTaskResultDataSummary(), taskId: {%s}, originId: {%s}",
			taskId, summarry.GetOriginId())
		return nil, err
	}

	localMetadata, err := s.carrier.carrierDB.QueryInternalMetadataById(summarry.GetMetadataId())
	if nil != err {
		log.WithError(err).Errorf("Failed query local metadata on CarrierAPIBackend.QueryTaskResultDataSummary(), taskId: {%s}, originId: {%s}, metadataId: {%s}",
			taskId, summarry.GetOriginId(), dataResourceDataUpload.GetMetadataId())
		return nil, err
	}

	// taskId, metadataId, originId, metadataName, dataHash, metadataOption, nodeId, extra string, dataType uint32
	return types.NewTaskResultDataSummary(
		summarry.GetTaskId(),
		dataResourceDataUpload.GetMetadataId(),
		dataResourceDataUpload.GetOriginId(),
		localMetadata.GetData().GetMetadataName(),
		dataResourceDataUpload.GetDataHash(),
		dataResourceDataUpload.GetMetadataOption(),
		dataResourceDataUpload.GetNodeId(),
		summarry.GetExtra(),
		dataResourceDataUpload.GetDataType(),
	), nil

}

func (s *CarrierAPIBackend) QueryTaskResultDataSummaryList() (types.TaskResultDataSummaryArr, error) {
	taskResultDataSummaryArr, err := s.carrier.carrierDB.QueryTaskUpResultDataList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed query all taskUpResultData on CarrierAPIBackend.QueryTaskResultDataSummaryList()")
		return nil, err
	}

	arr := make(types.TaskResultDataSummaryArr, 0)
	for _, summarry := range taskResultDataSummaryArr {
		dataResourceDataUpload, err := s.carrier.carrierDB.QueryDataResourceDataUpload(summarry.GetOriginId())
		if nil != err {
			log.WithError(err).Errorf("Failed query dataResourceDataUpload on CarrierAPIBackend.QueryTaskResultDataSummaryList(), taskId: {%s}, originId: {%s}",
				summarry.GetTaskId(), summarry.GetOriginId())
			continue
		}

		localMetadata, err := s.carrier.carrierDB.QueryInternalMetadataById(dataResourceDataUpload.GetMetadataId())
		if nil != err {
			log.WithError(err).Errorf("Failed query local metadata on CarrierAPIBackend.QueryTaskResultDataSummaryList(), taskId: {%s}, originId: {%s}, metadataId: {%s}",
				summarry.GetTaskId(), summarry.GetOriginId(), dataResourceDataUpload.GetMetadataId())
			continue
		}
		// taskId, metadataId, originId, metadataName, dataHash, metadataOption, nodeId, extra string, dataType uint32
		arr = append(arr, types.NewTaskResultDataSummary(
			summarry.GetTaskId(),
			dataResourceDataUpload.GetMetadataId(),
			dataResourceDataUpload.GetOriginId(),
			localMetadata.GetData().GetMetadataName(),
			dataResourceDataUpload.GetDataHash(),
			dataResourceDataUpload.GetMetadataOption(),
			dataResourceDataUpload.GetNodeId(),
			summarry.GetExtra(),
			dataResourceDataUpload.GetDataType(),
		))
	}

	return arr, nil
}

func (s *CarrierAPIBackend) EstimateTaskGas(taskSponsorAddress string, tokenItemList []*carrierapipb.TkItem) (gasLimit uint64, gasPrice *big.Int, err error) {
	gasLimit, gasPrice, err = s.carrier.datumPayManager.EstimateTaskGas(taskSponsorAddress, tokenItemList)
	if err != nil {
		log.WithError(err).Error("Failed to call EstimateTaskGas() on CarrierAPIBackend.EstimateTaskGas()")
	}
	return
}

// EstimateTaskGas
