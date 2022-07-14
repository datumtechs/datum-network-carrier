// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DatumPayMetaData contains all meta data concerning the DatumPay contract.
var DatumPayMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"PrepayEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"Agencyfee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refundOrAdd\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"SettleEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"addWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"authorize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAgency\",\"type\":\"address\"}],\"name\":\"deleteWhitelist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"getTaskInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"metisLat\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"tokenAddressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"tokenValueList\",\"type\":\"uint256[]\"}],\"name\":\"prepay\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"taskId\",\"type\":\"uint256\"}],\"name\":\"taskState\",\"outputs\":[{\"internalType\":\"int8\",\"name\":\"\",\"type\":\"int8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddress\",\"type\":\"address\"}],\"name\":\"whitelist\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506125dc806100206000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063b6a5d7de11610071578063b6a5d7de1461015f578063c4d66de814610172578063d1a1b99914610185578063f2fde38b146101aa578063f80f5dd5146101bd578063ff68f6af146101d057600080fd5b80630ef92166146100b957806326c6bee1146100e1578063715018a6146101075780638da5cb5b146101115780639a9c29f61461012c5780639b19251a1461013f575b600080fd5b6100cc6100c736600461210c565b6101e3565b60405190151581526020015b60405180910390f35b6100f46100ef3660046121f6565b610943565b60405160009190910b81526020016100d8565b61010f6109c1565b005b6033546040516001600160a01b0390911681526020016100d8565b6100cc61013a36600461220f565b6109d5565b61015261014d366004612231565b611397565b6040516100d89190612299565b6100cc61016d366004612231565b61146b565b61010f610180366004612231565b61147c565b6101986101933660046121f6565b6115a1565b6040516100d8969594939291906122dc565b61010f6101b8366004612231565b6117b4565b6100cc6101cb366004612231565b61182d565b6100cc6101de366004612231565b6119cd565b6000815183511461023b5760405162461bcd60e51b815260206004820152601960248201527f696e76616c696420746f6b656e20696e666f726d6174696f6e0000000000000060448201526064015b60405180910390fd5b606754600090815b8181101561029957876001600160a01b03166067828154811061026857610268612338565b6000918252602090912001546001600160a01b03160361028757600192505b8061029181612364565b915050610243565b50816102e75760405162461bcd60e51b815260206004820152601f60248201527f7573657227732077686974656c69737420646f6573206e6f74206578697374006044820152606401610232565b50506001600160a01b038516600090815260666020526040812054815b8181101561036a576001600160a01b038816600090815260666020526040902080543391908390811061033957610339612338565b6000918252602090912001546001600160a01b03160361035857600192505b8061036281612364565b915050610304565b50816103b85760405162461bcd60e51b815260206004820152601860248201527f6167656e6379206973206e6f7420617574686f72697a656400000000000000006044820152606401610232565b5050606954600090815b818110156104055788606982815481106103de576103de612338565b9060005260206000200154036103f357600192505b806103fd81612364565b9150506103c2565b50811561044d5760405162461bcd60e51b81526020600482015260166024820152757461736b20696420616c72656164792065786973747360501b6044820152606401610232565b606554600090610467906001600160a01b03168930611d95565b9050868110156104b95760405162461bcd60e51b815260206004820152601960248201527f776c617420696e73756666696369656e742062616c616e6365000000000000006044820152606401610232565b8551915060005b82811015610550576104ec8782815181106104dd576104dd612338565b60200260200101518a30611d95565b91508782101561053e5760405162461bcd60e51b815260206004820152601f60248201527f6461746120746f6b656e20696e73756666696369656e742062616c616e6365006044820152606401610232565b8061054881612364565b9150506104c0565b506065546040516001600160a01b038a81166024830152306044830152606482018a905260009260609284929091169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516105b8919061237d565b6000604051808303816000865af19150503d80600081146105f5576040519150601f19603f3d011682016040523d82523d6000602084013e6105fa565b606091505b5090935091508261061d5760405162461bcd60e51b8152600401610232906123b8565b8180602001905181019061063191906123ef565b9050806106505760405162461bcd60e51b815260040161023290612411565b60005b858110156107a05789818151811061066d5761066d612338565b60200260200101516001600160a01b03168c308b848151811061069257610692612338565b60209081029190910101516040516001600160a01b039384166024820152929091166044830152606482015260840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b179052516106f6919061237d565b6000604051808303816000865af19150503d8060008114610733576040519150601f19603f3d011682016040523d82523d6000602084013e610738565b606091505b5090945092508361075b5760405162461bcd60e51b8152600401610232906123b8565b8280602001905181019061076f91906123ef565b91508161078e5760405162461bcd60e51b815260040161023290612411565b8061079881612364565b915050610653565b506040518060c001604052808c6001600160a01b03168152602001336001600160a01b031681526020018b81526020018a8152602001898152602001600160000b815250606860008e815260200190815260200160002060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160a01b0302191690836001600160a01b03160217905550604082015181600201556060820151816003019080519060200190610877929190611f52565b5060808201518051610893916004840191602090910190611fb7565b5060a091909101516005909101805460ff191660ff909216919091179055606980546001810182556000919091527f7fb4302e8e91f9110a6554c2c0a24601252c2a42c2220ca988efcfe399914308018c905560405133906001600160a01b038d16908e907f3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f90610929908f908f908f90612456565b60405180910390a45060019b9a5050505050505050505050565b606954600090819081805b8281101561099157856069828154811061096a5761096a612338565b90600052602060002001540361097f57600191505b8061098981612364565b91505061094e565b50806109a15760001992506109b8565b600085815260686020526040812060050154900b92505b50909392505050565b6109c9611ea6565b6109d36000611f00565b565b6069546000908180805b83811015610a255786606982815481106109fb576109fb612338565b906000526020600020015403610a1357600192508091505b80610a1d81612364565b9150506109df565b5081610a655760405162461bcd60e51b815260206004820152600f60248201526e1a5b9d985b1a59081d185cdac81a59608a1b6044820152606401610232565b6000868152606860205260409020600101546001600160a01b03163314610ace5760405162461bcd60e51b815260206004820152601b60248201527f4f6e6c792075736572206167656e742063616e20646f207468697300000000006044820152606401610232565b600086815260686020526040812060050154900b600114610b3f5760405162461bcd60e51b815260206004820152602560248201527f707265706179206e6f7420636f6d706c65746564206f722072657065617420736044820152646574746c6560d81b6064820152608401610232565b6000868152606860205260408120600201546060908290881015610c7e57600089815260686020526040902060020154610b7a90899061248b565b60655460008b815260686020526040908190205490516001600160a01b03918216602482015260448101849052929350169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610be4919061237d565b6000604051808303816000865af19150503d8060008114610c21576040519150601f19603f3d011682016040523d82523d6000602084013e610c26565b606091505b50909350915082610c495760405162461bcd60e51b8152600401610232906124a2565b81806020019051810190610c5d91906123ef565b610c795760405162461bcd60e51b8152600401610232906124d0565b610db8565b600089815260686020526040902060020154881115610db857600089815260686020526040902060020154610cb3908961248b565b60655460008b815260686020526040908190205490516001600160a01b03918216602482015230604482015260648101849052929350169060840160408051601f198184030181529181526020820180516001600160e01b03166323b872dd60e01b17905251610d23919061237d565b6000604051808303816000865af19150503d8060008114610d60576040519150601f19603f3d011682016040523d82523d6000602084013e610d65565b606091505b50909350915082610d885760405162461bcd60e51b8152600401610232906123b8565b81806020019051810190610d9c91906123ef565b610db85760405162461bcd60e51b8152600401610232906124d0565b60655460008a815260686020526040908190206001015490516001600160a01b039182166024820152604481018b905291169060640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251610e23919061237d565b6000604051808303816000865af19150503d8060008114610e60576040519150601f19603f3d011682016040523d82523d6000602084013e610e65565b606091505b50909350915082610e885760405162461bcd60e51b8152600401610232906124a2565b81806020019051810190610e9c91906123ef565b610eb85760405162461bcd60e51b8152600401610232906124d0565b600089815260686020908152604080832060030180548251818502810185019093528083528493830182828015610f1857602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311610efa575b505050505090506000606860008d8152602001908152602001600020600401805480602002602001604051908101604052809291908181526020018280548015610f8157602002820191906000526020600020905b815481526020019060010190808311610f6d575b505085519394506000925050505b818110156111c857838181518110610fa957610fa9612338565b602090810291909101810151604080516004815260248101825292830180516001600160e01b03166303aa30b960e11b179052516001600160a01b0390911691610ff29161237d565b600060405180830381855afa9150503d806000811461102d576040519150601f19603f3d011682016040523d82523d6000602084013e611032565b606091505b5090985096508761107a5760405162461bcd60e51b815260206004820152601260248201527118d85b1b081b5a5b9d195c8819985a5b195960721b6044820152606401610232565b8680602001905181019061108e9190612511565b94508381815181106110a2576110a2612338565b60200260200101516001600160a01b0316858483815181106110c6576110c6612338565b60209081029190910101516040516001600160a01b039092166024830152604482015260640160408051601f198184030181529181526020820180516001600160e01b031663a9059cbb60e01b17905251611121919061237d565b6000604051808303816000865af19150503d806000811461115e576040519150601f19603f3d011682016040523d82523d6000602084013e611163565b606091505b509098509650876111865760405162461bcd60e51b8152600401610232906124a2565b8680602001905181019061119a91906123ef565b6111b65760405162461bcd60e51b8152600401610232906124d0565b806111c081612364565b915050610f8f565b50606860008e815260200190815260200160002060010160009054906101000a90046001600160a01b03166001600160a01b0316606860008f815260200190815260200160002060000160009054906101000a90046001600160a01b03166001600160a01b03168e7f44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef48f898888604051611265949392919061252e565b60405180910390a4875b61127a60018c61248b565b8110156112da57606961128e82600161255f565b8154811061129e5761129e612338565b9060005260206000200154606982815481106112bc576112bc612338565b600091825260209091200155806112d281612364565b91505061126f565b5060696112e860018c61248b565b815481106112f8576112f8612338565b6000918252602082200155606980548061131457611314612577565b6000828152602080822083016000199081018390559092019092558e8252606890526040812080546001600160a01b031990811682556001820180549091169055600281018290559061136a6003830182611ff2565b611378600483016000611ff2565b50600501805460ff191690555060019c9b505050505050505050505050565b60608060005b60675481101561146457836001600160a01b0316606782815481106113c4576113c4612338565b6000918252602090912001546001600160a01b031603611452576001600160a01b0384166000908152606660209081526040918290208054835181840281018401909452808452909183018282801561144657602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611428575b50505050509150611464565b8061145c81612364565b91505061139d565b5092915050565b60006114768261182d565b92915050565b600054610100900460ff161580801561149c5750600054600160ff909116105b806114b65750303b1580156114b6575060005460ff166001145b6115195760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610232565b6000805460ff19166001179055801561153c576000805461ff0019166101001790555b606580546001600160a01b0319166001600160a01b038416179055801561159d576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050565b600080600060608060008060698054905090506000805b828110156115fb5789606982815481106115d4576115d4612338565b9060005260206000200154036115e957600191505b806115f381612364565b9150506115b8565b50806116b1576040805160018082528183019092526000916020808301908036833701905050905060008160008151811061163857611638612338565b6001600160a01b03929092166020928302919091019091015260408051600180825281830190925260009181602001602082028036833701905050905060008160008151811061168a5761168a612338565b602090810291909101015260009950899850889750909550935060001992506117ab915050565b600089815260686020908152604080832080546001820154600283015460058401546003850180548751818a0281018a019098528088526001600160a01b0395861699949095169792969095600401949190930b929185919083018282801561174357602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611725575b505050505092508180548060200260200160405190810160405280929190818152602001828054801561179557602002820191906000526020600020905b815481526020019060010190808311611781575b5050505050915097509750975097509750975050505b91939550919395565b6117bc611ea6565b6001600160a01b0381166118215760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610232565b61182a81611f00565b50565b600080805b60675481101561188a57336001600160a01b03166067828154811061185957611859612338565b6000918252602090912001546001600160a01b03160361187857600191505b8061188281612364565b915050611832565b508015611950573360009081526066602052604081205490805b8281101561190a5733600090815260666020526040902080546001600160a01b0388169190839081106118d9576118d9612338565b6000918252602090912001546001600160a01b0316036118f857600191505b8061190281612364565b9150506118a4565b5080611949573360009081526066602090815260408220805460018101825590835291200180546001600160a01b0319166001600160a01b0387161790555b50506119c4565b6067805460018082019092557f9787eeb91fe3101235e4a76063c7023ecb40f923f97916639c598592fa30d6ae018054336001600160a01b031991821681179092556000918252606660209081526040832080549485018155835290912090910180549091166001600160a01b0385161790555b50600192915050565b60675460009081908190815b81811015611a3257336001600160a01b0316606782815481106119fe576119fe612338565b6000918252602090912001546001600160a01b031603611a2057600193508092505b80611a2a81612364565b9150506119d9565b5082611a6f5760405162461bcd60e51b815260206004820152600c60248201526b34b73b30b634b2103ab9b2b960a11b6044820152606401610232565b336000908152606660205260408120548190815b81811015611aec5733600090815260666020526040902080546001600160a01b038b16919083908110611ab857611ab8612338565b6000918252602090912001546001600160a01b031603611ada57600193508092505b80611ae481612364565b915050611a83565b5082611b2b5760405162461bcd60e51b815260206004820152600e60248201526d696e76616c6964206167656e637960901b6044820152606401610232565b6001811115611c6d57815b611b4160018361248b565b811015611be657336000908152606660205260409020611b6282600161255f565b81548110611b7257611b72612338565b60009182526020808320909101543383526066909152604090912080546001600160a01b039092169183908110611bab57611bab612338565b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905580611bde81612364565b915050611b36565b50336000908152606660205260409020611c0160018361248b565b81548110611c1157611c11612338565b6000918252602080832090910180546001600160a01b03191690553382526066905260409020805480611c4657611c46612577565b600082815260209020810160001990810180546001600160a01b0319169055019055611d87565b845b611c7a60018661248b565b811015611d05576067611c8e82600161255f565b81548110611c9e57611c9e612338565b600091825260209091200154606780546001600160a01b039092169183908110611cca57611cca612338565b600091825260209091200180546001600160a01b0319166001600160a01b039290921691909117905580611cfd81612364565b915050611c6f565b506067611d1360018661248b565b81548110611d2357611d23612338565b600091825260209091200180546001600160a01b03191690556067805480611d4d57611d4d612577565b60008281526020808220830160001990810180546001600160a01b03191690559092019092553382526066905260408120611d8791611ff2565b506001979650505050505050565b6040516001600160a01b0383811660248301528281166044830152600091829182919087169060640160408051601f198184030181529181526020820180516001600160e01b0316636eb1769f60e11b17905251611df3919061237d565b600060405180830381855afa9150503d8060008114611e2e576040519150601f19603f3d011682016040523d82523d6000602084013e611e33565b606091505b509150915081611e855760405162461bcd60e51b815260206004820152601b60248201527f73746174696363616c6c20616c6c6f77616e6365206661696c656400000000006044820152606401610232565b600081806020019051810190611e9b919061258d565b979650505050505050565b6033546001600160a01b031633146109d35760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610232565b603380546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b828054828255906000526020600020908101928215611fa7579160200282015b82811115611fa757825182546001600160a01b0319166001600160a01b03909116178255602090920191600190910190611f72565b50611fb392915061200c565b5090565b828054828255906000526020600020908101928215611fa7579160200282015b82811115611fa7578251825591602001919060010190611fd7565b508054600082559060005260206000209081019061182a91905b5b80821115611fb3576000815560010161200d565b6001600160a01b038116811461182a57600080fd5b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561207557612075612036565b604052919050565b600067ffffffffffffffff82111561209757612097612036565b5060051b60200190565b600082601f8301126120b257600080fd5b813560206120c76120c28361207d565b61204c565b82815260059290921b840181019181810190868411156120e657600080fd5b8286015b8481101561210157803583529183019183016120ea565b509695505050505050565b600080600080600060a0868803121561212457600080fd5b8535945060208087013561213781612021565b945060408701359350606087013567ffffffffffffffff8082111561215b57600080fd5b818901915089601f83011261216f57600080fd5b813561217d6120c28261207d565b81815260059190911b8301840190848101908c83111561219c57600080fd5b938501935b828510156121c35784356121b481612021565b825293850193908501906121a1565b9650505060808901359250808311156121db57600080fd5b50506121e9888289016120a1565b9150509295509295909350565b60006020828403121561220857600080fd5b5035919050565b6000806040838503121561222257600080fd5b50508035926020909101359150565b60006020828403121561224357600080fd5b813561224e81612021565b9392505050565b600081518084526020808501945080840160005b8381101561228e5781516001600160a01b031687529582019590820190600101612269565b509495945050505050565b60208152600061224e6020830184612255565b600081518084526020808501945080840160005b8381101561228e578151875295820195908201906001016122c0565b6001600160a01b038781168252861660208201526040810185905260c06060820181905260009061230f90830186612255565b828103608084015261232181866122ac565b9150508260000b60a0830152979650505050505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016123765761237661234e565b5060010190565b6000825160005b8181101561239e5760208186018101518583015201612384565b818111156123ad576000828501525b509190910192915050565b60208082526018908201527f63616c6c207472616e7366657246726f6d206661696c65640000000000000000604082015260600190565b60006020828403121561240157600080fd5b8151801515811461224e57600080fd5b60208082526025908201527f5468652072657475726e206f66207472616e7366657266726f6d206973206661604082015264696c75726560d81b606082015260800190565b83815260606020820152600061246f6060830185612255565b828103604084015261248181856122ac565b9695505050505050565b60008282101561249d5761249d61234e565b500390565b60208082526014908201527318d85b1b081d1c985b9cd9995c8819985a5b195960621b604082015260600190565b60208082526021908201527f5468652072657475726e206f66207472616e73666572206973206661696c75726040820152606560f81b606082015260800190565b60006020828403121561252357600080fd5b815161224e81612021565b84815283602082015260806040820152600061254d6080830185612255565b8281036060840152611e9b81856122ac565b600082198211156125725761257261234e565b500190565b634e487b7160e01b600052603160045260246000fd5b60006020828403121561259f57600080fd5b505191905056fea26469706673582212201576a3d28be5753bc6cc11a426e89a236aa799b4398b6927486000b64894aef764736f6c634300080d0033",
}

// DatumPayABI is the input ABI used to generate the binding from.
// Deprecated: Use DatumPayMetaData.ABI instead.
var DatumPayABI = DatumPayMetaData.ABI

// DatumPayBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DatumPayMetaData.Bin instead.
var DatumPayBin = DatumPayMetaData.Bin

// DeployDatumPay deploys a new Ethereum contract, binding an instance of DatumPay to it.
func DeployDatumPay(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DatumPay, error) {
	parsed, err := DatumPayMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DatumPayBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DatumPay{DatumPayCaller: DatumPayCaller{contract: contract}, DatumPayTransactor: DatumPayTransactor{contract: contract}, DatumPayFilterer: DatumPayFilterer{contract: contract}}, nil
}

// DatumPay is an auto generated Go binding around an Ethereum contract.
type DatumPay struct {
	DatumPayCaller     // Read-only binding to the contract
	DatumPayTransactor // Write-only binding to the contract
	DatumPayFilterer   // Log filterer for contract events
}

// DatumPayCaller is an auto generated read-only Go binding around an Ethereum contract.
type DatumPayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatumPayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DatumPayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatumPayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DatumPayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatumPaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DatumPaySession struct {
	Contract     *DatumPay         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DatumPayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DatumPayCallerSession struct {
	Contract *DatumPayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DatumPayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DatumPayTransactorSession struct {
	Contract     *DatumPayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DatumPayRaw is an auto generated low-level Go binding around an Ethereum contract.
type DatumPayRaw struct {
	Contract *DatumPay // Generic contract binding to access the raw methods on
}

// DatumPayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DatumPayCallerRaw struct {
	Contract *DatumPayCaller // Generic read-only contract binding to access the raw methods on
}

// DatumPayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DatumPayTransactorRaw struct {
	Contract *DatumPayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDatumPay creates a new instance of DatumPay, bound to a specific deployed contract.
func NewDatumPay(address common.Address, backend bind.ContractBackend) (*DatumPay, error) {
	contract, err := bindDatumPay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DatumPay{DatumPayCaller: DatumPayCaller{contract: contract}, DatumPayTransactor: DatumPayTransactor{contract: contract}, DatumPayFilterer: DatumPayFilterer{contract: contract}}, nil
}

// NewDatumPayCaller creates a new read-only instance of DatumPay, bound to a specific deployed contract.
func NewDatumPayCaller(address common.Address, caller bind.ContractCaller) (*DatumPayCaller, error) {
	contract, err := bindDatumPay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DatumPayCaller{contract: contract}, nil
}

// NewDatumPayTransactor creates a new write-only instance of DatumPay, bound to a specific deployed contract.
func NewDatumPayTransactor(address common.Address, transactor bind.ContractTransactor) (*DatumPayTransactor, error) {
	contract, err := bindDatumPay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DatumPayTransactor{contract: contract}, nil
}

// NewDatumPayFilterer creates a new log filterer instance of DatumPay, bound to a specific deployed contract.
func NewDatumPayFilterer(address common.Address, filterer bind.ContractFilterer) (*DatumPayFilterer, error) {
	contract, err := bindDatumPay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DatumPayFilterer{contract: contract}, nil
}

// bindDatumPay binds a generic wrapper to an already deployed contract.
func bindDatumPay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DatumPayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DatumPay *DatumPayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DatumPay.Contract.DatumPayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DatumPay *DatumPayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DatumPay.Contract.DatumPayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DatumPay *DatumPayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DatumPay.Contract.DatumPayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DatumPay *DatumPayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DatumPay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DatumPay *DatumPayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DatumPay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DatumPay *DatumPayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DatumPay.Contract.contract.Transact(opts, method, params...)
}

// GetTaskInfo is a free data retrieval call binding the contract method 0xd1a1b999.
//
// Solidity: function getTaskInfo(uint256 taskId) view returns(address, address, uint256, address[], uint256[], int8)
func (_DatumPay *DatumPayCaller) GetTaskInfo(opts *bind.CallOpts, taskId *big.Int) (common.Address, common.Address, *big.Int, []common.Address, []*big.Int, int8, error) {
	var out []interface{}
	err := _DatumPay.contract.Call(opts, &out, "getTaskInfo", taskId)

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
func (_DatumPay *DatumPaySession) GetTaskInfo(taskId *big.Int) (common.Address, common.Address, *big.Int, []common.Address, []*big.Int, int8, error) {
	return _DatumPay.Contract.GetTaskInfo(&_DatumPay.CallOpts, taskId)
}

// GetTaskInfo is a free data retrieval call binding the contract method 0xd1a1b999.
//
// Solidity: function getTaskInfo(uint256 taskId) view returns(address, address, uint256, address[], uint256[], int8)
func (_DatumPay *DatumPayCallerSession) GetTaskInfo(taskId *big.Int) (common.Address, common.Address, *big.Int, []common.Address, []*big.Int, int8, error) {
	return _DatumPay.Contract.GetTaskInfo(&_DatumPay.CallOpts, taskId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DatumPay *DatumPayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DatumPay.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DatumPay *DatumPaySession) Owner() (common.Address, error) {
	return _DatumPay.Contract.Owner(&_DatumPay.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DatumPay *DatumPayCallerSession) Owner() (common.Address, error) {
	return _DatumPay.Contract.Owner(&_DatumPay.CallOpts)
}

// TaskState is a free data retrieval call binding the contract method 0x26c6bee1.
//
// Solidity: function taskState(uint256 taskId) view returns(int8)
func (_DatumPay *DatumPayCaller) TaskState(opts *bind.CallOpts, taskId *big.Int) (int8, error) {
	var out []interface{}
	err := _DatumPay.contract.Call(opts, &out, "taskState", taskId)

	if err != nil {
		return *new(int8), err
	}

	out0 := *abi.ConvertType(out[0], new(int8)).(*int8)

	return out0, err

}

// TaskState is a free data retrieval call binding the contract method 0x26c6bee1.
//
// Solidity: function taskState(uint256 taskId) view returns(int8)
func (_DatumPay *DatumPaySession) TaskState(taskId *big.Int) (int8, error) {
	return _DatumPay.Contract.TaskState(&_DatumPay.CallOpts, taskId)
}

// TaskState is a free data retrieval call binding the contract method 0x26c6bee1.
//
// Solidity: function taskState(uint256 taskId) view returns(int8)
func (_DatumPay *DatumPayCallerSession) TaskState(taskId *big.Int) (int8, error) {
	return _DatumPay.Contract.TaskState(&_DatumPay.CallOpts, taskId)
}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address userAddress) view returns(address[])
func (_DatumPay *DatumPayCaller) Whitelist(opts *bind.CallOpts, userAddress common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _DatumPay.contract.Call(opts, &out, "whitelist", userAddress)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address userAddress) view returns(address[])
func (_DatumPay *DatumPaySession) Whitelist(userAddress common.Address) ([]common.Address, error) {
	return _DatumPay.Contract.Whitelist(&_DatumPay.CallOpts, userAddress)
}

// Whitelist is a free data retrieval call binding the contract method 0x9b19251a.
//
// Solidity: function whitelist(address userAddress) view returns(address[])
func (_DatumPay *DatumPayCallerSession) Whitelist(userAddress common.Address) ([]common.Address, error) {
	return _DatumPay.Contract.Whitelist(&_DatumPay.CallOpts, userAddress)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0xf80f5dd5.
//
// Solidity: function addWhitelist(address userAgency) returns(bool success)
func (_DatumPay *DatumPayTransactor) AddWhitelist(opts *bind.TransactOpts, userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "addWhitelist", userAgency)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0xf80f5dd5.
//
// Solidity: function addWhitelist(address userAgency) returns(bool success)
func (_DatumPay *DatumPaySession) AddWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.AddWhitelist(&_DatumPay.TransactOpts, userAgency)
}

// AddWhitelist is a paid mutator transaction binding the contract method 0xf80f5dd5.
//
// Solidity: function addWhitelist(address userAgency) returns(bool success)
func (_DatumPay *DatumPayTransactorSession) AddWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.AddWhitelist(&_DatumPay.TransactOpts, userAgency)
}

// Authorize is a paid mutator transaction binding the contract method 0xb6a5d7de.
//
// Solidity: function authorize(address userAgency) returns(bool success)
func (_DatumPay *DatumPayTransactor) Authorize(opts *bind.TransactOpts, userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "authorize", userAgency)
}

// Authorize is a paid mutator transaction binding the contract method 0xb6a5d7de.
//
// Solidity: function authorize(address userAgency) returns(bool success)
func (_DatumPay *DatumPaySession) Authorize(userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.Authorize(&_DatumPay.TransactOpts, userAgency)
}

// Authorize is a paid mutator transaction binding the contract method 0xb6a5d7de.
//
// Solidity: function authorize(address userAgency) returns(bool success)
func (_DatumPay *DatumPayTransactorSession) Authorize(userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.Authorize(&_DatumPay.TransactOpts, userAgency)
}

// DeleteWhitelist is a paid mutator transaction binding the contract method 0xff68f6af.
//
// Solidity: function deleteWhitelist(address userAgency) returns(bool success)
func (_DatumPay *DatumPayTransactor) DeleteWhitelist(opts *bind.TransactOpts, userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "deleteWhitelist", userAgency)
}

// DeleteWhitelist is a paid mutator transaction binding the contract method 0xff68f6af.
//
// Solidity: function deleteWhitelist(address userAgency) returns(bool success)
func (_DatumPay *DatumPaySession) DeleteWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.DeleteWhitelist(&_DatumPay.TransactOpts, userAgency)
}

// DeleteWhitelist is a paid mutator transaction binding the contract method 0xff68f6af.
//
// Solidity: function deleteWhitelist(address userAgency) returns(bool success)
func (_DatumPay *DatumPayTransactorSession) DeleteWhitelist(userAgency common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.DeleteWhitelist(&_DatumPay.TransactOpts, userAgency)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns()
func (_DatumPay *DatumPayTransactor) Initialize(opts *bind.TransactOpts, metisLat common.Address) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "initialize", metisLat)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns()
func (_DatumPay *DatumPaySession) Initialize(metisLat common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.Initialize(&_DatumPay.TransactOpts, metisLat)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address metisLat) returns()
func (_DatumPay *DatumPayTransactorSession) Initialize(metisLat common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.Initialize(&_DatumPay.TransactOpts, metisLat)
}

// Prepay is a paid mutator transaction binding the contract method 0x0ef92166.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList) returns(bool success)
func (_DatumPay *DatumPayTransactor) Prepay(opts *bind.TransactOpts, taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address, tokenValueList []*big.Int) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "prepay", taskId, user, fee, tokenAddressList, tokenValueList)
}

// Prepay is a paid mutator transaction binding the contract method 0x0ef92166.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList) returns(bool success)
func (_DatumPay *DatumPaySession) Prepay(taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address, tokenValueList []*big.Int) (*types.Transaction, error) {
	return _DatumPay.Contract.Prepay(&_DatumPay.TransactOpts, taskId, user, fee, tokenAddressList, tokenValueList)
}

// Prepay is a paid mutator transaction binding the contract method 0x0ef92166.
//
// Solidity: function prepay(uint256 taskId, address user, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList) returns(bool success)
func (_DatumPay *DatumPayTransactorSession) Prepay(taskId *big.Int, user common.Address, fee *big.Int, tokenAddressList []common.Address, tokenValueList []*big.Int) (*types.Transaction, error) {
	return _DatumPay.Contract.Prepay(&_DatumPay.TransactOpts, taskId, user, fee, tokenAddressList, tokenValueList)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DatumPay *DatumPayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DatumPay *DatumPaySession) RenounceOwnership() (*types.Transaction, error) {
	return _DatumPay.Contract.RenounceOwnership(&_DatumPay.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DatumPay *DatumPayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DatumPay.Contract.RenounceOwnership(&_DatumPay.TransactOpts)
}

// Settle is a paid mutator transaction binding the contract method 0x9a9c29f6.
//
// Solidity: function settle(uint256 taskId, uint256 fee) returns(bool success)
func (_DatumPay *DatumPayTransactor) Settle(opts *bind.TransactOpts, taskId *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "settle", taskId, fee)
}

// Settle is a paid mutator transaction binding the contract method 0x9a9c29f6.
//
// Solidity: function settle(uint256 taskId, uint256 fee) returns(bool success)
func (_DatumPay *DatumPaySession) Settle(taskId *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _DatumPay.Contract.Settle(&_DatumPay.TransactOpts, taskId, fee)
}

// Settle is a paid mutator transaction binding the contract method 0x9a9c29f6.
//
// Solidity: function settle(uint256 taskId, uint256 fee) returns(bool success)
func (_DatumPay *DatumPayTransactorSession) Settle(taskId *big.Int, fee *big.Int) (*types.Transaction, error) {
	return _DatumPay.Contract.Settle(&_DatumPay.TransactOpts, taskId, fee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DatumPay *DatumPayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DatumPay.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DatumPay *DatumPaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.TransferOwnership(&_DatumPay.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DatumPay *DatumPayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DatumPay.Contract.TransferOwnership(&_DatumPay.TransactOpts, newOwner)
}

// DatumPayInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DatumPay contract.
type DatumPayInitializedIterator struct {
	Event *DatumPayInitialized // Event containing the contract specifics and raw log

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
func (it *DatumPayInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatumPayInitialized)
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
		it.Event = new(DatumPayInitialized)
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
func (it *DatumPayInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatumPayInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatumPayInitialized represents a Initialized event raised by the DatumPay contract.
type DatumPayInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_DatumPay *DatumPayFilterer) FilterInitialized(opts *bind.FilterOpts) (*DatumPayInitializedIterator, error) {

	logs, sub, err := _DatumPay.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DatumPayInitializedIterator{contract: _DatumPay.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_DatumPay *DatumPayFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DatumPayInitialized) (event.Subscription, error) {

	logs, sub, err := _DatumPay.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatumPayInitialized)
				if err := _DatumPay.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_DatumPay *DatumPayFilterer) ParseInitialized(log types.Log) (*DatumPayInitialized, error) {
	event := new(DatumPayInitialized)
	if err := _DatumPay.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DatumPayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DatumPay contract.
type DatumPayOwnershipTransferredIterator struct {
	Event *DatumPayOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DatumPayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatumPayOwnershipTransferred)
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
		it.Event = new(DatumPayOwnershipTransferred)
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
func (it *DatumPayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatumPayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatumPayOwnershipTransferred represents a OwnershipTransferred event raised by the DatumPay contract.
type DatumPayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DatumPay *DatumPayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DatumPayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DatumPay.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DatumPayOwnershipTransferredIterator{contract: _DatumPay.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DatumPay *DatumPayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DatumPayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DatumPay.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatumPayOwnershipTransferred)
				if err := _DatumPay.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DatumPay *DatumPayFilterer) ParseOwnershipTransferred(log types.Log) (*DatumPayOwnershipTransferred, error) {
	event := new(DatumPayOwnershipTransferred)
	if err := _DatumPay.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DatumPayPrepayEventIterator is returned from FilterPrepayEvent and is used to iterate over the raw logs and unpacked data for PrepayEvent events raised by the DatumPay contract.
type DatumPayPrepayEventIterator struct {
	Event *DatumPayPrepayEvent // Event containing the contract specifics and raw log

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
func (it *DatumPayPrepayEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatumPayPrepayEvent)
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
		it.Event = new(DatumPayPrepayEvent)
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
func (it *DatumPayPrepayEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatumPayPrepayEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatumPayPrepayEvent represents a PrepayEvent event raised by the DatumPay contract.
type DatumPayPrepayEvent struct {
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
func (_DatumPay *DatumPayFilterer) FilterPrepayEvent(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*DatumPayPrepayEventIterator, error) {

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

	logs, sub, err := _DatumPay.contract.FilterLogs(opts, "PrepayEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &DatumPayPrepayEventIterator{contract: _DatumPay.contract, event: "PrepayEvent", logs: logs, sub: sub}, nil
}

// WatchPrepayEvent is a free log subscription operation binding the contract event 0x3a2d70b336733ffa4e257b7f871c733cc1cd2f1874bd1393a004e0ece7e3d55f.
//
// Solidity: event PrepayEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 fee, address[] tokenAddressList, uint256[] tokenValueList)
func (_DatumPay *DatumPayFilterer) WatchPrepayEvent(opts *bind.WatchOpts, sink chan<- *DatumPayPrepayEvent, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DatumPay.contract.WatchLogs(opts, "PrepayEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatumPayPrepayEvent)
				if err := _DatumPay.contract.UnpackLog(event, "PrepayEvent", log); err != nil {
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
func (_DatumPay *DatumPayFilterer) ParsePrepayEvent(log types.Log) (*DatumPayPrepayEvent, error) {
	event := new(DatumPayPrepayEvent)
	if err := _DatumPay.contract.UnpackLog(event, "PrepayEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DatumPaySettleEventIterator is returned from FilterSettleEvent and is used to iterate over the raw logs and unpacked data for SettleEvent events raised by the DatumPay contract.
type DatumPaySettleEventIterator struct {
	Event *DatumPaySettleEvent // Event containing the contract specifics and raw log

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
func (it *DatumPaySettleEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatumPaySettleEvent)
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
		it.Event = new(DatumPaySettleEvent)
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
func (it *DatumPaySettleEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatumPaySettleEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatumPaySettleEvent represents a SettleEvent event raised by the DatumPay contract.
type DatumPaySettleEvent struct {
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
func (_DatumPay *DatumPayFilterer) FilterSettleEvent(opts *bind.FilterOpts, taskId []*big.Int, user []common.Address, userAgency []common.Address) (*DatumPaySettleEventIterator, error) {

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

	logs, sub, err := _DatumPay.contract.FilterLogs(opts, "SettleEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return &DatumPaySettleEventIterator{contract: _DatumPay.contract, event: "SettleEvent", logs: logs, sub: sub}, nil
}

// WatchSettleEvent is a free log subscription operation binding the contract event 0x44b417d35556b6916cd7103ec4c060340e7f5e59b2f886680558e1c6d6366ef4.
//
// Solidity: event SettleEvent(uint256 indexed taskId, address indexed user, address indexed userAgency, uint256 Agencyfee, uint256 refundOrAdd, address[] tokenAddressList, uint256[] tokenValueList)
func (_DatumPay *DatumPayFilterer) WatchSettleEvent(opts *bind.WatchOpts, sink chan<- *DatumPaySettleEvent, taskId []*big.Int, user []common.Address, userAgency []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DatumPay.contract.WatchLogs(opts, "SettleEvent", taskIdRule, userRule, userAgencyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatumPaySettleEvent)
				if err := _DatumPay.contract.UnpackLog(event, "SettleEvent", log); err != nil {
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
func (_DatumPay *DatumPayFilterer) ParseSettleEvent(log types.Log) (*DatumPaySettleEvent, error) {
	event := new(DatumPaySettleEvent)
	if err := _DatumPay.contract.UnpackLog(event, "SettleEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
