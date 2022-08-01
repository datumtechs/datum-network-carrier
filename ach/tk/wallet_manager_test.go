package tk

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func Test_Address(t *testing.T) {

	if addr := getAddress(); addr == (common.Address{}) {
		t.Log("zero length")
	} else {
		t.Log("not zero length")
	}

}

func getAddress() common.Address {
	return common.Address{}
}
