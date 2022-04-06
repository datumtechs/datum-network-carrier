## KMS tool

### encrypt keystore file

./kmstool --kms-key-id 6a-4aa450a917 --kms-region-id cn-hangzhou --kms-access-key-id zEJNUex --kms-access-key-secret 7OuIh7 --src  ./conf/ori_keyfile.json encrypt

the output filename will be ./conf/ori_keyfile.json.enc

### decrypt the keystore file ciphertext

./kmstool --kms-key-id 6a-4aa450a917 --kms-region-id cn-hangzhou --kms-access-key-id zEJNUex --kms-access-key-secret 7OuIh7 --src  ./conf/ori_keyfile.json.enc decrypt

the output filename will be ./conf/ori_keyfile.json.enc.plain