/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

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
				Name:    "labelAuthor",
				Aliases: []string{"la"},
				Usage:   "set the label author",
			},
			&cli.StringFlag{
				Name:    "labelVersion",
				Aliases: []string{"lv"},
				Usage:   "set the label docker version",
			},
			&cli.StringFlag{
				Name:    "goVersion",
				Aliases: []string{"gv"},
				Usage:   "set the golang version",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "set the config file",
			},
			&cli.StringFlag{
				Name:    "remoteConfig",
				Aliases: []string{"rc"},
				Usage:   "set the remote config file",
			},
			&cli.StringFlag{
				Name:    "buildPath",
				Aliases: []string{"bp"},
				Usage:   "set the build path",
			},
			&cli.StringFlag{
				Name:    "project",
				Aliases: []string{"p"},
				Usage:   "set the project",
			},
			&cli.StringFlag{
				Name:    "memory",
				Aliases: []string{"m"},
				Usage:   "set the run memory",
			},
			&cli.StringFlag{
				Name:    "memory-swap",
				Aliases: []string{"ms"},
				Usage:   "set the run memory-swap",
			},
			&cli.StringFlag{
				Name:    "kernel-memory",
				Aliases: []string{"mk"},
				Usage:   "set the run kernel-memory",
			},
			&cli.StringFlag{
				Name:    "cpus",
				Aliases: []string{"cp"},
				Usage:   "set the run cpu num",
			},
			&cli.StringFlag{
				Name:    "cpuset-cpus",
				Aliases: []string{"cps"},
				Usage:   "set the run cpuset-cpus",
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
				Aliases: []string{"pf"},
				Usage:   "run prefix",
			},
			&cli.StringFlag{
				Name:    "registry",
				Aliases: []string{"r"},
				Usage:   "remote registry",
			},
			&cli.StringFlag{
				Name:    "exclude",
				Aliases: []string{"ex"},
				Usage:   "exclude server",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "appName",
				Usage: "view or set app name",
				Action: func(c *cli.Context) error {
					return service.appName(c)
				},
			},
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
				Name:  "remoteConfig",
				Usage: "view or set remote config path",
				Action: func(c *cli.Context) error {
					return service.remoteConfig(c)
				},
			},
			{
				Name:  "project",
				Usage: "view or set project",
				Action: func(c *cli.Context) error {
					return service.projectSet(c)
				},
			},
			{
				Name:  "labelAuthor",
				Usage: "view or set label author",
				Action: func(c *cli.Context) error {
					return service.labelAuthorSet(c)
				},
			},
			{
				Name:  "alpineVersion",
				Usage: "view or set alpine version",
				Action: func(c *cli.Context) error {
					return service.alpineVersionSet(c)
				},
			},
			{
				Name:  "goVersion",
				Usage: "view or set go version",
				Action: func(c *cli.Context) error {
					return service.goVersionSet(c)
				},
			},
			{
				Name:  "registry",
				Usage: "view or set registry addr",
				Action: func(c *cli.Context) error {
					return service.registry(c)
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
				Name:  "restart",
				Usage: "restart servers",
				Action: func(c *cli.Context) error {
					return service.restart(c)
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
				Name:  "prune",
				Usage: "rm none images",
				Action: func(c *cli.Context) error {
					return service.prune(c)
				},
			},
			{
				Name:  "clear",
				Usage: "clear servers",
				Action: func(c *cli.Context) error {
					return service.clear(c)
				},
			},
			{
				Name:  "remote",
				Usage: "remote registry operate",
				Before: func(c *cli.Context) error {
					return service.registryBefore(c)
				},
				Subcommands: []*cli.Command{
					{
						Name:  "push",
						Usage: "push local images to registry",
						Action: func(c *cli.Context) error {
							return service.registryPush(c)
						},
					},
					{
						Name:  "pull",
						Usage: "pull registry images to local",
						Action: func(c *cli.Context) error {
							return service.registryPull(c)
						},
					},
				},
			},
		}}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
