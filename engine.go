package cmd

import (
	"fmt"

	shellwords "github.com/mattn/go-shellwords"
)

// Handler is a function that gets invoked when a command is run.
type Handler func(ctx interface{}, args []string) error

type commandInfo struct {
	handler   Handler
	userLevel int
	help      string
}

// Engine is the core interface to the command interpreter.
type Engine struct {
	commands map[string]*commandInfo
	parser   *shellwords.Parser
}

// NewEngine creates a new command interpreter.
func NewEngine() *Engine {
	return &Engine{
		commands: make(map[string]*commandInfo),
		parser:   shellwords.NewParser(),
	}
}

// AddCommand adds a command to the engine.
func (e *Engine) AddCommand(name string, help string, handler Handler, userLevel int) error {
	if _, ok := e.commands[name]; ok {
		return fmt.Errorf("command \"%s\" already registered", name)
	}

	e.commands[name] = &commandInfo{
		handler:   handler,
		userLevel: userLevel,
		help:      help,
	}

	return nil
}

// RemoveCommand removes a command from the engine.
func (e *Engine) RemoveCommand(name string) error {
	if _, ok := e.commands[name]; !ok {
		return fmt.Errorf("command \"%s\" not registered", name)
	}

	delete(e.commands, name)
	return nil
}

// Exec executes a command.
func (e *Engine) Exec(ctx interface{}, userLevel int, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("command args need to be > 0")
	}

	cmd, ok := e.commands[args[0]]

	if !ok {
		return fmt.Errorf("command \"%s\" not found", args[0])
	}

	if cmd.userLevel > userLevel {
		return fmt.Errorf("user level %d not >= %d", userLevel, cmd.userLevel)
	}

	return cmd.handler(ctx, args[1:])
}

// ExecString parses the commandString and executes the command.
func (e *Engine) ExecString(ctx interface{}, userLevel int, commandString string) error {
	args, err := e.parser.Parse(commandString)
	if err != nil {
		return err
	}

	return e.Exec(ctx, userLevel, args)
}
