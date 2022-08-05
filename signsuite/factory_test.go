package signsuite

import (
	"bytes"
	"github.com/datumtechs/datum-network-carrier/common"
	"github.com/datumtechs/datum-network-carrier/common/rlputil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSender(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte("0xf795811af86e9f23a0c03de5115398b8d4778ed4"))
	v := rlputil.RlpHash(buf.Bytes())
	sign := common.Hex2Bytes("8741b27806ac5f62d9d43a6a5434e07d456e1b8d78f59ad51ad89acecfeb5f171fa7986924d2b92a94b0c6db7598bda67b69641a6c2c712b308bf8cf2679f2df1b")
	result, publicKey, err := Sender(1, v, sign)
	assert.Nil(t, err)
	t.Log(result)
	t.Log(publicKey)
}
