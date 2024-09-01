package internal

import (
	"encoding/json"
	"go-service/internal/entity"
	"go-service/pkg/types"
	"strings"
)

type IService interface {
	HandleDriverMsg(driverId string, client <-chan []byte)
}

type Service struct {
	queue types.IQueue
}

func NewService(queue types.IQueue) IService {

	s := &Service{
		queue: queue,
	}
	return s
}

func (s *Service) HandleDriverMsg(driveId string, clientMsg <-chan []byte) {
	dest := s.queue.BuildDefaultDestinationQueue()
	for msg := range clientMsg {

		jsonMsg := s.parserMsg(driveId, msg)

		s.queue.Notify(jsonMsg, dest)
	}
}

func (s *Service) parserMsg(driverId string, msg []byte) []byte {

	msgString := strings.Split(string(msg), ",")
	dest := entity.Destination{
		Order:     driverId,
		Latitude:  msgString[0],
		Longitude: msgString[1],
	}

	json, _ := json.Marshal(dest)

	return json

}
