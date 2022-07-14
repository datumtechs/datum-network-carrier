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

// Erc721MetaData contains all meta data concerning the Erc721 contract.
var Erc721MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cipherFlag\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"term\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"cipher_\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"tokenURI_\",\"type\":\"string\"}],\"name\":\"createToken\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getCharacter\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getExtInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"term\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"forEncryptAlg\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"proof_\",\"type\":\"string\"},{\"internalType\":\"uint8\",\"name\":\"cipherFlag_\",\"type\":\"uint8\"}],\"name\":\"initialize\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proof\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080604052600e805460ff60a81b191690553480156200001e57600080fd5b50604080518082018252600881526754656d706c61746560c01b60208083019182528351808501909452600e84526d15195b5c1b185d1954de5b589bdb60921b908401528151919291620000759160009162000094565b5080516200008b90600190602084019062000094565b50505062000176565b828054620000a2906200013a565b90600052602060002090601f016020900481019282620000c6576000855562000111565b82601f10620000e157805160ff191683800117855562000111565b8280016001018555821562000111579182015b8281111562000111578251825591602001919060010190620000f4565b506200011f92915062000123565b5090565b5b808211156200011f576000815560010162000124565b600181811c908216806200014f57607f821691505b6020821081036200017057634e487b7160e01b600052602260045260246000fd5b50919050565b6121fb80620001866000396000f3fe608060405234801561001057600080fd5b506004361061014d5760003560e01c80636eb5b233116100c3578063c0bb5c121161007c578063c0bb5c12146102c6578063c87b56dd146102e8578063dabb0531146102fb578063e985e9c51461031c578063f851a44014610358578063faf924cf1461036957600080fd5b80636eb5b2331461025357806370a082311461027257806378626f521461028557806395d89b4114610298578063a22cb465146102a0578063b88d4fde146102b357600080fd5b806318160ddd1161011557806318160ddd146101e257806323b872dd146101f45780632f745c591461020757806342842e0e1461021a5780634f6ccce71461022d5780636352211e1461024057600080fd5b806301ffc9a71461015257806306fdde031461017a578063081812fc1461018f578063095ea7b3146101ba57806313ebbf58146101cf575b600080fd5b610165610160366004611aec565b610371565b60405190151581526020015b60405180910390f35b610182610382565b6040516101719190611b61565b6101a261019d366004611b74565b610414565b6040516001600160a01b039091168152602001610171565b6101cd6101c8366004611ba9565b61043b565b005b6101656101dd366004611c1c565b610555565b6008545b604051908152602001610171565b6101cd610202366004611ce0565b610684565b6101e6610215366004611ba9565b6106b5565b6101cd610228366004611ce0565b61074b565b6101e661023b366004611b74565b610766565b6101a261024e366004611b74565b6107f9565b600e54600160a01b900460ff1660405160ff9091168152602001610171565b6101e6610280366004611d1c565b610859565b6101e6610293366004611df3565b6108df565b610182610af0565b6101cd6102ae366004611e67565b610aff565b6101cd6102c1366004611e9a565b610b0e565b6102d96102d4366004611b74565b610b46565b60405161017193929190611f16565b6101826102f6366004611b74565b610c08565b61030e610309366004611b74565b610c13565b604051610171929190611f4c565b61016561032a366004611f70565b6001600160a01b03918216600090815260056020908152604080832093909416825291909152205460ff1690565b600e546001600160a01b03166101a2565b610182610cc5565b600061037c82610cd4565b92915050565b6060600b805461039190611f9a565b80601f01602080910402602001604051908101604052809291908181526020018280546103bd90611f9a565b801561040a5780601f106103df5761010080835404028352916020019161040a565b820191906000526020600020905b8154815290600101906020018083116103ed57829003601f168201915b5050505050905090565b600061041f82610cf9565b506000908152600460205260409020546001600160a01b031690565b6000610446826107f9565b9050806001600160a01b0316836001600160a01b0316036104b85760405162461bcd60e51b815260206004820152602160248201527f4552433732313a20617070726f76616c20746f2063757272656e74206f776e656044820152603960f91b60648201526084015b60405180910390fd5b336001600160a01b03821614806104d457506104d4813361032a565b6105465760405162461bcd60e51b815260206004820152603e60248201527f4552433732313a20617070726f76652063616c6c6572206973206e6f7420746f60448201527f6b656e206f776e6572206e6f7220617070726f76656420666f7220616c6c000060648201526084016104af565b6105508383610d5b565b505050565b600e54600090600160a81b900460ff16156105cd5760405162461bcd60e51b815260206004820152603260248201527f45524337323154656d706c6174653a20746f6b656e20696e7374616e636520616044820152711b1c9958591e481a5b9a5d1a585b1a5e995960721b60648201526084016104af565b6106778989898080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050604080516020601f8d018190048102820181019092528b815292508b91508a908190840183828082843760009201919091525050604080516020601f8c018190048102820181019092528a815292508a9150899081908401838280828437600092019190915250899250610dc9915050565b9998505050505050505050565b61068e3382610ec0565b6106aa5760405162461bcd60e51b81526004016104af90611fd4565b610550838383610f3f565b60006106c083610859565b82106107225760405162461bcd60e51b815260206004820152602b60248201527f455243373231456e756d657261626c653a206f776e657220696e646578206f7560448201526a74206f6620626f756e647360a81b60648201526084016104af565b506001600160a01b03919091166000908152600660209081526040808320938352929052205490565b61055083838360405180602001604052806000815250610b0e565b600061077160085490565b82106107d45760405162461bcd60e51b815260206004820152602c60248201527f455243373231456e756d657261626c653a20676c6f62616c20696e646578206f60448201526b7574206f6620626f756e647360a01b60648201526084016104af565b600882815481106107e7576107e7612022565b90600052602060002001549050919050565b6000818152600260205260408120546001600160a01b03168061037c5760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b60448201526064016104af565b60006001600160a01b0382166108c35760405162461bcd60e51b815260206004820152602960248201527f4552433732313a2061646472657373207a65726f206973206e6f7420612076616044820152683634b21037bbb732b960b91b60648201526084016104af565b506001600160a01b031660009081526003602052604090205490565b600e546000906001600160a01b031633146109475760405162461bcd60e51b815260206004820152602260248201527f45524337323154656d706c6174653a20696e76616c6964206d73672073656e6460448201526132b960f11b60648201526084016104af565b82156109c657600e54600160a01b90046002908116146109c15760405162461bcd60e51b815260206004820152602f60248201527f6d6574616461746120646f6573206e6f7420737570706f72742063697068657260448201526e7465787420616c676f726974686d7360881b60648201526084016104af565b610a39565b600e54600160a01b9004600190811614610a395760405162461bcd60e51b815260206004820152602e60248201527f6d6574616461746120646f6573206e6f7420737570706f727420706c61696e7460448201526d65787420616c676f726974686d7360901b60648201526084016104af565b600f54604080518082018252868152851515602080830191909152600084815260108252929092208151805192939192610a769284920190611a3d565b50602091909101516001918201805460ff19169115159190911790556011805491820181556000527f31ecc21a745e3968a04e9570e4425bc18fa8019c68028196b546d1669c200c6801819055610acd33826110e6565b610ad78184611100565b600f54610ae590600161204e565b600f55949350505050565b6060600c805461039190611f9a565b610b0a33838361119a565b5050565b610b183383610ec0565b610b345760405162461bcd60e51b81526004016104af90611fd4565b610b4084848484611268565b50505050565b600060606000610b55846107f9565b60008581526010602052604090206001810154815460ff909116908290610b7b90611f9a565b80601f0160208091040260200160405190810160405280929190818152602001828054610ba790611f9a565b8015610bf45780601f10610bc957610100808354040283529160200191610bf4565b820191906000526020600020905b815481529060010190602001808311610bd757829003601f168201915b505050505091509250925092509193909250565b606061037c8261129b565b600081815260106020526040812060018101548154606093929160ff16908290610c3c90611f9a565b80601f0160208091040260200160405190810160405280929190818152602001828054610c6890611f9a565b8015610cb55780601f10610c8a57610100808354040283529160200191610cb5565b820191906000526020600020905b815481529060010190602001808311610c9857829003601f168201915b5050505050915091509150915091565b6060600d805461039190611f9a565b60006001600160e01b0319821663780e9d6360e01b148061037c575061037c826113a3565b6000818152600260205260409020546001600160a01b0316610d585760405162461bcd60e51b8152602060048201526018602482015277115490cdcc8c4e881a5b9d985b1a59081d1bdad95b88125160421b60448201526064016104af565b50565b600081815260046020526040902080546001600160a01b0319166001600160a01b0384169081179091558190610d90826107f9565b6001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45050565b60006001600160a01b038616610e355760405162461bcd60e51b815260206004820152602b60248201527f45524337323154656d706c6174653a20496e76616c69642061646d696e2c207a60448201526a65726f206164647265737360a81b60648201526084016104af565b8451610e4890600b906020880190611a3d565b508351610e5c90600c906020870190611a3d565b508251610e7090600d906020860190611a3d565b5050600e80546000600f55600160a81b6001600160a01b03979097166001600160a81b031990911617600160a01b60ff938416021760ff60a81b1916861790819055949094049093169392505050565b600080610ecc836107f9565b9050806001600160a01b0316846001600160a01b03161480610f1357506001600160a01b0380821660009081526005602090815260408083209388168352929052205460ff165b80610f375750836001600160a01b0316610f2c84610414565b6001600160a01b0316145b949350505050565b826001600160a01b0316610f52826107f9565b6001600160a01b031614610fb65760405162461bcd60e51b815260206004820152602560248201527f4552433732313a207472616e736665722066726f6d20696e636f72726563742060448201526437bbb732b960d91b60648201526084016104af565b6001600160a01b0382166110185760405162461bcd60e51b8152602060048201526024808201527f4552433732313a207472616e7366657220746f20746865207a65726f206164646044820152637265737360e01b60648201526084016104af565b6110238383836113f3565b61102e600082610d5b565b6001600160a01b0383166000908152600360205260408120805460019290611057908490612066565b90915550506001600160a01b038216600090815260036020526040812080546001929061108590849061204e565b909155505060008181526002602052604080822080546001600160a01b0319166001600160a01b0386811691821790925591518493918716917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4505050565b610b0a8282604051806020016040528060008152506113fe565b6000828152600260205260409020546001600160a01b031661117b5760405162461bcd60e51b815260206004820152602e60248201527f45524337323155524953746f726167653a2055524920736574206f66206e6f6e60448201526d32bc34b9ba32b73a103a37b5b2b760911b60648201526084016104af565b6000828152600a60209081526040909120825161055092840190611a3d565b816001600160a01b0316836001600160a01b0316036111fb5760405162461bcd60e51b815260206004820152601960248201527f4552433732313a20617070726f766520746f2063616c6c65720000000000000060448201526064016104af565b6001600160a01b03838116600081815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b611273848484610f3f565b61127f84848484611431565b610b405760405162461bcd60e51b81526004016104af9061207d565b60606112a682610cf9565b6000828152600a6020526040812080546112bf90611f9a565b80601f01602080910402602001604051908101604052809291908181526020018280546112eb90611f9a565b80156113385780601f1061130d57610100808354040283529160200191611338565b820191906000526020600020905b81548152906001019060200180831161131b57829003601f168201915b50505050509050600061135660408051602081019091526000815290565b90508051600003611368575092915050565b81511561139a5780826040516020016113829291906120cf565b60405160208183030381529060405292505050919050565b610f3784611532565b60006001600160e01b031982166380ac58cd60e01b14806113d457506001600160e01b03198216635b5e139f60e01b145b8061037c57506301ffc9a760e01b6001600160e01b031983161461037c565b6105508383836115a6565b611408838361165e565b6114156000848484611431565b6105505760405162461bcd60e51b81526004016104af9061207d565b60006001600160a01b0384163b1561152757604051630a85bd0160e11b81526001600160a01b0385169063150b7a02906114759033908990889088906004016120fe565b6020604051808303816000875af19250505080156114b0575060408051601f3d908101601f191682019092526114ad9181019061213b565b60015b61150d573d8080156114de576040519150601f19603f3d011682016040523d82523d6000602084013e6114e3565b606091505b5080516000036115055760405162461bcd60e51b81526004016104af9061207d565b805181602001fd5b6001600160e01b031916630a85bd0160e11b149050610f37565b506001949350505050565b606061153d82610cf9565b600061155460408051602081019091526000815290565b90506000815111611574576040518060200160405280600081525061159f565b8061157e846117ac565b60405160200161158f9291906120cf565b6040516020818303038152906040525b9392505050565b6001600160a01b038316611601576115fc81600880546000838152600960205260408120829055600182018355919091527ff3f7a9fe364faab93b216da50a3214154f22a0a2b415b23a84c8169e8b636ee30155565b611624565b816001600160a01b0316836001600160a01b0316146116245761162483826118ad565b6001600160a01b03821661163b576105508161194a565b826001600160a01b0316826001600160a01b0316146105505761055082826119f9565b6001600160a01b0382166116b45760405162461bcd60e51b815260206004820181905260248201527f4552433732313a206d696e7420746f20746865207a65726f206164647265737360448201526064016104af565b6000818152600260205260409020546001600160a01b0316156117195760405162461bcd60e51b815260206004820152601c60248201527f4552433732313a20746f6b656e20616c7265616479206d696e7465640000000060448201526064016104af565b611725600083836113f3565b6001600160a01b038216600090815260036020526040812080546001929061174e90849061204e565b909155505060008181526002602052604080822080546001600160a01b0319166001600160a01b03861690811790915590518392907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef908290a45050565b6060816000036117d35750506040805180820190915260018152600360fc1b602082015290565b8160005b81156117fd57806117e781612158565b91506117f69050600a83612187565b91506117d7565b60008167ffffffffffffffff81111561181857611818611d37565b6040519080825280601f01601f191660200182016040528015611842576020820181803683370190505b5090505b8415610f3757611857600183612066565b9150611864600a8661219b565b61186f90603061204e565b60f81b81838151811061188457611884612022565b60200101906001600160f81b031916908160001a9053506118a6600a86612187565b9450611846565b600060016118ba84610859565b6118c49190612066565b600083815260076020526040902054909150808214611917576001600160a01b03841660009081526006602090815260408083208584528252808320548484528184208190558352600790915290208190555b5060009182526007602090815260408084208490556001600160a01b039094168352600681528383209183525290812055565b60085460009061195c90600190612066565b6000838152600960205260408120546008805493945090928490811061198457611984612022565b9060005260206000200154905080600883815481106119a5576119a5612022565b60009182526020808320909101929092558281526009909152604080822084905585825281205560088054806119dd576119dd6121af565b6001900381819060005260206000200160009055905550505050565b6000611a0483610859565b6001600160a01b039093166000908152600660209081526040808320868452825280832085905593825260079052919091209190915550565b828054611a4990611f9a565b90600052602060002090601f016020900481019282611a6b5760008555611ab1565b82601f10611a8457805160ff1916838001178555611ab1565b82800160010185558215611ab1579182015b82811115611ab1578251825591602001919060010190611a96565b50611abd929150611ac1565b5090565b5b80821115611abd5760008155600101611ac2565b6001600160e01b031981168114610d5857600080fd5b600060208284031215611afe57600080fd5b813561159f81611ad6565b60005b83811015611b24578181015183820152602001611b0c565b83811115610b405750506000910152565b60008151808452611b4d816020860160208601611b09565b601f01601f19169290920160200192915050565b60208152600061159f6020830184611b35565b600060208284031215611b8657600080fd5b5035919050565b80356001600160a01b0381168114611ba457600080fd5b919050565b60008060408385031215611bbc57600080fd5b611bc583611b8d565b946020939093013593505050565b60008083601f840112611be557600080fd5b50813567ffffffffffffffff811115611bfd57600080fd5b602083019150836020828501011115611c1557600080fd5b9250929050565b60008060008060008060008060a0898b031215611c3857600080fd5b611c4189611b8d565b9750602089013567ffffffffffffffff80821115611c5e57600080fd5b611c6a8c838d01611bd3565b909950975060408b0135915080821115611c8357600080fd5b611c8f8c838d01611bd3565b909750955060608b0135915080821115611ca857600080fd5b50611cb58b828c01611bd3565b909450925050608089013560ff81168114611ccf57600080fd5b809150509295985092959890939650565b600080600060608486031215611cf557600080fd5b611cfe84611b8d565b9250611d0c60208501611b8d565b9150604084013590509250925092565b600060208284031215611d2e57600080fd5b61159f82611b8d565b634e487b7160e01b600052604160045260246000fd5b600067ffffffffffffffff80841115611d6857611d68611d37565b604051601f8501601f19908116603f01168101908282118183101715611d9057611d90611d37565b81604052809350858152868686011115611da957600080fd5b858560208301376000602087830101525050509392505050565b600082601f830112611dd457600080fd5b61159f83833560208501611d4d565b80358015158114611ba457600080fd5b600080600060608486031215611e0857600080fd5b833567ffffffffffffffff80821115611e2057600080fd5b611e2c87838801611dc3565b9450611e3a60208701611de3565b93506040860135915080821115611e5057600080fd5b50611e5d86828701611dc3565b9150509250925092565b60008060408385031215611e7a57600080fd5b611e8383611b8d565b9150611e9160208401611de3565b90509250929050565b60008060008060808587031215611eb057600080fd5b611eb985611b8d565b9350611ec760208601611b8d565b925060408501359150606085013567ffffffffffffffff811115611eea57600080fd5b8501601f81018713611efb57600080fd5b611f0a87823560208401611d4d565b91505092959194509250565b6001600160a01b0384168152606060208201819052600090611f3a90830185611b35565b90508215156040830152949350505050565b604081526000611f5f6040830185611b35565b905082151560208301529392505050565b60008060408385031215611f8357600080fd5b611f8c83611b8d565b9150611e9160208401611b8d565b600181811c90821680611fae57607f821691505b602082108103611fce57634e487b7160e01b600052602260045260246000fd5b50919050565b6020808252602e908201527f4552433732313a2063616c6c6572206973206e6f7420746f6b656e206f776e6560408201526d1c881b9bdc88185c1c1c9bdd995960921b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000821982111561206157612061612038565b500190565b60008282101561207857612078612038565b500390565b60208082526032908201527f4552433732313a207472616e7366657220746f206e6f6e20455243373231526560408201527131b2b4bb32b91034b6b83632b6b2b73a32b960711b606082015260800190565b600083516120e1818460208801611b09565b8351908301906120f5818360208801611b09565b01949350505050565b6001600160a01b038581168252841660208201526040810183905260806060820181905260009061213190830184611b35565b9695505050505050565b60006020828403121561214d57600080fd5b815161159f81611ad6565b60006001820161216a5761216a612038565b5060010190565b634e487b7160e01b600052601260045260246000fd5b60008261219657612196612171565b500490565b6000826121aa576121aa612171565b500690565b634e487b7160e01b600052603160045260246000fdfea26469706673582212209783756193f05996ef77e7d7ebae6efb0ff92282562eca59713a4048c33b60bb64736f6c634300080d0033",
}

// Erc721ABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc721MetaData.ABI instead.
var Erc721ABI = Erc721MetaData.ABI

// Erc721Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Erc721MetaData.Bin instead.
var Erc721Bin = Erc721MetaData.Bin

// DeployErc721 deploys a new Ethereum contract, binding an instance of Erc721 to it.
func DeployErc721(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Erc721, error) {
	parsed, err := Erc721MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Erc721Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc721{Erc721Caller: Erc721Caller{contract: contract}, Erc721Transactor: Erc721Transactor{contract: contract}, Erc721Filterer: Erc721Filterer{contract: contract}}, nil
}

// Erc721 is an auto generated Go binding around an Ethereum contract.
type Erc721 struct {
	Erc721Caller     // Read-only binding to the contract
	Erc721Transactor // Write-only binding to the contract
	Erc721Filterer   // Log filterer for contract events
}

// Erc721Caller is an auto generated read-only Go binding around an Ethereum contract.
type Erc721Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc721Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc721Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc721Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc721Session struct {
	Contract     *Erc721           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc721CallerSession struct {
	Contract *Erc721Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Erc721TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc721TransactorSession struct {
	Contract     *Erc721Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc721Raw is an auto generated low-level Go binding around an Ethereum contract.
type Erc721Raw struct {
	Contract *Erc721 // Generic contract binding to access the raw methods on
}

// Erc721CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc721CallerRaw struct {
	Contract *Erc721Caller // Generic read-only contract binding to access the raw methods on
}

// Erc721TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc721TransactorRaw struct {
	Contract *Erc721Transactor // Generic write-only contract binding to access the raw methods on
}

// NewErc721 creates a new instance of Erc721, bound to a specific deployed contract.
func NewErc721(address common.Address, backend bind.ContractBackend) (*Erc721, error) {
	contract, err := bindErc721(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc721{Erc721Caller: Erc721Caller{contract: contract}, Erc721Transactor: Erc721Transactor{contract: contract}, Erc721Filterer: Erc721Filterer{contract: contract}}, nil
}

// NewErc721Caller creates a new read-only instance of Erc721, bound to a specific deployed contract.
func NewErc721Caller(address common.Address, caller bind.ContractCaller) (*Erc721Caller, error) {
	contract, err := bindErc721(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721Caller{contract: contract}, nil
}

// NewErc721Transactor creates a new write-only instance of Erc721, bound to a specific deployed contract.
func NewErc721Transactor(address common.Address, transactor bind.ContractTransactor) (*Erc721Transactor, error) {
	contract, err := bindErc721(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc721Transactor{contract: contract}, nil
}

// NewErc721Filterer creates a new log filterer instance of Erc721, bound to a specific deployed contract.
func NewErc721Filterer(address common.Address, filterer bind.ContractFilterer) (*Erc721Filterer, error) {
	contract, err := bindErc721(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc721Filterer{contract: contract}, nil
}

// bindErc721 binds a generic wrapper to an already deployed contract.
func bindErc721(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc721ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721 *Erc721Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721.Contract.Erc721Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721 *Erc721Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721.Contract.Erc721Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721 *Erc721Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721.Contract.Erc721Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc721 *Erc721CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc721.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc721 *Erc721TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc721.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc721 *Erc721TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc721.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Erc721 *Erc721Caller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Erc721 *Erc721Session) Admin() (common.Address, error) {
	return _Erc721.Contract.Admin(&_Erc721.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Erc721 *Erc721CallerSession) Admin() (common.Address, error) {
	return _Erc721.Contract.Admin(&_Erc721.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Erc721 *Erc721Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Erc721 *Erc721Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Erc721.Contract.BalanceOf(&_Erc721.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_Erc721 *Erc721CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _Erc721.Contract.BalanceOf(&_Erc721.CallOpts, owner)
}

// CipherFlag is a free data retrieval call binding the contract method 0x6eb5b233.
//
// Solidity: function cipherFlag() view returns(uint8)
func (_Erc721 *Erc721Caller) CipherFlag(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "cipherFlag")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CipherFlag is a free data retrieval call binding the contract method 0x6eb5b233.
//
// Solidity: function cipherFlag() view returns(uint8)
func (_Erc721 *Erc721Session) CipherFlag() (uint8, error) {
	return _Erc721.Contract.CipherFlag(&_Erc721.CallOpts)
}

// CipherFlag is a free data retrieval call binding the contract method 0x6eb5b233.
//
// Solidity: function cipherFlag() view returns(uint8)
func (_Erc721 *Erc721CallerSession) CipherFlag() (uint8, error) {
	return _Erc721.Contract.CipherFlag(&_Erc721.CallOpts)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Caller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Session) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.GetApproved(&_Erc721.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721CallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.GetApproved(&_Erc721.CallOpts, tokenId)
}

// GetCharacter is a free data retrieval call binding the contract method 0xdabb0531.
//
// Solidity: function getCharacter(uint256 tokenId) view returns(string, bool)
func (_Erc721 *Erc721Caller) GetCharacter(opts *bind.CallOpts, tokenId *big.Int) (string, bool, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getCharacter", tokenId)

	if err != nil {
		return *new(string), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err

}

// GetCharacter is a free data retrieval call binding the contract method 0xdabb0531.
//
// Solidity: function getCharacter(uint256 tokenId) view returns(string, bool)
func (_Erc721 *Erc721Session) GetCharacter(tokenId *big.Int) (string, bool, error) {
	return _Erc721.Contract.GetCharacter(&_Erc721.CallOpts, tokenId)
}

// GetCharacter is a free data retrieval call binding the contract method 0xdabb0531.
//
// Solidity: function getCharacter(uint256 tokenId) view returns(string, bool)
func (_Erc721 *Erc721CallerSession) GetCharacter(tokenId *big.Int) (string, bool, error) {
	return _Erc721.Contract.GetCharacter(&_Erc721.CallOpts, tokenId)
}

// GetExtInfo is a free data retrieval call binding the contract method 0xc0bb5c12.
//
// Solidity: function getExtInfo(uint256 tokenId) view returns(address owner, string term, bool forEncryptAlg)
func (_Erc721 *Erc721Caller) GetExtInfo(opts *bind.CallOpts, tokenId *big.Int) (struct {
	Owner         common.Address
	Term          string
	ForEncryptAlg bool
}, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "getExtInfo", tokenId)

	outstruct := new(struct {
		Owner         common.Address
		Term          string
		ForEncryptAlg bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Term = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.ForEncryptAlg = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// GetExtInfo is a free data retrieval call binding the contract method 0xc0bb5c12.
//
// Solidity: function getExtInfo(uint256 tokenId) view returns(address owner, string term, bool forEncryptAlg)
func (_Erc721 *Erc721Session) GetExtInfo(tokenId *big.Int) (struct {
	Owner         common.Address
	Term          string
	ForEncryptAlg bool
}, error) {
	return _Erc721.Contract.GetExtInfo(&_Erc721.CallOpts, tokenId)
}

// GetExtInfo is a free data retrieval call binding the contract method 0xc0bb5c12.
//
// Solidity: function getExtInfo(uint256 tokenId) view returns(address owner, string term, bool forEncryptAlg)
func (_Erc721 *Erc721CallerSession) GetExtInfo(tokenId *big.Int) (struct {
	Owner         common.Address
	Term          string
	ForEncryptAlg bool
}, error) {
	return _Erc721.Contract.GetExtInfo(&_Erc721.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721 *Erc721Caller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721 *Erc721Session) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Erc721.Contract.IsApprovedForAll(&_Erc721.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_Erc721 *Erc721CallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _Erc721.Contract.IsApprovedForAll(&_Erc721.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721 *Erc721Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721 *Erc721Session) Name() (string, error) {
	return _Erc721.Contract.Name(&_Erc721.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Erc721 *Erc721CallerSession) Name() (string, error) {
	return _Erc721.Contract.Name(&_Erc721.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Caller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721Session) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.OwnerOf(&_Erc721.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_Erc721 *Erc721CallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _Erc721.Contract.OwnerOf(&_Erc721.CallOpts, tokenId)
}

// Proof is a free data retrieval call binding the contract method 0xfaf924cf.
//
// Solidity: function proof() view returns(string)
func (_Erc721 *Erc721Caller) Proof(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "proof")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Proof is a free data retrieval call binding the contract method 0xfaf924cf.
//
// Solidity: function proof() view returns(string)
func (_Erc721 *Erc721Session) Proof() (string, error) {
	return _Erc721.Contract.Proof(&_Erc721.CallOpts)
}

// Proof is a free data retrieval call binding the contract method 0xfaf924cf.
//
// Solidity: function proof() view returns(string)
func (_Erc721 *Erc721CallerSession) Proof() (string, error) {
	return _Erc721.Contract.Proof(&_Erc721.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721 *Erc721Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721 *Erc721Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Erc721.Contract.SupportsInterface(&_Erc721.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Erc721 *Erc721CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Erc721.Contract.SupportsInterface(&_Erc721.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721 *Erc721Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721 *Erc721Session) Symbol() (string, error) {
	return _Erc721.Contract.Symbol(&_Erc721.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Erc721 *Erc721CallerSession) Symbol() (string, error) {
	return _Erc721.Contract.Symbol(&_Erc721.CallOpts)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Erc721 *Erc721Caller) TokenByIndex(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "tokenByIndex", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Erc721 *Erc721Session) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _Erc721.Contract.TokenByIndex(&_Erc721.CallOpts, index)
}

// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
//
// Solidity: function tokenByIndex(uint256 index) view returns(uint256)
func (_Erc721 *Erc721CallerSession) TokenByIndex(index *big.Int) (*big.Int, error) {
	return _Erc721.Contract.TokenByIndex(&_Erc721.CallOpts, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Erc721 *Erc721Caller) TokenOfOwnerByIndex(opts *bind.CallOpts, owner common.Address, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "tokenOfOwnerByIndex", owner, index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Erc721 *Erc721Session) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _Erc721.Contract.TokenOfOwnerByIndex(&_Erc721.CallOpts, owner, index)
}

// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
//
// Solidity: function tokenOfOwnerByIndex(address owner, uint256 index) view returns(uint256)
func (_Erc721 *Erc721CallerSession) TokenOfOwnerByIndex(owner common.Address, index *big.Int) (*big.Int, error) {
	return _Erc721.Contract.TokenOfOwnerByIndex(&_Erc721.CallOpts, owner, index)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721 *Erc721Caller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721 *Erc721Session) TokenURI(tokenId *big.Int) (string, error) {
	return _Erc721.Contract.TokenURI(&_Erc721.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Erc721 *Erc721CallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Erc721.Contract.TokenURI(&_Erc721.CallOpts, tokenId)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Erc721 *Erc721Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Erc721.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Erc721 *Erc721Session) TotalSupply() (*big.Int, error) {
	return _Erc721.Contract.TotalSupply(&_Erc721.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Erc721 *Erc721CallerSession) TotalSupply() (*big.Int, error) {
	return _Erc721.Contract.TotalSupply(&_Erc721.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Transactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Session) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Approve(&_Erc721.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_Erc721 *Erc721TransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.Approve(&_Erc721.TransactOpts, to, tokenId)
}

// CreateToken is a paid mutator transaction binding the contract method 0x78626f52.
//
// Solidity: function createToken(string term, bool cipher_, string tokenURI_) returns(uint256)
func (_Erc721 *Erc721Transactor) CreateToken(opts *bind.TransactOpts, term string, cipher_ bool, tokenURI_ string) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "createToken", term, cipher_, tokenURI_)
}

// CreateToken is a paid mutator transaction binding the contract method 0x78626f52.
//
// Solidity: function createToken(string term, bool cipher_, string tokenURI_) returns(uint256)
func (_Erc721 *Erc721Session) CreateToken(term string, cipher_ bool, tokenURI_ string) (*types.Transaction, error) {
	return _Erc721.Contract.CreateToken(&_Erc721.TransactOpts, term, cipher_, tokenURI_)
}

// CreateToken is a paid mutator transaction binding the contract method 0x78626f52.
//
// Solidity: function createToken(string term, bool cipher_, string tokenURI_) returns(uint256)
func (_Erc721 *Erc721TransactorSession) CreateToken(term string, cipher_ bool, tokenURI_ string) (*types.Transaction, error) {
	return _Erc721.Contract.CreateToken(&_Erc721.TransactOpts, term, cipher_, tokenURI_)
}

// Initialize is a paid mutator transaction binding the contract method 0x13ebbf58.
//
// Solidity: function initialize(address admin_, string name_, string symbol_, string proof_, uint8 cipherFlag_) returns(bool)
func (_Erc721 *Erc721Transactor) Initialize(opts *bind.TransactOpts, admin_ common.Address, name_ string, symbol_ string, proof_ string, cipherFlag_ uint8) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "initialize", admin_, name_, symbol_, proof_, cipherFlag_)
}

// Initialize is a paid mutator transaction binding the contract method 0x13ebbf58.
//
// Solidity: function initialize(address admin_, string name_, string symbol_, string proof_, uint8 cipherFlag_) returns(bool)
func (_Erc721 *Erc721Session) Initialize(admin_ common.Address, name_ string, symbol_ string, proof_ string, cipherFlag_ uint8) (*types.Transaction, error) {
	return _Erc721.Contract.Initialize(&_Erc721.TransactOpts, admin_, name_, symbol_, proof_, cipherFlag_)
}

// Initialize is a paid mutator transaction binding the contract method 0x13ebbf58.
//
// Solidity: function initialize(address admin_, string name_, string symbol_, string proof_, uint8 cipherFlag_) returns(bool)
func (_Erc721 *Erc721TransactorSession) Initialize(admin_ common.Address, name_ string, symbol_ string, proof_ string, cipherFlag_ uint8) (*types.Transaction, error) {
	return _Erc721.Contract.Initialize(&_Erc721.TransactOpts, admin_, name_, symbol_, proof_, cipherFlag_)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Session) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721TransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721 *Erc721Transactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721 *Erc721Session) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom0(&_Erc721.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_Erc721 *Erc721TransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _Erc721.Contract.SafeTransferFrom0(&_Erc721.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721 *Erc721Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721 *Erc721Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721.Contract.SetApprovalForAll(&_Erc721.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Erc721 *Erc721TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Erc721.Contract.SetApprovalForAll(&_Erc721.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721Session) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.TransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_Erc721 *Erc721TransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _Erc721.Contract.TransferFrom(&_Erc721.TransactOpts, from, to, tokenId)
}

// Erc721ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Erc721 contract.
type Erc721ApprovalIterator struct {
	Event *Erc721Approval // Event containing the contract specifics and raw log

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
func (it *Erc721ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721Approval)
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
		it.Event = new(Erc721Approval)
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
func (it *Erc721ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721Approval represents a Approval event raised by the Erc721 contract.
type Erc721Approval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*Erc721ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721ApprovalIterator{contract: _Erc721.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Erc721Approval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721Approval)
				if err := _Erc721.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) ParseApproval(log types.Log) (*Erc721Approval, error) {
	event := new(Erc721Approval)
	if err := _Erc721.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Erc721 contract.
type Erc721ApprovalForAllIterator struct {
	Event *Erc721ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *Erc721ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721ApprovalForAll)
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
		it.Event = new(Erc721ApprovalForAll)
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
func (it *Erc721ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721ApprovalForAll represents a ApprovalForAll event raised by the Erc721 contract.
type Erc721ApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721 *Erc721Filterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*Erc721ApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &Erc721ApprovalForAllIterator{contract: _Erc721.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721 *Erc721Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Erc721ApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721ApprovalForAll)
				if err := _Erc721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_Erc721 *Erc721Filterer) ParseApprovalForAll(log types.Log) (*Erc721ApprovalForAll, error) {
	event := new(Erc721ApprovalForAll)
	if err := _Erc721.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc721TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Erc721 contract.
type Erc721TransferIterator struct {
	Event *Erc721Transfer // Event containing the contract specifics and raw log

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
func (it *Erc721TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc721Transfer)
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
		it.Event = new(Erc721Transfer)
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
func (it *Erc721TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc721TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc721Transfer represents a Transfer event raised by the Erc721 contract.
type Erc721Transfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*Erc721TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &Erc721TransferIterator{contract: _Erc721.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Erc721Transfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _Erc721.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc721Transfer)
				if err := _Erc721.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_Erc721 *Erc721Filterer) ParseTransfer(log types.Log) (*Erc721Transfer, error) {
	event := new(Erc721Transfer)
	if err := _Erc721.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
