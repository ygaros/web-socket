package query

import (
	"context"
	"sender/app/decorator"
	"sender/domain"
	"sender/domain/client"
)

// empty interface because this query doesnt need any input data
type ClientListHandler decorator.QueryHandler[interface{}, []domain.ClientDTO]

type clientListHandler struct {
	cRepo client.Repository
}

func (h clientListHandler) Handle(ctx context.Context, _ interface{}) ([]domain.ClientDTO, error) {
	clients, err := h.cRepo.GetAllClients(ctx)
	if err != nil {
		return nil, err
	}
	var parsedClients []domain.ClientDTO
	for _, c := range clients {
		dto := domain.ToClientDTO(c)
		parsedClients = append(parsedClients, dto)
	}
	return parsedClients, nil
}
func NewClientListHandler(cRepo client.Repository) ClientListHandler {
	if cRepo == nil {
		panic("empty client repository")
	}
	return decorator.NewQueryHandlerWithDefaultDecorators[interface{}, []domain.ClientDTO](clientListHandler{cRepo})
}
