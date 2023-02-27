package client

import (
	"errors"

	"github.com/google/uuid"
)

type Client struct {
	id   uuid.UUID
	name string
}

func NewClient(name string) (*Client, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	return &Client{
		id:   uuid.New(),
		name: name,
	}, nil
}
func NewClientWithId(id uuid.UUID, name string) (*Client, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	return &Client{
		id:   id,
		name: name,
	}, nil
}
func (c *Client) ID() uuid.UUID {
	return c.id
}
func (c *Client) Name() string {
	return c.name
}
