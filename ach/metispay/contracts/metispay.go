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
const MetisPayABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"PrepayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"SettleEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"metisLat\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"prepay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MetisPayBin is the compiled bytecode used for deploying new contracts.
var MetisPayBin = "0x608060405234801561001057600080fd5b506121ef806100206000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c8063b6a5d7de11610066578063b6a5d7de14610123578063c4d66de814610136578063f2fde38b14610149578063f80f5dd51461015c578063ff68f6af1461016f57600080fd5b80630ef92166146100a3578063715018a6146100cb5780638da5cb5b146100d55780639a9c29f6146100f05780639b19251a14610103575b600080fd5b6100b66100b1366004611dcb565b610182565b60405190151581526020015b60405180910390f35b6100d3610954565b005b6033546040516001600160a01b0390911681526020016100c2565b6100b66100fe366004611eb5565b6109ba565b610116610111366004611ed7565b6112c1565b6040516100c29190611f3f565b6100b6610131366004611ed7565b611395565b6100d3610144366004611ed7565b6113a6565b6100d3610157366004611ed7565b61147b565b6100b661016a366004611ed7565b611546565b6100b661017d366004611ed7565b6116e6565b600081518351146101da5760405162461bcd60e51b815260206004820152601960248201527f696e76616c696420746f6b656e20696e666f726d6174696f6e0000000000000060448201526064015b60405180910390fd5b606754600090815b8181101561023857876001600160a01b03166067828154811061020757610207611f52565b6000918252602090912001546001600160a01b03160361022657600192505b8061023081611f7e565b9150506101e2565b50816102865760405162461bcd60e51b815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f742065786973740060448201526064016101d1565b50506001600160a01b038516600090815260666020526040812054815b81811015610309576001600160a01b03881660009081526066602052604090208054339190839081106102d8576102d8611f52565b6000918252602090912001546001600160a01b0316036102f757600192505b8061030181611f7e565b9150506102a3565b50816103575760405162461bcd60e51b815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a6564000000000000000060448201526064016101d1565b5050606954600090815b818110156103a457886069828154811061037d5761037d611f52565b90600052602060002001540361039257600192505b8061039c81611f7e565b915050610361565b5081156103ec5760405162461bcd60e51b81526020600482015260166024820152757461736b20696420616c72656164792065786973747360501b60448201526064016101d1565b606554600090610406906001600160a01b03168930611aae565b9050868110156104585760405162461bcd60e51b815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e63650000000000000060448201526064016101d1565b8551915060005b828110156104ef5761048b87828151811061047c5761047c611f52565b60200260200101518a30611aae565b9150878210156104dd5760405162461bcd60e51b815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e63650060448201526064016101d1565b806104e781611f7e565b91505061045f565b506065546040516001600160a01b038a81166024830152306044830152606482018a905260009260609284929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516105579190611f97565b6000604051808303816000865af19150503d8060008114610594576040519150601f19603f3d011682016040523d82523d6000602084013e610599565b606091505b509093509150826105e75760405162461bcd60e51b815260206004820152601860248201527718d85b1b081d1c985b9cd9995c919c9bdb4819985a5b195960421b60448201526064016101d1565b818060200190518101906105fb9190611fd2565b90508061061a5760405162461bcd60e51b81526004016101d190611ff4565b60005b858110156107955789818151811061063757610637611f52565b60200260200101516001600160a01b03168c308b848151811061065c5761065c611f52565b60209081029190910101516040516001600160a01b039384166024820152929091166044830152606482015260840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516106c09190611f97565b6000604051808303816000865af19150503d80600081146106fd576040519150601f19603f3d011682016040523d82523d6000602084013e610702565b606091505b509094509250836107505760405162461bcd60e51b815260206004820152601860248201527718d85b1b081d1c985b9cd9995c919c9bdb4819985a5b195960421b60448201526064016101d1565b828060200190518101906107649190611fd2565b9150816107835760405162461bcd60e51b81526004016101d190611ff4565b8061078d81611f7e565b91505061061d565b506040518060e001604052808c6001600160a01b03168152602001336001600160a01b031681526020018b81526020018a815260200189815260200160011515815260200160001515815250606860008e815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550604082015181600201556060820151816003019080519060200190610874929190611c11565b5060808201518051610890916004840191602090910190611c76565b5060a08201516005909101805460c09093015115156101000261ff00199215159290921661ffff1990931692909217179055606980546001810182556000919091527f7fb4302e8e91f9110a6554c2c0a24601252c2a42c2220ca988efcfe399914308018c905560405133906001600160a01b038d16908e907f3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f9061093a908f908f908f90612069565b60405180910390a45060019b9a5050505050505050505050565b6033546001600160a01b031633146109ae5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016101d1565b6109b86000611bbf565b565b6069546000908180805b83811015610a0a5786606982815481106109e0576109e0611f52565b9060005260206000200154036109f857600192508091505b80610a0281611f7e565b9150506109c4565b5081610a4a5760405162461bcd60e51b815260206004820152600f60248201526e1a5b9d985b1a59081d185cdac81a59608a1b60448201526064016101d1565b6000868152606860205260409020600101546001600160a01b03163314610ab35760405162461bcd60e51b815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f2074686973000000000060448201526064016101d1565b60008681526068602052604090206005015460ff16610b0b5760405162461bcd60e51b81526020600482015260146024820152731c1c995c185e481b9bdd0818dbdb5c1b195d195960621b60448201526064016101d1565b600086815260686020526040902060050154610100900460ff1615610b625760405162461bcd60e51b815260206004820152600d60248201526c52657065617420736574746c6560981b60448201526064016101d1565b600086815260686020526040902060020154851115610bbb5760405162461bcd60e51b81526020600482015260156024820152741b5bdc99481d1a185b881c1c995c185a59081b185d605a1b60448201526064016101d1565b6065546000878152606860205260408082206001015490516001600160a01b03918216602482015260448101899052919260609291169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610c2a9190611f97565b6000604051808303816000865af19150503d8060008114610c67576040519150601f19603f3d011682016040523d82523d6000602084013e610c6c565b606091505b50909250905081610c8f5760405162461bcd60e51b81526004016101d19061209e565b80806020019051810190610ca39190611fd2565b610cbf5760405162461bcd60e51b81526004016101d1906120cc565b600088815260686020526040812060020154610cdc90899061210d565b90508015610de15760655460008a815260686020526040908190205490516001600160a01b0391821660248201526044810184905291169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610d4c9190611f97565b6000604051808303816000865af19150503d8060008114610d89576040519150601f19603f3d011682016040523d82523d6000602084013e610d8e565b606091505b50909350915082610db15760405162461bcd60e51b81526004016101d19061209e565b81806020019051810190610dc59190611fd2565b610de15760405162461bcd60e51b81526004016101d1906120cc565b600089815260686020908152604080832060030180548251818502810185019093528083528493830182828015610e4157602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610e23575b505050505090506000606860008d8152602001908152602001600020600401805480602002602001604051908101604052809291908181526020018280548015610eaa57602002820191906000526020600020905b815481526020019060010190808311610e96575b505085519394506000925050505b818110156110f157838181518110610ed257610ed2611f52565b602090810291909101810151604080516004815260248101825292830180516001600160e01b03166303aa30b960e11b179052516001600160a01b0390911691610f1b91611f97565b600060405180830381855afa9150503d8060008114610f56576040519150601f19603f3d011682016040523d82523d6000602084013e610f5b565b606091505b50909850965087610fa35760405162461bcd60e51b815260206004820152601260248201527118d85b1b081b5a5b9d195c8819985a5b195960721b60448201526064016101d1565b86806020019051810190610fb79190612124565b9450838181518110610fcb57610fcb611f52565b60200260200101516001600160a01b031685848381518110610fef57610fef611f52565b60209081029190910101516040516001600160a01b039092166024830152604482015260640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b1790525161104a9190611f97565b6000604051808303816000865af19150503d8060008114611087576040519150601f19603f3d011682016040523d82523d6000602084013e61108c565b606091505b509098509650876110af5760405162461bcd60e51b81526004016101d19061209e565b868060200190518101906110c39190611fd2565b6110df5760405162461bcd60e51b81526004016101d1906120cc565b806110e981611f7e565b915050610eb8565b50606860008e815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316606860008f815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168e7f44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef48f89888860405161118e9493929190612141565b60405180910390a4875b6111a360018c61210d565b8110156112035760696111b7826001612172565b815481106111c7576111c7611f52565b9060005260206000200154606982815481106111e5576111e5611f52565b600091825260209091200155806111fb81611f7e565b915050611198565b50606961121160018c61210d565b8154811061122157611221611f52565b6000918252602082200155606980548061123d5761123d61218a565b6000828152602080822083016000199081018390559092019092558e8252606890526040812080546001600160a01b03199081168255600182018054909116905560028101829055906112936003830182611cb1565b6112a1600483016000611cb1565b50600501805461ffff191690555060019c9b505050505050505050505050565b60608060005b60675481101561138e57836001600160a01b0316606782815481106112ee576112ee611f52565b6000918252602090912001546001600160a01b03160361137c576001600160a01b0384166000908152606660209081526040918290208054835181840281018401909452808452909183018282801561137057602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611352575b5050505050915061138e565b8061138681611f7e565b9150506112c7565b5092915050565b60006113a082611546565b92915050565b600054610100900460ff166113c15760005460ff16156113c5565b303b155b6114285760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016101d1565b600054610100900460ff1615801561144a576000805461ffff19166101011790555b606580546001600160a01b0319166001600160a01b0384161790558015611477576000805461ff00191690555b5050565b6033546001600160a01b031633146114d55760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016101d1565b6001600160a01b03811661153a5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016101d1565b61154381611bbf565b50565b600080805b6067548110156115a357336001600160a01b03166067828154811061157257611572611f52565b6000918252602090912001546001600160a01b03160361159157600191505b8061159b81611f7e565b91505061154b565b508015611669573360009081526066602052604081205490805b828110156116235733600090815260666020526040902080546001600160a01b0388169190839081106115f2576115f2611f52565b6000918252602090912001546001600160a01b03160361161157600191505b8061161b81611f7e565b9150506115bd565b5080611662573360009081526066602090815260408220805460018101825590835291200180546001600160a01b0319166001600160a01b0387161790555b50506116dd565b6067805460018082019092557f9787eeb91fe3101235e4a76063c7023ecb40f923f97916639c598592fa30d6ae018054336001600160a01b031991821681179092556000918252606660209081526040832080549485018155835290912090910180549091166001600160a01b0385161790555b50600192915050565b60675460009081908190815b8181101561174b57336001600160a01b03166067828154811061171757611717611f52565b6000918252602090912001546001600160a01b03160361173957600193508092505b8061174381611f7e565b9150506116f2565b50826117885760405162461bcd60e51b815260206004820152600c60248201526b34b73b30b634b2103ab9b2b960a11b60448201526064016101d1565b336000908152606660205260408120548190815b818110156118055733600090815260666020526040902080546001600160a01b038b169190839081106117d1576117d1611f52565b6000918252602090912001546001600160a01b0316036117f357600193508092505b806117fd81611f7e565b91505061179c565b50826118445760405162461bcd60e51b815260206004820152600e60248201526d696e76616c6964206167656e637960901b60448201526064016101d1565b600181111561198657815b61185a60018361210d565b8110156118ff5733600090815260666020526040902061187b826001612172565b8154811061188b5761188b611f52565b60009182526020808320909101543383526066909152604090912080546001600160a01b0390921691839081106118c4576118c4611f52565b600091825260209091200180546001600160a01b0319166001600160a01b0392909216919091179055806118f781611f7e565b91505061184f565b5033600090815260666020526040902061191a60018361210d565b8154811061192a5761192a611f52565b6000918252602080832090910180546001600160a01b0319169055338252606690526040902080548061195f5761195f61218a565b600082815260209020810160001990810180546001600160a01b0319169055019055611aa0565b845b61199360018661210d565b811015611a1e5760676119a7826001612172565b815481106119b7576119b7611f52565b600091825260209091200154606780546001600160a01b0390921691839081106119e3576119e3611f52565b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905580611a1681611f7e565b915050611988565b506067611a2c60018661210d565b81548110611a3c57611a3c611f52565b600091825260209091200180546001600160a01b03191690556067805480611a6657611a6661218a565b60008281526020808220830160001990810180546001600160a01b03191690559092019092553382526066905260408120611aa091611cb1565b506001979650505050505050565b6040516001600160a01b0383811660248301528281166044830152600091829182919087169060640160408051601f198184030181529181526020820180516001600160e01b0316636eb1769f60e11b17905251611b0c9190611f97565b600060405180830381855afa9150503d8060008114611b47576040519150601f19603f3d011682016040523d82523d6000602084013e611b4c565b606091505b509150915081611b9e5760405162461bcd60e51b815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c6564000000000060448201526064016101d1565b600081806020019051810190611bb491906121a0565b979650505050505050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b828054828255906000526020600020908101928215611c66579160200282015b82811115611c6657825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611c31565b50611c72929150611ccb565b5090565b828054828255906000526020600020908101928215611c66579160200282015b82811115611c66578251825591602001919060010190611c96565b508054600082559060005260206000209081019061154391905b5b80821115611c725760008155600101611ccc565b6001600160a01b038116811461154357600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715611d3457611d34611cf5565b604052919050565b600067ffffffffffffffff821115611d5657611d56611cf5565b5060051b60200190565b600082601f830112611d7157600080fd5b81356020611d86611d8183611d3c565b611d0b565b82815260059290921b84018101918181019086841115611da557600080fd5b8286015b84811015611dc05780358352918301918301611da9565b509695505050505050565b600080600080600060a08688031215611de357600080fd5b85359450602080870135611df681611ce0565b945060408701359350606087013567ffffffffffffffff80821115611e1a57600080fd5b818901915089601f830112611e2e57600080fd5b8135611e3c611d8182611d3c565b81815260059190911b8301840190848101908c831115611e5b57600080fd5b938501935b82851015611e82578435611e7381611ce0565b82529385019390850190611e60565b965050506080890135925080831115611e9a57600080fd5b5050611ea888828901611d60565b9150509295509295909350565b60008060408385031215611ec857600080fd5b50508035926020909101359150565b600060208284031215611ee957600080fd5b8135611ef481611ce0565b9392505050565b600081518084526020808501945080840160005b83811015611f345781516001600160a01b031687529582019590820190600101611f0f565b509495945050505050565b602081526000611ef46020830184611efb565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611f9057611f90611f68565b5060010190565b6000825160005b81811015611fb85760208186018101518583015201611f9e565b81811115611fc7576000828501525b509190910192915050565b600060208284031215611fe457600080fd5b81518015158114611ef457600080fd5b60208082526025908201527f5468652072657475726e206f66207472616e7366657266726f6d206973206661604082015264696c75726560d81b606082015260800190565b600081518084526020808501945080840160005b83811015611f345781518752958201959082019060010161204d565b8381526060602082015260006120826060830185611efb565b82810360408401526120948185612039565b9695505050505050565b60208082526014908201527318d85b1b081d1c985b9cd9995c8819985a5b195960621b604082015260600190565b60208082526021908201527f5468652072657475726e206f66207472616e73666572206973206661696c75726040820152606560f81b606082015260800190565b60008282101561211f5761211f611f68565b500390565b60006020828403121561213657600080fd5b8151611ef481611ce0565b8481528360208201526080604082015260006121606080830185611efb565b8281036060840152611bb48185612039565b6000821982111561218557612185611f68565b500190565b634e487b7160e01b600052603160045260246000fd5b6000602082840312156121b257600080fd5b505191905056fea26469706673582212206f5e484f4168e674fe03b8c941daa778c68d8d9f871c8a9c22f00b24fffca44a64736f6c634300080d0033"

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
