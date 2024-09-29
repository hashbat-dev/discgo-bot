# Commands
Commands are user triggered events via the !command format.

## Creating a New Command
1. Create a new file in Bot/Bang/ for the command code. There is an Interface template to copy below.
1. Add to the Command map in Bot/Commands/commands.go => func init(), referencing your new Command.

```
package bang

import (
	"github.com/bwmarrin/discordgo"
	config "github.com/hashbat-dev/discgo-bot/Config"
)

type MyCommand struct{}

func (s MyCommand) Name() string {
	return "MyCommand"
}

func (s MyCommand) PermissionRequirement() int {
	return config.CommandLevelUser
}

func (s MyCommand) ProcessPool() config.ProcessPool {
	return config.ProcessPools[config.ProcessPoolText]
}

func (s MyCommand) LockedByDefault() bool {
	return true
}

func (s MyCommand) Execute(message *discordgo.MessageCreate, command string) error {

	return nil
}


```