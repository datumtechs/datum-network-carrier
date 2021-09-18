package grpclient

type InternalResourceClientSet struct {
	// GRPC Client
	jobNodes  map[string]*JobNodeClient
	dataNodes map[string]*DataNodeClient
}

func NewInternalResourceNodeSet () *InternalResourceClientSet {
	return &InternalResourceClientSet{
		jobNodes:  make(map[string]*JobNodeClient),
		dataNodes: make(map[string]*DataNodeClient),
	}
}

func (nodeSet *InternalResourceClientSet) StoreJobNodeClient(nodeId string, client *JobNodeClient) {
	nodeSet.jobNodes[nodeId] = client
}

func (nodeSet *InternalResourceClientSet) QueryJobNodeClient(nodeId string) (*JobNodeClient, bool) {
	client, ok := nodeSet.jobNodes[nodeId]
	return client, ok
}

func (nodeSet *InternalResourceClientSet) QueryJobNodeClients() []*JobNodeClient {
	arr := make([]*JobNodeClient, 0)
	for _, client := range nodeSet.jobNodes {
		arr = append(arr, client)
	}
	return arr
}

func (nodeSet *InternalResourceClientSet) RemoveJobNodeClient(nodeId string)  {
	delete(nodeSet.jobNodes, nodeId)
}

func (nodeSet *InternalResourceClientSet) JobNodeClientSize() int {
	return len(nodeSet.jobNodes)
}

func (nodeSet *InternalResourceClientSet) StoreDataNodeClient(nodeId string, client *DataNodeClient) {
	nodeSet.dataNodes[nodeId] = client
}

func (nodeSet *InternalResourceClientSet) QueryDataNodeClient(nodeId string) (*DataNodeClient, bool) {
	client, ok := nodeSet.dataNodes[nodeId]
	return client, ok
}

func (nodeSet *InternalResourceClientSet) QueryDataNodeClients() []*DataNodeClient {
	arr := make([]*DataNodeClient, 0)
	for _, client := range nodeSet.dataNodes {
		arr = append(arr, client)
	}
	return arr
}

func (nodeSet *InternalResourceClientSet) RemoveDataNodeClient(nodeId string)  {
	delete(nodeSet.dataNodes, nodeId)
}

func (nodeSet *InternalResourceClientSet) DataNodeClientSize() int {
	return len(nodeSet.dataNodes)
}