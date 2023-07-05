package common

// Node 节点
type Node interface {
	GetUrl() URL
	IsAvailable() bool
	Destroy()
}
