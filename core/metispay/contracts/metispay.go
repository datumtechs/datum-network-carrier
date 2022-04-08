// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// MetisPayABI is the input ABI used to generate the binding from.
const MetisPayABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"tokenAddressList\",\"type\":\"address[]\"}],\"name\":\"prepay\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"metisLat\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"USAGE_FEE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"BASE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"metisLat\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddressList\",\"type\":\"address[]\"}],\"name\":\"PrepayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"refund\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddressList\",\"type\":\"address[]\"}],\"name\":\"SettleEvent\",\"type\":\"event\"}]"

// MetisPayFuncSigs maps the 4-byte function signature to its string representation.
var MetisPayFuncSigs = map[string]string{
	"ec342ad0": "BASE()",
	"d5b96c53": "USAGE_FEE()",
	"f80f5dd5": "addWhitelist(address)",
	"b6a5d7de": "authorize(address)",
	"ff68f6af": "deleteWhitelist(address)",
	"c4d66de8": "initialize(address)",
	"5d647525": "prepay(uint256,address,uint256,address[])",
	"9a9c29f6": "settle(uint256,uint256)",
	"9b19251a": "whitelist(address)",
}

// MetisPayBin is the compiled bytecode used for deploying new contracts.
var MetisPayBin = "0x60806040526003805460ff1916905534801561001a57600080fd5b50604051602080620020ac8339810180604052602081101561003b57600080fd5b505161004d81610054602090811b901c565b50506100eb565b60006001600160a01b0382166100b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001806200208a6022913960400191505060405180910390fd5b50600080546001600160a01b0383166001600160a01b03199091161790556003805460ff19166001179081905560ff16919050565b611f8f80620000fb6000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063c4d66de811610066578063c4d66de81461022a578063d5b96c5314610250578063ec342ad014610250578063f80f5dd51461026a578063ff68f6af1461029057610093565b80635d647525146100985780639a9c29f61461016b5780639b19251a1461018e578063b6a5d7de14610204575b600080fd5b610157600480360360808110156100ae57600080fd5b8135916001600160a01b0360208201351691604082013591908101906080810160608201356401000000008111156100e557600080fd5b8201836020820111156100f757600080fd5b8035906020019184602083028401116401000000008311171561011957600080fd5b9190808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152509295506102b6945050505050565b604080519115158252519081900360200190f35b6101576004803603604081101561018157600080fd5b5080359060200135610b1f565b6101b4600480360360208110156101a457600080fd5b50356001600160a01b0316611602565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156101f05781810151838201526020016101d8565b505050509050019250505060405180910390f35b6101576004803603602081101561021a57600080fd5b50356001600160a01b03166116c7565b6101576004803603602081101561024057600080fd5b50356001600160a01b03166116d8565b610258611729565b60408051918252519081900360200190f35b6101576004803603602081101561028057600080fd5b50356001600160a01b0316611735565b610157600480360360208110156102a657600080fd5b50356001600160a01b03166118b4565b6002546000908190815b8181101561030757866001600160a01b0316600282815481106102df57fe5b6000918252602090912001546001600160a01b031614156102ff57600192505b6001016102c0565b508161035d5760408051600160e51b62461bcd02815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f7420657869737400604482015290519081900360640190fd5b50506001600160a01b038416600090815260016020526040812054815b818110156103d1576001600160a01b03871660009081526001602052604090208054339190839081106103a957fe5b6000918252602090912001546001600160a01b031614156103c957600192505b60010161037a565b50816104275760408051600160e51b62461bcd02815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a65640000000000000000604482015290519081900360640190fd5b5050600554600090815b8181101561046557876005828154811061044757fe5b9060005260206000200154141561045d57600192505b600101610431565b5081156104bc5760408051600160e51b62461bcd02815260206004820152601660248201527f7461736b20696420616c72656164792065786973747300000000000000000000604482015290519081900360640190fd5b600080546104d4906001600160a01b03168830611bf4565b905085811161052d5760408051600160e51b62461bcd02815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e636500000000000000604482015290519081900360640190fd5b8451915060005b828110156105c35761055a86828151811061054b57fe5b60200260200101518930611bf4565b9150670de0b6b3a764000082116105bb5760408051600160e51b62461bcd02815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e636500604482015290519081900360640190fd5b600101610534565b5060008054604080516001600160a01b038b8116602483015230604483015260648083018c905283518084039091018152608490920183526020820180516001600160e01b0316600160e01b63b642fe57021781529251825160609587959316939282918083835b6020831061064a5780518252601f19909201916020918201910161062b565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d80600081146106ac576040519150601f19603f3d011682016040523d82523d6000602084013e6106b1565b606091505b5090935091508261070c5760408051600160e51b62461bcd02815260206004820152601860248201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604482015290519081900360640190fd5b81806020019051602081101561072157600080fd5b505190508061076457604051600160e51b62461bcd028152600401808060200182810382526025815260200180611f1d6025913960400191505060405180910390fd5b60005b858110156109335788818151811061077b57fe5b602090810291909101810151604080516001600160a01b038f81166024830152306044830152670de0b6b3a7640000606480840191909152835180840390910181526084909201835293810180516001600160e01b0316600160e01b63b642fe570217815291518151949093169390929182918083835b602083106108115780518252601f1990920191602091820191016107f2565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610873576040519150601f19603f3d011682016040523d82523d6000602084013e610878565b606091505b509094509250836108d35760408051600160e51b62461bcd02815260206004820152601860248201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604482015290519081900360640190fd5b8280602001905160208110156108e857600080fd5b505191508161092b57604051600160e51b62461bcd028152600401808060200182810382526025815260200180611f1d6025913960400191505060405180910390fd5b600101610767565b506040518060c001604052808b6001600160a01b03168152602001336001600160a01b031681526020018a815260200189815260200160011515815260200160001515815250600460008d815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550604082015181600201556060820151816003019080519060200190610a0c929190611ddc565b5060808201518160040160006101000a81548160ff02191690831515021790555060a08201518160040160016101000a81548160ff02191690831515021790555090505060058b9080600181540180825580915050906001820390600052602060002001600090919290919091505550336001600160a01b03168a6001600160a01b03168c7f2f05b3233a63a604771d18ab3e4bf65f10a488ea624dd4f1499d274b1bd7da598c8c6040518083815260200180602001828103825283818151815260200191508051906020019060200280838360005b83811015610afa578181015183820152602001610ae2565b50505050905001935050505060405180910390a45060019a9950505050505050505050565b6005546000908180805b83811015610b60578660058281548110610b3f57fe5b90600052602060002001541415610b5857600192508091505b600101610b29565b5081610bb65760408051600160e51b62461bcd02815260206004820152600f60248201527f696e76616c6964207461736b2069640000000000000000000000000000000000604482015290519081900360640190fd5b6000868152600460205260409020600101546001600160a01b03163314610c275760408051600160e51b62461bcd02815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f20746869730000000000604482015290519081900360640190fd5b6000868152600460208190526040909120015460ff16610c915760408051600160e51b62461bcd02815260206004820152601460248201527f707265706179206e6f7420636f6d706c65746564000000000000000000000000604482015290519081900360640190fd5b60008681526004602081905260409091200154610100900460ff1615610d015760408051600160e51b62461bcd02815260206004820152600d60248201527f52657065617420736574746c6500000000000000000000000000000000000000604482015290519081900360640190fd5b6000868152600460205260409020600201548510610d695760408051600160e51b62461bcd02815260206004820152601560248201527f6d6f7265207468616e2070726570616964206c61740000000000000000000000604482015290519081900360640190fd5b600080548782526004602090815260408084206001015481516001600160a01b03918216602482015260448082018c9052835180830390910181526064909101835292830180516001600160e01b0316600160e21b632758748d0217815291518351606095879593169382918083835b60208310610df85780518252601f199092019160209182019101610dd9565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610e5a576040519150601f19603f3d011682016040523d82523d6000602084013e610e5f565b606091505b50909350915082610eb45760408051600160e51b62461bcd0281526020600482015260146024820152600160621b7318d85b1b081d1c985b9cd9995c8819985a5b195902604482015290519081900360640190fd5b818060200190516020811015610ec957600080fd5b5051905080610f0c57604051600160e51b62461bcd028152600401808060200182810382526021815260200180611ecd6021913960400191505060405180910390fd5b600089815260046020908152604080832060028101549354905482516001600160a01b039182166024820152948d90036044808701829052845180880390910181526064909601845293850180516001600160e01b0316600160e21b632758748d021781529251855194959190921693909282918083835b60208310610fa35780518252601f199092019160209182019101610f84565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611005576040519150601f19603f3d011682016040523d82523d6000602084013e61100a565b606091505b5090945092508361105f5760408051600160e51b62461bcd0281526020600482015260146024820152600160621b7318d85b1b081d1c985b9cd9995c8819985a5b195902604482015290519081900360640190fd5b82806020019051602081101561107457600080fd5b50519150816110b757604051600160e51b62461bcd028152600401808060200182810382526021815260200180611ecd6021913960400191505060405180910390fd5b60008a8152600460209081526040808320600301805482518185028101850190935280835260609383018282801561111857602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116110fa575b505083519394506000925050505b8181101561143e5782818151811061113a57fe5b602090810291909101810151604080516004815260248101825292830180516001600160e01b0316600160e11b6303aa30b902178152905183516001600160a01b039093169392909182918083835b602083106111a85780518252601f199092019160209182019101611189565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461120a576040519150601f19603f3d011682016040523d82523d6000602084013e61120f565b606091505b5090985096508761126a5760408051600160e51b62461bcd02815260206004820152601260248201527f63616c6c206d696e746572206661696c65640000000000000000000000000000604482015290519081900360640190fd5b86806020019051602081101561127f57600080fd5b5051835190945083908290811061129257fe5b602090810291909101810151604080516001600160a01b038881166024830152670de0b6b3a7640000604480840191909152835180840390910181526064909201835293810180516001600160e01b0316600160e21b632758748d0217815291518151949093169390929182918083835b602083106113225780518252601f199092019160209182019101611303565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114611384576040519150601f19603f3d011682016040523d82523d6000602084013e611389565b606091505b509098509650876113de5760408051600160e51b62461bcd0281526020600482015260146024820152600160621b7318d85b1b081d1c985b9cd9995c8819985a5b195902604482015290519081900360640190fd5b8680602001905160208110156113f357600080fd5b505195508561143657604051600160e51b62461bcd028152600401808060200182810382526021815260200180611ecd6021913960400191505060405180910390fd5b600101611126565b50600460008e815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316600460008f815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168e7f0e5141db06f1fbaf3d754ea212eda4b7ad2bd1019d83809a2fd22bcb88e0c7568f88876040518084815260200183815260200180602001828103825283818151815260200191508051906020019060200280838360005b838110156115165781810151838201526020016114fe565b5050505090500194505050505060405180910390a4875b60018b03811015611574576005816001018154811061154857fe5b90600052602060002001546005828154811061156057fe5b60009182526020909120015560010161152d565b50600560018b038154811061158557fe5b600091825260208220015560058054906115a3906000198301611e41565b5060008d815260046020526040812080546001600160a01b03199081168255600182018054909116905560028101829055906115e26003830182611e6a565b50600401805461ffff191690555060019c9b505050505050505050505050565b60608060005b6002548110156116c057836001600160a01b03166002828154811061162957fe5b6000918252602090912001546001600160a01b031614156116b8576001600160a01b038416600090815260016020908152604091829020805483518184028101840190945280845290918301828280156116ac57602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161168e575b505050505091506116c0565b600101611608565b5092915050565b60006116d282611735565b92915050565b60035460009060ff161561172057604051600160e51b62461bcd02815260040180806020018281038252602f815260200180611eee602f913960400191505060405180910390fd5b6116d282611d5d565b670de0b6b3a764000081565b600080805b60025481101561178357336001600160a01b03166002828154811061175b57fe5b6000918252602090912001546001600160a01b0316141561177b57600191505b60010161173a565b508015611839573360009081526001602052604081205490805b828110156117f45733600090815260016020526040902080546001600160a01b0388169190839081106117cc57fe5b6000918252602090912001546001600160a01b031614156117ec57600191505b60010161179d565b50806118325733600090815260016020818152604083208054928301815583529091200180546001600160a01b0319166001600160a01b0387161790555b50506118ab565b6002805460018082019092557f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054336001600160a01b03199182168117909255600091825260208381526040832080549485018155835290912090910180549091166001600160a01b0385161790555b50600192915050565b60025460009081908190815b8181101561190a57336001600160a01b0316600282815481106118df57fe5b6000918252602090912001546001600160a01b0316141561190257600193508092505b6001016118c0565b50826119605760408051600160e51b62461bcd02815260206004820152600c60248201527f696e76616c696420757365720000000000000000000000000000000000000000604482015290519081900360640190fd5b336000908152600160205260408120548190815b818110156119ce5733600090815260016020526040902080546001600160a01b038b169190839081106119a357fe5b6000918252602090912001546001600160a01b031614156119c657600193508092505b600101611974565b5082611a245760408051600160e51b62461bcd02815260206004820152600e60248201527f696e76616c6964206167656e6379000000000000000000000000000000000000604482015290519081900360640190fd5b6001811115611b1c57815b60018203811015611abf57336000908152600160208190526040909120805490918301908110611a5b57fe5b60009182526020808320909101543383526001909152604090912080546001600160a01b039092169183908110611a8e57fe5b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055600101611a2f565b5033600090815260016020526040902080546000198301908110611adf57fe5b6000918252602080832090910180546001600160a01b03191690553382526001905260409020805490611b16906000198301611e41565b50611be6565b845b60018503811015611b905760028160010181548110611b3957fe5b600091825260209091200154600280546001600160a01b039092169183908110611b5f57fe5b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055600101611b1e565b5060026001850381548110611ba157fe5b600091825260209091200180546001600160a01b03191690556002805490611bcd906000198301611e41565b50336000908152600160205260408120611be691611e6a565b506001979650505050505050565b604080516001600160a01b03848116602483015283811660448084019190915283518084039091018152606490920183526020820180516001600160e01b0316600160e11b636eb1769f021781529251825160009485946060948a16939092909182918083835b60208310611c7a5780518252601f199092019160209182019101611c5b565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d8060008114611cda576040519150601f19603f3d011682016040523d82523d6000602084013e611cdf565b606091505b509150915081611d395760408051600160e51b62461bcd02815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c65640000000000604482015290519081900360640190fd5b6000818060200190516020811015611d5057600080fd5b5051979650505050505050565b60006001600160a01b038216611da757604051600160e51b62461bcd028152600401808060200182810382526022815260200180611f426022913960400191505060405180910390fd5b50600080546001600160a01b0383166001600160a01b03199091161790556003805460ff19166001179081905560ff16919050565b828054828255906000526020600020908101928215611e31579160200282015b82811115611e3157825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611dfc565b50611e3d929150611e8b565b5090565b815481835581811115611e6557600083815260209020611e65918101908301611eb2565b505050565b5080546000825590600052602060002090810190611e889190611eb2565b50565b611eaf91905b80821115611e3d5780546001600160a01b0319168155600101611e91565b90565b611eaf91905b80821115611e3d5760008155600101611eb856fe5468652072657475726e206f66207472616e73666572206973206661696c7572654d657469735061793a204d6574697350617920696e7374616e636520616c726561647920696e697469616c697a65645468652072657475726e206f66207472616e7366657266726f6d206973206661696c7572654d657469735061793a20496e76616c6964206d657469734c61742061646472657373a165627a7a72305820f67bd7671bc0b1dd3371a6c6835dd1f1e3f875edc43906c6a86c2cc8f7cf132d00294d657469735061793a20496e76616c6964206d657469734c61742061646472657373"

// DeployMetisPay deploys a new Ethereum contract, binding an instance of MetisPay to it.
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

// MetisPay is an auto generated Go binding around an Ethereum contract.
type MetisPay struct {
	MetisPayCaller     // Read-only binding to the contract
	MetisPayTransactor // Write-only binding to the contract
	MetisPayFilterer   // Log filterer for contract events
}

// MetisPayCaller is an auto generated read-only Go binding around an Ethereum contract.
type MetisPayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetisPayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MetisPayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetisPayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MetisPayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MetisPaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MetisPaySession struct {
	Contract     *MetisPay         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MetisPayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MetisPayCallerSession struct {
	Contract *MetisPayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// MetisPayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MetisPayTransactorSession struct {
	Contract     *MetisPayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// MetisPayRaw is an auto generated low-level Go binding around an Ethereum contract.
type MetisPayRaw struct {
	Contract *MetisPay // Generic contract binding to access the raw methods on
}

// MetisPayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MetisPayCallerRaw struct {
	Contract *MetisPayCaller // Generic read-only contract binding to access the raw methods on
}

// MetisPayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
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

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns(bool)
func (_MetisPay *MetisPayTransactor) Initialize(opts *bind.TransactOpts, metisLat common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "initialize", metisLat)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns(bool)
func (_MetisPay *MetisPaySession) Initialize(metisLat common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Initialize(&_MetisPay.TransactOpts, metisLat)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns(bool)
func (_MetisPay *MetisPayTransactorSession) Initialize(metisLat common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Initialize(&_MetisPay.TransactOpts, metisLat)
}

// Prepay is a paid mutator transaction binding the contract method 0x5d647525.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList) returns(bool success)
func (_MetisPay *MetisPayTransactor) Prepay(opts *bind.TransactOpts, taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "prepay", taskId, user, fee, tokenAddressList)
}

// Prepay is a paid mutator transaction binding the contract method 0x5d647525.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList) returns(bool success)
func (_MetisPay *MetisPaySession) Prepay(taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Prepay(&_MetisPay.TransactOpts, taskId, user, fee, tokenAddressList)
}

// Prepay is a paid mutator transaction binding the contract method 0x5d647525.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) Prepay(taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Prepay(&_MetisPay.TransactOpts, taskId, user, fee, tokenAddressList)
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

// MetisPayPrepayEventIterator is returned from FilterPrepayEvent and is used to iterate over the raw logs and unpacked data for PrepayEvent events raised by the MetisPay contract.
type MetisPayPrepayEventIterator struct {
	Event *MetisPayPrepayEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MetisPayPrepayEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetisPayPrepayEvent)
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
		it.Event = new(MetisPayPrepayEvent)
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
func (it *MetisPayPrepayEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetisPayPrepayEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetisPayPrepayEvent represents a PrepayEvent event raised by the MetisPay contract.
type MetisPayPrepayEvent struct {
	TaskId           *big.Int
	User             common.Address
	UserAgency       common.Address
	Fee              *big.Int
	TokenAddressList []common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPrepayEvent is a free log retrieval operation binding the contract event 0x2f05b3233a63a604771d18ab3e4bf65f10a488ea624dd4f1499d274b1bd7da59.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) FilterPrepayEvent(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*MetisPayPrepayEventIterator, error) {

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

	logs, sub, err := _MetisPay.contract.FilterLogs(opts, "PrepayEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &MetisPayPrepayEventIterator{contract: _MetisPay.contract, event: "PrepayEvent", logs: logs, sub: sub}, nil
}

// WatchPrepayEvent is a free log subscription operation binding the contract event 0x2f05b3233a63a604771d18ab3e4bf65f10a488ea624dd4f1499d274b1bd7da59.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) WatchPrepayEvent(opts *bind.WatchOpts, sink chan<- *MetisPayPrepayEvent, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _MetisPay.contract.WatchLogs(opts, "PrepayEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetisPayPrepayEvent)
				if err := _MetisPay.contract.UnpackLog(event, "PrepayEvent", log); err != nil {
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

// ParsePrepayEvent is a log parse operation binding the contract event 0x2f05b3233a63a604771d18ab3e4bf65f10a488ea624dd4f1499d274b1bd7da59.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) ParsePrepayEvent(log types.Log) (*MetisPayPrepayEvent, error) {
	event := new(MetisPayPrepayEvent)
	if err := _MetisPay.contract.UnpackLog(event, "PrepayEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MetisPaySettleEventIterator is returned from FilterSettleEvent and is used to iterate over the raw logs and unpacked data for SettleEvent events raised by the MetisPay contract.
type MetisPaySettleEventIterator struct {
	Event *MetisPaySettleEvent // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *MetisPaySettleEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetisPaySettleEvent)
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
		it.Event = new(MetisPaySettleEvent)
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
func (it *MetisPaySettleEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetisPaySettleEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetisPaySettleEvent represents a SettleEvent event raised by the MetisPay contract.
type MetisPaySettleEvent struct {
	TaskId           *big.Int
	User             common.Address
	UserAgency       common.Address
	Agencyfee        *big.Int
	Refund           *big.Int
	TokenAddressList []common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSettleEvent is a free log retrieval operation binding the contract event 0x0e5141db06f1fbaf3d754ea212eda4b7ad2bd1019d83809a2fd22bcb88e0c756.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) FilterSettleEvent(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*MetisPaySettleEventIterator, error) {

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

	logs, sub, err := _MetisPay.contract.FilterLogs(opts, "SettleEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &MetisPaySettleEventIterator{contract: _MetisPay.contract, event: "SettleEvent", logs: logs, sub: sub}, nil
}

// WatchSettleEvent is a free log subscription operation binding the contract event 0x0e5141db06f1fbaf3d754ea212eda4b7ad2bd1019d83809a2fd22bcb88e0c756.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) WatchSettleEvent(opts *bind.WatchOpts, sink chan<- *MetisPaySettleEvent, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _MetisPay.contract.WatchLogs(opts, "SettleEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetisPaySettleEvent)
				if err := _MetisPay.contract.UnpackLog(event, "SettleEvent", log); err != nil {
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

// ParseSettleEvent is a log parse operation binding the contract event 0x0e5141db06f1fbaf3d754ea212eda4b7ad2bd1019d83809a2fd22bcb88e0c756.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList)
func (_MetisPay *MetisPayFilterer) ParseSettleEvent(log types.Log) (*MetisPaySettleEvent, error) {
	event := new(MetisPaySettleEvent)
	if err := _MetisPay.contract.UnpackLog(event, "SettleEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
