package service

type GreeterService interface {
	SayHello(name string) (string, error)
}
