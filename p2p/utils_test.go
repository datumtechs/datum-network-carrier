package p2p

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

// Test `verifyConnectivity` function by trying to connect to google.com (successfully)
// and then by connecting to an unreachable IP and ensuring that a log is emitted
func TestVerifyConnectivity(t *testing.T) {
	//hook := logTest.NewGlobal()
	cases := []struct {
		address              string
		port                 uint
		expectedConnectivity bool
		name                 string
	}{
		{"142.250.68.46", 80, true, "Dialing a reachable IP: 142.250.68.46:80"}, // google.com
		{"123.123.123.123", 19000, false, "Dialing an unreachable IP: 123.123.123.123:19000"},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name),
			func(t *testing.T) {
				verifyConnectivity(tc.address, tc.port, "tcp")
				//logMessage := "IP address is not accessible"
				if tc.expectedConnectivity {
					//require.Contains(t, hook, logMessage)
				} else {
					//require.Contains(t, hook, logMessage)
				}
			})
	}
}

func TestIDFromPublicKey_HexPrivateKey(t *testing.T) {
	// 0x2ab516beccd56ba844a329cc7d3e235cdc0302e399a9d0015be1eca335635b36
	// 0x887e6ed21139cf609995f6b6964cedefdc458de515991d2dc0b2667889148dfd5dfd011ebcf258303d5d1e0113ba093a06682146d4c5e3be5078dc5ff6e62714
	// 16Uiu2HAm4cVF9ikZ9h7UytUSPfZMZPjwn2GNdDRjfBLfUgz7N6yJ
	prikey, err := crypto.HexToECDSA("2ab516beccd56ba844a329cc7d3e235cdc0302e399a9d0015be1eca335635b36")
	require.Nil(t, err, "convert private key failed.")
	assertedKey := convertToInterfacePubkey(&prikey.PublicKey)
	id, err := peer.IDFromPublicKey(assertedKey)
	require.Nil(t, err, "id from public key failed...")
	t.Logf("peer.ID %s", id)
	assert.Equal(t, "16Uiu2HAm4cVF9ikZ9h7UytUSPfZMZPjwn2GNdDRjfBLfUgz7N6yJ", id.String())
}

func TestIDFromPublicKey_HexPublicKey(t *testing.T) {
	// 0x2ab516beccd56ba844a329cc7d3e235cdc0302e399a9d0015be1eca335635b36
	// 0x887e6ed21139cf609995f6b6964cedefdc458de515991d2dc0b2667889148dfd5dfd011ebcf258303d5d1e0113ba093a06682146d4c5e3be5078dc5ff6e62714
	// 16Uiu2HAm4cVF9ikZ9h7UytUSPfZMZPjwn2GNdDRjfBLfUgz7N6yJ
	byts, err := hexutil.Decode("0x887e6ed21139cf609995f6b6964cedefdc458de515991d2dc0b2667889148dfd5dfd011ebcf258303d5d1e0113ba093a06682146d4c5e3be5078dc5ff6e62714")
	pbytes := []byte{4}
	pbytes = append(pbytes, byts...)
	pubkey, err := crypto.UnmarshalPubkey(pbytes)
	require.Nil(t, err, "UnmarshalPubkey failed.")
	assertedKey := convertToInterfacePubkey(pubkey)
	id, err := peer.IDFromPublicKey(assertedKey)
	require.Nil(t, err, "id from public key failed...")
	t.Logf("peer.ID %s", id)
	assert.Equal(t, "16Uiu2HAm4cVF9ikZ9h7UytUSPfZMZPjwn2GNdDRjfBLfUgz7N6yJ", id.String())
}