package consensus

import (
	"github.com/RosettaFlow/Carrier-Go/types"
)

//var (
//	ChainConsTyp = reflect.TypeOf(chaincons.Chaincons{})
//	TwoPcTyp = reflect.TypeOf(twopc.TwoPC{})
//)


type Engine interface {
	OnPrepare(task *types.ScheduleTask) error
	OnStart(task *types.ScheduleTask, result chan<- *types.TaskConsResult) error
	ValidateConsensusMsg(msg types.ConsensusMsg) error
	OnConsensusMsg(msg types.ConsensusMsg) error
	OnError() error
}



// Engine is an algorithm agnostic consensus engine.
//type Engine interface {
//	// Author retrieves the Ethereum address of the account that minted the given
//	// block, which may be different from the header's coinbase if a consensus
//	// engine is based on signatures.
//	Author(header *types.Header) (common.Address, error)
//
//	// VerifyHeader checks whether a header conforms to the consensus rules of a
//	// given engine. Verifying the seal may be done optionally here, or explicitly
//	// via the VerifySeal method.
//	VerifyHeader(chain ChainReader, header *types.Header, seal bool) error
//
//	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
//	// concurrently. The method returns a quit channel to abort the operations and
//	// a results channel to retrieve the async verifications (the order is that of
//	// the input slice).
//	VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)
//
//	// VerifySeal checks whether the crypto seal on a header is valid according to
//	// the consensus rules of the given engine.
//	VerifySeal(chain ChainReader, header *types.Header) error
//
//	// Prepare initializes the consensus fields of a block header according to the
//	// rules of a particular engine. The changes are executed inline.
//	Prepare(chain ChainReader, header *types.Header) error
//
//	// Finalize runs any post-transaction state modifications (e.g. block rewards)
//	// and assembles the final block.
//	// Note: The block header and state database might be updated to reflect any
//	// consensus rules that happen at finalization (e.g. block rewards).
//	Finalize(chain ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
//		receipts []*types.Receipt) (*types.Block, error)
//
//	// Seal generates a new sealing request for the given input block and pushes
//	// the result into the given channel.
//	//
//	// Note, the method returns immediately and will send the result async. More
//	// than one result may also be returned depending on the consensus algorithm.
//	Seal(chain ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}, complete chan<- struct{}) error
//
//	// SealHash returns the hash of a block prior to it being sealed.
//	SealHash(header *types.Header) common.Hash
//
//	// APIs returns the RPC APIs this consensus engine provides.
//	APIs(chain ChainReader) []rpc.API
//
//	Protocols() []p2p.Protocol
//
//	NextBaseBlock() *types.Block
//
//	InsertChain(block *types.Block) error
//
//	HasBlock(hash common.Hash, number uint64) bool
//
//	GetBlockByHash(hash common.Hash) *types.Block
//
//	GetBlockByHashAndNum(hash common.Hash, number uint64) *types.Block
//
//	CurrentBlock() *types.Block
//
//	FastSyncCommitHead(block *types.Block) error
//
//	// Close terminates any background threads maintained by the consensus engine.
//	Close() error
//
//	// Pause consensus
//	Pause()
//	// Resume consensus
//	Resume()
//
//	DecodeExtra(extra []byte) (common.Hash, uint64, error)
//}