package grpclient


type InternalResourceNodeSet struct {
	// GRPC Client
	jobNodes  map[string]*JobNodeClient
	dataNodes map[string]*DataNodeClient
}

func NewInternalResourceNodeSet () *InternalResourceNodeSet {
	return &InternalResourceNodeSet{
		jobNodes:  make(map[string]*JobNodeClient),
		dataNodes: make(map[string]*DataNodeClient),
	}
}

func (nodeSet *InternalResourceNodeSet) StoreJobNodeClient(nodeId string, client *JobNodeClient) {
	nodeSet.jobNodes[nodeId] = client
}
func (nodeSet *InternalResourceNodeSet) QueryJobNodeClient(nodeId string) (*JobNodeClient, bool) {
	client, ok := nodeSet.jobNodes[nodeId]
	return client, ok
}
func (nodeSet *InternalResourceNodeSet) QueryJobNodeClients() []*JobNodeClient {
	arr := make([]*JobNodeClient, 0)
	for _, client := range nodeSet.jobNodes {
		arr = append(arr, client)
	}
	return arr
}
func (nodeSet *InternalResourceNodeSet) RemoveJobNodeClient(nodeId string)  {
	delete(nodeSet.jobNodes, nodeId)
}


func (nodeSet *InternalResourceNodeSet) StoreDataNodeClient(nodeId string, client *DataNodeClient) {
	nodeSet.dataNodes[nodeId] = client
}
func (nodeSet *InternalResourceNodeSet) QueryDataNodeClient(nodeId string) (*DataNodeClient, bool) {
	client, ok := nodeSet.dataNodes[nodeId]
	return client, ok
}
func (nodeSet *InternalResourceNodeSet) QueryDataNodeClients() []*DataNodeClient {
	arr := make([]*DataNodeClient, 0)
	for _, client := range nodeSet.dataNodes {
		arr = append(arr, client)
	}
	return arr
}
func (nodeSet *InternalResourceNodeSet) RemoveDataNodeClient(nodeId string)  {
	delete(nodeSet.dataNodes, nodeId)
}