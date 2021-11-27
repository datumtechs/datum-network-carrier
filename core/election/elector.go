package election


type Elector interface {
	ElectionNode()
	ElectionOrganization()
}