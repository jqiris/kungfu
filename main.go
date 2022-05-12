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
				Aliases: []string{"d"},
				Usage:   "set the data dir",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "set the config file",
			},
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "build version",
			},
			&cli.StringFlag{
				Name:    "network",
				Aliases: []string{"n"},
				Usage:   "run network",
			},
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "run prefix",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "data",
				Usage: "view or set data dir",
				Action: func(c *cli.Context) error {
					return service.workDir(c)
				},
			},
			{
				Name:  "config",
				Usage: "view or set config path",
				Action: func(c *cli.Context) error {
					return service.config(c)
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
				Name:  "network",
				Usage: "view or set network",
				Action: func(c *cli.Context) error {
					return service.netView(c)
				},
			},
			{
				Name:  "prefix",
				Usage: "run prefix",
				Action: func(c *cli.Context) error {
					return service.runPrefix(c)
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
				Name:  "save",
				Usage: "save images",
				Action: func(c *cli.Context) error {
					return service.save(c)
				},
			},
			{
				Name:  "load",
				Usage: "load images",
				Action: func(c *cli.Context) error {
					return service.load(c)
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
