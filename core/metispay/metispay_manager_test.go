package metispay

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
	"time"
)

var (
	timespan = 4
)

func TestPrepayReceipt(t *testing.T) {
	ch := make(chan *prepayReceipt)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var receipt *prepayReceipt
LOOP:
	for {
		select {
		case receipt = <-ch:
			t.Log("get the receipt")
			break LOOP
		case <-ticker.C:
			t.Log("try to get the receipt")
			go getPrepayReceipt(common.HexToHash("0x00"), ch)
		}
	}
	t.Logf("receipt.success:%t", receipt.success)
	t.Logf("receipt.gasUsed:%d", receipt.gasUsed)
}

func getPrepayReceipt(txHash common.Hash, ch chan *prepayReceipt) {
	if timespan > 0 {
		timespan--
		return
	}
	prepayReceipt := new(prepayReceipt)
	prepayReceipt.success = true
	prepayReceipt.gasUsed = 100
	ch <- prepayReceipt
}
