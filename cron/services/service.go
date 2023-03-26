package services

type ServiceContainer struct {
	HelloService *HelloService
}

func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		HelloService: GetHelloService(),
	}
}
