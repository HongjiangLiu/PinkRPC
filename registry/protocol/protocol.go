package protocol

import (
	"pinkrpc/protocol"
	"sync"
)

var regProtocol *registryProtocol

type registryProtocol struct {
	invokers []protocol.Invoker
	// Registry  Map<RegistryAddress, Registry>
	registries sync.Map
	bounds     sync.Map
}
