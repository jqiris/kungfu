package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

var (
	storePath  = ".ini"
	dockerVer  = "DockerVer"
	dockerData = "DockerData"
	dockerPath = "Dockerfile"
	dockerTml  = `
FROM golang:1.18.1 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server

FROM alpine
ARG run_server
ARG client_port
ENV TZ Asia/Shanghai
ENV run_mode docker
ENV run_server ${run_server}

COPY --from=builder /src/server /app/

WORKDIR /app

EXPOSE ${client_port}
VOLUME /data
ENTRYPOINT ["/app/server", "-conf", "/data/conf/config.json"]	
	`
)

type MicroApp struct {
	ver   string
	data  string
	store *ini.File
}

func newMicroApp() *MicroApp {
	return &MicroApp{}
}

func (m *MicroApp) clear(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				m.clearServer(m.ver, server)
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			m.clearServer(m.ver, server)
		} else {
			log.Fatalf("can't find the server: %v\n", specialServer)
		}
	}
	return nil
}

func (m *MicroApp) rmi(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.rmiServer(m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("image rm result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.rmiServer(m.ver, server); err != nil {
				return err
			} else {
				fmt.Printf("image rm result:%v\n", string(bs))
			}
		} else {
			log.Fatalf("can't find the rm image: %v\n", specialServer)
		}
	}
	return nil
}

func (m *MicroApp) rm(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.rmServer(m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server rm result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.rmServer(m.ver, server); err != nil {
				return err
			} else {
				fmt.Printf("server rm result:%v\n", string(bs))
			}
		} else {
			log.Fatalf("can't find the rm server: %v\n", specialServer)
		}
	}
	return nil
}
func (m *MicroApp) stop(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.stopServer(m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server stop result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.stopServer(m.ver, server); err != nil {
				return err
			} else {
				fmt.Printf("server stop result:%v\n", string(bs))

			}
		} else {
			log.Fatalf("can't find the stop server: %v", specialServer)
		}
	}
	return nil
}

func (m *MicroApp) run(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.runServer(m.data, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server run result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.runServer(m.data, m.ver, server); err != nil {
				return err
			} else {
				fmt.Printf("server run result:%v\n", string(bs))
			}
		} else {
			log.Fatalf("can't find the run server: %v", specialServer)
		}
	}
	return nil
}

func (m *MicroApp) start(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.startServer(m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server start result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.startServer(m.ver, server); err != nil {
				return err
			} else {
				fmt.Printf("server start result:%v\n", string(bs))
			}
		} else {
			log.Fatalf("can't find the start server: %v", specialServer)
		}
	}
	return nil
}
func (m *MicroApp) build(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.buildServer(m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server build result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.buildServer(m.ver, server); err != nil {
				return err
			} else {
				fmt.Printf("server build result:%v\n", string(bs))
			}
		} else {
			log.Fatalf("can't find the build server: %v", specialServer)
		}
	}
	return nil
}

func (m *MicroApp) prepare(c *cli.Context) error {
	store, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	version := m.getVersion(store, c)
	if len(version) < 1 {
		return errors.New("please set the version first")
	}
	data := m.getData(store, c)
	if len(data) < 1 {
		return errors.New("can't find data dir")
	}
	m.readConf(data)
	m.store, m.ver, m.data = store, version, data
	return nil
}

func (m *MicroApp) version(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	ver := c.Args().Get(0)
	if len(ver) > 0 {
		m.setIniVar(cfg, dockerVer, ver)
		fmt.Println("版本设置成功")
	} else {
		ver = m.getIniVar(cfg, dockerVer)
		if len(ver) == 0 {
			fmt.Println("未设置版本")
		} else {
			fmt.Printf("当前版本为:%v\n", ver)
		}
	}
	return nil
}

func (m *MicroApp) workDir(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	dir := c.Args().Get(0)
	if len(dir) > 0 {
		m.setIniVar(cfg, dockerData, dir)
		fmt.Println("工作目录设置成功")
	} else {
		dir = m.getIniVar(cfg, dockerData)
		if len(dir) == 0 {
			fmt.Println("未设置工作目录")
		} else {
			fmt.Printf("当前工作目录为:%v\n", dir)
		}
	}
	return nil
}

func (m *MicroApp) before() error {
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
	if exist, _ := utils.PathExists(storePath); !exist {
		fp, err := os.Create(storePath)
		if err != nil {
			return err
		}
		defer fp.Close()
	}
	return nil
}

func (m *MicroApp) buildServer(ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"build", "--build-arg", fmt.Sprintf("run_server=%v", server.ServerId), "--build-arg", fmt.Sprintf(`client_port=%v`, server.ClientPort), "-t", fmt.Sprintf("%v:%v", server.ServerId, ver), "."}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) runServer(data, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"run", "-d", "-v", fmt.Sprintf("%v:/data", data), "-p", fmt.Sprintf("%v:%v", server.ClientPort, server.ClientPort), "--network=xg-net", fmt.Sprintf("--name=%v", server.ServerId), fmt.Sprintf("%v:%v", server.ServerId, ver)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) startServer(ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"start", server.ServerId}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) stopServer(ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"stop", server.ServerId}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) rmServer(ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"rm", server.ServerId}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) rmiServer(ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"rmi", fmt.Sprintf("%v:%v", server.ServerId, ver)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) clearServer(ver string, server *treaty.Server) {
	bs, err := m.stopServer(ver, server)
	fmt.Printf("stop server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
	bs, err = m.rmServer(ver, server)
	fmt.Printf("rm server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
	bs, err = m.rmiServer(ver, server)
	fmt.Printf("rmi server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
}

func (m *MicroApp) getVersion(cfg *ini.File, c *cli.Context) string {
	version := m.getIniVar(cfg, dockerVer)
	if ver := c.String("version"); len(ver) > 0 {
		version = ver
	}
	return version
}

func (m *MicroApp) getData(cfg *ini.File, c *cli.Context) string {
	dir := m.getIniVar(cfg, dockerData)
	if tmp := c.String("data"); len(tmp) > 0 {
		dir = tmp
	}
	return dir
}

func (m *MicroApp) getIniVar(cfg *ini.File, key string) string {
	return cfg.Section("").Key(key).String()
}

func (m *MicroApp) setIniVar(cfg *ini.File, key, val string) {
	cfg.Section("").Key(key).SetValue(val)
	cfg.SaveTo(storePath)
}

func (m *MicroApp) readConf(data string) {
	cfg := path.Join(data, "/conf/config.json")
	viper.SetConfigFile(cfg)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	//frame init
	frameCfg := viper.Get("frame")
	if err := config.InitFrameConf(frameCfg); err != nil {
		logger.Fatal(err)
	}
}
