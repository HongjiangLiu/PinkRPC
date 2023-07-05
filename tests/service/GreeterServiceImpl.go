package service

type GreeterServiceImpl struct {
}

func (g *GreeterServiceImpl) SayHello(name string) (string, error) {
	return "Hello, " + name, nil
}
