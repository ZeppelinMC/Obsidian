package core

import (
	"fmt"
	"obsidian/server/command"
	"strconv"
)

var help_cmd = command.Command{
	Name: "help",
	Execute: func(ctx command.CommandContext) {
		page := uint(1)
		if len(ctx.Arguments) > 0 {
			if n, err := strconv.ParseUint(ctx.Arguments[0], 10, 0); err == nil {
				page = uint(n)
			}
		}

		cmds := ctx.Manager.Paginate(int(page), 5)

		str := fmt.Sprintf("&fCommands (page %d)\n", page)

		for i, cmd := range cmds {
			str += fmt.Sprintf("&f/&e%s", cmd.Name)
			if i < len(cmds)-1 {
				str += "\n"
			}
		}

		ctx.Reply(str)
	},
}
