package node

// Config represents a small collection of configuration values for node.
// These values can be further extended by all registered services.
type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	DataDir string `json:"dataDir"`
}
