package service

import (
	"ScriptService/internal/model"
	"ScriptService/internal/repository"
)

type Command interface {
	CreateCommand(cmd model.Command) (model.Command, error)
	GetAllCommands() ([]string, error)
	ExecuteCommand(cmd model.Command) (model.Command, error)
	GetCommand(cmdId int) (string, error)
}

type Services struct {
	Command Command
}

func NewServices(Repos *repository.Repositories) (*Services, error) {
	return &Services{
		Command: NewCommandService(Repos),
	}, nil
}
