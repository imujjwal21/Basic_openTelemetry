package users

import "context"

type User struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type Storage interface {
	Create(ctx context.Context, name, password string) error
	Check(ctx context.Context, name, password string) error
}
