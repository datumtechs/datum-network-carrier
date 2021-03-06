package iputils_test

import (
	"github.com/RosettaFlow/Carrier-Go/common/iputils"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"net"
	"regexp"
	"testing"

)

func TestExternalIPv4(t *testing.T) {
	// Regular expression format for IPv4
	IPv4Format := `\.\d{1,3}\.\d{1,3}\b`
	test, err := iputils.ExternalIPv4()
	require.NoError(t, err)

	valid := regexp.MustCompile(IPv4Format)
	assert.Equal(t, true, valid.MatchString(test))
}

func TestRetrieveIP(t *testing.T) {
	ip, err := iputils.ExternalIP()
	if err != nil {
		t.Fatal(err)
	}
	retIP := net.ParseIP(ip)
	if retIP.To4() == nil && retIP.To16() == nil {
		t.Errorf("An invalid IP was retrieved: %s", ip)
	}
}

func TestSortAddresses(t *testing.T) {
	testAddresses := []net.IP{
		{0xff, 0x02, 0xAA, 0, 0x1F, 0, 0, 0, 0, 0, 0x02, 0x2E, 0, 0, 0x36, 0x45},
		{0xff, 0x02, 0xAA, 0, 0x1F, 0, 0x2E, 0, 0, 0x36, 0x45, 0, 0, 0, 0, 0x02},
		{0xAA, 0x11, 0x33, 0x19},
		{0x01, 0xBF, 0x33, 0x10},
		{0x03, 0x89, 0x33, 0x13},
	}

	sortedAddrs := iputils.SortAddresses(testAddresses)
	assert.Equal(t, true, sortedAddrs[0].To4() != nil, "expected ipv4 address")
	assert.Equal(t, true, sortedAddrs[1].To4() != nil, "expected ipv4 address")
	assert.Equal(t, true, sortedAddrs[2].To4() != nil, "expected ipv4 address")
	assert.Equal(t, true, sortedAddrs[3].To16() != nil && sortedAddrs[3].To4() == nil, "expected ipv6 address")
	assert.Equal(t, true, sortedAddrs[4].To16() != nil && sortedAddrs[4].To4() == nil, "expected ipv6 address")
}
