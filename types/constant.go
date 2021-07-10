package types

type ResourceDataStatus string

func (r ResourceDataStatus) String() string { return string(r) }

const (
	ResourceDataStatusN ResourceDataStatus = "N"
	ResourceDataStatusD ResourceDataStatus = "D"
)

type PowerState string

func (p PowerState) String() string { return string(p) }

const (
	PowerStateCreate  PowerState = "create"
	PowerStateRelease PowerState = "release"
	PowerStateRevoke  PowerState = "revoke"
)

type MetaDataState string

func (m MetaDataState) String() string { return string(m) }

const (
	MetaDataStateCreate  MetaDataState = "create"
	MetaDataStateRelease MetaDataState = "release"
	MetaDataStateRevoke  MetaDataState = "revoke"
)

type TaskState string

func (t TaskState) String() string { return string(t) }

// (pending: 等在中; running: 计算中; failed: 失败; success: 成功)
const (
	TaskStatePending TaskState = "pending"
	TaskStateRunning TaskState = "running"
	TaskStateFailed  TaskState = "failed"
	TaskStateSuccess TaskState = "success"
)
