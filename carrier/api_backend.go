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
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	"github.com/RosettaFlow/Carrier-Go/lib/fighter/computesvc"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/libp2p/go-libp2p-core/peer"
	"strings"
)

// CarrierAPIBackend implements rpc.Backend for Carrier
type CarrierAPIBackend struct {
	carrier *Service
}

func NewCarrierAPIBackend(carrier *Service) *CarrierAPIBackend {
	return &CarrierAPIBackend{carrier: carrier}
}

func (s *CarrierAPIBackend) SendMsg(msg types.Msg) error {
	return s.carrier.mempool.Add(msg)
}

// system (the yarn node self info)
func (s *CarrierAPIBackend) GetNodeInfo() (*pb.YarnNodeInfo, error) {

	seedNodes, err := s.GetSeedNodeList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to query all `seed nodes`, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}

	jobNodes, err := s.carrier.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to query all `job nodes`, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}
	dataNodes, err := s.carrier.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to query all `data nodes, on GetNodeInfo(), err: {%s}", err)
		return nil, err
	}
	jobsLen := len(jobNodes)
	datasLen := len(dataNodes)
	length := jobsLen + datasLen
	registerNodes := make([]*pb.YarnRegisteredPeer, length)
	if len(jobNodes) != 0 {
		for i, v := range jobNodes {
			client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(v.GetId())
			if ok && client.IsConnected(){
				v.ConnState = pb.ConnState_ConnState_Connected
			}
			n := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_JobNode,
				NodeDetail: v,
			}
			registerNodes[i] = n
		}
	}
	if len(dataNodes) != 0 {
		for i, v := range dataNodes {
			client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(v.GetId())
			if ok && client.IsConnected(){
				v.ConnState = pb.ConnState_ConnState_Connected
			}
			n := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_DataNode,
				NodeDetail: v,
			}
			registerNodes[jobsLen+i] = n
		}
	}

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
	enc, _ := rlp.EncodeToBytes(&enr) // always succeeds because record is valid
	b64 := base64.RawURLEncoding.EncodeToString(enc)
	enrStr := "enr:" + b64

	nodeInfo := &pb.YarnNodeInfo{
		NodeType:     pb.NodeType_NodeType_YarnNode,
		NodeId:       nodeId,
		IdentityType: types.IDENTITY_TYPE_DID, // default: DID
		IdentityId:   identityId,
		Name:         nodeName,
		Peers:        registerNodes,
		SeedPeers:    seedNodes,
		State:        pb.YarnNodeState_State_Active,
		RelatePeers:  uint32(len(s.carrier.config.P2P.Peers().Active())),
		LocalBootstrapNode: enrStr,
	}

	multiAddr := s.carrier.config.P2P.Host().Addrs()
	if len(multiAddr) != 0 {
		// /ip4/192.168.35.1/tcp/16788
		multiAddrParts := strings.Split(multiAddr[0].String(), "/")
		nodeInfo.ExternalIp = multiAddrParts[2]
		nodeInfo.ExternalPort = multiAddrParts[4]
	}

	selfMultiAddrs, _ := s.carrier.config.P2P.DiscoveryAddresses();
	if len(selfMultiAddrs) != 0 {
		nodeInfo.LocalMultiAddr = selfMultiAddrs[0].String()
	}

	return nodeInfo, nil
}

func (s *CarrierAPIBackend) GetRegisteredPeers(nodeType pb.NodeType) ([]*pb.YarnRegisteredPeer, error) {

	result := make([]*pb.YarnRegisteredPeer, 0)
	switch nodeType {
	case pb.NodeType_NodeType_Unknown:
		// all dataNodes on yarnNode
		dataNodes, err := s.carrier.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
		if nil != err {
			return nil, err
		}
		// handle dataNodes
		for _, v := range dataNodes {
			var duration uint64
			node, has := s.carrier.resourceClientSet.QueryDataNodeClient(v.GetId())
			if has {
				duration = uint64(node.RunningDuration())
			}
			v.Duration = duration // ms
			v.FileCount = 0
			v.FileTotalSize = 0
			registeredPeer := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_DataNode,
				NodeDetail: v,
			}
			result = append(result, registeredPeer)
		}

		// all jobNodes on yarnNode
		jobNodes, err := s.carrier.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
		if nil != err {
			return nil, err
		}

		// handle jobNodes
		for _, v := range jobNodes {
			var duration uint64
			node, has := s.carrier.resourceClientSet.QueryJobNodeClient(v.GetId())
			if has {
				duration = uint64(node.RunningDuration())
			}
			v.TaskCount, _ = s.carrier.carrierDB.QueryRunningTaskCountOnJobNode(v.GetId())
			v.TaskIdList, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(v.GetId())
			v.Duration = duration
			registeredPeer := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_JobNode,
				NodeDetail: v,
			}
			result = append(result, registeredPeer)
		}
	case pb.NodeType_NodeType_DataNode:
		// all dataNodes on yarnNode
		dataNodes, err := s.carrier.carrierDB.QueryRegisterNodeList(pb.PrefixTypeDataNode)
		if nil != err {
			return nil, err
		}
		// handle dataNodes
		for _, v := range dataNodes {
			var duration uint64
			node, has := s.carrier.resourceClientSet.QueryDataNodeClient(v.GetId())
			if has {
				duration = uint64(node.RunningDuration())
			}
			v.Duration = duration // ms
			v.FileCount = 0  // todo need implament this logic
			v.FileTotalSize = 0  // todo need implament this logic
			registeredPeer := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_DataNode,
				NodeDetail: v,
			}
			result = append(result, registeredPeer)
		}
	case pb.NodeType_NodeType_JobNode:
		// all jobNodes on yarnNode
		jobNodes, err := s.carrier.carrierDB.QueryRegisterNodeList(pb.PrefixTypeJobNode)
		if nil != err {
			return nil, err
		}
		// handle jobNodes
		for _, v := range jobNodes {
			var duration uint64
			node, has := s.carrier.resourceClientSet.QueryJobNodeClient(v.GetId())
			if has {
				duration = uint64(node.RunningDuration())
			}
			v.TaskCount, _ = s.carrier.carrierDB.QueryRunningTaskCountOnJobNode(v.GetId())
			v.TaskIdList, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(v.GetId())
			v.Duration = duration
			registeredPeer := &pb.YarnRegisteredPeer{
				NodeType:   pb.NodeType_NodeType_JobNode,
				NodeDetail: v,
			}
			result = append(result, registeredPeer)
		}
	default:
		return nil, fmt.Errorf("Invalid nodeType")
	}
	return result, nil
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
	addrs, err := s.carrier.config.P2P.PeerFromAddress([]string{ seed.GetAddr() })
	if err != nil || len(addrs) == 0{
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
			Addr: addr,
			IsDefault: true,
			ConnState: pb.ConnState_ConnState_UnConnected,
		}
	}

	// query seed node arr from db
	seeds , err := s.carrier.carrierDB.QuerySeedNodeList()
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

func (s *CarrierAPIBackend) storeLocalResource(identity *apicommonpb.Organization, jobNodeId string, jobNodeStatus *computesvc.GetStatusReply) error {

	// store into local db
	if err := s.carrier.carrierDB.InsertLocalResource(types.NewLocalResource(&libtypes.LocalResourcePB{
		IdentityId: identity.GetIdentityId(),
		NodeId:     identity.GetNodeId(),
		NodeName:   identity.GetNodeName(),
		JobNodeId:  jobNodeId,
		DataId:     "", // can not own powerId now, because power have not publish
		// the status of data, N means normal, D means deleted.
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		// resource status, eg: create/release/revoke
		State: apicommonpb.PowerState_PowerState_Created,
		// unit: byte
		TotalMem: jobNodeStatus.GetTotalMemory(),
		UsedMem: 0,
		// number of cpu cores.
		TotalProcessor: jobNodeStatus.GetTotalCpu(),
		UsedProcessor:  0,
		// unit: byte
		TotalBandwidth: jobNodeStatus.GetTotalBandwidth(),
		UsedBandwidth:  0,
	})); nil != err {
		log.Errorf("Failed to store power to local on MessageHandler with broadcast, jobNodeId: {%s}, err: {%s}",
			jobNodeId, err)
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
		return pb.ConnState_ConnState_UnConnected, errors.New("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)

		jobNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect jobNode query status failed, %s", err)
		}
		// add resource usage first, but not own power now (mem, proccessor, bandwidth)
		if err = s.storeLocalResource(identity, node.Id, jobNodeStatus); nil != err {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store jobNode local resource failed, %s", err)
		}
	}
	if typ == pb.PrefixTypeDataNode {
		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)

		dataNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect dataNode query status failed, %s", err)
		}
		// add data resource  (disk)
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk()))
		//err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
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
		return pb.ConnState_ConnState_UnConnected, errors.New("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {

		// The published jobNode cannot be updated directly
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on UpdateRegisterNode, %s", resourceTable.String())
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("still have the published computing power information on old jobNode failed, input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				node.Id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// First check whether there is a task being executed on jobNode
		runningTaskCount, err := s.carrier.carrierDB.QueryRunningTaskCountOnJobNode(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(node.Id); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceClientSet.RemoveJobNodeClient(node.Id)
		}

		// generate new client
		client, err := grpclient.NewJobNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new jobNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreJobNodeClient(node.Id, client)

		jobNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect jobNode query status failed, %s", err)
		}
		// add resource usage first, but not own power now (mem, proccessor, bandwidth)
		if err = s.storeLocalResource(identity, node.Id, jobNodeStatus); nil != err {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store jobNode local resource failed, %s", err)
		}

	}

	if typ == pb.PrefixTypeDataNode {

		// First verify whether the dataNode has been used
		dataNodeTable, err := s.carrier.carrierDB.QueryDataResourceTable(node.Id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("query disk used summary on old dataNode failed, %s", err)
		}
		if dataNodeTable.IsNotEmpty() && dataNodeTable.IsUsed() {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("the disk of old dataNode was used, don't remove it, totalDisk: {%d byte}, usedDisk: {%d byte}, remainDisk: {%d byte}",
				dataNodeTable.GetTotalDisk(), dataNodeTable.GetUsedDisk(), dataNodeTable.RemainDisk())
		}

		if client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(node.Id); ok {
			// remove old client instanse
			client.Close()
			s.carrier.resourceClientSet.RemoveDataNodeClient(node.Id)
		}

		// remove old data resource  (disk)
		if err := s.carrier.carrierDB.RemoveDataResourceTable(node.Id); rawdb.IsNoDBNotFoundErr(err) {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("remove disk summary of old dataNode, %s", err)
		}

		client, err := grpclient.NewDataNodeClientWithConn(s.carrier.ctx, fmt.Sprintf("%s:%s", node.InternalIp, node.InternalPort), node.Id)
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect new dataNode failed, %s", err)
		}
		s.carrier.resourceClientSet.StoreDataNodeClient(node.Id, client)

		dataNodeStatus, err := client.GetStatus()
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("connect dataNode query status failed, %s", err)
		}
		// add new data resource  (disk)
		err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, dataNodeStatus.GetTotalDisk(), dataNodeStatus.GetUsedDisk()))
		//err = s.carrier.carrierDB.StoreDataResourceTable(types.NewDataResourceTable(node.Id, types.DefaultDisk, 0))
		if err != nil {
			return pb.ConnState_ConnState_UnConnected, fmt.Errorf("store disk summary of new dataNode failed, %s", err)
		}
	}

	// remove  old jobNode from db
	if err := s.carrier.carrierDB.DeleteRegisterNode(typ, node.Id); nil != err {
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

	_, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return fmt.Errorf("query local identity failed, %s", err)
	}

	switch typ {
	case pb.PrefixTypeDataNode, pb.PrefixTypeJobNode:
	default:
		return errors.New("invalid nodeType")
	}
	if typ == pb.PrefixTypeJobNode {

		// The published jobNode cannot be updated directly
		resourceTable, err := s.carrier.carrierDB.QueryLocalResourceTable(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query local power resource on old jobNode failed, %s", err)
		}

		if nil != resourceTable {
			log.Debugf("still have the published computing power information on old jobNode on RemoveRegisterNode, %s", resourceTable.String())
			return fmt.Errorf("still have the published computing power information on old jobNode failed,input jobNodeId: {%s}, old jobNodeId: {%s}, old powerId: {%s}",
				id, resourceTable.GetNodeId(), resourceTable.GetPowerId())
		}

		// First check whether there is a task being executed on jobNode
		runningTaskCount, err := s.carrier.carrierDB.QueryRunningTaskCountOnJobNode(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query local running taskCount on old jobNode failed, %s", err)
		}
		if runningTaskCount > 0 {
			return fmt.Errorf("the old jobNode have been running {%d} task current, don't remove it", runningTaskCount)
		}

		if client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(id); ok {
			client.Close()
			s.carrier.resourceClientSet.RemoveJobNodeClient(id)
		}

		// remove jobNode local resource
		if err = s.carrier.carrierDB.RemoveLocalResource(id); nil != err {
			return fmt.Errorf("remove jobNode local resource failed, %s", err)
		}
	}

	if typ == pb.PrefixTypeDataNode {

		// First verify whether the dataNode has been used
		dataNodeTable, err := s.carrier.carrierDB.QueryDataResourceTable(id)
		if rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("query disk used summary on old dataNode failed, %s", err)
		}
		if dataNodeTable.IsNotEmpty() && dataNodeTable.IsUsed() {
			return fmt.Errorf("the disk of old dataNode was used, don't remove it, totalDisk: {%d byte}, usedDisk: {%d byte}, remainDisk: {%d byte}",
				dataNodeTable.GetTotalDisk(), dataNodeTable.GetUsedDisk(), dataNodeTable.RemainDisk())
		}

		if client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(id); ok {
			client.Close()
			s.carrier.resourceClientSet.RemoveDataNodeClient(id)
		}
		// remove data resource  (disk)
		if err := s.carrier.carrierDB.RemoveDataResourceTable(id); rawdb.IsNoDBNotFoundErr(err) {
			return fmt.Errorf("remove disk summary of old registerNode, %s", err)
		}
	}
	return s.carrier.carrierDB.DeleteRegisterNode(typ, id)
}

func (s *CarrierAPIBackend) GetRegisterNode(typ pb.RegisteredNodeType, id string) (*pb.YarnRegisteredPeerDetail, error) {
	return s.carrier.carrierDB.QueryRegisterNode(typ, id)
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

		var connState            pb.ConnState
		var duration             uint64
		var taskCount            uint32
		var taskIdList           []string
		var fileCount            uint32
		var fileTotalSize        uint32

		if typ == pb.PrefixTypeJobNode {
			node, has := s.carrier.resourceClientSet.QueryJobNodeClient(n.GetId())
			if has {
				duration = uint64(node.RunningDuration())
			}

			client, ok := s.carrier.resourceClientSet.QueryJobNodeClient(n.GetId())
			if !ok {
				connState = pb.ConnState_ConnState_UnConnected
			}
			if !client.IsConnected() {
				connState = pb.ConnState_ConnState_UnConnected
			} else {
				connState = pb.ConnState_ConnState_Connected
			}
			taskCount, _ = s.carrier.carrierDB.QueryRunningTaskCountOnJobNode(n.GetId())
			taskIdList, _ = s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(n.GetId())
			fileCount = 0  // todo need implament this logic
			fileTotalSize = 0  // todo need implament this logic
		} else {
			node, has := s.carrier.resourceClientSet.QueryDataNodeClient(n.GetId())
			if has {
				duration = uint64(node.RunningDuration())
			}
			client, ok := s.carrier.resourceClientSet.QueryDataNodeClient(n.GetId())
			if !ok {
				connState = pb.ConnState_ConnState_UnConnected
			}
			if !client.IsConnected() {
				connState = pb.ConnState_ConnState_UnConnected
			} else {
				connState = pb.ConnState_ConnState_Connected
			}
			taskCount = 0
			taskIdList = nil
			fileCount = 0  // todo need implament this logic
			fileTotalSize = 0  // todo need implament this logic
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

	oldUsage, err := s.carrier.carrierDB.QueryTaskResuorceUsage(usage.GetTaskId(), usage.GetPartyId())
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed to QueryTaskResuorceUsage on CarrierAPIBackend.ReportTaskResourceUsage(), taskId: {%s}, partyId: {%s}, nodeType: {%s}, ip:{%s}, port: {%s}, err: {%s}",
			usage.GetTaskId(), usage.GetPartyId(), nodeType.String(), ip, port, err)
		return err
	}

	var needSeedUsage bool
	// add
	if oldUsage == nil {
		oldUsage = usage
		err = s.carrier.carrierDB.StoreTaskResuorceUsage(oldUsage)
		needSeedUsage = true
	} else {
		// update
		if usage.GetUsedMem() > oldUsage.GetUsedMem() {
			oldUsage.SetUsedMem(usage.GetUsedMem())
			needSeedUsage = true
		}
		if usage.GetUsedProcessor() > oldUsage.GetUsedProcessor() {
			oldUsage.SetUsedProcessor(usage.GetUsedProcessor())
			needSeedUsage = true
		}
		if usage.GetUsedBandwidth() > oldUsage.GetUsedBandwidth() {
			oldUsage.SetUsedBandwidth(usage.GetUsedBandwidth())
			needSeedUsage = true
		}
		if usage.GetUsedDisk() > oldUsage.GetUsedDisk() {
			oldUsage.SetUsedDisk(usage.GetUsedDisk())
			needSeedUsage = true
		}
		if needSeedUsage {
			err = s.carrier.carrierDB.StoreTaskResuorceUsage(oldUsage)
		}
	}

	if nil != err {
		log.Errorf("Failed to StoreTaskResuorceUsage on CarrierAPIBackend.ReportTaskResourceUsage(), taskId: {%s}, partyId: {%s}, nodeType: {%s}, ip:{%s}, port: {%s}, err: {%s}",
			usage.GetTaskId(), usage.GetPartyId(), nodeType.String(), ip, port, err)
		return err
	}

	if needSeedUsage {
		if err := s.carrier.TaskManager.SendTaskResourceUsage(oldUsage); nil != err {
			log.Errorf("Failed to SendTaskResourceUsage on CarrierAPIBackend.ReportTaskResourceUsage(), taskId: {%s}, partyId: {%s}, nodeType: {%s}, ip:{%s}, port: {%s}, err: {%s}",
				usage.GetTaskId(), usage.GetPartyId(), nodeType.String(), ip, port, err)
			return err
		}
	}

	return nil
}

// metadata api
func (s *CarrierAPIBackend) IsInternalMetadata(metadataId string) (bool, error) {
	return s.carrier.carrierDB.IsInternalMetadataByDataId(metadataId)
}


func (s *CarrierAPIBackend) GetMetadataDetail(identityId, metadataId string) (*types.Metadata, error) {

	var (
		metadata *types.Metadata
		err      error
	)

	// find local metadata
	if "" == identityId {
		metadata, err = s.carrier.carrierDB.QueryInternalMetadataByDataId(metadataId)
		if rawdb.IsNoDBNotFoundErr(err) {
			return nil, errors.New("not found local metadata by special Id, " + err.Error())
		}
		if nil != metadata {
			return metadata, nil
		}
	}
	metadata, err = s.carrier.carrierDB.QueryMetadataByDataId(metadataId)
	if nil != err {
		return nil, errors.New("not found local metadata by special Id, " + err.Error())
	}

	if nil == metadata {
		return nil, errors.New("not found local metadata by special Id, metadata is empty")
	}
	return metadata, nil
}

// GetMetadataDetailList returns a list of all metadata details in the network.
func (s *CarrierAPIBackend) GetGlobalMetadataDetailList() ([]*pb.GetGlobalMetadataDetailResponse, error) {
	log.Debug("Invoke: GetGlobalMetadataDetailList executing...")
	var (
		arr []*pb.GetGlobalMetadataDetailResponse
		err error
	)

	publishMetadataArr, err := s.carrier.carrierDB.QueryMetadataList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, errors.New("found global metadata arr failed, " + err.Error())
	}
	if len(publishMetadataArr) != 0 {
		arr = append(arr, types.NewGlobalMetadataInfoArrayFromMetadataArray(publishMetadataArr)...)
	}
	//// set metadata used taskCount
	//for i, metadata := range arr {
	//	count, err := s.carrier.carrierDB.QueryMetadataUsedTaskIdCount(metadata.GetInformation().GetMetadataSummary().GetMetadataId())
	//	if nil != err {
	//		log.Warnf("Warn, query metadata used taskIdCount failed on CarrierAPIBackend.GetGlobalMetadataDetailList(), err: {%s}", err)
	//		continue
	//	}
	//	metadata.Information.TotalTaskCount = count
	//	arr[i] = metadata
	//}
	return arr, err
}

func (s *CarrierAPIBackend) GetLocalMetadataDetailList() ([]*pb.GetLocalMetadataDetailResponse, error) {
	log.Debug("Invoke: GetLocalMetadataDetailList executing...")

	var (
		arr []*pb.GetLocalMetadataDetailResponse
		err error
	)

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		if nil != err {
			return nil, errors.New("found local identity failed, " + err.Error())
		}
	}

	internalMetadataArr, err := s.carrier.carrierDB.QueryInternalMetadataList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, errors.New("found local metadata arr failed, " + err.Error())
	}

	globalMetadataArr, err := s.carrier.carrierDB.QueryMetadataList()
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, errors.New("found global metadata arr failed on query local metadata arr, " + err.Error())
	}

	publishMetadataArr := make(types.MetadataArray, 0)
	for _, metadata := range globalMetadataArr {
		if identity.GetIdentityId() == metadata.GetData().GetIdentityId() {
			publishMetadataArr = append(publishMetadataArr, metadata)
		}
	}

	arr = append(arr, types.NewLocalMetadataInfoArrayFromMetadataArray(internalMetadataArr, publishMetadataArr)...)

	// set metadata used taskCount
	for i, metadata := range arr {
		count, err := s.carrier.carrierDB.QueryMetadataUsedTaskIdCount(metadata.GetInformation().GetMetadataSummary().GetMetadataId())
		if nil != err {
			log.Warnf("Warn, query metadata used taskIdCount failed on CarrierAPIBackend.GetLocalMetadataDetailList(), err: {%s}", err)
			continue
		}
		metadata.Information.TotalTaskCount = count
		arr[i] = metadata
	}

	return arr, nil
}

func (s *CarrierAPIBackend) GetMetadataUsedTaskIdList(identityId, metadataId string) ([]string, error) {
	taskIds, err := s.carrier.carrierDB.QueryMetadataUsedTaskIds(metadataId)
	if nil != err {
		return nil, err
	}
	return taskIds, nil
}

// power api


func (s *CarrierAPIBackend) GetGlobalPowerSummaryList() ([]*pb.GetGlobalPowerSummaryResponse, error) {
	log.Debug("Invoke: GetGlobalPowerSummaryList executing...")
	resourceList, err := s.carrier.carrierDB.QueryGlobalResourceSummaryList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power summary list, len: {%d}", len(resourceList))
	powerList := make([]*pb.GetGlobalPowerSummaryResponse, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &pb.GetGlobalPowerSummaryResponse{
			Owner: &apicommonpb.Organization{
				NodeName:   resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentityId(),
			},
			Power: &libtypes.PowerUsageDetail{
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
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetGlobalPowerDetailList() ([]*pb.GetGlobalPowerDetailResponse, error) {
	log.Debug("Invoke: GetGlobalPowerDetailList executing...")
	resourceList, err := s.carrier.carrierDB.QueryGlobalResourceDetailList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Query all org's power detail list, len: {%d}", len(resourceList))
	powerList := make([]*pb.GetGlobalPowerDetailResponse, 0, resourceList.Len())
	for _, resource := range resourceList.To() {
		powerList = append(powerList, &pb.GetGlobalPowerDetailResponse{
			Owner: &apicommonpb.Organization{
				NodeName:   resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentityId(),
			},
			PowerId: resource.GetDataId(),
			Power: &libtypes.PowerUsageDetail{
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
			},
		})
	}
	return powerList, nil
}

func (s *CarrierAPIBackend) GetLocalPowerDetailList() ([]*pb.GetLocalPowerDetailResponse, error) {
	log.Debug("Invoke:GetLocalPowerDetailList executing...")
	// query local resource list from db.
	machineList, err := s.carrier.carrierDB.QueryLocalResourceList()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetLocalPowerDetailList, call QueryLocalResourceList, machineList: %s", machineList.String())

	// query slotUnit
	slotUnit, err := s.carrier.carrierDB.QueryNodeResourceSlotUnit()
	if err != nil {
		return nil, err
	}
	log.Debugf("Invoke:GetLocalPowerDetailList, call QueryNodeResourceSlotUnit, slotUint: %s",
		slotUnit.String())

	buildPowerTaskList := func(jobNodeId string) []*libtypes.PowerTask {

		powerTaskList := make([]*libtypes.PowerTask, 0)

		taskIds, err := s.carrier.carrierDB.QueryJobNodeRunningTaskIdList(jobNodeId)
		if err != nil {
			log.WithError(err).Errorf("Failed to query jobNode runningTaskIds on GetLocalPowerDetailList, jobNodeId: {%s}", jobNodeId)
			return powerTaskList
		}

		for _, taskId := range taskIds {
			// query local task information by taskId
			task, err := s.carrier.carrierDB.QueryLocalTask(taskId)
			if err != nil {
				log.WithError(err).Errorf("Failed to query local task on GetLocalPowerDetailList, taskId: {%s}", taskId)
				continue
			}

			partyIds, err := s.carrier.carrierDB.QueryJobNodeRunningTaskAllPartyIdList(jobNodeId, taskId)
			if err != nil {
				log.WithError(err).Errorf("Failed to query jobNode runningTask all partyIds on GetLocalPowerDetailList, jobNodeId: {%s}, taskId: {%s}", jobNodeId, taskId)
				continue
			}
			// build powerTask info
			powerTask := &libtypes.PowerTask{
				TaskId:   taskId,
				TaskName: task.GetTaskData().TaskName,
				Owner: &apicommonpb.Organization{
					NodeName:   task.GetTaskData().GetNodeName(),
					NodeId:     task.GetTaskData().GetNodeId(),
					IdentityId: task.GetTaskData().GetIdentityId(),
				},
				Receivers: make([]*apicommonpb.Organization, 0),
				OperationCost: &apicommonpb.TaskResourceCostDeclare{
					Processor: task.GetTaskData().GetOperationCost().GetProcessor(),
					Memory:    task.GetTaskData().GetOperationCost().GetMemory(),
					Bandwidth: task.GetTaskData().GetOperationCost().GetBandwidth(),
					Duration:  task.GetTaskData().GetOperationCost().GetDuration(),
				},
				OperationSpend: nil, // will be culculating after ...
				CreateAt:       task.GetTaskData().CreateAt,
			}
			// build dataSuppliers of task info
			for _, dataSupplier := range task.GetTaskData().GetDataSuppliers() {
				powerTask.Partners = append(powerTask.Partners, &apicommonpb.Organization{
					NodeName:   dataSupplier.GetOrganization().GetNodeName(),
					NodeId:     dataSupplier.GetOrganization().GetNodeId(),
					IdentityId: dataSupplier.GetOrganization().GetIdentityId(),
				})
			}
			// build receivers of task info
			for _, receiver := range task.GetTaskData().GetReceivers() {
				powerTask.Receivers = append(powerTask.Receivers, &apicommonpb.Organization{
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

			for _, powerSupplier := range task.GetTaskData().GetPowerSuppliers() {
				if _, ok := partyIdTmp[powerSupplier.GetOrganization().GetPartyId()]; ok  {
					processor += powerSupplier.GetResourceUsedOverview().GetUsedProcessor()
					memory += powerSupplier.GetResourceUsedOverview().GetUsedMem()
					bandwidth += powerSupplier.GetResourceUsedOverview().GetUsedBandwidth()
				}
			}

			if task.GetTaskData().GetStartAt() != 0 {
				duration = timeutils.UnixMsecUint64() - task.GetTaskData().GetStartAt()
			}
			powerTask.OperationSpend = &apicommonpb.TaskResourceCostDeclare{
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
			log.Errorf("Failed to query task totalCount with jobNodeId on GetLocalPowerDetailList, jobNodeId: {%s}, err: {%s}", jobNodeId, err)
			return 0
		}
		return count
	}
	// culculate task current running count on jobNode
	taskRunningCount := func(jobNodeId string) uint32 {
		count, err := s.carrier.carrierDB.QueryRunningTaskCountOnJobNode(jobNodeId)
		if err != nil {
			log.Errorf("Failed to query task runningCount with jobNodeId on GetLocalPowerDetailList, jobNodeId: {%s}, err: {%s}", jobNodeId, err)
			return 0
		}
		return count
	}

	resourceList := machineList.To()
	// handle  resource one by one
	result := make([]*pb.GetLocalPowerDetailResponse, len(resourceList))
	for i, resource := range resourceList {

		nodePowerDetail := &pb.GetLocalPowerDetailResponse{
			JobNodeId: resource.GetJobNodeId(),
			PowerId:   resource.GetDataId(),
			Owner: &apicommonpb.Organization{
				NodeName:   resource.GetNodeName(),
				NodeId:     resource.GetNodeId(),
				IdentityId: resource.GetIdentityId(),
			},
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
			},
		}
		nodePowerDetail.Power.Tasks = buildPowerTaskList(resource.GetJobNodeId())
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
	}), err
}

func (s *CarrierAPIBackend) GetIdentityList() ([]*types.Identity, error) {
	return s.carrier.carrierDB.QueryIdentityList()
}

// for metadataAuthority

func (s *CarrierAPIBackend) AuditMetadataAuthority(audit *types.MetadataAuthAudit) (apicommonpb.AuditMetadataOption, error) {
	return s.carrier.authEngine.AuditMetadataAuthority(audit)
}

func (s *CarrierAPIBackend) GetLocalMetadataAuthorityList() (types.MetadataAuthArray, error) {
	return s.carrier.authEngine.GetLocalMetadataAuthorityList()
}

func (s *CarrierAPIBackend) GetGlobalMetadataAuthorityList() (types.MetadataAuthArray, error) {
	return s.carrier.authEngine.GetGlobalMetadataAuthorityList()
}

func (s *CarrierAPIBackend) HasValidMetadataAuth(userType apicommonpb.UserType, user, identityId, metadataId string) (bool, error) {
	return s.carrier.authEngine.HasValidMetadataAuth(userType, user, identityId, metadataId)
}

// task api
func (s *CarrierAPIBackend) GetLocalTask(taskId string) (*pb.TaskDetailShow, error) {
	// the task is executing.
	localTask, err := s.carrier.carrierDB.QueryLocalTask(taskId)
	if nil != err {
		log.Errorf("Failed to query local task on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}", taskId)
		return nil, err
	}

	if nil == localTask {
		log.Errorf("Not found local task on `CarrierAPIBackend.GetLocalTask()`, taskId: {%s}", taskId)
		return nil, fmt.Errorf("not found local task")
	}

	detailShow :=  &pb.TaskDetailShow{
		TaskId:   localTask.GetTaskId(),
		TaskName: localTask.GetTaskData().GetTaskName(),
		UserType: localTask.GetTaskData().GetUserType(),
		User:     localTask.GetTaskData().GetUser(),
		Sender: &apicommonpb.TaskOrganization{
			PartyId:   localTask.GetTaskSender().GetPartyId(),
			NodeName:  localTask.GetTaskSender().GetNodeName(),
			NodeId:    localTask.GetTaskSender().GetNodeId(),
			IdentityId:localTask.GetTaskSender().GetIdentityId(),
		},
		AlgoSupplier: &apicommonpb.TaskOrganization{
			PartyId:    localTask.GetTaskData().GetAlgoSupplier().GetPartyId(),
			NodeName:   localTask.GetTaskData().GetAlgoSupplier().GetNodeName(),
			NodeId:     localTask.GetTaskData().GetAlgoSupplier().GetNodeId(),
			IdentityId: localTask.GetTaskData().GetAlgoSupplier().GetIdentityId(),
		},
		DataSuppliers:  make([]*pb.TaskDataSupplierShow, 0, len(localTask.GetTaskData().GetDataSuppliers())),
		PowerSuppliers: make([]*pb.TaskPowerSupplierShow, 0, len(localTask.GetTaskData().GetPowerSuppliers())),
		Receivers:      localTask.GetTaskData().GetReceivers(),
		CreateAt:       localTask.GetTaskData().GetCreateAt(),
		StartAt:        localTask.GetTaskData().GetStartAt(),
		EndAt:          localTask.GetTaskData().GetEndAt(),
		State:          localTask.GetTaskData().GetState(),
		OperationCost: &apicommonpb.TaskResourceCostDeclare{
			Processor: localTask.GetTaskData().GetOperationCost().GetProcessor(),
			Memory:    localTask.GetTaskData().GetOperationCost().GetMemory(),
			Bandwidth: localTask.GetTaskData().GetOperationCost().GetBandwidth(),
			Duration:  localTask.GetTaskData().GetOperationCost().GetDuration(),
		},
	}

	// DataSupplier
	for _, dataSupplier := range localTask.GetTaskData().GetDataSuppliers() {
		detailShow.DataSuppliers = append(detailShow.DataSuppliers, &pb.TaskDataSupplierShow{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    dataSupplier.GetOrganization().GetPartyId(),
				NodeName:   dataSupplier.GetOrganization().GetNodeName(),
				NodeId:     dataSupplier.GetOrganization().GetNodeId(),
				IdentityId: dataSupplier.GetOrganization().GetIdentityId(),
			},
			MetadataId:   dataSupplier.GetMetadataId(),
			MetadataName: dataSupplier.GetMetadataName(),
		})
	}
	// powerSupplier
	for _, data := range localTask.GetTaskData().GetPowerSuppliers() {
		detailShow.PowerSuppliers = append(detailShow.PowerSuppliers, &pb.TaskPowerSupplierShow{
			Organization: &apicommonpb.TaskOrganization{
				PartyId:    data.GetOrganization().GetPartyId(),
				NodeName:   data.GetOrganization().GetNodeName(),
				NodeId:     data.GetOrganization().GetNodeId(),
				IdentityId: data.GetOrganization().GetIdentityId(),
			},
			PowerInfo: &libtypes.ResourceUsageOverview{
				TotalMem:       data.GetResourceUsedOverview().GetTotalMem(),
				UsedMem:        data.GetResourceUsedOverview().GetUsedMem(),
				TotalProcessor: data.GetResourceUsedOverview().GetTotalProcessor(),
				UsedProcessor:  data.GetResourceUsedOverview().GetUsedProcessor(),
				TotalBandwidth: data.GetResourceUsedOverview().GetTotalBandwidth(),
				UsedBandwidth:  data.GetResourceUsedOverview().GetUsedBandwidth(),
			},
		})
	}

	return detailShow, nil
}

func (s *CarrierAPIBackend) GetTaskDetailList() ([]*pb.TaskDetailShow, error) {
	// the task is executing.
	localTaskArray, err := s.carrier.carrierDB.QueryLocalTaskList()

	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}
	localIdentityId, err := s.carrier.carrierDB.QueryIdentityId()
	if err != nil {
		return nil, fmt.Errorf("query local identityId failed, %s", err)
	}

	// the task has been executed.
	networkTaskList, err := s.carrier.carrierDB.QueryTaskListByIdentityId(localIdentityId)
	if rawdb.IsNoDBNotFoundErr(err) {
		return nil, err
	}

	makeTaskViewFn := func(task *types.Task) *pb.TaskDetailShow {
		return types.NewTaskDetailShowFromTaskData(task)
	}

	result := make([]*pb.TaskDetailShow, 0)
	for _, task := range localTaskArray {
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

	// 先查出 task 在本地的 eventList
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
			Owner: &apicommonpb.Organization{
				NodeName:   identity.GetNodeName(),
				NodeId:     identity.GetNodeId(),
				IdentityId: identity.GetIdentityId(),
			},
			PartyId: e.GetPartyId(),
		}
	}
	// 再查数据中心的.
	taskEvent, err := s.carrier.carrierDB.QueryTaskEventListByTaskId(taskId)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, types.NewTaskEventFromAPIEvent(taskEvent)...)
	return evenList, nil
}

func (s *CarrierAPIBackend) GetTaskEventListByTaskIds(taskIds []string) ([]*pb.TaskEventShow, error) {

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		return nil, err
	}

	evenList := make([]*pb.TaskEventShow, 0)

	// 先查出 task 在本地的 eventList
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
				Owner: &apicommonpb.Organization{
					NodeName:   identity.GetNodeName(),
					NodeId:     identity.GetNodeId(),
					IdentityId: identity.GetIdentityId(),
				},
				PartyId: e.GetPartyId(),
			})
		}
	}
	// 再查数据中心的.
	taskEvent, err := s.carrier.carrierDB.QueryTaskEventListByTaskIds(taskIds)
	if nil != err {
		return nil, err
	}
	evenList = append(evenList, types.NewTaskEventFromAPIEvent(taskEvent)...)
	return evenList, nil
}

func (s *CarrierAPIBackend) HasLocalTask () (bool, error) {
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

func (s *CarrierAPIBackend) StoreTaskResultFileSummary(taskId, originId, filePath, dataNodeId string) error {
	// generate metadataId
	var buf bytes.Buffer
	buf.Write([]byte(originId))
	buf.Write(bytesutil.Uint64ToBytes(timeutils.UnixMsecUint64()))
	originIdHash := rlputil.RlpHash(buf.Bytes())

	metadataId := types.PREFIX_METADATA_ID + originIdHash.Hex()

	identity, err := s.carrier.carrierDB.QueryIdentity()
	if nil != err {
		log.Errorf("Failed query local identity on CarrierAPIBackend.StoreTaskResultFileSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, filePath: {%s}, err: {%s}",
			taskId, dataNodeId, originId, metadataId, filePath, err)
		return err
	}

	// store local metadata (about task result file)
	s.carrier.carrierDB.StoreInternalMetadata(types.NewMetadata(&libtypes.MetadataPB{
		MetadataId:      metadataId,
		IdentityId:      identity.GetIdentityId(),
		NodeId:          identity.GetNodeId(),
		NodeName:        identity.GetNodeName(),
		DataId:          metadataId,
		OriginId:        originId,
		TableName:       fmt.Sprintf("task `%s` result file", taskId),
		FilePath:        filePath,
		FileType:        apicommonpb.OriginFileType_FileType_Unknown,
		Desc:            fmt.Sprintf("the task `%s` result file after executed", taskId),
		Rows:            0,
		Columns:         0,
		Size_:           0,
		HasTitle:        false,
		MetadataColumns: nil,
		Industry:        "Unknown",
		// the status of data, N means normal, D means deleted.
		DataStatus: apicommonpb.DataStatus_DataStatus_Normal,
		// metaData status, eg: create/release/revoke
		State: apicommonpb.MetadataState_MetadataState_Created,
	}))

	// todo whether need to store a dataResourceDiskUsed (metadataId. dataNodeId, diskUsed) ??? 后面需要上传 磁盘使用空间在弄吧

	// store dataResourceFileUpload (about task result file)
	err = s.carrier.carrierDB.StoreDataResourceFileUpload(types.NewDataResourceFileUpload(dataNodeId, originId, metadataId, filePath))
	if nil != err {
		log.Errorf("Failed store dataResourceFileUpload about task result file on CarrierAPIBackend.StoreTaskResultFileSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, filePath: {%s}, err: {%s}",
			taskId, dataNodeId, originId, metadataId, filePath, err)
		return err
	}
	// 记录原始数据占用资源大小   StoreDataResourceTable  todo 后续考虑是否加上, 目前不加 因为对于系统生成的元数据暂时不需要记录 disk 使用实况 ??
	// 单独记录 metaData 的 GetSize 和所在 dataNodeId   StoreDataResourceDiskUsed  todo 后续考虑是否加上, 目前不加 因为对于系统生成的元数据暂时不需要记录 disk 使用实况 ??

	// store taskId -> TaskUpResultFile (about task result file)
	err = s.carrier.carrierDB.StoreTaskUpResultFile(types.NewTaskUpResultFile(taskId, originId, metadataId))
	if nil != err {
		log.Errorf("Failed store taskUpResultFile on CarrierAPIBackend.StoreTaskResultFileSummary(), taskId: {%s}, dataNodeId: {%s}, originId: {%s}, metadataId: {%s}, filePath: {%s}, err: {%s}",
			taskId, dataNodeId, originId, metadataId, filePath, err)
		return err
	}
	return nil
}

func (s *CarrierAPIBackend) QueryTaskResultFileSummary(taskId string) (*types.TaskResultFileSummary, error) {
	taskUpResultFile, err := s.carrier.carrierDB.QueryTaskUpResultFile(taskId)
	if nil != err {
		log.Errorf("Failed query taskUpResultFile on CarrierAPIBackend.QueryTaskResultFileSummary(), taskId: {%s}, err: {%s}",
			taskId, err)
		return nil, err
	}
	dataResourceFileUpload, err := s.carrier.carrierDB.QueryDataResourceFileUpload(taskUpResultFile.GetOriginId())
	if nil != err {
		log.Errorf("Failed query dataResourceFileUpload on CarrierAPIBackend.QueryTaskResultFileSummary(), taskId: {%s}, originId: {%s}, err: {%s}",
			taskId, taskUpResultFile.GetOriginId(), err)
		return nil, err
	}

	localMetadata, err := s.carrier.carrierDB.QueryInternalMetadataByDataId(taskUpResultFile.GetMetadataId())
	if nil != err {
		log.Errorf("Failed query local metadata on CarrierAPIBackend.QueryTaskResultFileSummary(), taskId: {%s}, originId: {%s}, metadataId: {%s}, err: {%s}",
			taskId, taskUpResultFile.GetOriginId(), dataResourceFileUpload.GetMetadataId(), err)
		return nil, err
	}

	return types.NewTaskResultFileSummary(
		taskUpResultFile.GetTaskId(),
		localMetadata.GetData().GetTableName(),
		dataResourceFileUpload.GetMetadataId(),
		dataResourceFileUpload.GetOriginId(),
		dataResourceFileUpload.GetFilePath(),
		dataResourceFileUpload.GetNodeId(),
	), nil

}

func (s *CarrierAPIBackend) QueryTaskResultFileSummaryList() (types.TaskResultFileSummaryArr, error) {
	taskResultFileSummaryArr, err := s.carrier.carrierDB.QueryTaskUpResultFileList()
	if rawdb.IsNoDBNotFoundErr(err) {
		log.Errorf("Failed query all taskUpResultFile on CarrierAPIBackend.QueryTaskResultFileSummaryList(), err: {%s}", err)
		return nil, err
	}

	arr := make(types.TaskResultFileSummaryArr, 0)
	for _, summarry := range taskResultFileSummaryArr {
		dataResourceFileUpload, err := s.carrier.carrierDB.QueryDataResourceFileUpload(summarry.GetOriginId())
		if nil != err {
			log.Errorf("Failed query dataResourceFileUpload on CarrierAPIBackend.QueryTaskResultFileSummaryList(), taskId: {%s}, originId: {%s}, err: {%s}",
				summarry.GetTaskId(), summarry.GetOriginId(), err)
			continue
		}

		localMetadata, err := s.carrier.carrierDB.QueryInternalMetadataByDataId(dataResourceFileUpload.GetMetadataId())
		if nil != err {
			log.Errorf("Failed query local metadata on CarrierAPIBackend.QueryTaskResultFileSummaryList(), taskId: {%s}, originId: {%s}, metadataId: {%s}, err: {%s}",
				summarry.GetTaskId(), summarry.GetOriginId(), dataResourceFileUpload.GetMetadataId(), err)
			continue
		}

		arr = append(arr, types.NewTaskResultFileSummary(
			summarry.GetTaskId(),
			localMetadata.GetData().GetTableName(),
			dataResourceFileUpload.GetMetadataId(),
			dataResourceFileUpload.GetOriginId(),
			dataResourceFileUpload.GetFilePath(),
			dataResourceFileUpload.GetNodeId(),
		))
	}

	return arr, nil
}
