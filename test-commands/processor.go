package test

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// ProcessorFunc is a type which wraps functions for command processing
type ProcessorFunc func(string) error

func sleep(s string) error {
	log.Infof("sleep(%s)", s)
	t, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("sleep: %s ", err.Error())
	}
	time.Sleep(time.Second * time.Duration(t))
	return nil
}

// NewCommandsRegistry is a constructor for CommandsRegistry struct
func NewCommandsRegistry() (cR CommandsRegistry) {
	cR = CommandsRegistry{
		cmds: make(map[string]ProcessorFunc),
	}

	cR.cmds["sleep"] = sleep
	return

}

// CommandsRegistry is a struct with all registered commands' functions
type CommandsRegistry struct {
	cmds map[string]ProcessorFunc
}

// TODO: replace string signature to []slice
func (c *CommandsRegistry) execute(command, signature string) error {
	if _, ok := c.cmds[command]; !ok {
		return errors.New("No such command")
	}
	err := c.cmds[command](signature)
	if err != nil {
		return err
	}
	return nil
}

// RegisterCommand registers new command function
func (c *CommandsRegistry) RegisterCommand(name string, fn ProcessorFunc) {
	c.cmds[name] = fn
}

// CommandExecutor is an interface for CommandsRegistry; wraps execute function
type CommandExecutor interface {
	execute(command, signature string) error
}

// CommandRegistrar is an interface for CommandsRegistry; wraps RegisterCommand
type CommandRegistrar interface {
	RegisterCommand(string, ProcessorFunc)
}
