package types

type TaskEvent struct {
	TaskId   string     `json:"taskId"`
	Type     string     `json:"type"`
	CreateAt uint64     `json:"createAt"`
	Content  string     `json:"content"`
	Owner    *NodeAlias `json:"owner"`
}
