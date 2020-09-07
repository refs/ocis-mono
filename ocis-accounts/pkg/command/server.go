package command

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/flagset"
	"github.com/owncloud/ocis-accounts/pkg/metrics"
	"github.com/owncloud/ocis-accounts/pkg/server/grpc"
	"github.com/owncloud/ocis-accounts/pkg/server/http"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
)

// Server is the entry point for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Start ocis accounts service",
		Description: "uses an LDAP server as the storage backend",
		Flags:       flagset.ServerWithConfig(cfg),
		Before: func(ctx *cli.Context) error {
			if cfg.HTTP.Root != "/" {
				cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
			}

			// When running on single binary mode the before hook from the root command won't get called. We manually
			// call this before hook from ocis command, so the configuration can be loaded.
			return ParseConfig(ctx, cfg)
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				mtrcs       = metrics.New()
			)

			defer cancel()

			{
				server := http.Server(
					http.Logger(logger),
					http.Name(cfg.Server.Name),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					http.Flags(flagset.RootWithConfig(cfg)),
					http.Flags(flagset.ServerWithConfig(cfg)),
				)

				gr.Add(server.Run, func(_ error) {
					logger.Info().
						Str("server", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server := grpc.Server(
					grpc.Logger(logger),
					grpc.Name(cfg.Server.Name),
					grpc.Context(ctx),
					grpc.Config(cfg),
					grpc.Metrics(mtrcs),
				)

				gr.Add(func() error {
					logger.Info().Str("service", server.Name()).Msg("Reporting settings bundles to settings service")
					go svc.RegisterSettingsBundles(&logger)
					go svc.RegisterPermissions(&logger)
					return server.Run()
				}, func(_ error) {
					logger.Info().
						Str("server", "grpc").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
