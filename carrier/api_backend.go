package carrier

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/bytesutil"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/rawdb"
	"github.com/RosettaFlow/Carrier-Go/grpclient"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/computesvc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/RosettaFlow/Carrier-Go/policy"
	"github.com/RosettaFlow/Carrier-Go/types"
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

func (s *CarrierAPIBackend) GetCarrierChainConfig() *types.CarrierChainConfig {
	return params.CarrierConfig()
}

func (s *CarrierAPIBackend) SendMsg(msg types.Msg) error {
	return s.carrier.mempool.Add(msg)
}

// system (the yarn node self info)
func (s *CarrierAPIBackend) GetNodeInfo() (*pb.YarnNodeInfo, error) {

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

	nodeInfo := &pb.YarnNodeInfo{
		NodeType:           pb.NodeType_NodeType_YarnNode,
		NodeId:             nodeId,
		IdentityType:       types.IDENTITY_TYPE_DID, // default: DID
		IdentityId:         identityId,
		Name:               nodeName,
		State:              pb.YarnNodeState_State_Active,
		RelatePeers:        uint32(len(s.carrier.config.P2P.Peers().Active())),
		LocalBootstrapNode: enrStr,
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

	if addr, err := s.carrier.metisPayManager.QueryOrgWallet(); err == nil {
		nodeInfo.ObserverProxyWalletAddress = addr
	} else {
		log.WithError(err).Errorf("cannot load organization wallet of node info: %v", err)
		return nil, err
	}

	return nodeInfo, nil
}

func (s *CarrierAPIBackend) SetSeedNode(seed *pb.SeedPeer) (pb.ConnState, error) {
	// format: enr:-xxxxxx
	if err := s.carrier.config.P2P.AddPeer(seed.GetAddr()); nil != err {
		log.WithError(err).Errorf("Failed to call p2p.AddPeer() with seed.Addr on SetSeedNode(), addr: {%s}", seed.GetAddr())
		return pb.ConnState_ConnState_UnConnected, err
	}
	if err := s.carrier.carrierDB.SetSeedNode(seed); nil != err {
		log.WithError(err).Errorf("Failed to call SetSeedNode() to store seedNode on SetSeedNode(), seed: {%s}", seed.String())
		return pb.ConnState_ConnState_UnConnected, err
	}
	addrs, err := s.carrier.config.P2P.PeerFromAddress([]string{seed.GetAddr()})
	if err != nil || len(addrs) == 0 {
		log.WithError(err).Errorf("Failed to parse addr")
		return pb.ConnState_ConnState_UnConnected, err
	}
	addrInfo, err := peer.AddrInfoFromP2pAddr(addrs[0])
	if nil != err {
		log.WithError(err).Errorf("Failed to call peer.AddrInfoFromP2pAddr() with multiAddr on SetSeedNode(), addr: {%s}", addrs[0].String())
		return pb.ConnState_ConnState_UnConnected, err
	}
	for _, active := range s.carrier.config.P2P.Peers().Active() {
		if active.String() == addrInfo.ID.String() {
			return pb.ConnState_ConnState_Connected, nil
		}
	}
	return pb.ConnState_ConnState_UnConnected, nil
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

func (s *CarrierAPIBackend) GetSeedNodeList() ([]*pb.SeedPeer, error) {
	// load seed node from default bootstrap
	bootstrapNodeStrs, err := s.carrier.config.P2P.BootstrapAddresses()
	if nil != err {
		log.WithError(err).Errorf("Failed to call p2p.BootstrapAddresses() on GetSeedNodeList()")
		return nil, err
	}

	bootstrapNodes := make([]*pb.SeedPeer, len(bootstrapNodeStrs))
	for i, addr := range bootstrapNodeStrs {
		bootstrapNodes[i] = &pb.SeedPeer{
			Addr:      addr,
			IsDefault: true,
			ConnState: pb.ConnState_ConnState_UnConnected,
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
			seeds[i].ConnState = pb.ConnState_ConnState_Connected
		}
	}
	return seeds, nil
}

func (s *CarrierAPIBackend) storeLocalResource(identity *libtypes.Organization, jobNodeId string, jobNodeStatus *computesvc.GetStatusReply) error {

	// store into local db
	if err := s.carrier.carrierDB.StoreLocalResource(types.NewLocalResource(&libtypes.LocalResourcePB{
		Owner:      identity,
		JobNodeId:  jobNodeId,
		DataId:     "", // can not own powerId now, because power have not publish
		DataStatus: libtypes.DataStatus_DataStatus_Valid,
		// resource status, eg: create/release/revoke
		State: libtypes.PowerState_PowerState_Created,
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

func (s *CarrierAPIBackend) SetRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {
		client, err := grpclient.NewJobNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreJobNodeClient(node.GetId(), client)

		jobNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect jobNode query status failed, %s", err)
		}
		// add resource usage first, but not own power now (mem, proccessor, bandwidth)
		if err = s.storeLocalResource(identity, node.GetId(), jobNodeStatus); nil != err {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store jobNode local resource failed, %s", err)
		}
	}
	if typ == pb.PrefixTypeDataNode {
		client, err := grpclient.NewDataNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreDataNodeClient(node.GetId(), client)

		dataNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect dataNode query status failed, %s", err)
		}
		// add data resource  (disk)
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk(), true))
		//err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
		log.Debugf("Add dataNode resource table succeed, nodeId: {%s}, totalDisk: {%d}, usedDisk: {%d}", node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk())
	}
	node.ConnState = pb.ConnState_ConnState_Connected
	if err = s.carrier.carrierDB.SetRegisterNode(typ, node); err != nil {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("Store registerNode to db failed, %s", err)
	}
	return pb.ConnState_ConnState_Connected, nil
}

func (s *CarrierAPIBackend) UpdateRegisterNode(typ pb.RegisteredNodeType, node *pb.YarnRegisteredPeerDetail) (pb.ConnState, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {

		// The published jobNode cannot be updated directly
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(node.GetId())
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on UpdateRegisterNode, %s", resourceTable.String())
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("still have the published computing power information on old jobNode failed, input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				node.Id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// First check whether there is a task being executed on jobNode
		runningTaskCount, err := s.carrier.carrierDB.QueryJobNodeRunningTaskCount(node.GetId())
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceManager.QueryJobNodeClient(node.GetId()); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceManager.RemoveJobNodeClient(node.GetId())
		}

		// generate new client
		client, err := grpclient.NewJobNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreJobNodeClient(node.GetId(), client)

		jobNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect jobNode query status failed, %s", err)
		}
		// add resource usage first, but not own power now (mem, proccessor, bandwidth)
		if err = s.storeLocalResource(identity, node.GetId(), jobNodeStatus); nil != err {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store jobNode local resource failed, %s", err)
		}

	}

	if typ == pb.PrefixTypeDataNode {

		if client, ok := s.carrier.resourceManager.QueryDataNodeClient(node.GetId()); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceManager.RemoveDataNodeClient(node.GetId())
		}

		// remove old data resource  (disk)
		if err := s.carrier.carrierDB.RemoveDataResourceTable(node.GetId()); rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("remove disk summary of old dataNode, %s", err)
		}

		client, err := grpclient.NewDataNodeClient(s.carrier.ctx, fmt.Sprintf("%s:%s", node.GetInternalIp(), node.GetInternalPort()), node.GetId())
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceManager.StoreDataNodeClient(node.GetId(), client)

		dataNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect dataNode query status failed, %s", err)
		}
		// add new data resource  (disk)
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk(), true))
		//err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
		log.Debugf("Update dataNode resource table succeed, nodeId: {%s}, totalDisk: {%d}, usedDisk: {%d}", node.GetId(), dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk())
	}

	// remove  old jobNode from db
	if err := s.carrier.carrierDB.DeleteRegisterNode(typ, node.GetId()); nil != err {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("remove old registerNode from db failed, %s", err)
	}

	// add new node to db
	node.ConnState = pb.ConnState_ConnState_Connected
	if err = s.carrier.carrierDB.SetRegisterNode(typ, node); err != nil {
		return pb.ConnState_ConnState_UnConnected, fmt.Errorf("update registerNode to db failed, %s", err)
	}
	return pb.ConnState_ConnState_Connected, nil
}

func (s *CarrierAPIBackend) DeleteRegisterNode(typ pb.RegisteredNodeType, id string) error {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return fmt.Errorf("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {

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
			if err := s.carrier.carrierDB.RevokeResource(types.NewResource(&libtypes.ResourcePB{
				Owner:  identity,
				DataId: resourceTable.GetPowerId(),
				// the status of data, N means normal, D means deleted.
				DataStatus: libtypes.DataStatus_DataStatus_Invalid,
				// resource status, eg: create/release/revoke
				State:    libtypes.PowerState_PowerState_Revoked,
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

	if typ == pb.PrefixTypeDataNode {

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

func (s *CarrierAPIBackend) GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error) {
	node, err := s.carrier.carrierDB.QueryRegisterNode(typ, id)
	if nil != err {
		return nil, err
	}

	if typ == pb.PrefixTypeJobNode {

		client, ok := s.carrier.resourceManager.QueryJobNodeClient(id)
		if !ok {
			node.ConnState = pb.ConnState_ConnState_UnConnected
		} else {
			if !client.IsConnected() {
				node.ConnState = pb.ConnState_ConnState_UnConnected
			} else {
				node.ConnState = pb.ConnState_ConnState_Connected
			}
		}

	} else {

		client, ok := s.carrier.resourceManager.QueryDataNodeClient(id)
		if !ok {
			node.ConnState = pb.ConnState_ConnState_UnConnected
		} else {
			if !client.IsConnected() {
				node.ConnState = pb.ConnState_ConnState_UnConnected
			} else {
				node.ConnState = pb.ConnState_ConnState_Connected
			}
		}

	}
	return node, nil
}

func (s *CarrierAPIBackend) GetRegisterNodeList(typ pb.RegisteredNodeType) ([]*pb.YarnRegisteredPeerDetail, error) {

	if typ != pb.PrefixTypeJobNode &&
		typ != pb.PrefixTypeDataNode {

		return nil, fmt.Errorf("Invalid nodeType")
	}

	nodeList, err := s.carrier.carrierDB.QueryRegisterNodeList(typ)
	if nil != err {
		return nil, err
	}

	for i, n := range nodeList {

		var connState pb.ConnState
		var duration uint64
		var taskCount uint32
		var taskIdList []string
		var fileCount uint32
		var fileTotalSize uint32

		if typ == pb.PrefixTypeJobNode {

			client, ok := s.carrier.resourceManager.QueryJobNodeClient(n.GetId())
			if !ok {
				connState = pb.ConnState_ConnState_UnConnected
			} else {
				duration = uint64(client.RunningDuration())
				if !client.IsConnected() {
					connState = pb.ConnState_ConnState_UnConnected
				} else {
					connState = pb.ConnState_ConnState_Connected
				}
			}

			taskCount, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskCount(n.GetId())
			taskIdList, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(n.GetId())
			fileCount = 0     // todo need implament this logic
			fileTotalSize = 0 // todo need implament this logic
		} else {

			client, ok := s.carrier.resourceManager.QueryDataNodeClient(n.GetId())
			if !ok {
				connState = pb.ConnState_ConnState_UnConnected
			} else {
				duration = uint64(client.RunningDuration())
				if !client.IsConnected() {
					connState = pb.ConnState_ConnState_UnConnected
				} else {
					connState = pb.ConnState_ConnState_Connected
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

func (s *CarrierAPIBackend) SendTaskEvent(event *libtypes.TaskEvent) error {
	return s.carrier.TaskManager.SendTaskEvent(event)
}

func (s *CarrierAPIBackend) ReportTaskResourceUsage(nodeType pb.NodeType, ip, port string, usage *types.TaskResuorceUsage) error {

	if err := s.carrier.TaskManager.HandleReportResourceUsage(usage); nil != err {
		log.WithError(err).Errorf("Failed to call HandleReportResourceUsage() on CarrierAPIBackend.ReportTaskResourceUsage(), taskId: {%s}, partyId: {%s}, nodeType: {%s}, ip:{%s}, port: {%s}",
			usage.GetTaskId(), usage.GetPartyId(), nodeType.String(), ip, port)
		return err
	}
	return nil
}

func (s *CarrierAPIBackend) GenerateObServerProxyWalletAddress() (string, error) {
	if s.carrier.metisPayManager != nil {
		if addr, err := s.carrier.metisPayManager.GenerateOrgWallet(); nil != err {
			log.WithError(err).Error("Failed to call GenerateOrgWallet() on CarrierAPIBackend.GenerateObServerProxyWalletAddress()")
			return "", err
		} else {
			log.Debugf("Success to generate organization wallet %s", addr)
			return addr, nil
		}
	} else {
		return "", errors.New("MetisPay manager not initialized properly")
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
func (s *CarrierAPIBackend) GetGlobalMetadataDetailList(lastUpdate, pageSize uint64) ([]*pb.GetGlobalMetadataDetail, error) {
	log.Debug("Invoke: GetGlobalMetadataDetailList executing...")
	var (
		arr []*pb.GetGlobalMetadataDetail
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
func (s *CarrierAPIBackend) GetGlobalMetadataDetailListByIdentityId(identityId string, lastUpdate, pageSize uint64) ([]*pb.GetGlobalMetadataDetail, error) {
	log.Debug("Invoke: GetGlobalMetadataDetailListByIdentityId executing...")
	var (
		arr []*pb.GetGlobalMetadataDetail
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

func (s *CarrierAPIBackend) GetLocalMetadataDetailList(lastUpdate uint64, pageSize uint64) ([]*pb.GetLocalMetadataDetail, error) {
	log.Debug("Invoke: GetLocalMetadataDetailList executing...")

	var (
		arr []*pb.GetLocalMetadataDetail
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

func (s *CarrierAPIBackend) GetLocalInternalMetadataDetailList() ([]*pb.GetLocalMetadataDetail, error) {
	log.Debug("Invoke: GetLocalInternalMetadataDetailList executing...")

	var (
		arr []*pb.GetLocalMetadataDetail
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

func (s *CarrierAPIBackend) GetGlobalPowerSummaryList() ([]*pb.GetGlobalPowerSummary, error) {
	log.Debug("Invoke: GetGlobalPowerSummaryList executing...")
	resourceList, err := s.carrier.carrierDB.QueryGlobalResourceSummaryList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power summary list, len: {%d}", len(resourceList))
	powerList := make([]*pb.GetGlobalPowerSummary, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &pb.GetGlobalPowerSummary{
			Owner: resource.GetOwner(),
			Information: &libtypes.PowerUsageDetail{
				TotalTaskCount:   0,
				CurrentTaskCount: 0,
				Tasks:            make([]*libtypes.PowerTask, 0),
				Information: &libtypes.ResourceUsageOverview{
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

func (s *CarrierAPIBackend) GetGlobalPowerDetailList(lastUpdate uint64, pageSize uint64) ([]*pb.GetGlobalPowerDetail, error) {
	log.Debug("Invoke: GetGlobalPowerDetailList executing...")
	resourceList, err := s.carrier.carrierDB.QueryGlobalResourceDetailList(lastUpdate, pageSize)
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power detail list, len: {%d}", len(resourceList))
	powerList := make([]*pb.GetGlobalPowerDetail, 0, resourceList.Len())
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

		powerList = append(powerList, &pb.GetGlobalPowerDetail{
			Owner:   resource.GetOwner(),
			PowerId: resource.GetDataId(),
			Information: &libtypes.PowerUsageDetail{
				TotalTaskCount:   totalTaskCount,
				CurrentTaskCount: currentTaskCount,
				Tasks:            make([]*libtypes.PowerTask, 0),
				Information: &libtypes.ResourceUsageOverview{
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
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetLocalPowerDetailList() ([]*pb.GetLocalPowerDetail, error) {
	log.Debug("Invoke:GetLocalPowerDetailList executing...")
	// query local resource list from db.
	machineList, err := s.carrier.carrierDB.QueryLocalResourceList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetLocalPowerDetailList, call QueryLocalResourceList, machineList: %s", machineList.String())

	buildPowerTaskList := func(jobNodeId string) []*libtypes.PowerTask {

		powerTaskList := make([]*libtypes.PowerTask, 0)

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
			powerTask := &libtypes.PowerTask{
				TaskId:   taskId,
				TaskName: task.GetTaskData().TaskName,
				Owner: &libtypes.Organization{
					NodeName:   task.GetTaskSender().GetNodeName(),
					NodeId:     task.GetTaskSender().GetNodeId(),
					IdentityId: task.GetTaskSender().GetIdentityId(),
				},
				Receivers: make([]*libtypes.Organization, 0),
				OperationCost: &libtypes.TaskResourceCostDeclare{
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
				powerTask.Partners = append(powerTask.GetPartners(), &libtypes.Organization{
					NodeName:   dataSupplier.GetNodeName(),
					NodeId:     dataSupplier.GetNodeId(),
					IdentityId: dataSupplier.GetIdentityId(),
				})
			}
			// build receivers of task info
			for _, receiver := range task.GetTaskData().GetReceivers() {
				powerTask.Receivers = append(powerTask.GetReceivers(), &libtypes.Organization{
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
			powerTask.OperationSpend = &libtypes.TaskResourceCostDeclare{
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
	result := make([]*pb.GetLocalPowerDetail, len(resourceList))
	for i, resource := range resourceList {

		nodePowerDetail := &pb.GetLocalPowerDetail{
			JobNodeId: resource.GetJobNodeId(),
			PowerId:   resource.GetDataId(),
			Owner:     resource.GetOwner(),
			Power: &libtypes.PowerUsageDetail{
				TotalTaskCount:   taskTotalCount(resource.GetJobNodeId()),
				CurrentTaskCount: taskRunningCount(resource.GetJobNodeId()),
				Tasks:            make([]*libtypes.PowerTask, 0),
				Information: &libtypes.ResourceUsageOverview{
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
	return types.NewIdentity(&libtypes.IdentityPB{
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

func (s *CarrierAPIBackend) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (libtypes.AuditMetadataOption, error) {
	return s.carrier.authManager.AuditMetadataAuthority(audit)
}

func (s *CarrierAPIBackend) GetLocalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	return s.carrier.authManager.GetLocalMetadataAuthorityList(lastUpdate, pageSize)
}

func (s *CarrierAPIBackend) GetGlobalMetadataAuthorityList(lastUpdate, pageSize uint64) (types.MetadataAuthArray, error) {
	return s.carrier.authManager.GetGlobalMetadataAuthorityList(lastUpdate, pageSize)
}

func (s *CarrierAPIBackend) HasValidMetadataAuth(userType libtypes.UserType, user, identityId, metadataId string) (bool, error) {
	return s.carrier.authManager.HasValidMetadataAuth(userType, user, identityId, metadataId)
}

// task api
func (s *CarrierAPIBackend) GetLocalTask(taskId string) (*pb.TaskDetailShow, error) {
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

	detailShow := &pb.TaskDetailShow{
		TaskId:   localTask.GetTaskId(),
		TaskName: localTask.GetTaskData().GetTaskName(),
		UserType: localTask.GetTaskData().GetUserType(),
		User:     localTask.GetTaskData().GetUser(),
		Sender: &libtypes.TaskOrganization{
			PartyId:    localTask.GetTaskSender().GetPartyId(),
			NodeName:   localTask.GetTaskSender().GetNodeName(),
			NodeId:     localTask.GetTaskSender().GetNodeId(),
			IdentityId: localTask.GetTaskSender().GetIdentityId(),
		},
		//AlgoSupplier:
		DataSuppliers:  make([]*pb.TaskDataSupplierShow, 0, len(localTask.GetTaskData().GetDataSuppliers())),
		PowerSuppliers: make([]*pb.TaskPowerSupplierShow, 0, len(localTask.GetTaskData().GetPowerSuppliers())),
		Receivers:      localTask.GetTaskData().GetReceivers(),
		CreateAt:       localTask.GetTaskData().GetCreateAt(),
		StartAt:        localTask.GetTaskData().GetStartAt(),
		EndAt:          localTask.GetTaskData().GetEndAt(),
		State:          localTask.GetTaskData().GetState(),
		OperationCost: &libtypes.TaskResourceCostDeclare{
			Processor: localTask.GetTaskData().GetOperationCost().GetProcessor(),
			Memory:    localTask.GetTaskData().GetOperationCost().GetMemory(),
			Bandwidth: localTask.GetTaskData().GetOperationCost().GetBandwidth(),
			Duration:  localTask.GetTaskData().GetOperationCost().GetDuration(),
		},
	}

	//
	detailShow.AlgoSupplier = &pb.TaskAlgoSupplier{
		Organization: &libtypes.TaskOrganization{
			PartyId:    localTask.GetTaskData().GetAlgoSupplier().GetPartyId(),
			NodeName:   localTask.GetTaskData().GetAlgoSupplier().GetNodeName(),
			NodeId:     localTask.GetTaskData().GetAlgoSupplier().GetNodeId(),
			IdentityId: localTask.GetTaskData().GetAlgoSupplier().GetIdentityId(),
		},
		//MetaAlgorithmId: "", todo
		//MetaAlgorithmName: "", Organization
	}

	// DataSupplier
	for _, dataSupplier := range localTask.GetTaskData().GetDataSuppliers() {

		metadataId, err := policy.FetchMetedataIdByPartyId(dataSupplier.GetPartyId(), localTask.GetTaskData().GetDataPolicyType(), localTask.GetTaskData().GetDataPolicyOption())
		if nil != err {
			log.Errorf("not fetch metadataId of local task dataPolicy on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}, patyId: {%s}", taskId, dataSupplier.GetPartyId())
			return nil, fmt.Errorf("not fetch metadataId")
		}
		metadataName, err := policy.FetchMetedataNameByPartyId(dataSupplier.GetPartyId(), localTask.GetTaskData().GetDataPolicyType(), localTask.GetTaskData().GetDataPolicyOption())
		if nil != err {
			log.Errorf("not fetch metadataName of local task dataPolicy on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}, patyId: {%s}", taskId, dataSupplier.GetPartyId())
			return nil, fmt.Errorf("not fetch metadataName")
		}

		detailShow.DataSuppliers = append(detailShow.DataSuppliers, &pb.TaskDataSupplierShow{
			Organization: &libtypes.TaskOrganization{
				PartyId:    dataSupplier.GetPartyId(),
				NodeName:   dataSupplier.GetNodeName(),
				NodeId:     dataSupplier.GetNodeId(),
				IdentityId: dataSupplier.GetIdentityId(),
			},
			MetadataId:   metadataId,
			MetadataName: metadataName,
		})
	}
	// powerSupplier
	powerResourceCache := make(map[string]*libtypes.TaskPowerResourceOption, 0)
	for _, option := range localTask.GetTaskData().GetPowerResourceOptions() {
		powerResourceCache[option.GetPartyId()] = option
	}
	for _, data := range localTask.GetTaskData().GetPowerSuppliers() {

		option, ok := powerResourceCache[data.GetPartyId()]
		if !ok {
			log.Errorf("not found power resource option of local task dataPolicy on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}, patyId: {%s}", taskId, data.GetPartyId())
			return nil, fmt.Errorf("not found power resource option")
		}

		detailShow.PowerSuppliers = append(detailShow.PowerSuppliers, &pb.TaskPowerSupplierShow{
			Organization: &libtypes.TaskOrganization{
				PartyId:    data.GetPartyId(),
				NodeName:   data.GetNodeName(),
				NodeId:     data.GetNodeId(),
				IdentityId: data.GetIdentityId(),
			},
			PowerInfo: &libtypes.ResourceUsageOverview{
				TotalMem:       option.GetResourceUsedOverview().GetTotalMem(),
				UsedMem:        option.GetResourceUsedOverview().GetUsedMem(),
				TotalProcessor: option.GetResourceUsedOverview().GetTotalProcessor(),
				UsedProcessor:  option.GetResourceUsedOverview().GetUsedProcessor(),
				TotalBandwidth: option.GetResourceUsedOverview().GetTotalBandwidth(),
				UsedBandwidth:  option.GetResourceUsedOverview().GetUsedBandwidth(),
				TotalDisk:      option.GetResourceUsedOverview().GetTotalDisk(),
				UsedDisk:       option.GetResourceUsedOverview().GetUsedDisk(),
			},
		})
	}

	return detailShow, nil
}

func (s *CarrierAPIBackend) GetLocalTaskDetailList(lastUpdate, pageSize uint64) ([]*pb.TaskDetailShow, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("query local identity failed, %s", err)
	}
	// the task is executing.
	localTaskArray, err := s.carrier.carrierDB.QueryLocalTaskList()

	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.QueryTaskListByIdentityId(identity.GetIdentityId(), lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	makeTaskViewFn := func(task *types.Task) *pb.TaskDetailShow {
		return policy.NewTaskDetailShowFromTaskData(task)
	}

	result := make([]*pb.TaskDetailShow, 0)

next:
	for _, task := range localTaskArray {

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

		// For the initiator's local task, when the task has not started execution
		// (i.e. the status is still: pending), the 'powersuppliers' of the task should not be returned.
		if identity.GetIdentityId() == task.GetTaskSender().GetIdentityId() && task.GetTaskData().GetState() == libtypes.TaskState_TaskState_Pending {
			task.RemovePowerSuppliers() // clean powerSupplier when before return.
		}

		if taskView := makeTaskViewFn(task); nil != taskView {
			result = append(result, taskView)
		}
	}
	for _, networkTask := range networkTaskList {
		if taskView := makeTaskViewFn(networkTask); nil != taskView {
			result = append(result, taskView)
		}
	}
	return result, err
}

// v0.4.0
func (s *CarrierAPIBackend) GetGlobalTaskDetailList(lastUpdate, pageSize uint64) ([]*pb.TaskDetailShow, error) {

	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.QueryGlobalTaskList(lastUpdate, pageSize)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	makeTaskViewFn := func(task *types.Task) *pb.TaskDetailShow {
		return policy.NewTaskDetailShowFromTaskData(task)
	}
	result := make([]*pb.TaskDetailShow, 0)

	for _, networkTask := range networkTaskList {
		if taskView := makeTaskViewFn(networkTask); nil != taskView {
			result = append(result, taskView)
		}
	}
	return result, err
}

// v0.3.0
func (s *CarrierAPIBackend) GetTaskDetailListByTaskIds(taskIds []string) ([]*pb.TaskDetailShow, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, fmt.Errorf("query local identity failed, %s", err)
	}

	// the task status is pending and running.
	localTaskList, err := s.carrier.carrierDB.QueryLocalTaskList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	taskIdCache := make(map[string]struct{}, 0)
	for _, taskId := range taskIds {
		taskIdCache[taskId] = struct{}{}
	}

	for _, task := range localTaskList {
		if _, ok := taskIdCache[task.GetTaskId()]; ok {
			delete(taskIdCache, task.GetTaskId())
		}
	}

	var ids []string
	if len(taskIds) == len(taskIdCache) {
		ids = taskIds
	} else {
		for taskId, _ := range taskIdCache {
			ids = append(ids, taskId)
		}
	}

	// the task status was finished (failed or succeed).
	networkTaskList, err := s.carrier.carrierDB.QueryTaskListByTaskIds(ids)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	makeTaskViewFn := func(task *types.Task) *pb.TaskDetailShow {
		return policy.NewTaskDetailShowFromTaskData(task)
	}

	result := make([]*pb.TaskDetailShow, 0)

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
		if identity.GetIdentityId() == task.GetTaskSender().GetIdentityId() && task.GetTaskData().GetState() == libtypes.TaskState_TaskState_Pending {
			task.RemovePowerSuppliers() // clean powerSupplier when before return.
		}

		if taskView := makeTaskViewFn(task); nil != taskView {
			result = append(result, taskView)
		}
	}
	for _, networkTask := range networkTaskList {
		if taskView := makeTaskViewFn(networkTask); nil != taskView {
			result = append(result, taskView)
		}
	}
	return result, err
}

func (s *CarrierAPIBackend) GetTaskEventList(taskId string) ([]*pb.TaskEventShow, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, err
	}

	// 1、 If it is a local task, first find out the local eventList of the task
	localEventList, err := s.carrier.carrierDB.QueryTaskEventList(taskId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	evenList := make([]*pb.TaskEventShow, len(localEventList))
	for i, e := range localEventList {
		evenList[i] = &pb.TaskEventShow{
			TaskId:   e.GetTaskId(),
			Type:     e.GetType(),
			CreateAt: e.GetCreateAt(),
			Content:  e.GetContent(),
			Owner: &libtypes.Organization{
				NodeName:   identity.GetNodeName(),
				NodeId:     identity.GetNodeId(),
				IdentityId: identity.GetIdentityId(),
			},
			PartyId: e.GetPartyId(),
		}
	}
	// 2、Then find out the eventList of the task in the data center
	taskEvent, err := s.carrier.carrierDB.QueryTaskEventListByTaskId(taskId)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, policy.NewTaskEventFromAPIEvent(taskEvent)...)
	return evenList, nil
}

func (s *CarrierAPIBackend) GetTaskEventListByTaskIds(taskIds []string) ([]*pb.TaskEventShow, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, err
	}

	evenList := make([]*pb.TaskEventShow, 0)

	// 1、 If it is a local task, first find out the local eventList of the task
	for _, taskId := range taskIds {
		localEventList, err := s.carrier.carrierDB.QueryTaskEventList(taskId)
		if rawdb.IsNoDBNotFoundErr(err) {
			return nil, err
		}
		if rawdb.IsDBNotFoundErr(err) {
			continue
		}
		for _, e := range localEventList {
			evenList = append(evenList, &pb.TaskEventShow{
				TaskId:   e.GetTaskId(),
				Type:     e.GetType(),
				CreateAt: e.GetCreateAt(),
				Content:  e.GetContent(),
				Owner: &libtypes.Organization{
					NodeName:   identity.GetNodeName(),
					NodeId:     identity.GetNodeId(),
					IdentityId: identity.GetIdentityId(),
				},
				PartyId: e.GetPartyId(),
			})
		}
	}
	// 2、Then find out the eventList of the task in the data center
	taskEvent, err := s.carrier.carrierDB.QueryTaskEventListByTaskIds(taskIds)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, policy.NewTaskEventFromAPIEvent(taskEvent)...)
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

// about DataResourceFileUpload
func (s *CarrierAPIBackend) StoreDataResourceFileUpload(dataResourceDataUsed *types.DataResourceFileUpload) error {
	return s.carrier.carrierDB.StoreDataResourceFileUpload(dataResourceDataUsed)
}

func (s *CarrierAPIBackend) StoreDataResourceFileUploads(dataResourceDataUseds []*types.DataResourceFileUpload) error {
	return s.carrier.carrierDB.StoreDataResourceFileUploads(dataResourceDataUseds)
}

func (s *CarrierAPIBackend) RemoveDataResourceFileUpload(originId string) error {
	return s.carrier.carrierDB.RemoveDataResourceFileUpload(originId)
}

func (s *CarrierAPIBackend) QueryDataResourceFileUpload(originId string) (*types.DataResourceFileUpload, error) {
	return s.carrier.carrierDB.QueryDataResourceFileUpload(originId)
}

func (s *CarrierAPIBackend) QueryDataResourceFileUploads() ([]*types.DataResourceFileUpload, error) {
	return s.carrier.carrierDB.QueryDataResourceFileUploads()
}

func (s *CarrierAPIBackend) StoreTaskUpResultFile(turf *types.TaskUpResultFile) error {
	return s.carrier.carrierDB.StoreTaskUpResultFile(turf)
}

func (s *CarrierAPIBackend) QueryTaskUpResultFile(taskId string) (*types.TaskUpResultFile, error) {
	return s.carrier.carrierDB.QueryTaskUpResultFile(taskId)
}

func (s *CarrierAPIBackend) RemoveTaskUpResultFile(taskId string) error {
	return s.carrier.carrierDB.RemoveTaskUpResultFile(taskId)
}

func (s *CarrierAPIBackend) StoreTaskResultFileSummary(taskId, originId, dataHash, dataPath, dataNodeId, extra string) error {
	// generate metadataId
	var buf bytes.Buffer
	buf.Write([]byte(originId))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	originIdHash := rlputil.RlpHash(buf.Bytes())

	metadataId := types.PREFIX_METADATA_ID + originIdHash.Hex()

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		log.WithError(err).Errorf("Failed query local identity on CarrierAPIBackend.StoreTaskResultFileSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, dataPath: {%s}",
			taskId, dataNodeId, originId, metadataId, dataPath)
		return err
	}

	// store local metadata (about task result file)
	s.carrier.carrierDB.StoreInternalMetadata(types.NewMetadata(&libtypes.MetadataPB{
		MetadataId:   metadataId,
		Owner:        identity,
		DataId:       metadataId,
		DataStatus:   libtypes.DataStatus_DataStatus_Valid,
		MetadataName: fmt.Sprintf("task `%s` result file", taskId),
		MetadataType: 2,  // It means this is a module.
		DataHash:     "", // todo fill it.
		Desc:         fmt.Sprintf("the task `%s` result file after executed", taskId),
		DataType:     libtypes.OrigindataType_OrigindataType_Unknown,
		Industry:     "Unknown",
		// metaData status, eg: create/release/revoke
		State:          libtypes.MetadataState_MetadataState_Created,
		PublishAt:      0, // have not publish
		UpdateAt:       timeutils.UnixMsecUint64(),
		Nonce:          0,
		MetadataOption: "",
		TokenAddress:   "",
	}))

	// todo whether need to store a dataResourceDiskUsed (metadataId. dataNodeId, diskUsed) ??? 后面需要上传 磁盘使用空间在弄吧

	// store dataResourceFileUpload (about task result file)
	err = s.carrier.carrierDB.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(dataNodeId, originId, metadataId, dataPath, dataHash))
	if nil != err {
		log.WithError(err).Errorf("Failed store dataResourceFileUpload about task result file on CarrierAPIBackend.StoreTaskResultFileSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, dataPath: {%s}",
			taskId, dataNodeId, originId, metadataId, dataPath)
		return err
	}
	// 记录原始数据占用资源大小   StoreDataResourceTable  todo 后续考虑是否加上, 目前不加 因为对于系统生成的元数据暂时不需要记录 disk 使用实况 ??
	// 单独记录 metaData 的 GetSize 和所在 dataNodeId   StoreDataResourceDiskUsed  todo 后续考虑是否加上, 目前不加 因为对于系统生成的元数据暂时不需要记录 disk 使用实况 ??

	// store taskId -> TaskUpResultFile (about task result file)
	err = s.carrier.carrierDB.StoreTaskUpResultFile(types.NewTaskUpResultFile(taskId, originId, metadataId, extra))
	if nil != err {
		log.WithError(err).Errorf("Failed store taskUpResultFile on CarrierAPIBackend.StoreTaskResultFileSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, dataPath: {%s}",
			taskId, dataNodeId, originId, metadataId, dataPath)
		return err
	}
	return nil
}

func (s *CarrierAPIBackend) QueryTaskResultFileSummary(taskId string) (*types.TaskResultFileSummary, error) {
	summarry, err := s.carrier.carrierDB.QueryTaskUpResultFile(taskId)
	if nil != err {
		log.WithError(err).Errorf("Failed query taskUpResultFile on CarrierAPIBackend.QueryTaskResultFileSummary(), taskId: {%s}", taskId)
		return nil, err
	}
	dataResourceFileUpload, err := s.carrier.carrierDB.QueryDataResourceFileUpload(summarry.GetOriginId())
	if nil != err {
		log.WithError(err).Errorf("Failed query dataResourceFileUpload on CarrierAPIBackend.QueryTaskResultFileSummary(), taskId: {%s}, originId: {%s}",
			taskId, summarry.GetOriginId())
		return nil, err
	}

	localMetadata, err := s.carrier.carrierDB.QueryInternalMetadataById(summarry.GetMetadataId())
	if nil != err {
		log.WithError(err).Errorf("Failed query local metadata on CarrierAPIBackend.QueryTaskResultFileSummary(), taskId: {%s}, originId: {%s}, metadataId: {%s}",
			taskId, summarry.GetOriginId(), dataResourceFileUpload.GetMetadataId())
		return nil, err
	}

	return types.NewTaskResultFileSummary(
		summarry.GetTaskId(),
		dataResourceFileUpload.GetMetadataId(),
		dataResourceFileUpload.GetOriginId(),
		localMetadata.GetData().GetMetadataName(),
		dataResourceFileUpload.GetDataPath(),
		dataResourceFileUpload.GetNodeId(),
		summarry.GetExtra(),
	), nil

}

func (s *CarrierAPIBackend) QueryTaskResultFileSummaryList() (types.TaskResultFileSummaryArr, error) {
	taskResultFileSummaryArr, err := s.carrier.carrierDB.QueryTaskUpResultFileList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.WithError(err).Errorf("Failed query all taskUpResultFile on CarrierAPIBackend.QueryTaskResultFileSummaryList()")
		return nil, err
	}

	arr := make(types.TaskResultFileSummaryArr, 0)
	for _, summarry := range taskResultFileSummaryArr {
		dataResourceFileUpload, err := s.carrier.carrierDB.QueryDataResourceFileUpload(summarry.GetOriginId())
		if nil != err {
			log.WithError(err).Errorf("Failed query dataResourceFileUpload on CarrierAPIBackend.QueryTaskResultFileSummaryList(), taskId: {%s}, originId: {%s}",
				summarry.GetTaskId(), summarry.GetOriginId())
			continue
		}

		localMetadata, err := s.carrier.carrierDB.QueryInternalMetadataById(dataResourceFileUpload.GetMetadataId())
		if nil != err {
			log.WithError(err).Errorf("Failed query local metadata on CarrierAPIBackend.QueryTaskResultFileSummaryList(), taskId: {%s}, originId: {%s}, metadataId: {%s}",
				summarry.GetTaskId(), summarry.GetOriginId(), dataResourceFileUpload.GetMetadataId())
			continue
		}

		arr = append(arr, types.NewTaskResultFileSummary(
			summarry.GetTaskId(),
			dataResourceFileUpload.GetMetadataId(),
			dataResourceFileUpload.GetOriginId(),
			localMetadata.GetData().GetMetadataName(),
			dataResourceFileUpload.GetDataPath(),
			dataResourceFileUpload.GetNodeId(),
			summarry.GetExtra(),
		))
	}

	return arr, nil
}

func (s *CarrierAPIBackend) EstimateTaskGas(dataTokenTransferList []*pb.DataTokenTransferItem) (gasLimit uint64, gasPrice *big.Int, err error) {
	gasLimit, gasPrice, err = s.carrier.metisPayManager.EstimateTaskGas(dataTokenTransferList);
	if err != nil {
		log.WithError(err).Error("Failed to call EstimateTaskGas() on CarrierAPIBackend.EstimateTaskGas()")
	}
	return

	/*if gasLimit, gasPrice, err = s.carrier.metisPayManager.EstimateTaskGas(dataTokenTransferList); nil != err {
		log.WithError(err).Error("Failed to call EstimateTaskGas() on CarrierAPIBackend.EstimateTaskGas()")
		return 0, nil, err
	} else {
		return gasLimit, gasPrice, nil
	}*/
}

// EstimateTaskGas
