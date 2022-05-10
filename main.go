package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var (
	dockerPath = "Dockerfile"
	dockerTml  = `
FROM golang:1.18.1 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server

FROM alpine
ARG config_file
ARG run_server
ARG client_port
ENV TZ Asia/Shanghai
ENV run_mode = docker
ENV run_server = ${run_server}

COPY --from=builder /src/server /app

WORKDIR /app

EXPOSE ${client_port}
VOLUME /data/conf
VOLUME /data/logs
COPY ${config_file} /data/conf/config.json
ENTRYPOINT ["/app/server", "-conf", "/data/conf/config.json"]	
	`
)

func main() {
	app := &cli.App{
		Name: "kungfu",
		Before: func(c *cli.Context) error {
			if exist, _ := utils.PathExists(dockerPath); !exist {
				fp, err := os.Create(dockerPath)
				if err != nil {
					return err
				}
				defer fp.Close()
				_, err = fp.Write([]byte(dockerTml))
				if err != nil {
					return err
				}
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{"c"},
				Value:   "config.json",
				Usage:   "locate the config file",
			},
			&cli.StringFlag{
				Name:    "build_server",
				Aliases: []string{"bs"},
				Usage:   "build specify server",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "build servers",
				Action: func(c *cli.Context) error {
					cfg := c.String("conf")
					if len(cfg) < 1 {
						return errors.New("can't find config file")
					}
					viper.SetConfigFile(cfg)
					if err := viper.ReadInConfig(); err != nil {
						panic(err)
					}
					//frame init
					frameCfg := viper.Get("frame")
					if err := config.InitFrameConf(frameCfg); err != nil {
						logger.Fatal(err)
					}
					servers := config.GetServersConf()
					specialServer := c.String("build_server")
					if len(specialServer) < 1 {
						for _, server := range servers {
							if server.IsLaunch {
								if bs, err := buildServer(cfg, server); err != nil {
									return err
								} else {
									logger.Infof("server build result:%v", string(bs))
								}
							}
						}
					} else {
						if server, ok := servers[specialServer]; ok {
							if bs, err := buildServer(cfg, server); err != nil {
								return err
							} else {
								logger.Infof("server build result:%v", string(bs))
							}
						} else {
							logger.Fatalf("can't find the build server: %v", specialServer)
						}
					}
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func buildServer(cfg string, server *treaty.Server) ([]byte, error) {
	args := []string{"build", "--build-arg", fmt.Sprintf("config_file=%v", cfg), "--build-arg", fmt.Sprintf("run_server=%v", server.ServerId), "--build-arg", fmt.Sprintf(`client_port=%v`, server.ClientPort), "-t", server.ServerName, "."}
	cmd := exec.Command("docker", args...)
	logger.Info(cmd.String())
	return cmd.Output()
}
