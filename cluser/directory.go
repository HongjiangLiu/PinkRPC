package cluser

import (
	"pinkrpc/common"
	"pinkrpc/protocol"
)

// Directory 定义服务提供者（Provider）的目录信息
// 主要用于服务消费者（Consumer）在进行服务发现时，从注册中心或者服务目录中获取服务提供者的信息，以便进行服务调用
type Directory interface {
	common.Node
	List(invocation protocol.Invocation) []protocol.Invoker
}
