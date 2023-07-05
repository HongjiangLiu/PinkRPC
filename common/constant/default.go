package constant

const (
	DEFAULT_WEIGHT = 100
	DEFAULT_WARMUP = 10 * 60 //Java中是 10 * 60 * 1000 是因为System.currentTimeMillis()得到的是毫秒， Go中time.Unix()得到的是秒
)

const (
	DEFAULT_LOADBALANCE = "random"
	DEFAULT_RETRIES     = 2
	DEFAULT_PROTOCOL    = "pink"
	DEFAULT_VERSION     = "1.0.0"
	DEFAULT_REG_TIMEOUT = "10s" //TODO REG是REGISTRY的缩写吗？
	DEFAULT_CLUSTER     = "failover"
)

const ( //TODO?这些常量是干啥的？
	DEFAULT_KEY               = "default"
	DEFAULT_SERVICE_FILTERS   = "echo"
	DEFAULT_REFERENCE_FILTERS = ""
	ECHO                      = "$echo"
)
