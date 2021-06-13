package node

type Config struct {
	DataDir  string
	HTTPHost string
	HTTPPort uint64
	HTTPCors []string
	HTTPModules []string
}
