package types

type ResourceDataStatus string
func (r ResourceDataStatus) String() string {return string(r)}
const (
	ResourceDataStatusN =  "N"
	ResourceDataStatusD =  "D"
)


type PowerState string
func (p PowerState) String() string {return string(p)}
const (
	PowerStateCreate = "create"
	PowerStateRelease = "release"
	PowerStateRevoke = "revoke"
)


type MetaDataState string
func (m MetaDataState) String() string {return string(m)}
const (
	MetaDataStateCreate = "create"
	MetaDataStateRelease = "release"
	MetaDataStateRevoke = "revoke"
)