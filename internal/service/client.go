package service

import "github.com/tmrrwnxtsn/aero-table-booking-api/internal/store"

// ClientService представляет бизнес-логику работы с клиентами.
type ClientService interface {
}

type ClientServiceImpl struct {
	clientRepo store.ClientRepository
}

func NewClientService(clientRepo store.ClientRepository) ClientService {
	return &ClientServiceImpl{clientRepo: clientRepo}
}
