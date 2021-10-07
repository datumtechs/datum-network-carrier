package p2p

const (
	GossipTestDataTopicFormat = "/carrier/%x/gossip_test_data"

	// ````````````````` topic for consensus message `````````````````````
	// TwoPcPrepareMsgTopicFormat is the topic format for the prepare message of the 2pc consensus.
	TwoPcPrepareMsgTopicFormat    = "/carrier/%x/prepare_message"
	// TwoPcPrepareVoteTopicFormat is the topic format for the prepare vote of the 2pc consensus.
	TwoPcPrepareVoteTopicFormat   = "/carrier/%x/prepare_vote"
	// TwoPcConfirmMsgTopicFormat is the topic format for the confirm message of the 2pc consensus.
	TwoPcConfirmMsgTopicFormat    = "/carrier/%x/confirm_message"
	// TwoPcConfirmVoteTopicFormat is the topic format for the confirm vote of the 2pc consensus.
	TwoPcConfirmVoteTopicFormat   = "/carrier/%x/confirm_vote"
	// TwoPcCommitMsgTopicFormat is the topic format for the commit message of the 2pc consensus.
	TwoPcCommitMsgTopicFormat     = "/carrier/%x/commit_message"
	// TaskResultMsgTopicFormat is the topic format for the task result message of the task manager.
	TaskResultMsgTopicFormat = "/carrier/%x/task_result_message"
	// TaskResourceUsageMsgTopicFormat is the topic format for the task resource usage message of the task manager.
	TaskResourceUsageMsgTopicFormat = "/carrier/%x/task_resource_usage_message"
	// TaskTerminateMsgTopicFormat is the topic format for the task terminate message of the task manager.
	TaskTerminateMsgTopicFormat = "/carrier/%x/task_terminate_message"

	// AttestationSubnetTopicFormat is the topic format for the attestation subnet.
	AttestationSubnetTopicFormat = "/carrier/%x/beacon_attestation_%d"
	// BlockSubnetTopicFormat is the topic format for the block subnet.
	BlockSubnetTopicFormat = "/carrier/%x/carrier_block"
)
