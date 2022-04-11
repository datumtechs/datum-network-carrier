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
const MetisPayABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"fee\",\"type\":\"uint256\"},{\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"prepay\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"taskId\",\"type\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"metisLat\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"USAGE_FEE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"BASE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"PrepayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"refund\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"SettleEvent\",\"type\":\"event\"}]"

// MetisPayFuncSigs maps the 4-byte function signature to its string representation.
var MetisPayFuncSigs = map[string]string{
	"ec342ad0": "BASE()",
	"d5b96c53": "USAGE_FEE()",
	"f80f5dd5": "addWhitelist(address)",
	"b6a5d7de": "authorize(address)",
	"ff68f6af": "deleteWhitelist(address)",
	"c4d66de8": "initialize(address)",
	"0ef92166": "prepay(uint256,address,uint256,address[],uint256[])",
	"9a9c29f6": "settle(uint256,uint256)",
	"9b19251a": "whitelist(address)",
}

// MetisPayBin is the compiled bytecode used for deploying new contracts.
var MetisPayBin = "0x608060405234801561001057600080fd5b50611702806100206000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063c4d66de811610066578063c4d66de8146102ae578063d5b96c53146102d4578063ec342ad0146102d4578063f80f5dd5146102ee578063ff68f6af1461031457610093565b80630ef92166146100985780639a9c29f6146101ef5780639b19251a14610212578063b6a5d7de14610288575b600080fd5b6101db600480360360a08110156100ae57600080fd5b8135916001600160a01b0360208201351691604082013591908101906080810160608201356401000000008111156100e557600080fd5b8201836020820111156100f757600080fd5b8035906020019184602083028401116401000000008311171561011957600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250929594936020810193503591505064010000000081111561016957600080fd5b82018360208201111561017b57600080fd5b8035906020019184602083028401116401000000008311171561019d57600080fd5b91908080602002602001604051908101604052809392919081815260200183836020028082843760009201919091525092955061033a945050505050565b604080519115158252519081900360200190f35b6101db6004803603604081101561020557600080fd5b5080359060200135610c85565b6102386004803603602081101561022857600080fd5b50356001600160a01b0316610d4f565b60408051602080825283518183015283519192839290830191858101910280838360005b8381101561027457818101518382015260200161025c565b505050509050019250505060405180910390f35b6101db6004803603602081101561029e57600080fd5b50356001600160a01b0316610e14565b6101db600480360360208110156102c457600080fd5b50356001600160a01b0316610e25565b6102dc610e76565b60408051918252519081900360200190f35b6101db6004803603602081101561030457600080fd5b50356001600160a01b0316610e82565b6101db6004803603602081101561032a57600080fd5b50356001600160a01b0316611001565b600081518351146103955760408051600160e51b62461bcd02815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f7420657869737400604482015290519081900360640190fd5b600254600090815b818110156103e457876001600160a01b0316600282815481106103bc57fe5b6000918252602090912001546001600160a01b031614156103dc57600192505b60010161039d565b508161043a5760408051600160e51b62461bcd02815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f7420657869737400604482015290519081900360640190fd5b50506001600160a01b038516600090815260016020526040812054815b818110156104ae576001600160a01b038816600090815260016020526040902080543391908390811061048657fe5b6000918252602090912001546001600160a01b031614156104a657600192505b600101610457565b50816105045760408051600160e51b62461bcd02815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a65640000000000000000604482015290519081900360640190fd5b5050600554600090815b8181101561054257886005828154811061052457fe5b9060005260206000200154141561053a57600192505b60010161050e565b5081156105995760408051600160e51b62461bcd02815260206004820152601660248201527f7461736b20696420616c72656164792065786973747300000000000000000000604482015290519081900360640190fd5b600080546105b1906001600160a01b03168930611341565b905086811161060a5760408051600160e51b62461bcd02815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e636500000000000000604482015290519081900360640190fd5b8551915060005b828110156106ab5761063787828151811061062857fe5b60200260200101518a30611341565b915085818151811061064557fe5b602002602001015182116106a35760408051600160e51b62461bcd02815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e636500604482015290519081900360640190fd5b600101610611565b5060008054604080516001600160a01b038c8116602483015230604483015260648083018d905283518084039091018152608490920183526020820180516001600160e01b0316600160e01b63b642fe57021781529251825160609587959316939282918083835b602083106107325780518252601f199092019160209182019101610713565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610794576040519150601f19603f3d011682016040523d82523d6000602084013e610799565b606091505b509093509150826107f45760408051600160e51b62461bcd02815260206004820152601860248201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604482015290519081900360640190fd5b81806020019051602081101561080957600080fd5b505190508061084c57604051600160e51b62461bcd0281526004018080602001828103825260258152602001806116906025913960400191505060405180910390fd5b60005b85811015610a305789818151811061086357fe5b60200260200101516001600160a01b03168c308b848151811061088257fe5b602090810291909101810151604080516001600160a01b0395861660248201529390941660448401526064808401919091528351808403909101815260849092018352810180516001600160e01b0316600160e01b63b642fe570217815291518151919290918291908083835b6020831061090e5780518252601f1990920191602091820191016108ef565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114610970576040519150601f19603f3d011682016040523d82523d6000602084013e610975565b606091505b509094509250836109d05760408051600160e51b62461bcd02815260206004820152601860248201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604482015290519081900360640190fd5b8280602001905160208110156109e557600080fd5b5051915081610a2857604051600160e51b62461bcd0281526004018080602001828103825260258152602001806116906025913960400191505060405180910390fd5b60010161084f565b506040518060e001604052808c6001600160a01b03168152602001336001600160a01b031681526020018b81526020018a815260200189815260200160011515815260200160001515815250600460008e815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550604082015181600201556060820151816003019080519060200190610b0f929190611529565b5060808201518051610b2b91600484019160209091019061158e565b5060a08201518160050160006101000a81548160ff02191690831515021790555060c08201518160050160016101000a81548160ff02191690831515021790555090505060058c9080600181540180825580915050906001820390600052602060002001600090919290919091505550336001600160a01b03168b6001600160a01b03168d7f3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f8d8d8d604051808481526020018060200180602001838103835285818151815260200191508051906020019060200280838360005b83811015610c1e578181015183820152602001610c06565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610c5d578181015183820152602001610c45565b505050509050019550505050505060405180910390a45060019b9a5050505050505050505050565b6000828152600460205260408120600201548210610ced5760408051600160e51b62461bcd02815260206004820152601560248201527f6d6f7265207468616e2070726570616964206c61740000000000000000000000604482015290519081900360640190fd5b600083815260046020526040812080546001600160a01b0319908116825560018201805490911690556002810182905590610d2b60038301826115d5565b610d396004830160006115d5565b50600501805461ffff1916905550600192915050565b60608060005b600254811015610e0d57836001600160a01b031660028281548110610d7657fe5b6000918252602090912001546001600160a01b03161415610e05576001600160a01b03841660009081526001602090815260409182902080548351818402810184019094528084529091830182828015610df957602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610ddb575b50505050509150610e0d565b600101610d55565b5092915050565b6000610e1f82610e82565b92915050565b60035460009060ff1615610e6d57604051600160e51b62461bcd02815260040180806020018281038252602f815260200180611661602f913960400191505060405180910390fd5b610e1f826114aa565b670de0b6b3a764000081565b600080805b600254811015610ed057336001600160a01b031660028281548110610ea857fe5b6000918252602090912001546001600160a01b03161415610ec857600191505b600101610e87565b508015610f86573360009081526001602052604081205490805b82811015610f415733600090815260016020526040902080546001600160a01b038816919083908110610f1957fe5b6000918252602090912001546001600160a01b03161415610f3957600191505b600101610eea565b5080610f7f5733600090815260016020818152604083208054928301815583529091200180546001600160a01b0319166001600160a01b0387161790555b5050610ff8565b6002805460018082019092557f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace018054336001600160a01b03199182168117909255600091825260208381526040832080549485018155835290912090910180549091166001600160a01b0385161790555b50600192915050565b60025460009081908190815b8181101561105757336001600160a01b03166002828154811061102c57fe5b6000918252602090912001546001600160a01b0316141561104f57600193508092505b60010161100d565b50826110ad5760408051600160e51b62461bcd02815260206004820152600c60248201527f696e76616c696420757365720000000000000000000000000000000000000000604482015290519081900360640190fd5b336000908152600160205260408120548190815b8181101561111b5733600090815260016020526040902080546001600160a01b038b169190839081106110f057fe5b6000918252602090912001546001600160a01b0316141561111357600193508092505b6001016110c1565b50826111715760408051600160e51b62461bcd02815260206004820152600e60248201527f696e76616c6964206167656e6379000000000000000000000000000000000000604482015290519081900360640190fd5b600181111561126957815b6001820381101561120c573360009081526001602081905260409091208054909183019081106111a857fe5b60009182526020808320909101543383526001909152604090912080546001600160a01b0390921691839081106111db57fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560010161117c565b503360009081526001602052604090208054600019830190811061122c57fe5b6000918252602080832090910180546001600160a01b031916905533825260019052604090208054906112639060001983016115f6565b50611333565b845b600185038110156112dd576002816001018154811061128657fe5b600091825260209091200154600280546001600160a01b0390921691839081106112ac57fe5b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905560010161126b565b50600260018503815481106112ee57fe5b600091825260209091200180546001600160a01b0319169055600280549061131a9060001983016115f6565b50336000908152600160205260408120611333916115d5565b506001979650505050505050565b604080516001600160a01b03848116602483015283811660448084019190915283518084039091018152606490920183526020820180516001600160e01b0316600160e11b636eb1769f021781529251825160009485946060948a16939092909182918083835b602083106113c75780518252601f1990920191602091820191016113a8565b6001836020036101000a038019825116818451168082178552505050505050905001915050600060405180830381855afa9150503d8060008114611427576040519150601f19603f3d011682016040523d82523d6000602084013e61142c565b606091505b5091509150816114865760408051600160e51b62461bcd02815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c65640000000000604482015290519081900360640190fd5b600081806020019051602081101561149d57600080fd5b5051979650505050505050565b60006001600160a01b0382166114f457604051600160e51b62461bcd0281526004018080602001828103825260228152602001806116b56022913960400191505060405180910390fd5b50600080546001600160a01b0383166001600160a01b03199091161790556003805460ff19166001179081905560ff16919050565b82805482825590600052602060002090810192821561157e579160200282015b8281111561157e57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611549565b5061158a92915061161f565b5090565b8280548282559060005260206000209081019282156115c9579160200282015b828111156115c95782518255916020019190600101906115ae565b5061158a929150611646565b50805460008255906000526020600020908101906115f39190611646565b50565b81548183558181111561161a5760008381526020902061161a918101908301611646565b505050565b61164391905b8082111561158a5780546001600160a01b0319168155600101611625565b90565b61164391905b8082111561158a576000815560010161164c56fe4d657469735061793a204d6574697350617920696e7374616e636520616c726561647920696e697469616c697a65645468652072657475726e206f66207472616e7366657266726f6d206973206661696c7572654d657469735061793a20496e76616c6964206d657469734c61742061646472657373a165627a7a7230582055bfc42203f0cc498f8c06b21f28171f8a5c14422b493dfe96f1a77d019e53e50029"

// DeployMetisPay deploys a new Ethereum contract, binding an instance of MetisPay to it.
func DeployMetisPay(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MetisPay, error) {
	parsed, err := abi.JSON(strings.NewReader(MetisPayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MetisPayBin), backend)
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

// Prepay is a paid mutator transaction binding the contract method 0x0ef92166.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList) returns(bool success)
func (_MetisPay *MetisPayTransactor) Prepay(opts *bind.TransactOpts, taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address, tokenValueList []*big.Int) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "prepay", taskId, user, fee, tokenAddressList, tokenValueList)
}

// Prepay is a paid mutator transaction binding the contract method 0x0ef92166.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList) returns(bool success)
func (_MetisPay *MetisPaySession) Prepay(taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address, tokenValueList []*big.Int) (*types.Transaction, error) {
	return _MetisPay.Contract.Prepay(&_MetisPay.TransactOpts, taskId, user, fee, tokenAddressList, tokenValueList)
}

// Prepay is a paid mutator transaction binding the contract method 0x0ef92166.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList) returns(bool success)
func (_MetisPay *MetisPayTransactorSession) Prepay(taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address, tokenValueList []*big.Int) (*types.Transaction, error) {
	return _MetisPay.Contract.Prepay(&_MetisPay.TransactOpts, taskId, user, fee, tokenAddressList, tokenValueList)
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
	TokenValueList   []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPrepayEvent is a free log retrieval operation binding the contract event 0x3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList)
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

// WatchPrepayEvent is a free log subscription operation binding the contract event 0x3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList)
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

// ParsePrepayEvent is a log parse operation binding the contract event 0x3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList)
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
	TokenValueList   []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSettleEvent is a free log retrieval operation binding the contract event 0x44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef4.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList, uint256[] tokenValueList)
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

// WatchSettleEvent is a free log subscription operation binding the contract event 0x44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef4.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList, uint256[] tokenValueList)
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

// ParseSettleEvent is a log parse operation binding the contract event 0x44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef4.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refund, address[] tokenAddressList, uint256[] tokenValueList)
func (_MetisPay *MetisPayFilterer) ParseSettleEvent(log types.Log) (*MetisPaySettleEvent, error) {
	event := new(MetisPaySettleEvent)
	if err := _MetisPay.contract.UnpackLog(event, "SettleEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
