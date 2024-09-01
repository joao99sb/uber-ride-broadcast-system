package types

type IController interface {
	HandleDriver(client IClient)
}
