package blacklist_test

import (
	"encoding/json"
	"fmt"
	"github.com/datumtechs/datum-network-carrier/blacklist"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/consensus/twopc"
	pbtypes "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
	"gotest.tools/assert"
	"os"
	"testing"
)

type Test struct {
	SavePath string
	Cache    int
	Handles  int
}

func generateJsonFile() {
	ConsensusStateFile := &Test{
		SavePath: "./tests",
		Cache:    32,
		Handles:  32,
	}
	filePtr, err := os.Create("test.json")
	if err != nil {
		fmt.Println("Create file failed", err.Error())
		return
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)

	err = encoder.Encode(ConsensusStateFile)
	if err != nil {
		fmt.Println("Encoder failed", err.Error())

	} else {
		fmt.Println("Encoder success")
	}
}
func GenerateObg() *blacklist.IdentityBackListCache {
	identityBlackListCache := blacklist.NewIdentityBackListCache()
	_, _ = twopc.New(&twopc.Config{
		PeerMsgQueueSize:   12,
		ConsensusStateFile: "test.json",
	},
		nil,
		nil,
		nil,
		nil,
		identityBlackListCache,
	)
	return identityBlackListCache
}

func TestFindBlackOrgByWalPrefix(t *testing.T) {
	taskOrgDatas := []*types.Task{
		types.NewTask(&pbtypes.TaskPB{
			DataSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p0",
					NodeName:   "",
					NodeId:     "NodeId01",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p1",
					NodeName:   "",
					NodeId:     "NodeId02",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc189",
				},
			},
			PowerSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p2",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p3",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
			},
			Receivers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p4",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
			},
		}),
		types.NewTask(&pbtypes.TaskPB{
			DataSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p0",
					NodeName:   "",
					NodeId:     "NodeId01",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p1",
					NodeName:   "",
					NodeId:     "NodeId02",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc189",
				},
			},
			PowerSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p2",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p3",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
			},
			Receivers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p4",
					NodeName:   "",
					NodeId:     "NodeIdK",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18k",
				},
			},
		}),
		types.NewTask(&pbtypes.TaskPB{
			DataSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p0",
					NodeName:   "",
					NodeId:     "NodeId01",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p1",
					NodeName:   "",
					NodeId:     "NodeId02",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc189",
				},
			},
			PowerSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p2",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p3",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
			},
			Receivers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p4",
					NodeName:   "",
					NodeId:     "NodeIdK",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18j",
				},
			},
		}),
		types.NewTask(&pbtypes.TaskPB{
			DataSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p0",
					NodeName:   "",
					NodeId:     "NodeId01",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18q",
				},
				{
					PartyId:    "p1",
					NodeName:   "",
					NodeId:     "NodeId02",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18u",
				},
			},
			PowerSuppliers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p2",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
				{
					PartyId:    "p3",
					NodeName:   "",
					NodeId:     "NodeId03",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
				},
			},
			Receivers: []*pbtypes.TaskOrganization{
				{
					PartyId:    "p4",
					NodeName:   "",
					NodeId:     "NodeIdK",
					IdentityId: "identity:4d7b5f1f114b43b682d9c73d6d2bc18j",
				},
			},
		}),
	}
	proposalIds := []string{
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac99",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac98",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac97",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac96",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac95",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac94",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac93",
		"0x66fa46e75bf5c0f161ab62f520da5f3d710fe9da87b6ed44a307176c8aedac92",
	}
	identityIds := []string{
		"identity:4d7b5f1f114b43b682d9c73d6d2bc18u",
		"identity:4d7b5f1f114b43b682d9c73d6d2bc189",
		"identity:4d7b5f1f114b43b682d9c73d6d2bc18e",
		"identity:4d7b5f1f114b43b682d9c73d6d2bc18j",
		"identity:4d7b5f1f114b43b682d9c73d6d2bc18k",
		"identity:4d7b5f1f114b43b682d9c73d6d2bc18q",
	}
	generateJsonFile()
	obj := GenerateObg()
	for _, taskOrg := range taskOrgDatas {
		for _, proposalId := range proposalIds {
			obj.CheckConsensusResultOfNoVote(common.HexToHash(proposalId), taskOrg)
		}
	}

	temp := 0
	for _, _ = range obj.FindBlackOrgByWalPrefix() {
		temp += 1
	}
	assert.Equal(t, temp, 6)
	obj.RemoveBlackOrgByIdentity("identity:4d7b5f1f114b43b682d9c73d6d2bc18e")
	temp = 0
	for _, _ = range obj.FindBlackOrgByWalPrefix() {
		temp += 1
	}
	assert.Equal(t, temp, 5)
	for _, identityId := range identityIds {
		obj.RemoveBlackOrgByIdentity(identityId)
	}
	os.Remove("test.json")
}
