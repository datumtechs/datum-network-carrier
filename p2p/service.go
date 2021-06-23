package p2p

import (
	"context"
	"crypto/ecdsa"
	"github.com/RosettaFlow/Carrier-Go/common/slotutil"
	pb "github.com/RosettaFlow/Carrier-Go/lib/p2p/v1"
	"github.com/RosettaFlow/Carrier-Go/p2p/peers"
	"github.com/RosettaFlow/Carrier-Go/params"
	"github.com/kevinms/leakybucket-go"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"sync"
	"time"
)

//var _ shared.Service = (*Service)(nil)

// In the event that we are at our peer limit, we stop looking for new peers
// and instead poll for the current peer limit status
// for the time period defined below.
var pollingPeriod = 6 * time.Second

// Refresh rate of ENR set at twice per slot.
var refreshRate = slotutil.DivideSlotBy(2)

// maxBadResponses is the maximum number of bad responses from a peer before we stop talking to it.
const maxBadResponses = 5

// maxDialTimeout is the timeout for a single peer dial.
var maxDialTimeout = params.CarrierNetworkConfig().RespTimeout

// Service for managing peer to peer (p2p) networking.
type Service struct {
	started               bool
	isPreGenesis          bool
	currentForkDigest     [4]byte
	pingMethod            func(ctx context.Context, id peer.ID) error
	cancel                context.CancelFunc
	cfg                   *Config
	peers                 *peers.Status
	addrFilter            *multiaddr.Filters
	ipLimiter             *leakybucket.Collector
	privKey               *ecdsa.PrivateKey
	metaData              *pb.MetaData
	pubsub                *pubsub.PubSub
	joinedTopics          map[string]*pubsub.Topic
	joinedTopicsLock      sync.Mutex
	subnetsLock           map[uint64]*sync.RWMutex
	subnetsLockLock       sync.Mutex // Lock access to subnetsLock
	initializationLock    sync.Mutex
	dv5Listener           Listener
	startupErr            error
	//stateNotifier         statefeed.Notifier
	ctx                   context.Context
	host                  host.Host
	genesisTime           time.Time
	genesisValidatorsRoot []byte
	activeValidatorCount  uint64
}
