package signsuite

import (
	"bytes"
	"encoding/hex"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/bytesutil"
	"github.com/datumtechs/datum-network-carrier/common/rlputil"
	"github.com/datumtechs/datum-network-carrier/signsuite/eip712"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testTypedData = &eip712.TypedData{
	Domain: eip712.TypedDataDomain{
		Name: "Datum",
		//Version: "1.0",
	},
	Message: eip712.TypedDataMessage{
		"key":  "1",
		"desc": "Welcome to Datum!",
	},
	PrimaryType: "Login",
	Types: eip712.Types{
		"EIP712Domain": {
			{
				Name: "name",
				Type: "string",
			},
		},
		"Login": {
			{
				Name: "key",
				Type: "string",
			},
			{
				Name: "desc",
				Type: "string",
			},
		},
	},
}

func TestSender(t *testing.T) {
	{
		expected, err := hex.DecodeString("f93142f65cf427b089322f31f687373201fe59c289b9252e395931a22650c2035c1eaf90cfc822b1df745ae507b226b3a4189d22cedccbb0b25f3f3dc585491d1b")

		if err != nil {
			t.Fatal(err)
		}
		address, _, err := RecoverEIP712(expected, testTypedData)
		assert.Equal(t, address, "0xc115ceadf9e5923330e5f42903fe7f926dda65d2")

	}

	var buf bytes.Buffer
	buf.Write([]byte("111111111111111111111111"))
	buf.Write(bytesutil.Uint32ToBytes(uint32(1)))
	buf.Write([]byte("7e7afd35ff1ebe5db3a06759a2d42293b8fc73f2b1dfd7e16527bddca9ed9584"))
	buf.Write([]byte("1")) // null
	buf.Write(bytesutil.Uint32ToBytes(uint32(1)))
	buf.Write(bytesutil.Uint32ToBytes(uint32(1)))
	buf.Write([]byte("1"))
	buf.Write([]byte("1"))
	buf.Write([]byte("{\"metadataColumns\":[{\"index\":1,\"type\":\"string\",\"size\":0,\"name\":\"CLIENT_ID\"},{\"index\":2,\"type\":\"string\",\"size\":0,\"name\":\"DEFAULT\"},{\"index\":3,\"type\":\"string\",\"size\":0,\"name\":\"HOUSING\"},{\"index\":4,\"type\":\"string\",\"size\":0,\"name\":\"LOAN\"},{\"index\":5,\"type\":\"string\",\"size\":0,\"name\":\"CONTACT\"},{\"index\":6,\"type\":\"string\",\"size\":0,\"name\":\"CAMPAIGN\"},{\"index\":7,\"type\":\"string\",\"size\":0,\"name\":\"PDAYS\"},{\"index\":8,\"type\":\"string\",\"size\":0,\"name\":\"PREVIOUS\"},{\"index\":9,\"type\":\"string\",\"size\":0,\"name\":\"EURIBOR3M\"},{\"index\":10,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_APR\"},{\"index\":11,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_AUG\"},{\"index\":12,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_DEC\"},{\"index\":13,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_JUL\"},{\"index\":14,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_JUN\"},{\"index\":15,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_MAR\"},{\"index\":16,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_MAY\"},{\"index\":17,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_NOV\"},{\"index\":18,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_OCT\"},{\"index\":19,\"type\":\"string\",\"size\":0,\"name\":\"MONTH_SEP\"},{\"index\":20,\"type\":\"string\",\"size\":0,\"name\":\"DAY_OF_WEEK_FRI\"},{\"index\":21,\"type\":\"string\",\"size\":0,\"name\":\"DAY_OF_WEEK_MON\"},{\"index\":22,\"type\":\"string\",\"size\":0,\"name\":\"DAY_OF_WEEK_THU\"},{\"index\":23,\"type\":\"string\",\"size\":0,\"name\":\"DAY_OF_WEEK_TUE\"},{\"index\":24,\"type\":\"string\",\"size\":0,\"name\":\"DAY_OF_WEEK_WED\"},{\"index\":25,\"type\":\"string\",\"size\":0,\"name\":\"POUTCOME_FAILURE\"},{\"index\":26,\"type\":\"string\",\"size\":0,\"name\":\"POUTCOME_NONEXISTENT\"},{\"index\":27,\"type\":\"string\",\"size\":0,\"name\":\"POUTCOME_SUCCESS\"},{\"index\":28,\"type\":\"string\",\"size\":0,\"name\":\"Y\"}],\"columns\":28,\"rows\":45036,\"consumeOptions\":[],\"dataPath\":\"/home/user1/data/data_root/bank_train_partyA_20220805-084455.csv\",\"condition\":3,\"originId\":\"4c3b974ac80f6cda3f0654b28212606c91e0fcecac4c88ba3183cd51fd0f3f41\",\"size\":3792331,\"consumeTypes\":[],\"hasTitle\":true}"))
	v := rlputil.RlpHash(buf.Bytes())
	assert.Equal(t, v.Hex(), "0x6a536d83476b9aa73c8abf72e61b4acb980530652e7c0ef9afcf6a5ac9005a00")

	sign := common.Hex2Bytes("44a0f4174bf8c9dc313ffa3eea2d7fc158acc01ab74f6c65aa0730b556fd07c71c36003a110118c1a7ac6d1a2c912b60cb2064666ea48987cb013a92e4bb508a1b")
	hashStr := "0x6ba5a559cce87a060cbae660818fc2b0159f6522812b825d6df2497fe2e3ec8b"
	walletAddress, _, err := Sender(1, common.HexToHash(hashStr), sign)
	assert.Nil(t, err)
	assert.Equal(t, walletAddress, "0xc115ceadf9e5923330e5f42903fe7f926dda65d2")
}
