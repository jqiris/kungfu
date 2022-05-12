package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"gopkg.in/ini.v1"
)

var (
	Null          = "null"
	storePath     = ".ini"
	dockerVer     = "DockerVer"
	dockerPrefix  = "DockerPrefix"
	dockerData    = "DockerData"
	dockerConfig  = "DockerConfig"
	dockerNetwork = "DockerNetwork"
	dockerPath    = "Dockerfile"
	dockerTml     = `
FROM golang:1.18.1 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server

FROM alpine
ARG run_server
ARG client_port
# 时区控制
ENV TZ Asia/Shanghai
RUN echo "http://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories \
    && apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone
ENV run_mode docker
ENV run_server ${run_server}

COPY --from=builder /src/server /app/

WORKDIR /app

EXPOSE ${client_port}
VOLUME /data
ENTRYPOINT ["/app/server", "-conf", "/data/%s"]	
	`
)

type MicroApp struct {
	ver     string
	data    string
	cfg     string
	network string
	prefix  string
	store   *ini.File
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
				m.clearServer(m.prefix, m.ver, server)
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			m.clearServer(m.prefix, m.ver, server)
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
				if bs, err := m.rmiServer(m.prefix, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("image rm result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.rmiServer(m.prefix, m.ver, server); err != nil {
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
				if bs, err := m.rmServer(m.prefix, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server rm result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.rmServer(m.prefix, m.ver, server); err != nil {
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
				if bs, err := m.stopServer(m.prefix, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server stop result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.stopServer(m.prefix, m.ver, server); err != nil {
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
				if bs, err := m.runServer(m.data, m.network, m.prefix, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server run result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.runServer(m.data, m.network, m.prefix, m.ver, server); err != nil {
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
				if bs, err := m.startServer(m.prefix, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server start result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.startServer(m.prefix, m.ver, server); err != nil {
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

func (m *MicroApp) save(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	saveList, imageList := []string{}, []string{}
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				saveList = append(saveList, server.ServerId)
				imageList = append(imageList, m.runImage(m.prefix, m.ver, server))
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			saveList = append(saveList, server.ServerId)
			imageList = append(imageList, m.runImage(m.prefix, m.ver, server))
		} else {
			log.Fatalf("can't find the start server: %v", specialServer)
		}
	}
	saveName := m.saveName(m.prefix, m.ver, saveList)
	args := append([]string{"save", "-o", saveName}, imageList...)
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	bs, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(bs))
	return nil
}

func (m *MicroApp) load(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	loadTar := c.Args().Get(0)
	if len(loadTar) == 0 || !strings.HasSuffix(loadTar, ".tar") {
		return errors.New("请输入正确加载镜像目录")
	}
	cmd := exec.Command("docker", "load", "-i", loadTar)
	fmt.Println(cmd.String())
	bs, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(bs))
	return nil
}

func (m *MicroApp) build(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	if exist, _ := utils.PathExists(dockerPath); !exist {
		fp, err := os.Create(dockerPath)
		if err != nil {
			return err
		}
		defer fp.Close()
		dockerFile := fmt.Sprintf(dockerTml, m.cfg)
		_, err = fp.Write([]byte(dockerFile))
		if err != nil {
			return err
		}
	}
	servers := config.GetServersConf()
	specialServer := c.Args().Get(0)
	if len(specialServer) < 1 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.buildServer(m.prefix, m.ver, server); err != nil {
					return err
				} else {
					fmt.Printf("server build result:%v\n", string(bs))
				}
			}
		}
	} else {
		if server, ok := servers[specialServer]; ok {
			if bs, err := m.buildServer(m.prefix, m.ver, server); err != nil {
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
	cfg := m.getConfig(store, c)
	if len(cfg) < 1 {
		return errors.New("can't find config path")
	}
	network := m.getNetwork(store, c)
	if len(network) < 1 {
		return errors.New("not set network")
	}
	data := m.getData(store, c)
	if len(data) < 1 {
		return errors.New("can't find data dir")
	}
	m.readConf(data, cfg)
	prefix := m.getPrefix(store, c)
	m.store, m.ver, m.data, m.cfg, m.network, m.prefix = store, version, data, cfg, network, prefix
	fmt.Printf("ver:%v,data:%v,cfg:%v,network:%v,prefix:%v \n", version, data, cfg, network, prefix)
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
func (m *MicroApp) runPrefix(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	prefix := c.Args().Get(0)
	if len(prefix) > 0 {
		m.setIniVar(cfg, dockerPrefix, prefix)
		fmt.Println("运行前缀设置成功")
	} else {
		prefix = m.getIniVar(cfg, dockerPrefix)
		if len(prefix) == 0 {
			fmt.Println("未设置运行前缀")
		} else {
			fmt.Printf("当前运行前缀为:%v\n", prefix)
		}
	}
	return nil
}

func (m *MicroApp) config(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	configure := c.Args().Get(0)
	if len(configure) > 0 {
		m.setIniVar(cfg, dockerConfig, configure)
		fmt.Println("工作目录下配置文件位置设置成功")
	} else {
		configure = m.getIniVar(cfg, dockerConfig)
		if len(configure) == 0 {
			fmt.Println("工作目录下未设置配置文件位置")
		} else {
			fmt.Printf("当前工作目录下配置文件位置为:%v\n", configure)
		}
	}
	return nil
}

func (m *MicroApp) netView(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	network := c.Args().Get(0)
	if len(network) > 0 {
		m.setIniVar(cfg, dockerNetwork, network)
		fmt.Println("运行网络设置成功")
	} else {
		network = m.getIniVar(cfg, dockerNetwork)
		if len(network) == 0 {
			fmt.Println("未设置运行网络")
		} else {
			fmt.Printf("当前运行网络为:%v\n", network)
		}
	}
	return nil
}

func (m *MicroApp) before() error {
	if exist, _ := utils.PathExists(storePath); !exist {
		fp, err := os.Create(storePath)
		if err != nil {
			return err
		}
		defer fp.Close()
	}
	return nil
}

func (m *MicroApp) runName(prefix, ver string, server *treaty.Server) string {
	if len(prefix) == 0 {
		return server.ServerId
	}
	return prefix + "_" + server.ServerId
}

func (m *MicroApp) runImage(prefix, ver string, server *treaty.Server) string {
	if len(prefix) == 0 {
		return fmt.Sprintf("%v:%v", server.ServerId, ver)
	}
	return fmt.Sprintf("%v_%v:%v", prefix, server.ServerId, ver)
}
func (m *MicroApp) saveName(prefix, ver string, list []string) string {
	item := strings.Join(list, "-")
	if len(prefix) == 0 {
		return fmt.Sprintf("%v_%v.tar", item, ver)
	}
	return fmt.Sprintf("%v_%v_%v.tar", prefix, item, ver)
}

func (m *MicroApp) buildServer(prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"build", "--build-arg", fmt.Sprintf("run_server=%v", server.ServerId), "--build-arg", fmt.Sprintf(`client_port=%v`, server.ClientPort), "-t", m.runImage(prefix, ver, server), "."}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) runServer(data, network, prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"run", "-d", "-v", fmt.Sprintf("%v:/data", data), "-p", fmt.Sprintf("%v:%v", server.ClientPort, server.ClientPort), fmt.Sprintf("--network=%s", network), fmt.Sprintf("--name=%v", m.runName(prefix, ver, server)), m.runImage(prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) startServer(prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"start", m.runName(prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) stopServer(prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"stop", m.runName(prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) rmServer(prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"rm", m.runName(prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) rmiServer(prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"rmi", m.runImage(prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) clearServer(prefix, ver string, server *treaty.Server) {
	bs, err := m.stopServer(prefix, ver, server)
	fmt.Printf("stop server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
	bs, err = m.rmServer(prefix, ver, server)
	fmt.Printf("rm server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
	bs, err = m.rmiServer(prefix, ver, server)
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

func (m *MicroApp) getConfig(cfg *ini.File, c *cli.Context) string {
	configure := m.getIniVar(cfg, dockerConfig)
	if tmp := c.String("config"); len(tmp) > 0 {
		configure = tmp
	}
	return configure
}

func (m *MicroApp) getNetwork(cfg *ini.File, c *cli.Context) string {
	network := m.getIniVar(cfg, dockerNetwork)
	if tmp := c.String("network"); len(tmp) > 0 {
		network = tmp
	}
	return network
}

func (m *MicroApp) getPrefix(cfg *ini.File, c *cli.Context) string {
	prefix := m.getIniVar(cfg, dockerPrefix)
	if tmp := c.String("prefix"); len(tmp) > 0 {
		prefix = tmp
		if prefix == Null {
			prefix = ""
		}
	}
	return prefix
}

func (m *MicroApp) getIniVar(cfg *ini.File, key string) string {
	return cfg.Section("").Key(key).String()
}

func (m *MicroApp) setIniVar(cfg *ini.File, key, val string) {
	if val == Null {
		val = ""
	}
	cfg.Section("").Key(key).SetValue(val)
	cfg.SaveTo(storePath)
}

func (m *MicroApp) readConf(data, cfgPath string) {
	cfg := path.Join(data, cfgPath)
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
