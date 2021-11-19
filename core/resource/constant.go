package resource


type ReleaseResourceOption uint32

const (
	UnlockLocalResorce    ReleaseResourceOption = 1<< iota // 0001     1
	RemoveLocalTask                                        // 0010     2
	RemoveLocalTaskEvents                                  // 0100     4
)

// 1
func SetUnlockLocalResorce () ReleaseResourceOption {
	return UnlockLocalResorce
}

// 2
func SetRemoveLocalTask () ReleaseResourceOption {
	return RemoveLocalTask
}

// 4
func SetRemoveLocalTaskEvents () ReleaseResourceOption {
	return RemoveLocalTaskEvents
}

// 3
func SetUnlockLocalResorceAndRemoveLocalTask () ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTask
}

// 6
func SetRemoveLocalTaskAndEvents () ReleaseResourceOption {
	return RemoveLocalTask | RemoveLocalTaskEvents
}

// 5
func SetUnlockLocalResorceAndRemoveEvents () ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTaskEvents
}

// 7
func SetAllReleaseResourceOption() ReleaseResourceOption {
	return UnlockLocalResorce | RemoveLocalTask | RemoveLocalTaskEvents
}

// 1
func (option ReleaseResourceOption) IsUnlockLocalResorce() bool {
	return option&UnlockLocalResorce == UnlockLocalResorce
}

// 2
func (option ReleaseResourceOption) IsRemoveLocalTask() bool {
	return option&RemoveLocalTask == RemoveLocalTask
}

// 4
func (option ReleaseResourceOption) IsRemoveLocalTaskEvents() bool {
	return option&RemoveLocalTaskEvents == RemoveLocalTaskEvents
}


