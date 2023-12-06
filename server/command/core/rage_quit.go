package core

import (
	"obsidian/server/command"
	"obsidian/server/player"
)

var ragequit_cmd = command.Command{
	Name:       "ragequit",
	Aliases:    []string{"rq"},
	PlayerOnly: true,
	Execute: func(ctx command.CommandContext) {
		ctx.Executor.(*player.Player).Disconnect("RAGE QUIT!!!")
	},
}
