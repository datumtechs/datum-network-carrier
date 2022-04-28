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
const MetisPayABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"PrepayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"SettleEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"metisLat\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"prepay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"taskState\",\"outputs\":[{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MetisPayBin is the compiled bytecode used for deploying new contracts.
var MetisPayBin = "0x608060405234801561001057600080fd5b5061225c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80639b19251a116100715780639b19251a14610134578063b6a5d7de14610154578063c4d66de814610167578063f2fde38b1461017a578063f80f5dd51461018d578063ff68f6af146101a057600080fd5b80630ef92166146100ae57806326c6bee1146100d6578063715018a6146100fc5780638da5cb5b146101065780639a9c29f614610121575b600080fd5b6100c16100bc366004611e1f565b6101b3565b60405190151581526020015b60405180910390f35b6100e96100e4366004611f09565b610969565b60405160009190910b81526020016100cd565b6101046109e7565b005b6033546040516001600160a01b0390911681526020016100cd565b6100c161012f366004611f22565b610a4d565b610147610142366004611f44565b611315565b6040516100cd9190611fac565b6100c1610162366004611f44565b6113e9565b610104610175366004611f44565b6113fa565b610104610188366004611f44565b6114cf565b6100c161019b366004611f44565b61159a565b6100c16101ae366004611f44565b61173a565b6000815183511461020b5760405162461bcd60e51b815260206004820152601960248201527f696e76616c696420746f6b656e20696e666f726d6174696f6e0000000000000060448201526064015b60405180910390fd5b606754600090815b8181101561026957876001600160a01b03166067828154811061023857610238611fbf565b6000918252602090912001546001600160a01b03160361025757600192505b8061026181611feb565b915050610213565b50816102b75760405162461bcd60e51b815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f74206578697374006044820152606401610202565b50506001600160a01b038516600090815260666020526040812054815b8181101561033a576001600160a01b038816600090815260666020526040902080543391908390811061030957610309611fbf565b6000918252602090912001546001600160a01b03160361032857600192505b8061033281611feb565b9150506102d4565b50816103885760405162461bcd60e51b815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a656400000000000000006044820152606401610202565b5050606954600090815b818110156103d55788606982815481106103ae576103ae611fbf565b9060005260206000200154036103c357600192505b806103cd81611feb565b915050610392565b50811561041d5760405162461bcd60e51b81526020600482015260166024820152757461736b20696420616c72656164792065786973747360501b6044820152606401610202565b606554600090610437906001600160a01b03168930611b02565b9050868110156104895760405162461bcd60e51b815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e6365000000000000006044820152606401610202565b8551915060005b82811015610520576104bc8782815181106104ad576104ad611fbf565b60200260200101518a30611b02565b91508782101561050e5760405162461bcd60e51b815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e6365006044820152606401610202565b8061051881611feb565b915050610490565b506065546040516001600160a01b038a81166024830152306044830152606482018a905260009260609284929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516105889190612004565b6000604051808303816000865af19150503d80600081146105c5576040519150601f19603f3d011682016040523d82523d6000602084013e6105ca565b606091505b509093509150826106185760405162461bcd60e51b815260206004820152601860248201527718d85b1b081d1c985b9cd9995c919c9bdb4819985a5b195960421b6044820152606401610202565b8180602001905181019061062c919061203f565b90508061064b5760405162461bcd60e51b815260040161020290612061565b60005b858110156107c65789818151811061066857610668611fbf565b60200260200101516001600160a01b03168c308b848151811061068d5761068d611fbf565b60209081029190910101516040516001600160a01b039384166024820152929091166044830152606482015260840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516106f19190612004565b6000604051808303816000865af19150503d806000811461072e576040519150601f19603f3d011682016040523d82523d6000602084013e610733565b606091505b509094509250836107815760405162461bcd60e51b815260206004820152601860248201527718d85b1b081d1c985b9cd9995c919c9bdb4819985a5b195960421b6044820152606401610202565b82806020019051810190610795919061203f565b9150816107b45760405162461bcd60e51b815260040161020290612061565b806107be81611feb565b91505061064e565b506040518060c001604052808c6001600160a01b03168152602001336001600160a01b031681526020018b81526020018a8152602001898152602001600160000b815250606860008e815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160020155606082015181600301908051906020019061089d929190611c65565b50608082015180516108b9916004840191602090910190611cca565b5060a091909101516005909101805460ff191660ff909216919091179055606980546001810182556000919091527f7fb4302e8e91f9110a6554c2c0a24601252c2a42c2220ca988efcfe399914308018c905560405133906001600160a01b038d16908e907f3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f9061094f908f908f908f906120d6565b60405180910390a45060019b9a5050505050505050505050565b606954600090819081805b828110156109b757856069828154811061099057610990611fbf565b9060005260206000200154036109a557600191505b806109af81611feb565b915050610974565b50806109c75760001992506109de565b600085815260686020526040812060050154900b92505b50909392505050565b6033546001600160a01b03163314610a415760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610202565b610a4b6000611c13565b565b6069546000908180805b83811015610a9d578660698281548110610a7357610a73611fbf565b906000526020600020015403610a8b57600192508091505b80610a9581611feb565b915050610a57565b5081610add5760405162461bcd60e51b815260206004820152600f60248201526e1a5b9d985b1a59081d185cdac81a59608a1b6044820152606401610202565b6000868152606860205260409020600101546001600160a01b03163314610b465760405162461bcd60e51b815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f207468697300000000006044820152606401610202565b600086815260686020526040812060050154900b600114610bb75760405162461bcd60e51b815260206004820152602560248201527f707265706179206e6f7420636f6d706c65746564206f722072657065617420736044820152646574746c6560d81b6064820152608401610202565b600086815260686020526040902060020154851115610c105760405162461bcd60e51b81526020600482015260156024820152741b5bdc99481d1a185b881c1c995c185a59081b185d605a1b6044820152606401610202565b6065546000878152606860205260408082206001015490516001600160a01b03918216602482015260448101899052919260609291169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610c7f9190612004565b6000604051808303816000865af19150503d8060008114610cbc576040519150601f19603f3d011682016040523d82523d6000602084013e610cc1565b606091505b50909250905081610ce45760405162461bcd60e51b81526004016102029061210b565b80806020019051810190610cf8919061203f565b610d145760405162461bcd60e51b815260040161020290612139565b600088815260686020526040812060020154610d3190899061217a565b90508015610e365760655460008a815260686020526040908190205490516001600160a01b0391821660248201526044810184905291169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610da19190612004565b6000604051808303816000865af19150503d8060008114610dde576040519150601f19603f3d011682016040523d82523d6000602084013e610de3565b606091505b50909350915082610e065760405162461bcd60e51b81526004016102029061210b565b81806020019051810190610e1a919061203f565b610e365760405162461bcd60e51b815260040161020290612139565b600089815260686020908152604080832060030180548251818502810185019093528083528493830182828015610e9657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610e78575b505050505090506000606860008d8152602001908152602001600020600401805480602002602001604051908101604052809291908181526020018280548015610eff57602002820191906000526020600020905b815481526020019060010190808311610eeb575b505085519394506000925050505b8181101561114657838181518110610f2757610f27611fbf565b602090810291909101810151604080516004815260248101825292830180516001600160e01b03166303aa30b960e11b179052516001600160a01b0390911691610f7091612004565b600060405180830381855afa9150503d8060008114610fab576040519150601f19603f3d011682016040523d82523d6000602084013e610fb0565b606091505b50909850965087610ff85760405162461bcd60e51b815260206004820152601260248201527118d85b1b081b5a5b9d195c8819985a5b195960721b6044820152606401610202565b8680602001905181019061100c9190612191565b945083818151811061102057611020611fbf565b60200260200101516001600160a01b03168584838151811061104457611044611fbf565b60209081029190910101516040516001600160a01b039092166024830152604482015260640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b1790525161109f9190612004565b6000604051808303816000865af19150503d80600081146110dc576040519150601f19603f3d011682016040523d82523d6000602084013e6110e1565b606091505b509098509650876111045760405162461bcd60e51b81526004016102029061210b565b86806020019051810190611118919061203f565b6111345760405162461bcd60e51b815260040161020290612139565b8061113e81611feb565b915050610f0d565b50606860008e815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316606860008f815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168e7f44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef48f8988886040516111e394939291906121ae565b60405180910390a4875b6111f860018c61217a565b81101561125857606961120c8260016121df565b8154811061121c5761121c611fbf565b90600052602060002001546069828154811061123a5761123a611fbf565b6000918252602090912001558061125081611feb565b9150506111ed565b50606961126660018c61217a565b8154811061127657611276611fbf565b60009182526020822001556069805480611292576112926121f7565b6000828152602080822083016000199081018390559092019092558e8252606890526040812080546001600160a01b03199081168255600182018054909116905560028101829055906112e86003830182611d05565b6112f6600483016000611d05565b50600501805460ff191690555060019c9b505050505050505050505050565b60608060005b6067548110156113e257836001600160a01b03166067828154811061134257611342611fbf565b6000918252602090912001546001600160a01b0316036113d0576001600160a01b038416600090815260666020908152604091829020805483518184028101840190945280845290918301828280156113c457602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116113a6575b505050505091506113e2565b806113da81611feb565b91505061131b565b5092915050565b60006113f48261159a565b92915050565b600054610100900460ff166114155760005460ff1615611419565b303b155b61147c5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610202565b600054610100900460ff1615801561149e576000805461ffff19166101011790555b606580546001600160a01b0319166001600160a01b03841617905580156114cb576000805461ff00191690555b5050565b6033546001600160a01b031633146115295760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610202565b6001600160a01b03811661158e5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610202565b61159781611c13565b50565b600080805b6067548110156115f757336001600160a01b0316606782815481106115c6576115c6611fbf565b6000918252602090912001546001600160a01b0316036115e557600191505b806115ef81611feb565b91505061159f565b5080156116bd573360009081526066602052604081205490805b828110156116775733600090815260666020526040902080546001600160a01b03881691908390811061164657611646611fbf565b6000918252602090912001546001600160a01b03160361166557600191505b8061166f81611feb565b915050611611565b50806116b6573360009081526066602090815260408220805460018101825590835291200180546001600160a01b0319166001600160a01b0387161790555b5050611731565b6067805460018082019092557f9787eeb91fe3101235e4a76063c7023ecb40f923f97916639c598592fa30d6ae018054336001600160a01b031991821681179092556000918252606660209081526040832080549485018155835290912090910180549091166001600160a01b0385161790555b50600192915050565b60675460009081908190815b8181101561179f57336001600160a01b03166067828154811061176b5761176b611fbf565b6000918252602090912001546001600160a01b03160361178d57600193508092505b8061179781611feb565b915050611746565b50826117dc5760405162461bcd60e51b815260206004820152600c60248201526b34b73b30b634b2103ab9b2b960a11b6044820152606401610202565b336000908152606660205260408120548190815b818110156118595733600090815260666020526040902080546001600160a01b038b1691908390811061182557611825611fbf565b6000918252602090912001546001600160a01b03160361184757600193508092505b8061185181611feb565b9150506117f0565b50826118985760405162461bcd60e51b815260206004820152600e60248201526d696e76616c6964206167656e637960901b6044820152606401610202565b60018111156119da57815b6118ae60018361217a565b811015611953573360009081526066602052604090206118cf8260016121df565b815481106118df576118df611fbf565b60009182526020808320909101543383526066909152604090912080546001600160a01b03909216918390811061191857611918611fbf565b600091825260209091200180546001600160a01b0319166001600160a01b03929092169190911790558061194b81611feb565b9150506118a3565b5033600090815260666020526040902061196e60018361217a565b8154811061197e5761197e611fbf565b6000918252602080832090910180546001600160a01b031916905533825260669052604090208054806119b3576119b36121f7565b600082815260209020810160001990810180546001600160a01b0319169055019055611af4565b845b6119e760018661217a565b811015611a725760676119fb8260016121df565b81548110611a0b57611a0b611fbf565b600091825260209091200154606780546001600160a01b039092169183908110611a3757611a37611fbf565b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905580611a6a81611feb565b9150506119dc565b506067611a8060018661217a565b81548110611a9057611a90611fbf565b600091825260209091200180546001600160a01b03191690556067805480611aba57611aba6121f7565b60008281526020808220830160001990810180546001600160a01b03191690559092019092553382526066905260408120611af491611d05565b506001979650505050505050565b6040516001600160a01b0383811660248301528281166044830152600091829182919087169060640160408051601f198184030181529181526020820180516001600160e01b0316636eb1769f60e11b17905251611b609190612004565b600060405180830381855afa9150503d8060008114611b9b576040519150601f19603f3d011682016040523d82523d6000602084013e611ba0565b606091505b509150915081611bf25760405162461bcd60e51b815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c656400000000006044820152606401610202565b600081806020019051810190611c08919061220d565b979650505050505050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b828054828255906000526020600020908101928215611cba579160200282015b82811115611cba57825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611c85565b50611cc6929150611d1f565b5090565b828054828255906000526020600020908101928215611cba579160200282015b82811115611cba578251825591602001919060010190611cea565b508054600082559060005260206000209081019061159791905b5b80821115611cc65760008155600101611d20565b6001600160a01b038116811461159757600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611d8857611d88611d49565b604052919050565b600067ffffffffffffffff821115611daa57611daa611d49565b5060051b60200190565b600082601f830112611dc557600080fd5b81356020611dda611dd583611d90565b611d5f565b82815260059290921b84018101918181019086841115611df957600080fd5b8286015b84811015611e145780358352918301918301611dfd565b509695505050505050565b600080600080600060a08688031215611e3757600080fd5b85359450602080870135611e4a81611d34565b945060408701359350606087013567ffffffffffffffff80821115611e6e57600080fd5b818901915089601f830112611e8257600080fd5b8135611e90611dd582611d90565b81815260059190911b8301840190848101908c831115611eaf57600080fd5b938501935b82851015611ed6578435611ec781611d34565b82529385019390850190611eb4565b965050506080890135925080831115611eee57600080fd5b5050611efc88828901611db4565b9150509295509295909350565b600060208284031215611f1b57600080fd5b5035919050565b60008060408385031215611f3557600080fd5b50508035926020909101359150565b600060208284031215611f5657600080fd5b8135611f6181611d34565b9392505050565b600081518084526020808501945080840160005b83811015611fa15781516001600160a01b031687529582019590820190600101611f7c565b509495945050505050565b602081526000611f616020830184611f68565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611ffd57611ffd611fd5565b5060010190565b6000825160005b81811015612025576020818601810151858301520161200b565b81811115612034576000828501525b509190910192915050565b60006020828403121561205157600080fd5b81518015158114611f6157600080fd5b60208082526025908201527f5468652072657475726e206f66207472616e7366657266726f6d206973206661604082015264696c75726560d81b606082015260800190565b600081518084526020808501945080840160005b83811015611fa1578151875295820195908201906001016120ba565b8381526060602082015260006120ef6060830185611f68565b828103604084015261210181856120a6565b9695505050505050565b60208082526014908201527318d85b1b081d1c985b9cd9995c8819985a5b195960621b604082015260600190565b60208082526021908201527f5468652072657475726e206f66207472616e73666572206973206661696c75726040820152606560f81b606082015260800190565b60008282101561218c5761218c611fd5565b500390565b6000602082840312156121a357600080fd5b8151611f6181611d34565b8481528360208201526080604082015260006121cd6080830185611f68565b8281036060840152611c0881856120a6565b600082198211156121f2576121f2611fd5565b500190565b634e487b7160e01b600052603160045260246000fd5b60006020828403121561221f57600080fd5b505191905056fea2646970667358221220ad302b6ae1b8e2666b35fccb745c230943537c0af78d3300892d3e8f8f67bd9164736f6c634300080d0033"

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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MetisPay *MetisPayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _MetisPay.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MetisPay *MetisPaySession) Owner() (common.Address, error) {
	return _MetisPay.Contract.Owner(&_MetisPay.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_MetisPay *MetisPayCallerSession) Owner() (common.Address, error) {
	return _MetisPay.Contract.Owner(&_MetisPay.CallOpts)
}

// TaskState is a free data retrieval call binding the contract method 0x26c6bee1.
//
// Solidity: function taskState(uint256 taskId) view returns(int8)
func (_MetisPay *MetisPayCaller) TaskState(opts *bind.CallOpts, taskId *big.Int) (int8, error) {
	var out []interface{}
	err := _MetisPay.contract.Call(opts, &out, "taskState", taskId)

	if err != nil {
		return *new(int8), err
	}

	out0 := *abi.ConvertType(out[0], new(int8)).(*int8)

	return out0, err

}

// TaskState is a free data retrieval call binding the contract method 0x26c6bee1.
//
// Solidity: function taskState(uint256 taskId) view returns(int8)
func (_MetisPay *MetisPaySession) TaskState(taskId *big.Int) (int8, error) {
	return _MetisPay.Contract.TaskState(&_MetisPay.CallOpts, taskId)
}

// TaskState is a free data retrieval call binding the contract method 0x26c6bee1.
//
// Solidity: function taskState(uint256 taskId) view returns(int8)
func (_MetisPay *MetisPayCallerSession) TaskState(taskId *big.Int) (int8, error) {
	return _MetisPay.Contract.TaskState(&_MetisPay.CallOpts, taskId)
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
// Solidity: function initialize(address metisLat) returns()
func (_MetisPay *MetisPayTransactor) Initialize(opts *bind.TransactOpts, metisLat common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "initialize", metisLat)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns()
func (_MetisPay *MetisPaySession) Initialize(metisLat common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.Initialize(&_MetisPay.TransactOpts, metisLat)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns()
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MetisPay *MetisPayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MetisPay *MetisPaySession) RenounceOwnership() (*types.Transaction, error) {
	return _MetisPay.Contract.RenounceOwnership(&_MetisPay.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_MetisPay *MetisPayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _MetisPay.Contract.RenounceOwnership(&_MetisPay.TransactOpts)
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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MetisPay *MetisPayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MetisPay.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MetisPay *MetisPaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.TransferOwnership(&_MetisPay.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_MetisPay *MetisPayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MetisPay.Contract.TransferOwnership(&_MetisPay.TransactOpts, newOwner)
}

// MetisPayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MetisPay contract.
type MetisPayOwnershipTransferredIterator struct {
	Event *MetisPayOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MetisPayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MetisPayOwnershipTransferred)
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
		it.Event = new(MetisPayOwnershipTransferred)
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
func (it *MetisPayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MetisPayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MetisPayOwnershipTransferred represents a OwnershipTransferred event raised by the MetisPay contract.
type MetisPayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MetisPay *MetisPayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MetisPayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MetisPay.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MetisPayOwnershipTransferredIterator{contract: _MetisPay.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MetisPay *MetisPayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MetisPayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MetisPay.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MetisPayOwnershipTransferred)
				if err := _MetisPay.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_MetisPay *MetisPayFilterer) ParseOwnershipTransferred(log types.Log) (*MetisPayOwnershipTransferred, error) {
	event := new(MetisPayOwnershipTransferred)
	if err := _MetisPay.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
