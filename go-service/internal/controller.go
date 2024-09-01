package internal

import "go-service/pkg/types"

type Controller struct {
	Service IService
}

func NewController(service IService) *Controller {

	return &Controller{
		Service: service,
	}
}

func (c *Controller) HandleDriver(client types.IClient) {
	msgChann := client.GetMessageChann()

	c.Service.HandleDriverMsg(client.GetId(), msgChann)
}
