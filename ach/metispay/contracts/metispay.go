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
const MetisPayABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"PrepayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refundOrAdd\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"SettleEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"getTaskInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"metisLat\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"prepay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"taskState\",\"outputs\":[{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// MetisPayBin is the compiled bytecode used for deploying new contracts.
var MetisPayBin = "0x608060405234801561001057600080fd5b506125d6806100206000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063b6a5d7de11610071578063b6a5d7de1461015f578063c4d66de814610172578063d1a1b99914610185578063f2fde38b146101aa578063f80f5dd5146101bd578063ff68f6af146101d057600080fd5b80630ef92166146100b957806326c6bee1146100e1578063715018a6146101075780638da5cb5b146101115780639a9c29f61461012c5780639b19251a1461013f575b600080fd5b6100cc6100c7366004612106565b6101e3565b60405190151581526020015b60405180910390f35b6100f46100ef3660046121f0565b610943565b60405160009190910b81526020016100d8565b61010f6109c1565b005b6033546040516001600160a01b0390911681526020016100d8565b6100cc61013a366004612209565b610a27565b61015261014d36600461222b565b6113e9565b6040516100d89190612293565b6100cc61016d36600461222b565b6114bd565b61010f61018036600461222b565b6114ce565b6101986101933660046121f0565b6115a3565b6040516100d8969594939291906122d6565b61010f6101b836600461222b565b6117b6565b6100cc6101cb36600461222b565b611881565b6100cc6101de36600461222b565b611a21565b6000815183511461023b5760405162461bcd60e51b815260206004820152601960248201527f696e76616c696420746f6b656e20696e666f726d6174696f6e0000000000000060448201526064015b60405180910390fd5b606754600090815b8181101561029957876001600160a01b03166067828154811061026857610268612332565b6000918252602090912001546001600160a01b03160361028757600192505b806102918161235e565b915050610243565b50816102e75760405162461bcd60e51b815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f74206578697374006044820152606401610232565b50506001600160a01b038516600090815260666020526040812054815b8181101561036a576001600160a01b038816600090815260666020526040902080543391908390811061033957610339612332565b6000918252602090912001546001600160a01b03160361035857600192505b806103628161235e565b915050610304565b50816103b85760405162461bcd60e51b815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a656400000000000000006044820152606401610232565b5050606954600090815b818110156104055788606982815481106103de576103de612332565b9060005260206000200154036103f357600192505b806103fd8161235e565b9150506103c2565b50811561044d5760405162461bcd60e51b81526020600482015260166024820152757461736b20696420616c72656164792065786973747360501b6044820152606401610232565b606554600090610467906001600160a01b03168930611de9565b9050868110156104b95760405162461bcd60e51b815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e6365000000000000006044820152606401610232565b8551915060005b82811015610550576104ec8782815181106104dd576104dd612332565b60200260200101518a30611de9565b91508782101561053e5760405162461bcd60e51b815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e6365006044820152606401610232565b806105488161235e565b9150506104c0565b506065546040516001600160a01b038a81166024830152306044830152606482018a905260009260609284929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516105b89190612377565b6000604051808303816000865af19150503d80600081146105f5576040519150601f19603f3d011682016040523d82523d6000602084013e6105fa565b606091505b5090935091508261061d5760405162461bcd60e51b8152600401610232906123b2565b8180602001905181019061063191906123e9565b9050806106505760405162461bcd60e51b81526004016102329061240b565b60005b858110156107a05789818151811061066d5761066d612332565b60200260200101516001600160a01b03168c308b848151811061069257610692612332565b60209081029190910101516040516001600160a01b039384166024820152929091166044830152606482015260840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516106f69190612377565b6000604051808303816000865af19150503d8060008114610733576040519150601f19603f3d011682016040523d82523d6000602084013e610738565b606091505b5090945092508361075b5760405162461bcd60e51b8152600401610232906123b2565b8280602001905181019061076f91906123e9565b91508161078e5760405162461bcd60e51b81526004016102329061240b565b806107988161235e565b915050610653565b506040518060c001604052808c6001600160a01b03168152602001336001600160a01b031681526020018b81526020018a8152602001898152602001600160000b815250606860008e815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550604082015181600201556060820151816003019080519060200190610877929190611f4c565b5060808201518051610893916004840191602090910190611fb1565b5060a091909101516005909101805460ff191660ff909216919091179055606980546001810182556000919091527f7fb4302e8e91f9110a6554c2c0a24601252c2a42c2220ca988efcfe399914308018c905560405133906001600160a01b038d16908e907f3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f90610929908f908f908f90612450565b60405180910390a45060019b9a5050505050505050505050565b606954600090819081805b8281101561099157856069828154811061096a5761096a612332565b90600052602060002001540361097f57600191505b806109898161235e565b91505061094e565b50806109a15760001992506109b8565b600085815260686020526040812060050154900b92505b50909392505050565b6033546001600160a01b03163314610a1b5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610232565b610a256000611efa565b565b6069546000908180805b83811015610a77578660698281548110610a4d57610a4d612332565b906000526020600020015403610a6557600192508091505b80610a6f8161235e565b915050610a31565b5081610ab75760405162461bcd60e51b815260206004820152600f60248201526e1a5b9d985b1a59081d185cdac81a59608a1b6044820152606401610232565b6000868152606860205260409020600101546001600160a01b03163314610b205760405162461bcd60e51b815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f207468697300000000006044820152606401610232565b600086815260686020526040812060050154900b600114610b915760405162461bcd60e51b815260206004820152602560248201527f707265706179206e6f7420636f6d706c65746564206f722072657065617420736044820152646574746c6560d81b6064820152608401610232565b6000868152606860205260408120600201546060908290881015610cd057600089815260686020526040902060020154610bcc908990612485565b60655460008b815260686020526040908190205490516001600160a01b03918216602482015260448101849052929350169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610c369190612377565b6000604051808303816000865af19150503d8060008114610c73576040519150601f19603f3d011682016040523d82523d6000602084013e610c78565b606091505b50909350915082610c9b5760405162461bcd60e51b81526004016102329061249c565b81806020019051810190610caf91906123e9565b610ccb5760405162461bcd60e51b8152600401610232906124ca565b610e0a565b600089815260686020526040902060020154881115610e0a57600089815260686020526040902060020154610d059089612485565b60655460008b815260686020526040908190205490516001600160a01b03918216602482015230604482015260648101849052929350169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b17905251610d759190612377565b6000604051808303816000865af19150503d8060008114610db2576040519150601f19603f3d011682016040523d82523d6000602084013e610db7565b606091505b50909350915082610dda5760405162461bcd60e51b8152600401610232906123b2565b81806020019051810190610dee91906123e9565b610e0a5760405162461bcd60e51b8152600401610232906124ca565b60655460008a815260686020526040908190206001015490516001600160a01b039182166024820152604481018b905291169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610e759190612377565b6000604051808303816000865af19150503d8060008114610eb2576040519150601f19603f3d011682016040523d82523d6000602084013e610eb7565b606091505b50909350915082610eda5760405162461bcd60e51b81526004016102329061249c565b81806020019051810190610eee91906123e9565b610f0a5760405162461bcd60e51b8152600401610232906124ca565b600089815260686020908152604080832060030180548251818502810185019093528083528493830182828015610f6a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610f4c575b505050505090506000606860008d8152602001908152602001600020600401805480602002602001604051908101604052809291908181526020018280548015610fd357602002820191906000526020600020905b815481526020019060010190808311610fbf575b505085519394506000925050505b8181101561121a57838181518110610ffb57610ffb612332565b602090810291909101810151604080516004815260248101825292830180516001600160e01b03166303aa30b960e11b179052516001600160a01b039091169161104491612377565b600060405180830381855afa9150503d806000811461107f576040519150601f19603f3d011682016040523d82523d6000602084013e611084565b606091505b509098509650876110cc5760405162461bcd60e51b815260206004820152601260248201527118d85b1b081b5a5b9d195c8819985a5b195960721b6044820152606401610232565b868060200190518101906110e0919061250b565b94508381815181106110f4576110f4612332565b60200260200101516001600160a01b03168584838151811061111857611118612332565b60209081029190910101516040516001600160a01b039092166024830152604482015260640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b179052516111739190612377565b6000604051808303816000865af19150503d80600081146111b0576040519150601f19603f3d011682016040523d82523d6000602084013e6111b5565b606091505b509098509650876111d85760405162461bcd60e51b81526004016102329061249c565b868060200190518101906111ec91906123e9565b6112085760405162461bcd60e51b8152600401610232906124ca565b806112128161235e565b915050610fe1565b50606860008e815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316606860008f815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168e7f44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef48f8988886040516112b79493929190612528565b60405180910390a4875b6112cc60018c612485565b81101561132c5760696112e0826001612559565b815481106112f0576112f0612332565b90600052602060002001546069828154811061130e5761130e612332565b600091825260209091200155806113248161235e565b9150506112c1565b50606961133a60018c612485565b8154811061134a5761134a612332565b6000918252602082200155606980548061136657611366612571565b6000828152602080822083016000199081018390559092019092558e8252606890526040812080546001600160a01b03199081168255600182018054909116905560028101829055906113bc6003830182611fec565b6113ca600483016000611fec565b50600501805460ff191690555060019c9b505050505050505050505050565b60608060005b6067548110156114b657836001600160a01b03166067828154811061141657611416612332565b6000918252602090912001546001600160a01b0316036114a4576001600160a01b0384166000908152606660209081526040918290208054835181840281018401909452808452909183018282801561149857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161147a575b505050505091506114b6565b806114ae8161235e565b9150506113ef565b5092915050565b60006114c882611881565b92915050565b600054610100900460ff166114e95760005460ff16156114ed565b303b155b6115505760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610232565b600054610100900460ff16158015611572576000805461ffff19166101011790555b606580546001600160a01b0319166001600160a01b038416179055801561159f576000805461ff00191690555b5050565b600080600060608060008060698054905090506000805b828110156115fd5789606982815481106115d6576115d6612332565b9060005260206000200154036115eb57600191505b806115f58161235e565b9150506115ba565b50806116b3576040805160018082528183019092526000916020808301908036833701905050905060008160008151811061163a5761163a612332565b6001600160a01b03929092166020928302919091019091015260408051600180825281830190925260009181602001602082028036833701905050905060008160008151811061168c5761168c612332565b602090810291909101015260009950899850889750909550935060001992506117ad915050565b600089815260686020908152604080832080546001820154600283015460058401546003850180548751818a0281018a019098528088526001600160a01b0395861699949095169792969095600401949190930b929185919083018282801561174557602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611727575b505050505092508180548060200260200160405190810160405280929190818152602001828054801561179757602002820191906000526020600020905b815481526020019060010190808311611783575b5050505050915097509750975097509750975050505b91939550919395565b6033546001600160a01b031633146118105760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610232565b6001600160a01b0381166118755760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610232565b61187e81611efa565b50565b600080805b6067548110156118de57336001600160a01b0316606782815481106118ad576118ad612332565b6000918252602090912001546001600160a01b0316036118cc57600191505b806118d68161235e565b915050611886565b5080156119a4573360009081526066602052604081205490805b8281101561195e5733600090815260666020526040902080546001600160a01b03881691908390811061192d5761192d612332565b6000918252602090912001546001600160a01b03160361194c57600191505b806119568161235e565b9150506118f8565b508061199d573360009081526066602090815260408220805460018101825590835291200180546001600160a01b0319166001600160a01b0387161790555b5050611a18565b6067805460018082019092557f9787eeb91fe3101235e4a76063c7023ecb40f923f97916639c598592fa30d6ae018054336001600160a01b031991821681179092556000918252606660209081526040832080549485018155835290912090910180549091166001600160a01b0385161790555b50600192915050565b60675460009081908190815b81811015611a8657336001600160a01b031660678281548110611a5257611a52612332565b6000918252602090912001546001600160a01b031603611a7457600193508092505b80611a7e8161235e565b915050611a2d565b5082611ac35760405162461bcd60e51b815260206004820152600c60248201526b34b73b30b634b2103ab9b2b960a11b6044820152606401610232565b336000908152606660205260408120548190815b81811015611b405733600090815260666020526040902080546001600160a01b038b16919083908110611b0c57611b0c612332565b6000918252602090912001546001600160a01b031603611b2e57600193508092505b80611b388161235e565b915050611ad7565b5082611b7f5760405162461bcd60e51b815260206004820152600e60248201526d696e76616c6964206167656e637960901b6044820152606401610232565b6001811115611cc157815b611b95600183612485565b811015611c3a57336000908152606660205260409020611bb6826001612559565b81548110611bc657611bc6612332565b60009182526020808320909101543383526066909152604090912080546001600160a01b039092169183908110611bff57611bff612332565b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905580611c328161235e565b915050611b8a565b50336000908152606660205260409020611c55600183612485565b81548110611c6557611c65612332565b6000918252602080832090910180546001600160a01b03191690553382526066905260409020805480611c9a57611c9a612571565b600082815260209020810160001990810180546001600160a01b0319169055019055611ddb565b845b611cce600186612485565b811015611d59576067611ce2826001612559565b81548110611cf257611cf2612332565b600091825260209091200154606780546001600160a01b039092169183908110611d1e57611d1e612332565b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905580611d518161235e565b915050611cc3565b506067611d67600186612485565b81548110611d7757611d77612332565b600091825260209091200180546001600160a01b03191690556067805480611da157611da1612571565b60008281526020808220830160001990810180546001600160a01b03191690559092019092553382526066905260408120611ddb91611fec565b506001979650505050505050565b6040516001600160a01b0383811660248301528281166044830152600091829182919087169060640160408051601f198184030181529181526020820180516001600160e01b0316636eb1769f60e11b17905251611e479190612377565b600060405180830381855afa9150503d8060008114611e82576040519150601f19603f3d011682016040523d82523d6000602084013e611e87565b606091505b509150915081611ed95760405162461bcd60e51b815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c656400000000006044820152606401610232565b600081806020019051810190611eef9190612587565b979650505050505050565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b828054828255906000526020600020908101928215611fa1579160200282015b82811115611fa157825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611f6c565b50611fad929150612006565b5090565b828054828255906000526020600020908101928215611fa1579160200282015b82811115611fa1578251825591602001919060010190611fd1565b508054600082559060005260206000209081019061187e91905b5b80821115611fad5760008155600101612007565b6001600160a01b038116811461187e57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561206f5761206f612030565b604052919050565b600067ffffffffffffffff82111561209157612091612030565b5060051b60200190565b600082601f8301126120ac57600080fd5b813560206120c16120bc83612077565b612046565b82815260059290921b840181019181810190868411156120e057600080fd5b8286015b848110156120fb57803583529183019183016120e4565b509695505050505050565b600080600080600060a0868803121561211e57600080fd5b853594506020808701356121318161201b565b945060408701359350606087013567ffffffffffffffff8082111561215557600080fd5b818901915089601f83011261216957600080fd5b81356121776120bc82612077565b81815260059190911b8301840190848101908c83111561219657600080fd5b938501935b828510156121bd5784356121ae8161201b565b8252938501939085019061219b565b9650505060808901359250808311156121d557600080fd5b50506121e38882890161209b565b9150509295509295909350565b60006020828403121561220257600080fd5b5035919050565b6000806040838503121561221c57600080fd5b50508035926020909101359150565b60006020828403121561223d57600080fd5b81356122488161201b565b9392505050565b600081518084526020808501945080840160005b838110156122885781516001600160a01b031687529582019590820190600101612263565b509495945050505050565b602081526000612248602083018461224f565b600081518084526020808501945080840160005b83811015612288578151875295820195908201906001016122ba565b6001600160a01b038781168252861660208201526040810185905260c0606082018190526000906123099083018661224f565b828103608084015261231b81866122a6565b9150508260000b60a0830152979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b60006001820161237057612370612348565b5060010190565b6000825160005b81811015612398576020818601810151858301520161237e565b818111156123a7576000828501525b509190910192915050565b60208082526018908201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604082015260600190565b6000602082840312156123fb57600080fd5b8151801515811461224857600080fd5b60208082526025908201527f5468652072657475726e206f66207472616e7366657266726f6d206973206661604082015264696c75726560d81b606082015260800190565b838152606060208201526000612469606083018561224f565b828103604084015261247b81856122a6565b9695505050505050565b60008282101561249757612497612348565b500390565b60208082526014908201527318d85b1b081d1c985b9cd9995c8819985a5b195960621b604082015260600190565b60208082526021908201527f5468652072657475726e206f66207472616e73666572206973206661696c75726040820152606560f81b606082015260800190565b60006020828403121561251d57600080fd5b81516122488161201b565b848152836020820152608060408201526000612547608083018561224f565b8281036060840152611eef81856122a6565b6000821982111561256c5761256c612348565b500190565b634e487b7160e01b600052603160045260246000fd5b60006020828403121561259957600080fd5b505191905056fea2646970667358221220186c1135dce9634cb95ea47f88f2e008a7a2a88eee4bf652840b6b06fe5d914564736f6c634300080d0033"

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

// GetTaskInfo is a free data retrieval call binding the contract method 0xd1a1b999.
//
// Solidity: function getTaskInfo(uint256 taskId) view returns(address, address, uint256, address[], uint256[], int8)
func (_MetisPay *MetisPayCaller) GetTaskInfo(opts *bind.CallOpts, taskId *big.Int) (common.Address, common.Address, *big.Int, []common.Address, []*big.Int, int8, error) {
	var out []interface{}
	err := _MetisPay.contract.Call(opts, &out, "getTaskInfo", taskId)

	if err != nil {
		return *new(common.Address), *new(common.Address), *new(*big.Int), *new([]common.Address), *new([]*big.Int), *new(int8), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)
	out4 := *abi.ConvertType(out[4], new([]*big.Int)).(*[]*big.Int)
	out5 := *abi.ConvertType(out[5], new(int8)).(*int8)

	return out0, out1, out2, out3, out4, out5, err

}

// GetTaskInfo is a free data retrieval call binding the contract method 0xd1a1b999.
//
// Solidity: function getTaskInfo(uint256 taskId) view returns(address, address, uint256, address[], uint256[], int8)
func (_MetisPay *MetisPaySession) GetTaskInfo(taskId *big.Int) (common.Address, common.Address, *big.Int, []common.Address, []*big.Int, int8, error) {
	return _MetisPay.Contract.GetTaskInfo(&_MetisPay.CallOpts, taskId)
}

// GetTaskInfo is a free data retrieval call binding the contract method 0xd1a1b999.
//
// Solidity: function getTaskInfo(uint256 taskId) view returns(address, address, uint256, address[], uint256[], int8)
func (_MetisPay *MetisPayCallerSession) GetTaskInfo(taskId *big.Int) (common.Address, common.Address, *big.Int, []common.Address, []*big.Int, int8, error) {
	return _MetisPay.Contract.GetTaskInfo(&_MetisPay.CallOpts, taskId)
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
	RefundOrAdd      *big.Int
	TokenAddressList []common.Address
	TokenValueList   []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterSettleEvent is a free log retrieval operation binding the contract event 0x44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef4.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refundOrAdd, address[] tokenAddressList, uint256[] tokenValueList)
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
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refundOrAdd, address[] tokenAddressList, uint256[] tokenValueList)
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
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refundOrAdd, address[] tokenAddressList, uint256[] tokenValueList)
func (_MetisPay *MetisPayFilterer) ParseSettleEvent(log types.Log) (*MetisPaySettleEvent, error) {
	event := new(MetisPaySettleEvent)
	if err := _MetisPay.contract.UnpackLog(event, "SettleEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
