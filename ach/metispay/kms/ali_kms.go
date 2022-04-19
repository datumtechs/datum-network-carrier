package kms

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

const domainPattern = "kms.%s.aliyuncs.com"

type AliKms struct {
	Config *Config
}

func (k *AliKms) Encrypt(plain string) (string, error) {
	client, err := sdk.NewClientWithAccessKey(k.Config.RegionId, k.Config.AccessKeyId, k.Config.AccessKeySecret)
	if err != nil {
		return "", errors.Wrapf(err, "new KMS client")
	}

	domain := fmt.Sprintf(domainPattern, k.Config.RegionId)

	request := requests.NewCommonRequest() // Make a common request
	request.Method = "GET"                 // Set request method
	request.Product = "KMS"                // Specify product
	request.Domain = domain                // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2016-01-20"         // Specify product version
	request.Scheme = "https"               // Set request scheme. Default: http
	request.ApiName = "Encrypt"            // Specify product interface
	request.AcceptFormat = "json"
	request.QueryParams["Action"] = "Encrypt"     // Assign values to parameters in the path
	request.QueryParams["KeyId"] = k.Config.KeyId // Assign values to parameters in the path

	//request.QueryParams["Plaintext"] = base64.StdEncoding.EncodeToString(jsonBytes)
	request.QueryParams["Plaintext"] = plain // Specify the requested regionId, if not specified, use the client regionId, then default regionId

	request.QueryParams["RequestId"] = uuid.New() // Specify the requested regionId, if not specified, use the client regionId, then default regionId
	commonResponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return "", errors.Wrapf(err, "call KMS")
	}

	resp := new(kms.EncryptResponse)
	err = json.Unmarshal(commonResponse.GetHttpContentBytes(), &resp)

	if err != nil {
		return "", errors.Wrapf(err, "unmarshal KMS response")
	}

	return resp.CiphertextBlob, nil
}

func (k *AliKms) Decrypt(encryptedKeyStore string) (string, error) {
	domain := fmt.Sprintf(domainPattern, k.Config.RegionId)

	client, err := sdk.NewClientWithAccessKey(k.Config.RegionId, k.Config.AccessKeyId, k.Config.AccessKeySecret)
	if err != nil {
		return "", errors.Wrapf(err, "new KMS client")
	}
	request := requests.NewCommonRequest() // Make a common request
	request.Method = "GET"                 // Set request method
	request.Product = "KMS"                // Specify product
	request.Domain = domain                // Location Service will not be enabled if the host is specified. For example, service with a Certification type-Bearer Token should be specified
	request.Version = "2016-01-20"         // Specify product version
	request.Scheme = "https"               // Set request scheme. Default: http
	request.ApiName = "Decrypt"            // Specify product interface
	request.AcceptFormat = "json"
	request.QueryParams["Action"] = "Decrypt"     // Assign values to parameters in the path
	request.QueryParams["KeyId"] = k.Config.KeyId // Assign values to parameters in the path
	//request.QueryParams["keyVersionId"] = k.Config.KeyVersionId
	request.QueryParams["CiphertextBlob"] = encryptedKeyStore // Specify the requested regionId, if not specified, use the client regionId, then default regionId
	request.QueryParams["RequestId"] = uuid.New()             // Specify the requested regionId, if not specified, use the client regionId, then default regionId
	commonResponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		return "", errors.Wrapf(err, "call KMS")
	}

	resp := new(kms.DecryptResponse)
	err = json.Unmarshal(commonResponse.GetHttpContentBytes(), &resp)

	if err != nil {
		return "", errors.Wrapf(err, "unmarshal KMS response")
	}

	return resp.Plaintext, nil
}
