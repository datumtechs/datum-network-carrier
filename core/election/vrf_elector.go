package election

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/rlputil"
	"github.com/RosettaFlow/Carrier-Go/common/timeutils"
	"github.com/RosettaFlow/Carrier-Go/core/resource"
	"github.com/RosettaFlow/Carrier-Go/crypto/vrf"
	pb "github.com/RosettaFlow/Carrier-Go/lib/api"
	apicommonpb "github.com/RosettaFlow/Carrier-Go/lib/common"
	libtypes "github.com/RosettaFlow/Carrier-Go/lib/types"
	"github.com/RosettaFlow/Carrier-Go/p2p"
	"github.com/RosettaFlow/Carrier-Go/rpc/backend"
	"github.com/RosettaFlow/Carrier-Go/types"
	"math/big"
	"sort"
	"strings"
)

var (
	ErrEnoughResourceOrgCountLessCalculateCount = fmt.Errorf("the enough resource org count is less calculate count")
	ErrEnoughInternalResourceCount              = fmt.Errorf("has not enough internal resource count")
)

type VrfElector struct {
	privateKey      *ecdsa.PrivateKey // privateKey of current node
	resourceMng     *resource.Manager
}

func NewVrfElector(privateKey *ecdsa.PrivateKey, resourceMng *resource.Manager) *VrfElector {
	return &VrfElector{
		privateKey:      privateKey,
		resourceMng:     resourceMng,
	}
}

func (s *VrfElector) ElectionOrganization(
	taskId string,
	powerPartyIds []string,
	skipIdentityIdCache map[string]struct{},
	mem, bandwidth, disk uint64, processor uint32,
	extra []byte,
) ([]*libtypes.TaskPowerSupplier, []byte, [][]byte, error) {

	calculateCount := len(powerPartyIds)

	// Find global identitys
	identityInfoArr, err := s.resourceMng.GetDB().QueryIdentityList(timeutils.BeforeYearUnixMsecUint64(), backend.DefaultMaxPageSize)
	if nil != err {
		return nil, nil, nil, err
	}

	if len(identityInfoArr) < calculateCount {
		return nil, nil, nil, fmt.Errorf("query identityList count less calculate count")
	}

	log.Debugf("QueryIdentityList by dataCenter on VrfElector.ElectionOrganization(), taskId: {%s}, len: {%d}, identityList: %s", taskId, len(identityInfoArr), identityInfoArr.String())
	identityInfoTmp := make(map[string]*types.Identity, calculateCount)
	for _, identityInfo := range identityInfoArr {

		// Skip the mock identityId
		if s.resourceMng.IsMockIdentityId(identityInfo.GetIdentityId()) {
			continue
		}

		identityInfoTmp[identityInfo.GetIdentityId()] = identityInfo
	}

	if len(identityInfoTmp) < calculateCount {
		return nil, nil, nil, fmt.Errorf("find valid identityIds count less calculate count")
	}

	// Find global power resources
	globalResources, err := s.resourceMng.GetDB().QueryGlobalResourceSummaryList()
	if nil != err {
		return nil, nil, nil, err
	}
	log.Debugf("GetRemoteResouceTables on VrfElector.ElectionOrganization(), taskId: {%s}, len: {%d}, globalResources: %s", taskId, len(globalResources), globalResources.String())

	if len(globalResources) < calculateCount {
		return nil, nil, nil, fmt.Errorf("query org's power resource count less calculate count")
	}

	nonce, err := s.vrfNonce(extra)
	if nil != err {
		return nil, nil, nil, err
	}
	queue, weights := s.vrfElectionOrganizationResourceQueue(globalResources, nonce, calculateCount)

	orgs := make([]*libtypes.TaskPowerSupplier, 0)
	i := 0
	for _, r := range queue {

		if i == calculateCount {
			break
		}

		// skip
		if len(skipIdentityIdCache) != 0 {
			if _, ok := skipIdentityIdCache[r.GetIdentityId()]; ok {
				continue
			}
		}

		// append one, if it enouph
		if info, ok := identityInfoTmp[r.GetIdentityId()]; ok {
			orgs = append(orgs, &libtypes.TaskPowerSupplier{
				Organization: &apicommonpb.TaskOrganization{
					PartyId:    powerPartyIds[i],
					NodeName:   info.GetName(),
					NodeId:     info.GetNodeId(),
					IdentityId: info.GetIdentityId(),
				},
				ResourceUsedOverview: &libtypes.ResourceUsageOverview{
					TotalMem:       r.GetTotalMem(), // total resource value of org.
					UsedMem:        0,               // used resource of this task (real time max used)
					TotalBandwidth: r.GetTotalBandWidth(),
					UsedBandwidth:  0, // used resource of this task (real time max used)
					TotalDisk:      r.GetTotalDisk(),
					UsedDisk:       0,
					TotalProcessor: r.GetTotalProcessor(),
					UsedProcessor:  0, // used resource of this task (real time max used)
				},
			})
			i++
		}
	}

	if len(orgs) < calculateCount {
		return nil, nil, nil, ErrEnoughResourceOrgCountLessCalculateCount
	}
	return orgs, nonce, weights, nil
}

func (s *VrfElector) ElectionNode(taskId string, mem, bandwidth, disk uint64, processor uint32, extra string) (*pb.YarnRegisteredPeerDetail, error) {

	if s.resourceMng.HasNotInternalJobNodeClientSet() {
		return nil, fmt.Errorf("not found alive jobNode")
	}

	resourceNodeIdArr := make([]string, 0)

	tables, err := s.resourceMng.QueryLocalResourceTables()
	if nil != err {
		return nil, fmt.Errorf("query local resource tables failed, %s", err)
	}
	log.Debugf("QueryLocalResourceTables on electionJobNode, taskId: {%s}, localResources: %s", taskId, types.UtilLocalResourceArrString(tables))
	for _, r := range tables {
		isEnough := r.IsEnough(mem, bandwidth, disk, processor)
		if isEnough {
			log.Debugf("Call electionJobNode find resource enough jobNode, taskId: {%s}, resource: %s, r.RemainMem(): %d, r.RemainBandwidth(): %d, r.RemainDisk(): %d, r.RemainProcessor(): %d, needMem: %d, needBandwidth: %d, needDisk: %d, needProcessor: %d, isEnough: %v",
				taskId, r.String(), r.RemainMem(), r.RemainBandwidth(), r.RemainDisk(), r.RemainProcessor(), mem, bandwidth, disk, processor, isEnough)
			jobNodeClient, find := s.resourceMng.QueryJobNodeClient(r.GetNodeId())
			if find && jobNodeClient.IsConnected() {
				resourceNodeIdArr = append(resourceNodeIdArr, r.GetNodeId())
				log.Debugf("Call electionJobNode collect jobNode, append jobNodeId: %s", r.GetNodeId())
			}
		} else {
			taskIds, _ := s.resourceMng.GetDB().QueryJobNodeRunningTaskIdList(r.GetNodeId())
			log.Debugf("Call electionJobNode it is a not enough resource jobNode, taskId: {%s}, resource: %s, r.RemainMem(): %d, r.RemainBandwidth(): %d, r.RemainDisk(): %d, r.RemainProcessor(): %d, needMem: %d, needBandwidth: %d, needDisk: %d, needProcessor: %d, isEnough: %v, was running taskIds: %v",
				taskId, r.String(), r.RemainMem(), r.RemainBandwidth(), r.RemainDisk(), r.RemainProcessor(), mem, bandwidth, disk, processor, isEnough, "[" + strings.Join(taskIds, ",") + "]")
		}
	}

	if len(resourceNodeIdArr) == 0 {
		return nil, ErrEnoughInternalResourceCount
	}

	resourceId := resourceNodeIdArr[len(resourceNodeIdArr)-1]
	jobNode, err := s.resourceMng.GetDB().QueryRegisterNode(pb.PrefixTypeJobNode, resourceId)
	if nil != err {
		return nil, fmt.Errorf("query jobNode info failed, %s", err)
	}
	if nil == jobNode {
		return nil, fmt.Errorf("not found jobNode information")
	}
	return jobNode, nil
}

func (s *VrfElector) EnoughAvailableOrganization(taskId string, calculateCount int, mem, bandwidth, disk uint64, processor uint32) (bool, error) {

	// Find global power resources
	globalResources, err := s.resourceMng.GetDB().QueryGlobalResourceSummaryList()
	if nil != err {
		return false, err
	}
	if len(globalResources) < calculateCount {
		return false, fmt.Errorf("query org's power resource count less calculate count")
	}

	i := 0
	for _, r := range globalResources {

		if i == calculateCount {
			break
		}

		// Find one, if have enough resource
		rMem, rBandwidth, rProcessor := r.GetTotalMem()-r.GetUsedMem(), r.GetTotalBandWidth()-r.GetUsedBandWidth(), r.GetTotalProcessor()-r.GetUsedProcessor()
		if rMem < mem {
			continue
		}
		if rProcessor < processor {
			continue
		}
		if rBandwidth < bandwidth {
			continue
		}
		// ignore disk for power resource.

		i++
	}
	if i < calculateCount {
		return false, nil
	}
	return true, nil
}

func (s *VrfElector) VerifyElectionOrganization(taskId string, powerSuppliers []*libtypes.TaskPowerSupplier, nodeIdStr string, extra, nonce []byte, weights [][]byte) (bool, error) {

	if len(powerSuppliers) != len(weights) {
		return false, fmt.Errorf("powerSuppliers count is invalid, powerSuppliers count : %d, weights count: %d", len(powerSuppliers), len(weights))
	}

	if len(nonce) == 0 {
		return false, fmt.Errorf("empty vrf nonce <proof + rand>")
	}

	nodeId, err := p2p.HexID(nodeIdStr)
	if nil != err {
		return false, fmt.Errorf("convert nodeId from hex failed, %s", err)
	}
	pubKey, err := nodeId.Pubkey()
	if nil != err {
		return false, fmt.Errorf("fetch publicKey from nodeId failed, %s", err)
	}

	input := rlputil.RlpHash(extra) // extra just is a taskId + electionAt

	//log.Debugf("Verify vrt nonce, taskId: %s, nodeId: %s, pubKey: %s%s, input: %s, nonce: %v",
	//taskId, nodeId.String(), pubKey.X.String(), pubKey.Y.String(), input.String(), nonce)

	flag, err := vrf.Verify(pubKey, nonce, input.Bytes())
	if nil != err {
		return false, fmt.Errorf("verify vrf nonce <proof + rand> failed, %s", err)
	}
	if !flag {
		return false, fmt.Errorf("verify vrf nonce <proof + rand> result is %v", flag)
	}

	weightMap := make(map[string]struct{}, len(weights))
	for _, weight := range weights {
		weightMap[new(big.Int).SetBytes(weight).String()] = struct{}{}
	}

	rand := vrf.ProofToHash(nonce) // nonce == proof + rand , len(rand) == 32

	identityIdMap := make(map[string]struct{}, len(powerSuppliers))
	for _, powerSupplier := range powerSuppliers {
		dh := rlputil.RlpHash(powerSupplier.GetOrganization().GetIdentityId()) // len(dh) == 32
		value := new(big.Int).Xor(new(big.Int).SetBytes(dh.Bytes()), new(big.Int).SetBytes(rand)).String()
		if _, ok := weightMap[value]; !ok {
			return false, fmt.Errorf("not found vrf xor weight value of powerSupplier, identity: %s, weight: %s", powerSupplier.GetOrganization().GetIdentityId(), value)
		}
		identityIdMap[powerSupplier.GetOrganization().GetIdentityId()] = struct{}{}
	}

	// Find global power resources
	globalResources, err := s.resourceMng.GetDB().QueryGlobalResourceSummaryList()
	if nil != err {
		return false, fmt.Errorf("query global resource summary list failed, %s", err)
	}
	log.Debugf("GetRemoteResouceTables on VrfElector.VerifyElectionOrganization(), taskId: {%s}, len: {%d}, globalResources: %s", taskId, len(globalResources), globalResources.String())

	if len(globalResources) < len(powerSuppliers) {
		return false, fmt.Errorf("query org's power resource count less calculate count")
	}
	queue, reweights := s.vrfElectionOrganizationResourceQueue(globalResources, nonce, len(powerSuppliers))
	for _, powerSupplier := range queue {
		if _, ok := identityIdMap[powerSupplier.GetIdentityId()]; !ok {
			return false, fmt.Errorf("not found identityId of powerSupplier when reElectionOrganizationResource, identity: %s", powerSupplier.GetIdentityId())
		}
	}
	for _, weight := range reweights {
		value := new(big.Int).SetBytes(weight).String()
		if _, ok := weightMap[value]; !ok {
			return false, fmt.Errorf("not found reweight value of powerSupplier when reElectionOrganizationResource, weight: %s", value)
		}
	}
	return true, nil
}

// data is taskId
func (s *VrfElector) vrfNonce(data []byte) ([]byte, error) {
	if nil == s.privateKey {
		return nil, fmt.Errorf("not found privateKey of current node")
	}
	input := rlputil.RlpHash(data)
	nonce, err := vrf.Prove(s.privateKey, input.Bytes()) // nonce == proof + rand
	if nil != err {
		return nil, fmt.Errorf("Failed to generate vrf proof, %s", err)
	}
	//log.Debugf("Generate vrt nonce, privateKey: %s, input: %s, nonce: %v", s.privateKey.D.String(), input.String(), nonce)
	return nonce, nil
}

type randomIden struct {
	data            *types.Resource
	value           *big.Int
	totalMem       uint64
	totalBandwidth uint64
	totalProcessor uint32
}

func newRandomIden(resource *types.Resource, value *big.Int, totalMem, totalBandwidth uint64, totalProcessor uint32) *randomIden {
	return &randomIden{
		data:  resource,
		value: value,
		totalMem: totalMem,
		totalBandwidth: totalBandwidth,
		totalProcessor: totalProcessor,
	}
}

type randomIdenQueue []*randomIden

func (r randomIdenQueue) Len() int {
	return len(r)
}

func (r randomIdenQueue) Less(i, j int) bool { // from max to min
	a, b := r[i], r[j]
	if a.value.Cmp(b.value) < 0 {
		return false
	} else if  a.value.Cmp(b.value) > 0 {
		return true
	} else {
		return new(big.Int).SetBytes(rlputil.RlpHash(a.data.GetIdentityId()).Bytes()).Cmp(new(big.Int).SetBytes(rlputil.RlpHash(b.data.GetIdentityId()).Bytes())) >= 0
	}
	//if a.value.Cmp(b.value) < 0 {
	//	return false
	//} else {  // >= 0
	//	flag := 1
	//	if a.totalMem >= b.totalMem {
	//		flag &= 1
	//	} else {
	//		flag &= 0
	//	}
	//	if a.totalBandwidth >= b.totalBandwidth {
	//		flag &= 1
	//	} else {
	//		flag &= 0
	//	}
	//	if a.totalProcessor >= b.totalProcessor {
	//		flag &= 1
	//	} else {
	//		flag &= 0
	//	}
	//	if flag != 1 {
	//		return false
	//	} else {
	//		return true
	//	}
	//}
}

func (r randomIdenQueue) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (s *VrfElector) vrfElectionOrganizationResourceQueue(resources types.ResourceArray, nonce []byte, count int) (types.ResourceArray, [][]byte) {

	rand := vrf.ProofToHash(nonce) // nonce == proof + rand , len(rand) == 32

	queue := make(randomIdenQueue, len(resources))
	for i, resource := range resources {
		dh := rlputil.RlpHash(resource.GetIdentityId()) // len(dh) == 32
		value := new(big.Int).Xor(new(big.Int).SetBytes(dh.Bytes()), new(big.Int).SetBytes(rand))
		queue[i] = newRandomIden(resource, value,  resource.GetTotalMem(), resource.GetTotalBandWidth(), resource.GetTotalProcessor())
	}
	sort.Sort(queue)

	res := make(types.ResourceArray, count)
	bs := make([][]byte, count)
	for i, riden := range queue {
		if i == count {
			break
		}
		res[i] = riden.data
		bs[i] = riden.value.Bytes()
	}
	return res, bs
}
