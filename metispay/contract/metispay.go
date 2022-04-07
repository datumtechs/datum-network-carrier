// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"math/big"
	"strings"

	platon "github.com/PlatONnetwork/PlatON-Go"
	"github.com/PlatONnetwork/PlatON-Go/accounts/abi"
	"github.com/PlatONnetwork/PlatON-Go/accounts/abi/bind"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = platon.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MetisPayABI is the input ABI used to generate the binding from.
const MetisPayABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"tokenAddressList\",\"type\":\"address[]\"}],\"name\":\"prepay\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"createTask\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"USAGE_FEE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"BASE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"metisLat\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"TaskCreate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddressList\",\"type\":\"address[]\"}],\"name\":\"TaskPrepay\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"refund\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddressList\",\"type\":\"address[]\"}],\"name\":\"TaskSettle\",\"type\":\"event\"}]"

// MetisPayFuncSigs maps the 4-byte function signature to its string representation.
var MetisPayFuncSigs = map[string]string{
	"ec342ad0": "BASE()",
	"d5b96c53": "USAGE_FEE()",
	"f80f5dd5": "addWhitelist(address)",
	"b6a5d7de": "authorize(address)",
	"cbba9f18": "createTask(uint256,address,address)",
	"ff68f6af": "deleteWhitelist(address)",
	"94886421": "prepay(uint256,uint256,address[])",
	"9a9c29f6": "settle(uint256,uint256)",
	"9b19251a": "whitelist(address)",
}

// MetisPayBin is the compiled bytecode used for deploying new contracts.
var MetisPayBin = "0x608060405234801561001057600080fd5b5060405160208061200d8339810180604052602081101561003057600080fd5b50516001600160a01b038116610091576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180611feb6022913960400191505060405180910390fd5b600080546001600160a01b039092166001600160a01b0319909216919091179055611f2a806100c16000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063cbba9f1811610066578063cbba9f181461021a578063d5b96c531461024e578063ec342ad01461024e578063f80f5dd514610268578063ff68f6af1461028e57610093565b806394886421146100985780639a9c29f61461015b5780639b19251a1461017e578063b6a5d7de146101f4575b600080fd5b610147600480360360608110156100ae57600080fd5b8135916020810135918101906060810160408201356401000000008111156100d557600080fd5b8201836020820111156100e757600080fd5b8035906020019184602083028401116401000000008311171561010957600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295506102b4945050505050565b604080519115158252519081900360200190f35b6101476004803603604081101561017157600080fd5b5080359060200135610a24565b6101a46004803603602081101561019457600080fd5b50356001600160a01b03166114b1565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156101e05781810151838201526020016101c8565b505050509050019250505060405180910390f35b6101476004803603602081101561020a57600080fd5b50356001600160a01b0316611586565b6101476004803603606081101561023057600080fd5b508035906001600160a01b0360208201358116916040013516611597565b61025661191b565b60408051918252519081900360200190f35b6101476004803603602081101561027e57600080fd5b50356001600160a01b0316611927565b610147600480360360208110156102a457600080fd5b50356001600160a01b0316611aa9565b60045460009081805b828110156102f15786600482815481106102d357fe5b906000526020600020015414156102e957600191505b6001016102bd565b508061033c5760408051600160e51b62461bcd02815260206004820152600f60248201526001608a1b6e1a5b9d985b1a59081d185cdac81a5902604482015290519081900360640190fd5b6000868152600360205260409020600101546001600160a01b031633146103ad5760408051600160e51b62461bcd02815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f20746869730000000000604482015290519081900360640190fd5b60008681526003602052604090206004015460ff16156104175760408051600160e51b62461bcd02815260206004820152600d60248201527f5265706561742070726570617900000000000000000000000000000000000000604482015290519081900360640190fd5b600080548782526003602052604082205461043f916001600160a01b03908116911630611c5f565b90508581116104985760408051600160e51b62461bcd02815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e636500000000000000604482015290519081900360640190fd5b8451925060005b8381101561054a576104e18682815181106104b657fe5b60209081029190910181015160008b815260039092526040909120546001600160a01b031630611c5f565b9150670de0b6b3a764000082116105425760408051600160e51b62461bcd02815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e636500604482015290519081900360640190fd5b60010161049f565b50600080548882526003602090815260408084205481516001600160a01b03918216602482015230604482015260648082018d9052835180830390910181526084909101835292830180516001600160e01b0316600160e01b63b642fe570217815291518351606095879593169382918083835b602083106105dd5780518252601f1990920191602091820191016105be565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461063f576040519150601f19603f3d011682016040523d82523d6000602084013e610644565b606091505b5090935091508261069f5760408051600160e51b62461bcd02815260206004820152601860248201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604482015290519081900360640190fd5b8180602001905160208110156106b457600080fd5b50519050806106f757604051600160e51b62461bcd028152600401808060200182810382526025815260200180611eda6025913960400191505060405180910390fd5b60005b868110156108d35788818151811061070e57fe5b60209081029190910181015160008d8152600383526040908190205481516001600160a01b039182166024820152306044820152670de0b6b3a7640000606480830191909152835180830390910181526084909101835293840180516001600160e01b0316600160e01b63b642fe5702178152915184519190931693929182918083835b602083106107b15780518252601f199092019160209182019101610792565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610813576040519150601f19603f3d011682016040523d82523d6000602084013e610818565b606091505b509094509250836108735760408051600160e51b62461bcd02815260206004820152601860248201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604482015290519081900360640190fd5b82806020019051602081101561088857600080fd5b50519150816108cb57604051600160e51b62461bcd028152600401808060200182810382526025815260200180611eda6025913960400191505060405180910390fd5b6001016106fa565b5060008a8152600360208181526040909220600281018c90558a51610900939190920191908b0190611dc8565b506001600360008c815260200190815260200160002060040160006101000a81548160ff021916908315150217905550600360008b815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316600360008c815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168b7f7d751fe0754678cea44f8bc0a6faf6baa657d920c8a9b9e8721072bbd0260cf18c8c6040518083815260200180602001828103825283818151815260200191508051906020019060200280838360005b83811015610a005781810151838201526020016109e8565b50505050905001935050505060405180910390a45060019998505050505050505050565b6004546000908180805b83811015610a65578660048281548110610a4457fe5b90600052602060002001541415610a5d57600192508091505b600101610a2e565b5081610ab05760408051600160e51b62461bcd02815260206004820152600f60248201526001608a1b6e1a5b9d985b1a59081d185cdac81a5902604482015290519081900360640190fd5b6000868152600360205260409020600101546001600160a01b03163314610b215760408051600160e51b62461bcd02815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f20746869730000000000604482015290519081900360640190fd5b60008681526003602052604090206004015460ff16610b8a5760408051600160e51b62461bcd02815260206004820152601460248201527f707265706179206e6f7420636f6d706c65746564000000000000000000000000604482015290519081900360640190fd5b600086815260036020526040902060040154610100900460ff1615610bf95760408051600160e51b62461bcd02815260206004820152600d60248201527f52657065617420736574746c6500000000000000000000000000000000000000604482015290519081900360640190fd5b6000868152600360205260409020600201548510610c615760408051600160e51b62461bcd02815260206004820152601560248201527f6d6f7265207468616e2070726570616964206c61740000000000000000000000604482015290519081900360640190fd5b600080548782526003602090815260408084206001015481516001600160a01b03918216602482015260448082018c9052835180830390910181526064909101835292830180516001600160e01b0316600160e21b632758748d0217815291518351606095879593169382918083835b60208310610cf05780518252601f199092019160209182019101610cd1565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610d52576040519150601f19603f3d011682016040523d82523d6000602084013e610d57565b606091505b50909350915082610dac5760408051600160e51b62461bcd0281526020600482015260146024820152600160621b7318d85b1b081d1c985b9cd9995c8819985a5b195902604482015290519081900360640190fd5b818060200190516020811015610dc157600080fd5b5051905080610e0457604051600160e51b62461bcd028152600401808060200182810382526021815260200180611eb96021913960400191505060405180910390fd5b600089815260036020908152604080832060028101549354905482516001600160a01b039182166024820152948d90036044808701829052845180880390910181526064909601845293850180516001600160e01b0316600160e21b632758748d021781529251855194959190921693909282918083835b60208310610e9b5780518252601f199092019160209182019101610e7c565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610efd576040519150601f19603f3d011682016040523d82523d6000602084013e610f02565b606091505b50909450925083610f575760408051600160e51b62461bcd0281526020600482015260146024820152600160621b7318d85b1b081d1c985b9cd9995c8819985a5b195902604482015290519081900360640190fd5b828060200190516020811015610f6c57600080fd5b5051915081610faf57604051600160e51b62461bcd028152600401808060200182810382526021815260200180611eb96021913960400191505060405180910390fd5b60008a815260036020818152604080842090920180548351818402810184019094528084526060939283018282801561101157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610ff3575b505083519c509293506000925050505b898110156113395781818151811061103557fe5b602090810291909101810151604080516004815260248101825292830180516001600160e01b0316600160e11b6303aa30b902178152905183516001600160a01b039093169392909182918083835b602083106110a35780518252601f199092019160209182019101611084565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611105576040519150601f19603f3d011682016040523d82523d6000602084013e61110a565b606091505b509097509550866111655760408051600160e51b62461bcd02815260206004820152601260248201527f63616c6c206d696e746572206661696c65640000000000000000000000000000604482015290519081900360640190fd5b85806020019051602081101561117a57600080fd5b5051825190935082908290811061118d57fe5b602090810291909101810151604080516001600160a01b038781166024830152670de0b6b3a7640000604480840191909152835180840390910181526064909201835293810180516001600160e01b0316600160e21b632758748d0217815291518151949093169390929182918083835b6020831061121d5780518252601f1990920191602091820191016111fe565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461127f576040519150601f19603f3d011682016040523d82523d6000602084013e611284565b606091505b509097509550866112d95760408051600160e51b62461bcd0281526020600482015260146024820152600160621b7318d85b1b081d1c985b9cd9995c8819985a5b195902604482015290519081900360640190fd5b8580602001905160208110156112ee57600080fd5b505194508461133157604051600160e51b62461bcd028152600401808060200182810382526021815260200180611eb96021913960400191505060405180910390fd5b600101611021565b50600360008d815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316600360008e815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168d7fa6c1cad6a3fe81c601fc1d9aed900bdfa5d1c584370e07028037d667fd07ca9b8e87866040518084815260200183815260200180602001828103825283818151815260200191508051906020019060200280838360005b838110156114115781810151838201526020016113f9565b5050505090500194505050505060405180910390a46004878154811061143357fe5b60009182526020822001556004805490611451906000198301611e2d565b5060008c8152600360208190526040822080546001600160a01b03199081168255600182018054909116905560028101839055919061149290830182611e56565b50600401805461ffff191690555060019b9a5050505050505050505050565b60025460609081906000805b8281101561150457856001600160a01b0316600282815481106114dc57fe5b6000918252602090912001546001600160a01b031614156114fc57600191505b6001016114bd565b50801561157b576001600160a01b0385166000908152600160209081526040918290208054835181840281018401909452808452909183018282801561157357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611555575b505050505092505b50909150505b919050565b600061159182611927565b92915050565b6002546000908190815b818110156115e857336001600160a01b0316600282815481106115c057fe5b6000918252602090912001546001600160a01b031614156115e057600192505b6001016115a1565b508161163e5760408051600160e51b62461bcd02815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f7420657869737400604482015290519081900360640190fd5b33600090815260016020526040812054815b818110156116a75733600090815260016020526040902080546001600160a01b03891691908390811061167f57fe5b6000918252602090912001546001600160a01b0316141561169f57600192505b600101611650565b50816116fd5760408051600160e51b62461bcd02815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a65640000000000000000604482015290519081900360640190fd5b600454600090815b81811015611739578a6004828154811061171b57fe5b9060005260206000200154141561173157600192505b600101611705565b5081156117905760408051600160e51b62461bcd02815260206004820152601660248201527f7461736b20696420616c72656164792065786973747300000000000000000000604482015290519081900360640190fd5b60606040518060c001604052808b6001600160a01b031681526020018a6001600160a01b031681526020016000815260200182815260200160001515815260200160001515815250600360008d815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160020155606082015181600301908051906020019061186b929190611dc8565b5060808201516004918201805460a09094015115156101000261ff001992151560ff19909516949094179190911692909217909155805460018101825560009182527f8a35acfbc15ff81a39ae7d344fd709f28e8600b4aa8c65c6b64bfe7fe36bd19b018c90556040516001600160a01b03808c1692908d16918e917f05de3a4f956df2caba4856353901924f16f116b62eb25e7842c926e02288610591a45060019a9950505050505050505050565b670de0b6b3a764000081565b60025460009081805b8281101561197757336001600160a01b03166002828154811061194f57fe5b6000918252602090912001546001600160a01b0316141561196f57600191505b600101611930565b508015611a2d573360009081526001602052604081205490805b828110156119e85733600090815260016020526040902080546001600160a01b0389169190839081106119c057fe5b6000918252602090912001546001600160a01b031614156119e057600191505b600101611991565b5080611a265733600090815260016020818152604083208054928301815583529091200180546001600160a01b0319166001600160a01b0388161790555b5050611a9f565b6002805460018082019092557f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054336001600160a01b03199182168117909255600091825260208381526040832080549485018155835290912090910180549091166001600160a01b0386161790555b5060019392505050565b6002546000908180805b83811015611afd57336001600160a01b031660028281548110611ad257fe5b6000918252602090912001546001600160a01b03161415611af557600192508091505b600101611ab3565b508115611c4857336000908152600160205260408120549080805b83811015611b725733600090815260016020526040902080546001600160a01b038b16919083908110611b4757fe5b6000918252602090912001546001600160a01b03161415611b6a57600192508091505b600101611b18565b508115611c31576001831115611bda57336000908152600160205260409020805482908110611b9d57fe5b6000918252602080832090910180546001600160a01b03191690553382526001905260409020805490611bd4906000198301611e2d565b50611c2c565b60028481548110611be757fe5b600091825260209091200180546001600160a01b03191690556002805490611c13906000198301611e2d565b50336000908152600160205260408120611c2c91611e56565b611c40565b60009650505050505050611581565b505050611c54565b60009350505050611581565b506001949350505050565b604080516001600160a01b03848116602483015283811660448084019190915283518084039091018152606490920183526020820180516001600160e01b0316600160e11b636eb1769f021781529251825160009485946060948a16939092909182918083835b60208310611ce55780518252601f199092019160209182019101611cc6565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d8060008114611d45576040519150601f19603f3d011682016040523d82523d6000602084013e611d4a565b606091505b509150915081611da45760408051600160e51b62461bcd02815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c65640000000000604482015290519081900360640190fd5b6000818060200190516020811015611dbb57600080fd5b5051979650505050505050565b828054828255906000526020600020908101928215611e1d579160200282015b82811115611e1d57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611de8565b50611e29929150611e77565b5090565b815481835581811115611e5157600083815260209020611e51918101908301611e9e565b505050565b5080546000825590600052602060002090810190611e749190611e9e565b50565b611e9b91905b80821115611e295780546001600160a01b0319168155600101611e7d565b90565b611e9b91905b80821115611e295760008155600101611ea456fe5468652072657475726e206f66207472616e73666572206973206661696c7572655468652072657475726e206f66207472616e7366657266726f6d206973206661696c757265a165627a7a723058206a084c917505d4b0135d69d61f56446705f4b3df96394c97d7b4b83cbc4e489800294d657469735061793a20496e76616c6964206d657469734c61742061646472657373"

// DeployMetisPay deploys a new platon contract, binding an instance of MetisPay to it.
func DeployMetisPay(auth *bind.TransactOpts, backend bind.ContractBackend, metisLat common.Address) (common.Address, *types.Transaction, *MetisPay, error) {
	parsed, err := abi.JSON(strings.NewReader(MetisPayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MetisPayBin), backend, metisLat)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MetisPay{MetisPayCaller: MetisPayCaller{contract: contract}, MetisPayTransactor: MetisPayTransactor{contract: contract}, MetisPayFilterer: MetisPayFilterer{contract: contract}}, nil
}

// MetisPay is an auto generated Go binding around an platon contract.
type MetisPay struct {
	MetisPayCaller     // Read-only binding to the contract
	MetisPayTransactor // Write-only binding to the contract
	MetisPayFilterer   // Log filterer for contract events
}

// MetisPayCaller is an auto generated read-only Go binding around an platon contract.
type MetisPayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetisPayTransactor is an auto generated write-only Go binding around an platon contract.
type MetisPayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetisPayFilterer is an auto generated log filtering Go binding around an platon contract events.
type MetisPayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetisPaySession is an auto generated Go binding around an platon contract,
// with pre-set call and transact options.
type MetisPaySession struct {
	Contract     *MetisPay         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MetisPayCallerSession is an auto generated read-only Go binding around an platon contract,
// with pre-set call options.
type MetisPayCallerSession struct {
	Contract *MetisPayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MetisPayTransactorSession is an auto generated write-only Go binding around an platon contract,
// with pre-set transact options.
type MetisPayTransactorSession struct {
	Contract     *MetisPayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MetisPayRaw is an auto generated low-level Go binding around an platon contract.
type MetisPayRaw struct {
	Contract *MetisPay // Generic contract binding to access the raw methods on
}

// MetisPayCallerRaw is an auto generated low-level read-only Go binding around an platon contract.
type MetisPayCallerRaw struct {
	Contract *MetisPayCaller // Generic read-only contract binding to access the raw methods on
}

// MetisPayTransactorRaw is an auto generated low-level write-only Go binding around an platon contract.
type MetisPayTransactorRaw struct {
	Contract *MetisPayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMetisPay creates a new instance of MetisPay, bound to a specific deployed contract.
func NewMetisPay(address common.Address, backend bind.ContractBackend) (*MetisPay, error) {
	contract, err := bindMetisPay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MetisPay{MetisPayCaller: MetisPayCaller{contract: contract}, MetisPayTransactor: MetisPayTransactor{contract: contract}, MetisPayFilterer: MetisPayFilterer{contract: contract}}, nil
}

// NewMetisPayCaller creates a new read-only instance of MetisPay, bound to a specific deployed contract.
func NewMetisPayCaller(address common.Address, caller bind.ContractCaller) (*MetisPayCaller, error) {
	contract, err := bindMetisPay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MetisPayCaller{contract: contract}, nil
}

// NewMetisPayTransactor creates a new write-only instance of MetisPay, bound to a specific deployed contract.
func NewMetisPayTransactor(address common.Address, transactor bind.ContractTransactor) (*MetisPayTransactor, error) {
	contract, err := bindMetisPay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MetisPayTransactor{contract: contract}, nil
}

// NewMetisPayFilterer creates a new log filterer instance of MetisPay, bound to a specific deployed contract.
func NewMetisPayFilterer(address common.Address, filterer bind.ContractFilterer) (*MetisPayFilterer, error) {
	contract, err := bindMetisPay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MetisPayFilterer{contract: contract}, nil
}

// bindMetisPay binds a generic wrapper to an already deployed contract.
func bindMetisPay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MetisPayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MetisPay *MetisPayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MetisPay.Contract.MetisPayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MetisPay *MetisPayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MetisPay.Contract.MetisPayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MetisPay *MetisPayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MetisPay.Contract.MetisPayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MetisPay *MetisPayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MetisPay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MetisPay *MetisPayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MetisPay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MetisPay *MetisPayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MetisPay.Contract.contract.Transact(opts, method, params...)
}

// BASE is a free data retrieval call binding the contract method 0xec342ad0.
//
// Solidity: function BASE() view returns(uint256)
func (_MetisPay *MetisPayCaller) BASE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MetisPay.contract.Call(opts, &out, "BASE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BASE is a free data retrieval call binding the contract method 0xec342ad0.
//
// Solidity: function BASE() view returns(uint256)
func (_MetisPay *MetisPaySession) BASE() (*big.Int, error) {
	return _MetisPay.Contract.BASE(&_MetisPay.CallOpts)
}

// BASE is a free data retrieval call binding the contract method 0xec342ad0.
//
// Solidity: function BASE() view returns(uint256)
func (_MetisPay *MetisPayCallerSession) BASE() (*big.Int, error) {
	return _MetisPay.Contract.BASE(&_MetisPay.CallOpts)
}

// USAGEFEE is a free data retrieval call binding the contract method 0xd5b96c53.
//
// Solidity: function USAGE_FEE() view returns(uint256)
func (_MetisPay *MetisPayCaller) USAGEFEE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MetisPay.contract.Call(opts, &out, "USAGE_FEE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// USAGEFEE is a free data retrieval call binding the contract method 0xd5b96c53.
//
// Solidity: function USAGE_FEE() view returns(uint256)
func (_MetisPay *MetisPaySession) USAGEFEE() (*big.Int, error) {
	return _MetisPay.Contract.USAGEFEE(&_MetisPay.CallOpts)
}

// USAGEFEE is a free data retrieval call binding the contract method 0xd5b96c53.
//
// Solidity: function USAGE_FEE() view returns(uint256)
func (_MetisPay *MetisPayCallerSession) USAGEFEE() (*big.Int, error) {
	return _MetisPay.Contract.USAGEFEE(&_MetisPay.CallOpts)
}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address userAddress) view returns(address[])
func (_MetisPay *MetisPayCaller) Whitelist(opts *bind.CallOpts, userAddress common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _MetisPay.contract.Call(opts, &out, "whitelist", userAddress)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address userAddress) view returns(address[])
func (_MetisPay *MetisPaySession) Whitelist(userAddress common.Address) ([]common.Address, error) {
	return _MetisPay.Contract.Whitelist(&_MetisPay.CallOpts, userAddress)
}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address userAddress) view returns(address[])
func (_MetisPay *MetisPayCallerSession) Whitelist(userAddress common.Address) ([]common.Address, error) {
	return _MetisPay.Contract.Whitelist(&_MetisPay.CallOpts, userAddress)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0xf80f5dd5.
//
// Solidity: function addWhitelist(address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactor) AddWhitelist(opts *bind.TransactOpts, userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "addWhitelist", userAgency)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0xf80f5dd5.
//
// Solidity: function addWhitelist(address userAgency) returns(bool success)
func (_MetisPay *MetisPaySession) AddWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.AddWhitelist(&_MetisPay.TransactOpts, userAgency)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0xf80f5dd5.
//
// Solidity: function addWhitelist(address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) AddWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.AddWhitelist(&_MetisPay.TransactOpts, userAgency)
}

// Authorize is a paid mutator transaction binding the contract method 0xb6a5d7de.
//
// Solidity: function authorize(address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactor) Authorize(opts *bind.TransactOpts, userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "authorize", userAgency)
}

// Authorize is a paid mutator transaction binding the contract method 0xb6a5d7de.
//
// Solidity: function authorize(address userAgency) returns(bool success)
func (_MetisPay *MetisPaySession) Authorize(userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Authorize(&_MetisPay.TransactOpts, userAgency)
}

// Authorize is a paid mutator transaction binding the contract method 0xb6a5d7de.
//
// Solidity: function authorize(address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) Authorize(userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Authorize(&_MetisPay.TransactOpts, userAgency)
}

// CreateTask is a paid mutator transaction binding the contract method 0xcbba9f18.
//
// Solidity: function createTask(uint256 taskId, address user, address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactor) CreateTask(opts *bind.TransactOpts, taskId *big.Int, user common.Address, userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "createTask", taskId, user, userAgency)
}

// CreateTask is a paid mutator transaction binding the contract method 0xcbba9f18.
//
// Solidity: function createTask(uint256 taskId, address user, address userAgency) returns(bool success)
func (_MetisPay *MetisPaySession) CreateTask(taskId *big.Int, user common.Address, userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.CreateTask(&_MetisPay.TransactOpts, taskId, user, userAgency)
}

// CreateTask is a paid mutator transaction binding the contract method 0xcbba9f18.
//
// Solidity: function createTask(uint256 taskId, address user, address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) CreateTask(taskId *big.Int, user common.Address, userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.CreateTask(&_MetisPay.TransactOpts, taskId, user, userAgency)
}

// DeleteWhitelist is a paid mutator transaction binding the contract method 0xff68f6af.
//
// Solidity: function deleteWhitelist(address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactor) DeleteWhitelist(opts *bind.TransactOpts, userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "deleteWhitelist", userAgency)
}

// DeleteWhitelist is a paid mutator transaction binding the contract method 0xff68f6af.
//
// Solidity: function deleteWhitelist(address userAgency) returns(bool success)
func (_MetisPay *MetisPaySession) DeleteWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.DeleteWhitelist(&_MetisPay.TransactOpts, userAgency)
}

// DeleteWhitelist is a paid mutator transaction binding the contract method 0xff68f6af.
//
// Solidity: function deleteWhitelist(address userAgency) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) DeleteWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.DeleteWhitelist(&_MetisPay.TransactOpts, userAgency)
}

// Prepay is a paid mutator transaction binding the contract method 0x94886421.
//
// Solidity: function prepay(uint256 taskId, uint256 fee, address[] tokenAddressList) returns(bool success)
func (_MetisPay *MetisPayTransactor) Prepay(opts *bind.TransactOpts, taskId *big.Int, fee *big.Int, tokenAddressList []common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "prepay", taskId, fee, tokenAddressList)
}

// Prepay is a paid mutator transaction binding the contract method 0x94886421.
//
// Solidity: function prepay(uint256 taskId, uint256 fee, address[] tokenAddressList) returns(bool success)
func (_MetisPay *MetisPaySession) Prepay(taskId *big.Int, fee *big.Int, tokenAddressList []common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Prepay(&_MetisPay.TransactOpts, taskId, fee, tokenAddressList)
}

// Prepay is a paid mutator transaction binding the contract method 0x94886421.
//
// Solidity: function prepay(uint256 taskId, uint256 fee, address[] tokenAddressList) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) Prepay(taskId *big.Int, fee *big.Int, tokenAddressList []common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Prepay(&_MetisPay.TransactOpts, taskId, fee, tokenAddressList)
}

// Settle is a paid mutator transaction binding the contract method 0x9a9c29f6.
//
// Solidity: function settle(uint256 taskId, uint256 fee) returns(bool success)
func (_MetisPay *MetisPayTransactor) Settle(opts *bind.TransactOpts, taskId *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "settle", taskId, fee)
}

// Settle is a paid mutator transaction binding the contract method 0x9a9c29f6.
//
// Solidity: function settle(uint256 taskId, uint256 fee) returns(bool success)
func (_MetisPay *MetisPaySession) Settle(taskId *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _MetisPay.Contract.Settle(&_MetisPay.TransactOpts, taskId, fee)
}

// Settle is a paid mutator transaction binding the contract method 0x9a9c29f6.
//
// Solidity: function settle(uint256 taskId, uint256 fee) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) Settle(taskId *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _MetisPay.Contract.Settle(&_MetisPay.TransactOpts, taskId, fee)
}

// MetisPayTaskCreateIterator is returned from FilterTaskCreate and is used to iterate over the raw logs and unpacked data for TaskCreate events raised by the MetisPay contract.
type MetisPayTaskCreateIterator struct {
	Event *MetisPayTaskCreate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  platon.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MetisPayTaskCreateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetisPayTaskCreate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MetisPayTaskCreate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MetisPayTaskCreateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetisPayTaskCreateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetisPayTaskCreate represents a TaskCreate event raised by the MetisPay contract.
type MetisPayTaskCreate struct {
	TaskId     *big.Int
	User       common.Address
	UserAgency common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTaskCreate is a free log retrieval operation binding the contract event 0x05de3a4f956df2caba4856353901924f16f116b62eb25e7842c926e022886105.
//
// Solidity: event TaskCreate(uint256 indexed taskId, address indexed user, address indexed userAgency)
func (_MetisPay *MetisPayFilterer) FilterTaskCreate(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*MetisPayTaskCreateIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var userAgencyRule []interface{}
	for _, userAgencyItem := range userAgency {
		userAgencyRule = append(userAgencyRule, userAgencyItem)
	}

	logs, sub, err := _MetisPay.contract.FilterLogs(opts, "TaskCreate", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &MetisPayTaskCreateIterator{contract: _MetisPay.contract, event: "TaskCreate", logs: logs, sub: sub}, nil
}

// WatchTaskCreate is a free log subscription operation binding the contract event 0x05de3a4f956df2caba4856353901924f16f116b62eb25e7842c926e022886105.
//
// Solidity: event TaskCreate(uint256 indexed taskId, address indexed user, address indexed userAgency)
func (_MetisPay *MetisPayFilterer) WatchTaskCreate(opts *bind.WatchOpts, sink chan<- *MetisPayTaskCreate, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var userAgencyRule []interface{}
	for _, userAgencyItem := range userAgency {
		userAgencyRule = append(userAgencyRule, userAgencyItem)
	}

	logs, sub, err := _MetisPay.contract.WatchLogs(opts, "TaskCreate", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetisPayTaskCreate)
				if err := _MetisPay.contract.UnpackLog(event, "TaskCreate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskCreate is a log parse operation binding the contract event 0x05de3a4f956df2caba4856353901924f16f116b62eb25e7842c926e022886105.
//
// Solidity: event TaskCreate(uint256 indexed taskId, address indexed user, address indexed userAgency)
func (_MetisPay *MetisPayFilterer) ParseTaskCreate(log types.Log) (*MetisPayTaskCreate, error) {
	event := new(MetisPayTaskCreate)
	if err := _MetisPay.contract.UnpackLog(event, "TaskCreate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MetisPayTaskPrepayIterator is returned from FilterTaskPrepay and is used to iterate over the raw logs and unpacked data for TaskPrepay events raised by the MetisPay contract.
type MetisPayTaskPrepayIterator struct {
	Event *MetisPayTaskPrepay // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  platon.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MetisPayTaskPrepayIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetisPayTaskPrepay)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MetisPayTaskPrepay)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MetisPayTaskPrepayIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetisPayTaskPrepayIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetisPayTaskPrepay represents a TaskPrepay event raised by the MetisPay contract.
type MetisPayTaskPrepay struct {
	TaskId           *big.Int
	User             common.Address
	UserAgency       common.Address
	Fee              *big.Int
	TokenAddressList []common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterTaskPrepay is a free log retrieval operation binding the contract event 0x7d751fe0754678cea44f8bc0a6faf6baa657d920c8a9b9e8721072bbd0260cf1.
//
// Solidity: event TaskPrepay(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) FilterTaskPrepay(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*MetisPayTaskPrepayIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var userAgencyRule []interface{}
	for _, userAgencyItem := range userAgency {
		userAgencyRule = append(userAgencyRule, userAgencyItem)
	}

	logs, sub, err := _MetisPay.contract.FilterLogs(opts, "TaskPrepay", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &MetisPayTaskPrepayIterator{contract: _MetisPay.contract, event: "TaskPrepay", logs: logs, sub: sub}, nil
}

// WatchTaskPrepay is a free log subscription operation binding the contract event 0x7d751fe0754678cea44f8bc0a6faf6baa657d920c8a9b9e8721072bbd0260cf1.
//
// Solidity: event TaskPrepay(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) WatchTaskPrepay(opts *bind.WatchOpts, sink chan<- *MetisPayTaskPrepay, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var userAgencyRule []interface{}
	for _, userAgencyItem := range userAgency {
		userAgencyRule = append(userAgencyRule, userAgencyItem)
	}

	logs, sub, err := _MetisPay.contract.WatchLogs(opts, "TaskPrepay", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetisPayTaskPrepay)
				if err := _MetisPay.contract.UnpackLog(event, "TaskPrepay", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskPrepay is a log parse operation binding the contract event 0x7d751fe0754678cea44f8bc0a6faf6baa657d920c8a9b9e8721072bbd0260cf1.
//
// Solidity: event TaskPrepay(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) ParseTaskPrepay(log types.Log) (*MetisPayTaskPrepay, error) {
	event := new(MetisPayTaskPrepay)
	if err := _MetisPay.contract.UnpackLog(event, "TaskPrepay", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MetisPayTaskSettleIterator is returned from FilterTaskSettle and is used to iterate over the raw logs and unpacked data for TaskSettle events raised by the MetisPay contract.
type MetisPayTaskSettleIterator struct {
	Event *MetisPayTaskSettle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  platon.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MetisPayTaskSettleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetisPayTaskSettle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(MetisPayTaskSettle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *MetisPayTaskSettleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetisPayTaskSettleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetisPayTaskSettle represents a TaskSettle event raised by the MetisPay contract.
type MetisPayTaskSettle struct {
	TaskId           *big.Int
	User             common.Address
	UserAgency       common.Address
	Agencyfee        *big.Int
	Refund           *big.Int
	TokenAddressList []common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterTaskSettle is a free log retrieval operation binding the contract event 0xa6c1cad6a3fe81c601fc1d9aed900bdfa5d1c584370e07028037d667fd07ca9b.
//
// Solidity: event TaskSettle(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) FilterTaskSettle(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*MetisPayTaskSettleIterator, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var userAgencyRule []interface{}
	for _, userAgencyItem := range userAgency {
		userAgencyRule = append(userAgencyRule, userAgencyItem)
	}

	logs, sub, err := _MetisPay.contract.FilterLogs(opts, "TaskSettle", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &MetisPayTaskSettleIterator{contract: _MetisPay.contract, event: "TaskSettle", logs: logs, sub: sub}, nil
}

// WatchTaskSettle is a free log subscription operation binding the contract event 0xa6c1cad6a3fe81c601fc1d9aed900bdfa5d1c584370e07028037d667fd07ca9b.
//
// Solidity: event TaskSettle(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) WatchTaskSettle(opts *bind.WatchOpts, sink chan<- *MetisPayTaskSettle, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

	var taskIdRule []interface{}
	for _, taskIdItem := range taskId {
		taskIdRule = append(taskIdRule, taskIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var userAgencyRule []interface{}
	for _, userAgencyItem := range userAgency {
		userAgencyRule = append(userAgencyRule, userAgencyItem)
	}

	logs, sub, err := _MetisPay.contract.WatchLogs(opts, "TaskSettle", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetisPayTaskSettle)
				if err := _MetisPay.contract.UnpackLog(event, "TaskSettle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTaskSettle is a log parse operation binding the contract event 0xa6c1cad6a3fe81c601fc1d9aed900bdfa5d1c584370e07028037d667fd07ca9b.
//
// Solidity: event TaskSettle(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) ParseTaskSettle(log types.Log) (*MetisPayTaskSettle, error) {
	event := new(MetisPayTaskSettle)
	if err := _MetisPay.contract.UnpackLog(event, "TaskSettle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
