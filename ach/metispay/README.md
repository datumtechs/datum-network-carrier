## KMS tool


### encrypt input key string

./kmstool --kms.key-id 6a-4aa450a917 --kms.region-id cn-hangzhou --kms.access-key-id zEJNUex --kms.access-key-secret 7OuIh7 --input 14754086774b445e8d0a631d2ade1e97a8ff0df33910cb0fcd93c8089553d4a1 encrypt

### decrypt input key ciphertext

./kmstool --kms.key-id 6a-4aa450a917 --kms.region-id cn-hangzhou --kms.access-key-id zEJNUex --kms.access-key-secret 7OuIh7 --input NTY3MWE1ZWQtZWFkOS00YjVkLTlmZDctYTlkODRjMTdlZTljgwN7xHHej7+Ht7RdUMEwjE8LhU1QqxNZduoNJvmB7QGi3nKK0UzYq+Oeb8zlMpxTg5RFOBMpZsArZORjwZaVjEmRgN3VWe2MhfivrtDTYo3UniSGpWFNDb213us= decrypt


### encrypt keystore file

./kmstool --kms.key-id 6a-4aa450a917 --kms.region-id cn-hangzhou --kms.access-key-id zEJNUex --kms.access-key-secret 7OuIh7 --keystore  ./conf/ori_keyfile.json encrypt

the output filename will be ./conf/ori_keyfile.json.enc

### decrypt the keystore file ciphertext

./kmstool --kms.key-id 6a-4aa450a917 --kms.region-id cn-hangzhou --kms.access-key-id zEJNUex --kms.access-key-secret 7OuIh7 --keystore  ./conf/ori_keyfile.json.enc decrypt

the output filename will be ./conf/ori_keyfile.json.enc.plain


