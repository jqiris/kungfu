package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	service := newMicroApp()
	app := &cli.App{
		Name: "kungfu",
		Before: func(c *cli.Context) error {
			return service.before()
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "data",
				Aliases: []string{"c"},
				Usage:   "set the data dir",
			},
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "build version",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "data",
				Usage: "set data dir",
				Action: func(c *cli.Context) error {
					return service.workDir(c)
				},
			},
			{
				Name:  "version",
				Usage: "view or set build version",
				Action: func(c *cli.Context) error {
					return service.version(c)
				},
			},
			{
				Name:  "build",
				Usage: "build servers",
				Action: func(c *cli.Context) error {
					return service.build(c)
				},
			},
			{
				Name:  "run",
				Usage: "run servers",
				Action: func(c *cli.Context) error {
					return service.run(c)
				},
			},
			{
				Name:  "start",
				Usage: "start servers",
				Action: func(c *cli.Context) error {
					return service.start(c)
				},
			},
			{
				Name:  "stop",
				Usage: "stop servers",
				Action: func(c *cli.Context) error {
					return service.stop(c)
				},
			},
			{
				Name:  "rm",
				Usage: "rm servers",
				Action: func(c *cli.Context) error {
					return service.rm(c)
				},
			},
			{
				Name:  "rmi",
				Usage: "rm servers images",
				Action: func(c *cli.Context) error {
					return service.rmi(c)
				},
			},
			{
				Name:  "clear",
				Usage: "clear servers",
				Action: func(c *cli.Context) error {
					return service.clear(c)
				},
			},
		}}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
