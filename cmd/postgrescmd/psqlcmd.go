package postgrescmd

import (
	"fmt"
	"github.com/nais/cli/pkg/postgres"
	"github.com/urfave/cli/v2"
)

func psqlCommand() *cli.Command {
	return &cli.Command{
		Name:        "psql",
		Usage:       "Connect to the database using psql",
		Description: "Create a shell to the postgres instance by opening a proxy on a random port (see the proxy command for more info) and opening a psql shell.",
		ArgsUsage:   "appname",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
			},
		},
		Before: func(context *cli.Context) error {
			if context.Args().Len() != 1 {
				return fmt.Errorf("missing name of app")
			}

			return nil
		},
		Action: func(context *cli.Context) error {
			appName := context.Args().First()

			namespace := context.String("namespace")
			cluster := context.String("context")
			database := context.String("database")
			verbose := context.Bool("verbose")

			return postgres.RunPSQL(context.Context, appName, cluster, namespace, database, verbose)
		},
	}
}
