package router

//TODO to pkg
import (
	"errors"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrCommandNotFound    = errors.New("command was not found")
	ErrInvalidMessage     = errors.New("invalid message")
	ErrNotEnoughArguments = errors.New("not enouth arguments")
	ErrSoManyArguments    = errors.New("so many arguments")
)

type CommandHandler func(message *tgbotapi.Message)

type Command struct {
	ArgumentsCount int
	CommandHandler CommandHandler
}

type CommandRouter struct {
	commands map[string]Command
}

func NewRouter() *CommandRouter {
	return &CommandRouter{commands: make(map[string]Command)}
}

func (c *CommandRouter) AddCommand(prefix string, command Command) {
	c.commands[prefix] = command
}

func (c *CommandRouter) Route(message *tgbotapi.Message) error {
	args := strings.Split(message.Text, " ")
	if len(args) == 0 {
		return ErrInvalidMessage
	}

	command, ok := c.commands[args[0]]
	if !ok {
		return ErrCommandNotFound
	}
	if len(args) > command.ArgumentsCount+1 {
		return ErrSoManyArguments
	} else if len(args) < command.ArgumentsCount+1 {
		return ErrNotEnoughArguments
	}

	command.CommandHandler(message)
	return nil
}
