// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/refs/ocis-mono/ocis-migration/pkg/command"
	toolconfig "github.com/refs/ocis-mono/ocis-migration/pkg/config"
	"github.com/refs/ocis-mono/ocis-migration/pkg/flagset"
	"github.com/refs/ocis-mono/pkg/config"
	"github.com/refs/ocis-mono/pkg/register"
)

// ImportCommand is the entrypoint for the accounts command.
func ImportCommand(cfg *config.Config) *cli.Command {
	tc := toolconfig.New()
	return &cli.Command{
		Name:  "import",
		Usage: "Import a user exported by owncloud/data_exporter",
		Flags: flagset.ImportWithConfig(tc),
		Action: func(c *cli.Context) error {
			importCommand := command.Import(tc)
			return cli.HandleAction(importCommand.Action, c)
		},
	}
}

func init() {
	register.AddCommand(ImportCommand)
}
