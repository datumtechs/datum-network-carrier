package kms

type Config struct {
	KeyId           string `json:"KeyId"`
	RegionId        string `json:"RegionId"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
}

type KmsService interface {
	Encrypt(jsonBytes []byte) ([]byte, error)
	Decrypt(encryptedKeyStore []byte) ([]byte, error)
}
