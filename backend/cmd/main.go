package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/emochka2007/block-accounting/cmd/commands"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "biocom-ioannes",
		Version:  "0.0.1a",
		Commands: commands.Commands(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Value: "debug",
			},
			&cli.StringFlag{
				Name:  "address",
				Value: "localhost:8080",
			},
			&cli.StringFlag{
				Name: "db-host",
			},
			&cli.StringFlag{
				Name: "db-user",
			},
			&cli.StringFlag{
				Name: "db-password",
			},
		},
		Action: func(c *cli.Context) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()

			// todo build config

			// todo build service

			// todo run service

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
