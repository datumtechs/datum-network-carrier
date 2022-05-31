package kms

type Config struct {
	KeyId           string `json:"KeyId"`
	RegionId        string `json:"RegionId"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
}

type KmsService interface {
	Encrypt(hexPriKey string) (string, error)
	Decrypt(hexEncryptedPriKey string) (string, error)
}
