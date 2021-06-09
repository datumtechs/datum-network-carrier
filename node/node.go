package node


type Node struct {

	// p2p


	serviceFuncs []ServiceConstructor     // Service constructors (in dependency order)

}


func New() *Node {

	return nil
}