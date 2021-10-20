package resource


type ReleaseResourceOption uint32

const (
	UnlockLocalResorce    ReleaseResourceOption = 1<< iota // 0001
	RemoveLocalTask                                        // 0010
	RemoveLocalTaskEvents                                  // 0100
)


func SetUnlockLocalResorce () ReleaseResourceOption {
	return UnlockLocalResorce
}

func SetRemoveLocalTask () ReleaseResourceOption {
	return RemoveLocalTask
}

func SetRemoveLocalTaskEvents () ReleaseResourceOption {
	return RemoveLocalTaskEvents
}

func SetUnlockLocalResorceAndRemoveLocalTask () ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTask
}

func SetRemoveLocalTaskAndEvents () ReleaseResourceOption {
	return RemoveLocalTask | RemoveLocalTaskEvents
}

func SetUnlockLocalResorceAndRemoveEvents () ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTaskEvents
}

func SetAllReleaseResourceOption() ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTask | RemoveLocalTaskEvents
}


func (option ReleaseResourceOption) IsUnlockLocalResorce() bool {
	return option&UnlockLocalResorce == UnlockLocalResorce
}

func (option ReleaseResourceOption) IsRemoveLocalTask() bool {
	return option&RemoveLocalTask == RemoveLocalTask
}

func (option ReleaseResourceOption) IsRemoveLocalTaskEvents() bool {
	return option&RemoveLocalTaskEvents == RemoveLocalTaskEvents
}


