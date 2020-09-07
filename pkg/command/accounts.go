// +build !simple

package command

import (
	"github.com/micro/cli/v2"
	"github.com/refs/ocis-mono/ocis-accounts/pkg/command"
	svcconfig "github.com/refs/ocis-mono/ocis-accounts/pkg/config"
	"github.com/refs/ocis-mono/ocis-accounts/pkg/flagset"
	"github.com/refs/ocis-mono/pkg/config"
	"github.com/refs/ocis-mono/pkg/register"
)

// AccountsCommand is the entrypoint for the accounts command.
func AccountsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "accounts",
		Usage:    "Start accounts server",
		Category: "Extensions",
		Flags:    flagset.ServerWithConfig(cfg.Accounts),
		Subcommands: []*cli.Command{
			command.ListAccounts(cfg.Accounts),
			command.AddAccount(cfg.Accounts),
			command.UpdateAccount(cfg.Accounts),
			command.RemoveAccount(cfg.Accounts),
			command.InspectAccount(cfg.Accounts),
		},
		Action: func(c *cli.Context) error {
			accountsCommand := command.Server(configureAccounts(cfg))
			if err := accountsCommand.Before(c); err != nil {
				return err
			}

			return cli.HandleAction(accountsCommand.Action, c)
		},
	}
}

func configureAccounts(cfg *config.Config) *svcconfig.Config {
	cfg.Accounts.Log.Level = cfg.Log.Level
	cfg.Accounts.Log.Pretty = cfg.Log.Pretty
	cfg.Accounts.Log.Color = cfg.Log.Color

	return cfg.Accounts
}

func init() {
	register.AddCommand(AccountsCommand)
}
