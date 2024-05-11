package service

import (
	"BannerService/internal/model"
	"BannerService/internal/repository"
)

type Command interface {
	CreateCommand(cmd model.Command) (model.Command, error)
	executeCommand(cmdStr string) (string, error)
	GetAllCommands() ([]string, error)
	ExecuteCommandAsync(cmd model.Command)
}

type Services struct {
	Command Command
}

func NewServices(Repos *repository.Repositories) (*Services, error) {
	return &Services{
		Command: NewCommandService(Repos),
	}, nil
}
