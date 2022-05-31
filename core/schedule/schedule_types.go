package schedule

import (
	carriertypespb "github.com/datumtechs/datum-network-carrier/pb/carrier/types"
	"github.com/datumtechs/datum-network-carrier/types"
)

// for schedule powerSuppliers
type ScheduleCollecter interface {
	Name() string
	String() string
}

type ScheduleWithSymbolRandomElectionPower struct {
	partyIds []string
}
func (s *ScheduleWithSymbolRandomElectionPower) Name() string { return "ScheduleWithSymbolRandomElectionPower" }
func (s *ScheduleWithSymbolRandomElectionPower) String() string        { return "" }
func (s *ScheduleWithSymbolRandomElectionPower) GetPartyIds() []string { return s.partyIds }
func (s *ScheduleWithSymbolRandomElectionPower) AppendPartyId(partyId string) {
	s.partyIds = append(s.partyIds, partyId)
}


type ScheduleWithDataNodeProvidePower struct {
	provides []*types.TaskPowerPolicyDataNodeProvide
}
func (s *ScheduleWithDataNodeProvidePower) Name() string   { return "ScheduleWithDataNodeProvidePower" }
func (s *ScheduleWithDataNodeProvidePower) String() string { return "" }
func (s *ScheduleWithDataNodeProvidePower) Getprovides() []*types.TaskPowerPolicyDataNodeProvide { return s.provides }
func (s *ScheduleWithDataNodeProvidePower) AppendProvide(provide *types.TaskPowerPolicyDataNodeProvide) {
	s.provides = append(s.provides, provide)
}



// for reschedule powerSuppliers
type ReScheduleCollecter interface {
	Name() string
	String() string
}

type ReScheduleWithSymbolRandomElectionPower struct {
	suppliers []*carriertypespb.TaskOrganization
	resources []*carriertypespb.TaskPowerResourceOption
}
func (r *ReScheduleWithSymbolRandomElectionPower) Name() string { return "ReScheduleWithSymbolRandomElectionPower" }
func (r *ReScheduleWithSymbolRandomElectionPower) String() string        { return "" }
func (r *ReScheduleWithSymbolRandomElectionPower) GetSuppliers() []*carriertypespb.TaskOrganization { return r.suppliers }
func (r *ReScheduleWithSymbolRandomElectionPower) GetResources() []*carriertypespb.TaskPowerResourceOption { return r.resources }
func (r *ReScheduleWithSymbolRandomElectionPower) AppendSupplier(supplier *carriertypespb.TaskOrganization) {
	r.suppliers = append(r.suppliers, supplier)
}
func (r *ReScheduleWithSymbolRandomElectionPower) AppendResource(resource *carriertypespb.TaskPowerResourceOption) {
	r.resources = append(r.resources, resource)
}


type ReScheduleWithDataNodeProvidePower struct {
	suppliers []*carriertypespb.TaskOrganization
	//resources  []*carriertypespb.TaskPowerResourceOption
	provides []*types.TaskPowerPolicyDataNodeProvide
}
func (r *ReScheduleWithDataNodeProvidePower) Name() string { return "ReScheduleWithDataNodeProvidePower" }
func (r *ReScheduleWithDataNodeProvidePower) String() string        { return "" }
func (r *ReScheduleWithDataNodeProvidePower) GetSuppliers() []*carriertypespb.TaskOrganization { return r.suppliers }
func (r *ReScheduleWithDataNodeProvidePower) GetProvides() []*types.TaskPowerPolicyDataNodeProvide { return r.provides }
func (r *ReScheduleWithDataNodeProvidePower) AppendSupplier(supplier *carriertypespb.TaskOrganization) {
	r.suppliers = append(r.suppliers, supplier)
}
func (r *ReScheduleWithDataNodeProvidePower) AppendProvide(provide *types.TaskPowerPolicyDataNodeProvide) {
	r.provides = append(r.provides, provide)
}