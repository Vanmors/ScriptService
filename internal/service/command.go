package service

import (
	"BannerService/internal/model"
	"BannerService/internal/repository"
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
	//result, err := cs.executeCommand(cmd.Command)
	//if err != nil {
	//	return model.Command{}, err
	//}

	//cmd.Result = result

	cmd.Result = ""
	id, err := cs.Repos.Command.CreateCommand(cmd)
	cmd.ID = id
	if err != nil {
		log.Fatal(err)
		return model.Command{}, err
	}

	return cmd, nil
}

func (cs *CommandService) executeCommand(cmdStr string) (string, error) {
	cmd := exec.Command("bash", "-c", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (cs *CommandService) ExecuteCommandAsync(cmd model.Command) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		outChan := make(chan string)
		errChan := make(chan error)

		cmdId := cmd.ID
		cmdStr := cmd.Command
		cmdResult := ""

		cmd := exec.Command("bash", "-c", cmdStr)
		cmd.Stderr = nil

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("Error creating stdout pipe:", err)
			errChan <- err
			return
		}

		if err := cmd.Start(); err != nil {
			log.Println("Error starting command:", err)
			errChan <- err
			return
		}

		go func() {
			buf := make([]byte, 1024)
			output := ""
			for {
				n, err := stdout.Read(buf)
				if err != nil {
					close(outChan)
					return
				}
				fmt.Println(string(buf[:n]))
				output += string(buf[:n])
				cs.Repos.Command.UpdateCommand(cmdId, output)
				log.Println("updated", cmdId)
				outChan <- string(buf[:n])
			}
		}()

		for {
			select {
			case output, ok := <-outChan:
				if !ok {
					//cmd.Wait()
					cmdResult += "Command execution completed."
					//cs.Repos.Command.UpdateCommand(1, cmd.String())
					//wg.Done()
					return
				}
				cmdResult += output
				//cs.Repos.Command.UpdateCommand(1, cmd.String())
			case err := <-errChan:
				cmdResult += fmt.Sprintf("Error: %v", err)
				//cs.Repos.Command.UpdateCommand(1, cmd.String())
				//wg.Done()
				return
			}
		}
	}()

	wg.Wait()
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
