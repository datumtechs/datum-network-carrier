solc version: 0.8.13

abigen version: go-ethereum branch v1.10.3

in go-ethereum/build/bin folder

npm install @openzeppelin/contracts-upgradeable
npm install @openzeppelin/contracts
 
1. generate abi/bin:
solc --include-path ./node_modules/ --base-path . --optimize --bin --abi --overwrite -o ./ Token20Pay.sol

2. generate contract wrapper in golang
./abigen --abi=Token20Pay.abi --bin=Token20Pay.bin --pkg=contracts --type=Token20Pay --out=Token20Pay.go
