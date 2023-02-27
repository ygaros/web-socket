package client

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetAllClients(ctx context.Context) ([]Client, error)
	//In case of client get should also create if client doesnt exists
	GetClient(ctx context.Context, clientId uuid.UUID) (*Client, error)
	GetClientByName(ctx context.Context, clientName string) (*Client, error)
	UpdateClient(
		ctx context.Context,
		clientId uuid.UUID,
		updateFn func(c *Client) (*Client, error),
	) error
	SaveClient(ctx context.Context, client *Client) (*Client, error)
}
