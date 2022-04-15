package params

// add by v 0.4.0
type TaskManagerConfig struct {
	ConsumOption                uint32 // Task consumption options
	NeedReplayScheduleTaskChLen uint32
	NeedExecuteTaskChanLen      uint32
	TaskConsResultChanLen       uint32
}
