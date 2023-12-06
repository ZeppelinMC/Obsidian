package command

import "log"

type Command struct {
	Name       string
	Aliases    []string
	PlayerOnly bool
	Execute    func(ctx CommandContext)
}

type CommandContext struct {
	Arguments []string
	Executor  any
}

func (c CommandContext) Reply(msg string) {
	if exe, ok := c.Executor.(interface{ SendMessage(string) }); ok {
		exe.SendMessage(msg)
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
