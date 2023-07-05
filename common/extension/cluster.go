package extension

import cluster "pinkrpc/cluser"

var (
	clusters = make(map[string]func() cluster.Cluster)
)
