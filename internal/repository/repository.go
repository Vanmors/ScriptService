package repository

import (
	"ScriptService/internal/model"
	"ScriptService/pkg/store"
)

type Repositories struct {
	Command Command
}

type Command interface {
	CreateCommand(cmd model.Command) (int, error)
	GetAllCommand() ([]model.Command, error)
	UpdateCommand(cmdId int, cmdResult string) error
	DeleteCommand(cmdId int) error
	GetCommand(cmdId int) (model.Command, error)
}

func NewRepositories(dbname, username, password, host, port string) (*Repositories, error) {
	db, err := store.NewClient(dbname, username, password, host, port)
	if err != nil {
		return nil, err
	}
	return &Repositories{Command: NewCommandRepository(db)}, nil
}
