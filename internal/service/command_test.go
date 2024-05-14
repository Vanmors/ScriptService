package service

import (
	"ScriptService/internal/model"
	"ScriptService/internal/repository"
	mock_repository "ScriptService/internal/repository/mock"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestCommandService_CreateCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	command := mock_repository.NewMockCommand(ctrl)
	repos := &repository.Repositories{Command: command}
	commandService := &Services{Command: NewCommandService(repos)}

	t.Run("Success", func(t *testing.T) {
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		command.EXPECT().CreateCommand(gomock.Any()).Return(1, nil)
		resultCmd, err := commandService.Command.CreateCommand(cmd)
		expectedCmd := model.Command{ID: 1, Command: "echo 'Hello, world!'", Result: ""}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, resultCmd, expectedCmd)
	})
	t.Run("Error_UpdateCommand", func(t *testing.T) {
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		command.EXPECT().CreateCommand(gomock.Any()).Return(0, errors.New("db error"))
		expectedCmd := model.Command{ID: 0, Command: "", Result: ""}
		resultCmd, err := commandService.Command.CreateCommand(cmd)
		if err == nil {
			t.Error("expected error, got nil")
		}

		require.Equal(t, resultCmd, expectedCmd)
	})
}

func TestCommandService_ExecuteCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	command := mock_repository.NewMockCommand(ctrl)
	repos := &repository.Repositories{Command: command}
	commandService := &Services{Command: NewCommandService(repos)}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Success", func(t *testing.T) {
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		command.EXPECT().UpdateCommand(gomock.Any(), gomock.Any()).Return(nil)
		resultCmd, err := commandService.Command.ExecuteCommand(model.ContextCommand{Ctx: ctx, Cancel: cancel}, cmd)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resultCmd.Command != cmd.Command {
			t.Errorf("expected command: %s, got: %s", cmd.Command, resultCmd.Command)
		}
	})
	t.Run("Error_UpdateCommand", func(t *testing.T) {
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		command.EXPECT().UpdateCommand(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
		resultCmd, err := commandService.Command.ExecuteCommand(model.ContextCommand{Ctx: ctx, Cancel: cancel}, cmd)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if resultCmd.ID != 0 || resultCmd.Command != "" {
			t.Error("expected empty command result, got non-empty")
		}
	})

	t.Run("Context_Cancel", func(t *testing.T) {
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		resultCmd, err := commandService.Command.ExecuteCommand(model.ContextCommand{Ctx: ctx, Cancel: cancel}, cmd)
		if err != nil {
			t.Error("expected nil, got nil")
		}
		if resultCmd.ID != 0 || resultCmd.Command != "" {
			t.Error("expected empty command result, got non-empty")
		}
	})
}

func TestCommandService_GetCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	command := mock_repository.NewMockCommand(ctrl)
	repos := &repository.Repositories{Command: command}
	commandService := &Services{Command: NewCommandService(repos)}

	t.Run("Success", func(t *testing.T) {
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		command.EXPECT().GetCommand(gomock.Any()).Return(cmd, nil)
		resultCmd, err := commandService.Command.GetCommand(1)
		expectedCmd := "echo 'Hello, world!'"
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, resultCmd, expectedCmd)
	})
	t.Run("Error_GetCommand", func(t *testing.T) {
		command.EXPECT().GetCommand(gomock.Any()).Return(model.Command{}, errors.New("db error"))
		resultCmd, err := commandService.Command.GetCommand(1)
		expectedCmd := ""
		if err == nil {
			t.Error("expected error, got nil")
		}

		require.Equal(t, resultCmd, expectedCmd)
	})
}

func TestCommandService_GetAllCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	command := mock_repository.NewMockCommand(ctrl)
	repos := &repository.Repositories{Command: command}
	commandService := &Services{Command: NewCommandService(repos)}

	t.Run("Success", func(t *testing.T) {
		cmd := []model.Command{{ID: 1, Command: "echo 'Hello, world!'"}, {ID: 2, Command: "ls -la"}}
		command.EXPECT().GetAllCommand().Return(cmd, nil)
		resultCmd, err := commandService.Command.GetAllCommands()
		expectedCmd := []string{"echo 'Hello, world!'", "ls -la"}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		require.Equal(t, resultCmd, expectedCmd)
	})
	t.Run("Error_GetAllCommand", func(t *testing.T) {
		command.EXPECT().GetAllCommand().Return(nil, errors.New("bd error"))
		resultCmd, err := commandService.Command.GetAllCommands()
		expectedCmd := []string(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}
		require.Equal(t, resultCmd, expectedCmd)
	})
}

func TestCommandService_CancelCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	command := mock_repository.NewMockCommand(ctrl)
	repos := &repository.Repositories{Command: command}
	commandService := &Services{Command: NewCommandService(repos)}

	var ContextMap sync.Map

	t.Run("Success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		contextCommand := model.ContextCommand{Ctx: ctx, Cancel: cancel}
		ContextMap.Store(cmd.ID, contextCommand)

		command.EXPECT().DeleteCommand(gomock.Any()).Return(nil)
		err := commandService.Command.CancelCommand(&ContextMap, cmd.ID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("Error_DeleteCommand", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		cmd := model.Command{ID: 1, Command: "echo 'Hello, world!'"}
		contextCommand := model.ContextCommand{Ctx: ctx, Cancel: cancel}
		ContextMap.Store(cmd.ID, contextCommand)

		command.EXPECT().DeleteCommand(gomock.Any()).Return(errors.New("db error"))
		err := commandService.Command.CancelCommand(&ContextMap, cmd.ID)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
