package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTwopcMsgOption(t *testing.T) {
	start := TwopcMsgStart
	unknown := TwopcMsgUnknown
	stop := TwopcMsgStop
	require.Equal(t, start.String(), "START")
	require.Equal(t, stop.String(), "STOP")
	require.Equal(t, unknown.String(), "UNKNOWN")

	require.Equal(t, TwopcMsgOptionFromUint8(0), unknown)
	require.Equal(t, TwopcMsgOptionFromUint8(1), start)
	require.Equal(t, TwopcMsgOptionFromUint8(2), stop)

	require.Equal(t, TwopcMsgOptionFromBytes([]byte{0}), unknown)
	require.Equal(t, TwopcMsgOptionFromBytes([]byte{1}), start)
	require.Equal(t, TwopcMsgOptionFromBytes([]byte{2}), stop)
}
