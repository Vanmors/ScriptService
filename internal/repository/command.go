package repository

import (
	"ScriptService/internal/model"
	"database/sql"
)

type CommandRepository struct {
	conn *sql.DB
}

func NewCommandRepository(db *sql.DB) *CommandRepository {
	return &CommandRepository{conn: db}
}

func (cr *CommandRepository) CreateCommand(cmd model.Command) (int, error) {
	var id int
	err := cr.conn.QueryRow("INSERT INTO commands (command, result) VALUES ($1, $2) RETURNING id", cmd.Command, cmd.Result).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (cr *CommandRepository) UpdateCommand(cmdId int, cmdResult string) error {
	_, err := cr.conn.Exec("UPDATE commands SET result=$1 WHERE id=$2", cmdResult, cmdId)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CommandRepository) GetAllCommand() ([]model.Command, error) {
	rows, err := cr.conn.Query("SELECT * FROM commands")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var commands []model.Command
	for rows.Next() {
		command := model.Command{}
		err = rows.Scan(&command.ID, &command.Command, &command.Result)
		if err != nil {
			return nil, err
		}
		commands = append(commands, command)
	}
	return commands, nil
}

func (cr *CommandRepository) DeleteCommand(cmdId int) error {
	_, err := cr.conn.Exec("DELETE FROM commands WHERE id=$1", cmdId)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CommandRepository) GetCommand(cmdId int) (model.Command, error) {
	var command model.Command
	err := cr.conn.QueryRow("SELECT * FROM commands WHERE id=$1", cmdId).Scan(&command.ID, &command.Command, &command.Result)
	if err != nil {
		return model.Command{}, err
	}
	return command, nil
}
