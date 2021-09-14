package resource


type ReleaseResourceOption uint32

const (
	UnlockLocalResorce ReleaseResourceOption = 1<< iota // 0001
	RemoveLocalTask                                     // 0010
	CleanTaskEvents                                     // 0100
)

func SetAllReleaseResourceOption() ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTask | CleanTaskEvents
}

func (option ReleaseResourceOption) IsUnlockLocalResorce() bool {
	return option&UnlockLocalResorce == UnlockLocalResorce
}

func (option ReleaseResourceOption) IsRemoveLocalTask() bool {
	return option&RemoveLocalTask == RemoveLocalTask
}

func (option ReleaseResourceOption) IsCleanTaskEvents() bool {
	return option&CleanTaskEvents == CleanTaskEvents
}

