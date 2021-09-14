package types

import (
	"fmt"
	"github.com/RosettaFlow/Carrier-Go/common"
	"github.com/RosettaFlow/Carrier-Go/crypto/sha3"
	"github.com/RosettaFlow/Carrier-Go/lib/types"
	"math/big"
	"sort"
	"sync/atomic"
	"time"
)

// Header represents a block header in the RosettaNet.
type Header types.HeaderPb

// Hash returns the block hash of the header, which is simply the keccak256 hash of its ProtoBuf encoding.
func (h *Header) Hash() common.Hash {
	header := h.GetHeaderPb()
	data, _ := header.Marshal()
	return protoBufHash(data)
}

func (h *Header) GetHeaderPb() *types.HeaderPb {
	header := &types.HeaderPb{
		ParentHash: h.ParentHash,
		Timestamp: h.Timestamp,
		Version: h.Version,
		Extra: h.Extra,
	}
	return header
}

func protoBufHash(x []byte) (h common.Hash) {
	hw := sha3.NewKeccak256()
	hw.Write(x)
	hw.Sum(h[:0])
	return h
}

// Signature returns the signature of seal hash from extra.
func (h *Header) Signature() []byte {
	if len(h.Extra) < 32 {
		return []byte{}
	}
	return h.Extra[32:]
}

func (h *Header) _sealHash() (hash common.Hash) {
	extra := h.Extra
	hasher := sha3.NewKeccak256()
	if len(h.Extra) > 32 {
		extra = h.Extra[0:32]
	}
	header := &types.HeaderPb{
		ParentHash: h.ParentHash,
		Timestamp: h.Timestamp,
		Version: h.Version,
		Extra: extra,
	}
	data, _ := header.Marshal()
	hasher.Write(data)
	hasher.Sum(hash[:0])
	return hash
}

// Block represents an entire block in the RosettaNet.
type Block struct {
	header       *Header
	metadatas 	 MetadataArray
	resources    ResourceArray
	identities 	 IdentityArray
	taskDatas	 TaskDataArray

	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package rosetta to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
	extraData    []byte
}

func NewBlock(header *Header, metadataList []*Metadata, resourcesList []*Resource,
					identityList []*Identity, taskDataList []*Task) *Block {
	b := &Block{header: CopyHeader(header)}
	if len(metadataList) != 0 {
		b.metadatas = make(MetadataArray, len(metadataList))
		copy(b.metadatas, metadataList)
	}
	if len(resourcesList) != 0 {
		b.resources = make(ResourceArray, len(resourcesList))
		copy(b.resources, resourcesList)
	}
	if len(identityList) != 0 {
		b.identities = make(IdentityArray, len(identityList))
		copy(b.identities, identityList)
	}
	if len(taskDataList) != 0 {
		b.taskDatas = make(TaskDataArray, len(taskDataList))
		copy(b.taskDatas, taskDataList)
	}
	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	cpy.Timestamp = h.Timestamp
	cpy.Version = h.Version
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	return &cpy
}

// DecodePb decodes the RosettaNet
func (b *Block) DecodePb(data []byte) error {
	var blockData types.BlockData
	blockData.Unmarshal(data)
	b.header = (*Header)(blockData.Header)
	//
	b.metadatas = NewMetadataArray(blockData.Metadata)
	b.resources = NewResourceArray(blockData.Resourcedata)
	b.identities = NewIdentityArray(blockData.Identitydata)
	b.taskDatas = NewTaskDataArray(blockData.Taskdata)
	return nil
}


// EncodePb serializes b into the Ethereum RLP block format.
func (b *Block) EncodePb() ([]byte, error) {
	blockData := &types.BlockData{
		Header: b.header.GetHeaderPb(),
		Metadata: b.metadatas.To(),
		Resourcedata: b.resources.To(),
		Taskdata: b.taskDatas.To(),
		Identitydata: b.identities.To(),
	}
	return blockData.Marshal()
}

func (b *Block) Metadatas() MetadataArray { return b.metadatas }
func (b *Block) Resources() ResourceArray { return b.resources }
func (b *Block) Identities() IdentityArray { return b.identities }
func (b *Block) TaskDatas() TaskDataArray { return b.taskDatas }

func (b *Block) Metadata(hash common.Hash) *Metadata {
	for _, metadata := range b.metadatas {
		if metadata.Hash() == hash {
			return metadata
		}
	}
	return nil
}
func (b *Block) SetExtraData(extraData []byte) { b.extraData = extraData }
func (b *Block) ExtraData() []byte             { return b.extraData }
func (b *Block) Time() *big.Int                { return new(big.Int).SetUint64(b.header.Timestamp) }

func (b *Block) NumberU64() uint64        { return b.header.Version }
func (b *Block) ParentHash() common.Hash  { return common.BytesToHash(b.header.ParentHash) }
func (b *Block) Extra() []byte            { return common.CopyBytes(b.header.Extra) }

func (b *Block) Header() *Header { return CopyHeader(b.header) }

// WithBody returns a new block with the given data(metadata/identity/resource/taskdata).
func (b *Block) WithBody(metadataList []*Metadata, resourcesList []*Resource,
	identityList []*Identity, taskDataList []*Task, extraData []byte) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		metadatas: 	  make(MetadataArray, len(metadataList)),
		resources: 	  make(ResourceArray, len(resourcesList)),
		identities:   make(IdentityArray, len(identityList)),
		taskDatas:    make(TaskDataArray, len(taskDataList)),
		extraData:    make([]byte, len(extraData)),
	}
	copy(block.metadatas, metadataList)
	copy(block.resources, resourcesList)
	copy(block.identities, identityList)
	copy(block.taskDatas, taskDataList)
	copy(block.extraData, extraData)
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type Blocks []*Block

func (b Blocks) String() string {
	s := "["
	for _, v := range b {
		s += fmt.Sprintf("[hash:%s, number:%d]", v.Hash().TerminalString(), v.NumberU64())
	}
	s += "]"
	return s
}

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.header.Version < b2.header.Version }

