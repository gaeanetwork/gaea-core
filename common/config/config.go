package config

// variables
var (
	// TODO - read in config
	ListenAddr     = ":12666"
	GRPCAddr       = ":12667"
	PProfAddr      = ":12668"
	ProfileEnabled = false

	// Max send and receive bytes for grpc clients and servers
	MaxRecvMsgSize = 100 * 1024 * 1024
	MaxSendMsgSize = 100 * 1024 * 1024
)
