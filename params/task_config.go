package params

// add by v 0.4.0
type TaskManagerConfig struct {
	MetadataConsumeOption          int // Metadata of task consumption options (0: nothing, 1: metadataAuth, 2: datatoken)
	NeedReplayScheduleTaskChanSize int
	NeedExecuteTaskChanSize        int
	TaskConsResultChanSize         int
}
