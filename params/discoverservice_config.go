package params


type DiscoverServiceConfig struct {
	DiscoveryServerIP     string
	DiscoveryServerPort   int
	DiscoveryServiceId  string
	DiscoveryServiceName string
	DiscoveryServerTags  []string
	DiscoveryServiceHealthCheckInterval int
	DiscoveryServiceHealthCheckDeregister int
}