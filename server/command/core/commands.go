package core

import "obsidian/server/command"

var Manager = command.NewManager(ragequit_cmd, help_cmd)
