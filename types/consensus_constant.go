package types

const (
	MsgSignLength = 65
	MsgEpochMaxNumber = 2
)
// The proposalMsg signature and the voteMsg signature
type MsgSign [MsgSignLength]byte

func (s MsgSign) Bytes() []byte { return s[:] }

type VoteOption uint8

func (v VoteOption) Bytes() []byte { return []byte{byte(v)} }
func (v VoteOption) String() string {
	switch v {
	case Yes:
		return "Yes"
	case No:
		return "No"
	default:
		return "Abstention"
	}
}

const (
	VoteUnknown VoteOption = 0x00
	Yes         VoteOption = 0x01
	No          VoteOption = 0x02
	Abstention  VoteOption = 0x03
)

func VoteOptionFromUint8(option uint8) VoteOption {
	switch option {
	case 0x01:
		return Yes
	case 0x02:
		return No
	case 0x03:
		return Abstention
	default:
		return VoteUnknown
	}
}
func VoteOptionFromStr(option string) VoteOption {
	switch option {
	case "Yes":
		return Yes
	case "No":
		return No
	case "Abstention":
		return Abstention
	default:
		return VoteUnknown
	}
}
func VoteOptionFromBytes(option []byte) VoteOption {
	if len(option) != 1 {
		return VoteUnknown
	}
	return VoteOptionFromUint8(uint8(option[0]))
}




//type TaskRole uint8
//
//func (t TaskRole) Bytes() []byte { return []byte{byte(t)} }
//func (t TaskRole) String() string {
//	switch t {
//	case DataSupplier:
//		return "DataSupplier"
//	case PowerSupplier:
//		return "PowerSupplier"
//	case ResultSupplier:
//		return "ResultSupplier"
//	case TaskOwner:
//		return "TaskOwner"
//	default:
//		return "TaskRoleUnknown"
//	}
//}
//
//const (
//	TaskRoleUnknown TaskRole = 0x00
//	DataSupplier    TaskRole = 0x01
//	PowerSupplier   TaskRole = 0x02
//	ResultSupplier  TaskRole = 0x03
//	TaskOwner       TaskRole = 0x04
//)
//
//func TaskRoleFromUint8(role uint8) TaskRole {
//	switch role {
//	case 0x01:
//		return DataSupplier
//	case 0x02:
//		return PowerSupplier
//	case 0x03:
//		return ResultSupplier
//	case 0x04:
//		return TaskOwner
//	default:
//		return TaskRoleUnknown
//	}
//}
//
//func TaskRoleFromStr(role string) TaskRole {
//	switch role {
//	case "DataSupplier":
//		return DataSupplier
//	case "PowerSupplier":
//		return PowerSupplier
//	case "ResultSupplier":
//		return ResultSupplier
//	case "TaskOwner":
//		return TaskOwner
//	default:
//		return TaskRoleUnknown
//	}
//}
//func TaskRoleFromBytes(role []byte) TaskRole {
//	if len(role) != 1 {
//		return TaskRoleUnknown
//	}
//	return TaskRoleFromUint8(role[0])
//}
