package command

import (
	"obsidian/log"
	"strings"
)

type Command struct {
	Name         string
	Aliases      []string
	PlayerOnly   bool
	OperatorOnly bool
	Execute      func(ctx CommandContext)
}

type CommandContext struct {
	Arguments []string
	Manager   *CommandManager
	Executor  any
}

func (c CommandContext) Reply(msg string) {
	if exe, ok := c.Executor.(interface{ SendMessage(string) }); ok {
		msgs := strings.Split(msg, "\n")
		for _, msg := range msgs {
			exe.SendMessage(msg)
		}
	} else {
		log.Print(msg)
	}
}

func NewManager(cmds ...Command) *CommandManager {
	return &CommandManager{commands: cmds}
}

type CommandManager struct {
	commands []Command
}

func (c *CommandManager) Search(cmd string) (Command, bool) {
	for _, command := range c.commands {
		if command.Name == cmd {
			return command, true
		}
		for _, a := range command.Aliases {
			if a == cmd {
				return command, true
			}
		}
	}
	return Command{}, false
}

func (c *CommandManager) Paginate(pageNumber, pageSize int) []Command {
	startIndex := (pageNumber - 1) * pageSize
	endIndex := startIndex + pageSize

	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex > len(c.commands) {
		startIndex = 0
	}
	if endIndex > len(c.commands) {
		endIndex = len(c.commands)
	}

	return c.commands[startIndex:endIndex]
}
