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
	return VoteOptionFromUint8(option[0])
}




type TwopcMsgOption uint8

func (c TwopcMsgOption) Bytes() []byte { return []byte{byte(c)} }
func (c TwopcMsgOption) String() string {
	switch c {
	case TwopcMsgStart:
		return "Start"
	case TwopcMsgStop:
		return "Stop"
	default:
		return "Unknown"
	}
}

const (
	TwopcMsgUnknown TwopcMsgOption = 0x00
	TwopcMsgStart   TwopcMsgOption = 0x01
	TwopcMsgStop    TwopcMsgOption = 0x02
)

func TwopcMsgOptionFromUint8(option uint8) TwopcMsgOption {
	switch option {
	case 0x01:
		return TwopcMsgStart
	case 0x02:
		return TwopcMsgStop
	default:
		return TwopcMsgUnknown
	}
}
func TwopcMsgOptionFromStr(option string) TwopcMsgOption {
	switch option {
	case "Start":
		return TwopcMsgStart
	case "Stop":
		return TwopcMsgStop
	default:
		return TwopcMsgUnknown
	}
}
func TwopcMsgOptionFromBytes(option []byte) TwopcMsgOption {
	if len(option) != 1 {
		return TwopcMsgUnknown
	}
	return TwopcMsgOptionFromUint8(option[0])
}



