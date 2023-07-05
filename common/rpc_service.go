package common

type RPCService interface {
	Service() string // Path InterfaceName
	Version() string
}
