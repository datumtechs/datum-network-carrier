package types

// The proposalMsg signature and the voteMsg signature
type MsgSign [MsgSignLength]byte

func (s MsgSign) Bytes() []byte { return s[:] }

type VoteOption uint8

func (v VoteOption) Byte() byte { return byte(v) }
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
	MsgSignLength = 65
)
const (
	Yes        VoteOption = 0x01
	No         VoteOption = 0x02
	Abstention VoteOption = 0x03
)

func ParseVoteOption(option uint8) VoteOption {
	switch option {
	case 0x01:
		return Yes
	case 0x02:
		return No
	case 0x03:
		return Abstention
	}
	return Abstention
}

type TaskRole uint8
func (t TaskRole)Bytes() []byte {return []byte{byte(t)}}
func (t TaskRole) String() string {
	switch t {
	case DataSupplier:
		return "DataSupplier"
	case PowerSupplier:
		return "PowerSupplier"
	case ResultSupplier:
		return "ResultSupplier"
	default:
		return "Unknown"
	}
}

const (
	Unknown        TaskRole = 0x00
	DataSupplier   TaskRole = 0x01
	PowerSupplier  TaskRole = 0x02
	ResultSupplier TaskRole = 0x03
)

func TaskRoleFromStr(role string) TaskRole {
	switch role {
	case "DataSupplier":
		return DataSupplier
	case "PowerSupplier":
		return PowerSupplier
	case "ResultSupplier":
		return ResultSupplier
	default:
		return Unknown
	}
}
func TaskRoleFromBytes(role []byte) TaskRole {
	if len(role) != 1 {
		return Unknown
	}
	switch role[0] {
	case 0x01:
		return DataSupplier
	case 0x02:
		return PowerSupplier
	case 0x03:
		return ResultSupplier
	default:
		return Unknown
	}
}
