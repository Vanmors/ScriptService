package transport

import (
	"ScriptService/internal/model"
	"ScriptService/internal/service"
	mock_service "ScriptService/internal/service/mock"
	"bytes"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestHandler_GetCommandById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCommand, commandId int)

	testTable := []struct {
		name                string
		inputBody           string
		inputId             int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{}`,
			inputId:   1,
			mockBehavior: func(s *mock_service.MockCommand, commandId int) {
				s.EXPECT().GetCommand(commandId).Return(model.Command{ID: 1, Command: "echo 'Hello, world!'", Result: "Hello, world!", Status: "Done"}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "{\"id\":1,\"command\":\"echo 'Hello, world!'\",\"result\":\"Hello, world!\",\"status\":\"Done\"}\n",
		},
		{
			name:      "Element does not exist",
			inputBody: `{}`,
			inputId:   1,
			mockBehavior: func(s *mock_service.MockCommand, commandId int) {
				s.EXPECT().GetCommand(commandId).Return(model.Command{}, errors.New("element does not exist"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			command := mock_service.NewMockCommand(c)
			testCase.mockBehavior(command, testCase.inputId)

			services := &service.Services{Command: command}
			handler := NewHandler(services)

			rr := httptest.NewRecorder()

			reqBody := bytes.NewBufferString(testCase.inputBody)
			req := httptest.NewRequest(http.MethodGet, "/command/1", reqBody)

			handler.GetCommandById(rr, req)

			require.Equal(t, testCase.expectedStatusCode, rr.Code)
			require.Equal(t, testCase.expectedRequestBody, rr.Body.String())
		})
	}
}

func TestHandler_GetAllCommands(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCommand)

	testTable := []struct {
		name                string
		inputBody           string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{}`,
			mockBehavior: func(s *mock_service.MockCommand) {
				s.EXPECT().GetAllCommands().Return([]string{"ls", "echo 'hello'"}, nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "[\"ls\",\"echo 'hello'\"]\n",
		},
		{
			name:      "Error",
			inputBody: `{}`,
			mockBehavior: func(s *mock_service.MockCommand) {
				s.EXPECT().GetAllCommands().Return(nil, errors.New("error"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			command := mock_service.NewMockCommand(c)
			testCase.mockBehavior(command)

			services := &service.Services{Command: command}
			handler := NewHandler(services)

			rr := httptest.NewRecorder()

			reqBody := bytes.NewBufferString(testCase.inputBody)
			req := httptest.NewRequest(http.MethodGet, "/commands", reqBody)

			handler.GetAllCommands(rr, req)

			require.Equal(t, testCase.expectedStatusCode, rr.Code)
			require.Equal(t, testCase.expectedRequestBody, rr.Body.String())
		})
	}
}

func TestHandler_CancelCommand(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCommand, commandId int, contextMap sync.Map)

	testTable := []struct {
		name                string
		inputBody           string
		inputId             int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{}`,
			inputId:   1,
			mockBehavior: func(s *mock_service.MockCommand, commandId int, contextMap sync.Map) {
				s.EXPECT().CancelCommand(&contextMap, commandId).Return(nil)
			},
			expectedStatusCode:  http.StatusOK,
			expectedRequestBody: "",
		},
		{
			name:      "Element does not exist",
			inputBody: `{}`,
			inputId:   1,
			mockBehavior: func(s *mock_service.MockCommand, commandId int, contextMap sync.Map) {
				s.EXPECT().CancelCommand(&contextMap, commandId).Return(errors.New("element does not exist"))
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			command := mock_service.NewMockCommand(c)
			testCase.mockBehavior(command, testCase.inputId, sync.Map{})

			services := &service.Services{Command: command}
			handler := NewHandler(services)

			rr := httptest.NewRecorder()

			reqBody := bytes.NewBufferString(testCase.inputBody)
			req := httptest.NewRequest(http.MethodDelete, "/cancelCommand/1", reqBody)

			handler.CancelCommand(rr, req)

			require.Equal(t, testCase.expectedStatusCode, rr.Code)
			require.Equal(t, testCase.expectedRequestBody, rr.Body.String())
		})
	}
}

func TestHandler_CreateCommand(t *testing.T) {
	type mockBehavior func(s *mock_service.MockCommand, command model.Command, contextCommand model.ContextCommand)

	testTable := []struct {
		name                string
		inputBody           string
		inputCommand        model.Command
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:         "OK",
			inputBody:    `{"command": "echo 'Hello, world!' && sleep 10 && echo 'Goodbye, world!'"}`,
			inputCommand: model.Command{Command: "echo 'Hello, world!' && sleep 10 && echo 'Goodbye, world!'"},
			mockBehavior: func(s *mock_service.MockCommand, command model.Command, contextCommand model.ContextCommand) {
				s.EXPECT().CreateCommand(command).Return(model.Command{ID: 1,
					Command: "echo 'Hello, world!' && sleep 10 && echo 'Goodbye, world!'",
					Result:  "",
					Status:  "In progress"}, nil)
				s.EXPECT().ExecuteCommand(gomock.Any(), gomock.Any()).Return(model.Command{}, nil)
			},
			expectedStatusCode:  http.StatusCreated,
			expectedRequestBody: "1\n",
		},
		{
			name:         "Element does not exist",
			inputBody:    `{"command": "echo 'Hello, world!' && sleep 10 && echo 'Goodbye, world!'"}`,
			inputCommand: model.Command{Command: "echo 'Hello, world!' && sleep 10 && echo 'Goodbye, world!'"},
			mockBehavior: func(s *mock_service.MockCommand, command model.Command, contextCommand model.ContextCommand) {
				s.EXPECT().CreateCommand(command).Return(
					model.Command{ID: 1,
						Command: "echo 'Hello, world!' && sleep 10 && echo 'Goodbye, world!'",
						Result:  "",
						Status:  "In progress"},
					errors.New("element"))
				s.EXPECT().ExecuteCommand(gomock.Any(), gomock.Any()).Return(model.Command{}, nil)
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedRequestBody: "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			command := mock_service.NewMockCommand(c)
			testCase.mockBehavior(command, testCase.inputCommand, model.ContextCommand{Ctx: ctx, Cancel: cancel})

			services := &service.Services{Command: command}
			handler := NewHandler(services)

			rr := httptest.NewRecorder()

			reqBody := bytes.NewBufferString(testCase.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/command", reqBody)

			handler.CreateCommand(rr, req)
			time.Sleep(1 * time.Second)
			require.Equal(t, testCase.expectedStatusCode, rr.Code)
			require.Equal(t, testCase.expectedRequestBody, rr.Body.String())
		})
	}
}
