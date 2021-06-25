package types

type ResourceDataStatus string
func (r ResourceDataStatus) String() string {return string(r)}
const (
	ResourceDataStatusN ResourceDataStatus =  "N"
	ResourceDataStatusD ResourceDataStatus =  "D"
)


type PowerState string
func (p PowerState) String() string {return string(p)}
const (
	PowerStateCreate PowerState = "create"
	PowerStateRelease PowerState = "release"
	PowerStateRevoke PowerState = "revoke"
)


type MetaDataState string
func (m MetaDataState) String() string {return string(m)}
const (
	MetaDataStateCreate MetaDataState = "create"
	MetaDataStateRelease MetaDataState = "release"
	MetaDataStateRevoke MetaDataState = "revoke"
)