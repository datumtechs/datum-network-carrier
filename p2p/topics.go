package p2p

const (
	//TODO: to define more topics....
	GossipTestDataTopicFormat = "/carrier/%x/gossip_test_data"
	// AttestationSubnetTopicFormat is the topic format for the attestation subnet.
	AttestationSubnetTopicFormat = "/carrier/%x/beacon_attestation_%d"
	// BlockSubnetTopicFormat is the topic format for the block subnet.
	BlockSubnetTopicFormat = "/carrier/%x/carrier_block"
	// ExitSubnetTopicFormat is the topic format for the voluntary exit subnet.
	ExitSubnetTopicFormat = "/carrier/%x/voluntary_exit"
	// ProposerSlashingSubnetTopicFormat is the topic format for the proposer slashing subnet.
	ProposerSlashingSubnetTopicFormat = "/carrier/%x/proposer_slashing"
	// AttesterSlashingSubnetTopicFormat is the topic format for the attester slashing subnet.
	AttesterSlashingSubnetTopicFormat = "/carrier/%x/attester_slashing"
	// AggregateAndProofSubnetTopicFormat is the topic format for the aggregate and proof subnet.
	AggregateAndProofSubnetTopicFormat = "/carrier/%x/beacon_aggregate_and_proof"
)
