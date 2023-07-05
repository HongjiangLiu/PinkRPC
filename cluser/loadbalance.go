package cluser

import "pinkrpc/protocol"

type LoadBalance interface {
	Select([]protocol.Invoker, protocol.Invocation) protocol.Invoker
}
