package main

import (
	"pinkrpc/config"
)

// go-client
func main() {
	//TODO 实现测试用例  ref:https://github.com/apache/dubbo-go-samples/blob/master/helloworld/go-client/cmd/client.go
	//config.SetConsumerService(service.GreeterServiceImpl{})
	//if err := config.Load(); err != nil {
	//	panic(err)
	//}
	config.SetConsumerService(nil)
	config.Load()

	//req := &api.HelloRequest{
	//	Name: "laurence",
	//}
	//reply, err := grpcGreeterImpl.SayHello(context.Background(), req)

}
