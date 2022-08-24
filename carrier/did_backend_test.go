package carrier

import (
	"context"
	"github.com/datumtechs/datum-network-carrier/common/hashutil"
	"github.com/datumtechs/datum-network-carrier/common/hexutil"
	"github.com/datumtechs/datum-network-carrier/grpclient"
	"github.com/datumtechs/datum-network-carrier/pb/carrier/api"
	"github.com/datumtechs/datum-network-carrier/service/discovery"
	didsdkgocrypto "github.com/datumtechs/did-sdk-go/crypto"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func Test_DialRemoteCarrier(t *testing.T) {

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	issuerUrl := "192.168.9.154:10033"
	conn, err := grpclient.DialContext(ctx, issuerUrl, true)
	if err != nil {
		t.Fatalf("cannot dial up url: %+v", err)
	}
	defer conn.Close()

	client := api.NewVcServiceClient(conn)
	t.Log(client)

}

func Test_ApplyLocalVC(t *testing.T) {
	/*s := CarrierAPIBackend{}
	s.ApplyVCLocal()*/
}

func Test_findAdminService(t *testing.T) {
	consulManager := newConsulManager()
	//查找本地admin服务端口
	adminService, err := consulManager.QueryAdminService()
	if err != nil {
		t.Fatalf("cannot find local admin RPC service: %+v", err)
		return
	}
	if adminService == nil {
		t.Fatalf("cannot find local admin RPC service: %+v", err)
		return
	}
	log.Debugf("adminService info:%+v", adminService)

	ctx, cancelFn := context.WithTimeout(context.Background(), grpclient.DefaultGrpcDialTimeout)
	defer cancelFn()

	adminServiceUrl := adminService.Address + ":" + strconv.Itoa(adminService.Port)
	conn, err := grpclient.DialContext(ctx, adminServiceUrl, true)
	defer conn.Close()

	if err != nil {
		t.Fatalf("failed to dial admin service: %+v", err)
		return
	}
	t.Log(conn)
}

func newConsulManager() *discovery.ConnectConsul {
	consulManager := discovery.NewConsulClient(&discovery.ConsulService{
		ServiceIP:   "192.168.9.156",
		ServicePort: "8500",
		Tags:        DefaultConfig.DiscoverServiceConfig.DiscoveryServerTags,
		Name:        DefaultConfig.DiscoverServiceConfig.DiscoveryServiceName,
		Id:          DefaultConfig.DiscoverServiceConfig.DiscoveryServiceId,
		Interval:    DefaultConfig.DiscoverServiceConfig.DiscoveryServiceHealthCheckInterval,
		Deregister:  DefaultConfig.DiscoverServiceConfig.DiscoveryServiceHealthCheckDeregister,
	},
		"192.168.9.156",
		8500,
	)

	return consulManager
}

func Test_HexEncoding(t *testing.T) {
	rawData := "applicantDid + claim"
	reqHash := hashutil.HashSHA256([]byte(rawData))

	t.Logf("hex1:%s", hexutil.MustDecode(hexutil.Encode(reqHash)))

}

func Test_verifySign(t *testing.T) {
	digest := "0xebd5186e34a8a67b541cf27b17106bae3c77d11154d6e72c2a2f6f243bc04533"
	sign := "0x7f709a6633ecbe51bcbed198d35ec499383cb72edf1e4326b907b5385640af95025a7ea1b16060b6606244e1fc05578580f78396f82e508231778d27916c999300"

	pubKey := "042c1691dd1f657ae994dedc99bbde5f1a95c687ab47df70bd5e556a7d6138d3ddfb4848703d2050b1a23674ad29a28c56e90d5c3b1ac1ca7961c678e3c00eb588"

	publicKey, err := crypto.UnmarshalPubkey(ethcommon.FromHex(pubKey))
	if err != nil {
		t.Fatalf("failed to unmarshal pubkey: %+v", err)
		return
	}

	ok := didsdkgocrypto.VerifySecp256k1Signature(ethcommon.FromHex(digest), sign, publicKey)
	a := assert.New(t)
	if !a.Equal(true, ok) {
		t.Fatalf("failed to verify signature")
	}
}
