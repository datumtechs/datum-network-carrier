package node

// Config represents a small collection of configuration values for node.
// These values can be further extended by all registered services.
type Config struct {
	DataDir     string
	HTTPHost    string
	HTTPPort    uint64
	HTTPCors    []string
	HTTPModules []string
}
