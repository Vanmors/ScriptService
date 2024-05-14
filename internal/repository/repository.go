package repository

import (
	"ScriptService/internal/model"
	"ScriptService/pkg/store"
)

//go:generate mockgen -source=repository.go -destination=mock/mock.go

type Command interface {
	CreateCommand(cmd model.Command) (int, error)
	GetAllCommand() ([]model.Command, error)
	UpdateCommand(cmdId int, cmdResult string) error
	DeleteCommand(cmdId int) error
	GetCommand(cmdId int) (model.Command, error)
}

type Repositories struct {
	Command Command
}

func NewRepositories(dbname, username, password, host, port string) (*Repositories, error) {
	db, err := store.NewClient(dbname, username, password, host, port)
	if err != nil {
		return nil, err
	}
	return &Repositories{Command: NewCommandRepository(db)}, nil
}
