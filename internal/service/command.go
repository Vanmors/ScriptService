package service

import (
	"ScriptService/internal/model"
	"ScriptService/internal/repository"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type CommandService struct {
	Repos *repository.Repositories
}

func NewCommandService(repos *repository.Repositories) *CommandService {
	return &CommandService{Repos: repos}
}

func (cs *CommandService) CreateCommand(cmd model.Command) (model.Command, error) {
	cmd.Result = ""
	id, err := cs.Repos.Command.CreateCommand(cmd)
	cmd.ID = id
	if err != nil {
		return model.Command{}, err
	}

	return cmd, nil
}

func (cs *CommandService) ExecuteCommand(contextCommand model.ContextCommand, cmd model.Command) (model.Command, error) {
	command := exec.Command("bash", "-c", cmd.Command)
	command.Stderr = nil

	stdout, err := command.StdoutPipe()
	if err != nil {
		log.Println("Error creating stdout pipe:", err)
		return model.Command{}, err
	}

	if err := command.Start(); err != nil {
		log.Println("Error starting command:", err)
		return model.Command{}, err
	}

	buf := make([]byte, 1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			break
		}
		fmt.Println(string(buf[:n]))
		cmd.Result += string(buf[:n])

		select {
		case <-contextCommand.Ctx.Done():
			return model.Command{}, err
		default:
			err = cs.Repos.Command.UpdateCommand(cmd.ID, cmd.Result)
			if err != nil {
				return model.Command{}, err
			}

			log.Println("updated", cmd.ID)
		}
	}
	return cmd, err
}

func (cs *CommandService) GetAllCommands() ([]string, error) {
	commands, err := cs.Repos.Command.GetAllCommand()
	if err != nil {
		return nil, err
	}
	var nameCommands []string
	for i := range commands {
		nameCommands = append(nameCommands, commands[i].Command)
	}

	return nameCommands, nil
}

func (cs *CommandService) GetCommand(cmdId int) (string, error) {
	command, err := cs.Repos.Command.GetCommand(cmdId)
	if err != nil {
		return "", err
	}
	return command.Command, nil
}

func (cs *CommandService) CancelCommand(contextMap *sync.Map, cmdId int) error {
	val, ok := contextMap.LoadAndDelete(cmdId)
	if ok {
		ctx, ok := val.(model.ContextCommand)
		if ok {
			ctx.Cancel()
			err := cs.Repos.Command.DeleteCommand(cmdId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
