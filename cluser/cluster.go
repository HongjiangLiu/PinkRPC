package cluser

import "pinkrpc/protocol"

type Cluster interface {
	Join(Directory) protocol.Invoker
}
