package types

import "context"

type IStreamServer interface {
	Run(ctx context.Context)
}
type IClient interface {
	GetId() string
	GetMessageChann() <-chan []byte
}
