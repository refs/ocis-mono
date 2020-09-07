package command

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/micro/cli/v2"
	"github.com/refs/ocis-mono/pkg/config"
	"github.com/refs/ocis-mono/pkg/register"
)

// KillCommand is the entrypoint for the kill command.
func KillCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "kill",
		Usage:    "Kill an extension by name",
		Category: "Runtime",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "hostname",
				Value:       "localhost",
				EnvVars:     []string{"OCIS_RUNTIME_HOSTNAME"},
				Destination: &cfg.Runtime.Hostname,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       "10666",
				EnvVars:     []string{"OCIS_RUNTIME_PORT"},
				Destination: &cfg.Runtime.Port,
			},
		},
		Action: func(c *cli.Context) error {
			client, err := rpc.DialHTTP("tcp", net.JoinHostPort(cfg.Runtime.Hostname, cfg.Runtime.Port))
			if err != nil {
				log.Fatal("dialing:", err)
			}

			var arg1 int

			if err := client.Call("Service.Kill", os.Args[2], &arg1); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("process %v terminated", os.Args[2])

			return nil
		},
	}
}

func init() {
	register.AddCommand(KillCommand)
}
