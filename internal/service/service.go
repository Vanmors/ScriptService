package service

import (
	"ScriptService/internal/model"
	"ScriptService/internal/repository"
	"sync"
)

//go:generate mockgen -source=service.go -destination=mock/mock.go

type CommandServ interface {
	CreateCommand(cmd model.Command) (model.Command, error)
	GetAllCommands() ([]string, error)
	ExecuteCommand(contextCommand model.ContextCommand, cmd model.Command) (model.Command, error)
	GetCommand(cmdId int) (string, error)
	CancelCommand(contextMap *sync.Map, cmdId int) error
}

type Services struct {
	Command CommandServ
}

func NewServices(Repos *repository.Repositories) (*Services, error) {
	return &Services{
		Command: NewCommandService(Repos),
	}, nil
}
