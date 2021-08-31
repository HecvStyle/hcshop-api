package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
)

type RegisterClient struct {
	Host string
	Port int
}

func NewRegisterClient(host string, port int) RegisterClient {
	return RegisterClient{
		host,
		port,
	}
}

type ClientRegister interface {
	RegisterService(host, name string, port int, tags []string)
	DeRegisterService(name string)
}

func (r *RegisterClient) RegisterService(host, name string, port int, tags []string) (string, error) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID, _ = uuid.GenerateUUID()
	registration.Port = port
	registration.Tags = tags
	registration.Address = host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return registration.ID, err
}

func (r *RegisterClient) DeRegisterService(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	return client.Agent().ServiceDeregister(serviceId)
}
