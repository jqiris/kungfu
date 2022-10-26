package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
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
	Null                = "null"
	storePath           = ".ini"
	dockerVer           = "DockerVer"
	dockerPrefix        = "DockerPrefix"
	dockerData          = "DockerData"
	dockerConfig        = "DockerConfig"
	dockerRemoteConfig  = "DockerRemoteConfig"
	dockerNetwork       = "DockerNetwork"
	dockerPath          = "Dockerfile"
	dockerRegistry      = "DockerRegistry"
	dockerProject       = "DockerProject"
	defaultProject      = "default"
	labelAuthor         = "LabelAuthor"
	defaultLabelAuthor  = "jqiris"
	labelVersion        = "LabelVersion"
	defaultLabelVersion = "1.0.0"
	goVersion           = "goVersion"
	defaultGoVersion    = "1.19.2"
	dockerTml           = `
FROM golang:%s AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o buildsrv

FROM alpine
MAINTAINER %s
LABEL VERSION="%s"
ARG run_server
ARG client_port
# 时区控制
ENV TZ Asia/Shanghai
RUN echo "http://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories \
	&& apk update \
	&& apk add curl \
    && apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone
ENV run_mode docker
ENV run_server ${run_server}

COPY --from=builder /src/buildsrv /app/

WORKDIR /app

EXPOSE ${client_port}
VOLUME /data
ENTRYPOINT ["/app/buildsrv", "-conf", "/data/%s"]	
	`
)

type MicroApp struct {
	ver          string
	data         string
	cfg          string
	remoteCfg    string
	network      string
	prefix       string
	depot        string
	project      string
	labelAuthor  string
	labelVersion string
	goVersion    string
	store        *ini.File
}

func newMicroApp() *MicroApp {
	return &MicroApp{}
}

func (m *MicroApp) clear(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	var launchArr []*treaty.Server
	if c.NArg() == 0 {
		for _, cfg := range servers {
			if cfg.IsLaunch {
				launchArr = append(launchArr, cfg)
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				launchArr = append(launchArr, server)
			} else {
				log.Fatalf("can't find the clear server: %v", specialServer)
			}
		}
	}
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].ShutWeight < launchArr[j].ShutWeight
	})
	for _, server := range launchArr {
		m.clearServer(server)
	}
	return nil
}

func (m *MicroApp) rmi(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	if c.NArg() == 0 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.rmiServer(server); err != nil {
					return err
				} else {
					fmt.Printf("image rm result:%v\n", string(bs))
				}
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				if bs, err := m.rmiServer(server); err != nil {
					return err
				} else {
					fmt.Printf("image rm result:%v\n", string(bs))
				}
			} else {
				log.Fatalf("can't find the rm image: %v\n", specialServer)
			}
		}
	}
	return nil
}
func (m *MicroApp) prune(c *cli.Context) error {
	//ocker rmi $(docker images -q -f dangling=true)
	args := []string{"images", "-q", "-f", "dangling=true"}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	pidStr := string(out)
	if len(pidStr) > 0 {
		pids := strings.Split(pidStr, "\n")
		for _, pid := range pids {
			if len(pid) == 0 {
				continue
			}
			cmd := exec.Command("docker", "rmi", pid)
			fmt.Println(cmd.String())
			out, err := cmd.Output()
			if err != nil {
				return err
			}
			fmt.Println("prune:", string(out))
		}
	}
	return nil
}

func (m *MicroApp) rm(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	var launchArr []*treaty.Server
	if c.NArg() == 0 {
		for _, cfg := range servers {
			if cfg.IsLaunch {
				launchArr = append(launchArr, cfg)
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				launchArr = append(launchArr, server)
			} else {
				log.Fatalf("can't find the rm server: %v", specialServer)
			}
		}
	}
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].ShutWeight < launchArr[j].ShutWeight
	})
	for _, server := range launchArr {
		if bs, err := m.rmServer(server); err != nil {
			return err
		} else {
			fmt.Printf("server rm result:%v\n", string(bs))
		}
	}
	return nil
}
func (m *MicroApp) stop(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	var launchArr []*treaty.Server
	if c.NArg() == 0 {
		for _, cfg := range servers {
			if cfg.IsLaunch {
				launchArr = append(launchArr, cfg)
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				launchArr = append(launchArr, server)
			} else {
				log.Fatalf("can't find the stop server: %v", specialServer)
			}
		}
	}
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].ShutWeight < launchArr[j].ShutWeight
	})
	for _, server := range launchArr {
		if bs, err := m.stopServer(server); err != nil {
			return err
		} else {
			fmt.Printf("server stop result:%v\n", string(bs))
		}
	}
	return nil
}

func (m *MicroApp) run(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	mem, memSwap, memKernel := c.String("memory"), c.String("memory-swap"), c.String("kernel-memory")
	cpu, cpuSet := c.String("cpus"), c.String("cpuset-cpus")
	var launchArr []*treaty.Server
	if c.NArg() == 0 {
		for _, cfg := range servers {
			if cfg.IsLaunch {
				launchArr = append(launchArr, cfg)
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				launchArr = append(launchArr, server)
			} else {
				log.Fatalf("can't find the run server: %v", specialServer)
			}
		}
	}
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].LaunchWeight < launchArr[j].LaunchWeight
	})
	for _, server := range launchArr {
		if bs, err := m.runServer(mem, memSwap, memKernel, cpu, cpuSet, server); err != nil {
			return err
		} else {
			fmt.Printf("server run result:%v\n", string(bs))
		}
	}
	return nil
}

func (m *MicroApp) start(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	var launchArr []*treaty.Server
	if c.NArg() == 0 {
		for _, cfg := range servers {
			if cfg.IsLaunch {
				launchArr = append(launchArr, cfg)
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				launchArr = append(launchArr, server)
			} else {
				log.Fatalf("can't find the start server: %v", specialServer)
			}
		}
	}
	sort.Slice(launchArr, func(i, j int) bool {
		return launchArr[i].LaunchWeight < launchArr[j].LaunchWeight
	})
	for _, server := range launchArr {
		if bs, err := m.startServer(server); err != nil {
			return err
		} else {
			fmt.Printf("server start result:%v\n", string(bs))
		}
	}
	return nil
}

func (m *MicroApp) registryBefore(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	registry := m.getRegistry(cfg, c)
	if len(registry) == 0 {
		return errors.New("未设置仓库地址")
	}
	m.depot = registry
	return nil
}

func (m *MicroApp) registryPush(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	if c.NArg() == 0 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.pushServer(server); err != nil {
					return err
				} else {
					fmt.Printf("registry push result:%v\n", string(bs))
				}
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				if bs, err := m.pushServer(server); err != nil {
					return err
				} else {
					fmt.Printf("registry push result:%v\n", string(bs))
				}
			} else {
				log.Fatalf("can't find the registry push server: %v", specialServer)
			}
		}
	}
	return nil
}

func (m *MicroApp) registryPull(c *cli.Context) error {
	if err := m.prepare(c); err != nil {
		return err
	}
	servers := config.GetServersConf()
	if c.NArg() == 0 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.pullServer(server); err != nil {
					return err
				} else {
					fmt.Printf("registry pull result:%v\n", string(bs))
				}
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				if bs, err := m.pullServer(server); err != nil {
					return err
				} else {
					fmt.Printf("registry pull result:%v\n", string(bs))
				}
			} else {
				log.Fatalf("can't find the registry pull server: %v", specialServer)
			}
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

func (m *MicroApp) removeDockerFile() {
	err := os.Remove(dockerPath)
	if err != nil {
		fmt.Println(err)
	}
}

func (m *MicroApp) serverDockerName(server *treaty.Server) string {
	return fmt.Sprintf("%v_%v", server.ServerId, "dockerFile")
}

func (m *MicroApp) createServerDocker(server *treaty.Server) error {
	file, err := os.Open(dockerPath)
	if err != nil {
		return err
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	serverName := "server_" + server.ServerId
	if len(m.prefix) > 0 {
		serverName = fmt.Sprintf("%v_%v", m.prefix, server.ServerId)
	}
	ntext := strings.Replace(string(bs), "buildsrv", serverName, -1)
	fp, err := os.Create(m.serverDockerName(server))
	if err != nil {
		return err
	}
	defer fp.Close()
	_, err = fp.Write([]byte(ntext))
	if err != nil {
		return err
	}
	return nil
}

func (m *MicroApp) removeServerDocker(server *treaty.Server) {
	err := os.Remove(m.serverDockerName(server))
	if err != nil {
		fmt.Println(err)
	}
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
		dockerFile := fmt.Sprintf(dockerTml, m.goVersion, m.labelAuthor, m.labelVersion, m.remoteCfg)
		_, err = fp.Write([]byte(dockerFile))
		if err != nil {
			return err
		}
	}
	buildPath := "."
	if v := c.String("buildPath"); len(v) > 0 {
		buildPath = v
	}
	servers := config.GetServersConf()
	if c.NArg() == 0 {
		for _, server := range servers {
			if server.IsLaunch {
				if bs, err := m.buildServer(buildPath, server); err != nil {
					return err
				} else {
					fmt.Printf("server build result:%v\n", string(bs))
				}
			}
		}
	} else {
		for _, specialServer := range c.Args().Slice() {
			if server, ok := servers[specialServer]; ok {
				if bs, err := m.buildServer(buildPath, server); err != nil {
					return err
				} else {
					fmt.Printf("server build result:%v\n", string(bs))
				}
			} else {
				log.Fatalf("can't find the build server: %v", specialServer)
			}
		}
	}
	return nil
}

func (m *MicroApp) prepare(c *cli.Context) error {
	store, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	project := m.getProject(store, c)
	if len(project) < 1 {
		return errors.New("please set the project first")
	}
	version := m.getVersion(store, c)
	if len(version) < 1 {
		return errors.New("please set the version first")
	}
	cfg := m.getConfig(store, c)
	if len(cfg) < 1 {
		return errors.New("can't find config path")
	}
	remoteCfg := m.getRemoteConfig(store, c)
	if len(remoteCfg) < 1 {
		return errors.New("can't find remote config path")
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
	la := m.getLabelAuthor(store, c)
	lv := m.getLabelVersion(store, c)
	gv := m.getGoVersion(store, c)
	m.store, m.ver, m.data, m.cfg, m.network, m.prefix, m.remoteCfg, m.project, m.labelAuthor, m.labelVersion, m.goVersion = store, version, data, cfg, network, prefix, remoteCfg, project, la, lv, gv
	fmt.Printf("project:%v,ver:%v,data:%v,cfg:%v,remoteCfg:%v,network:%v,labelAuthor:%v,labelVersion:%v,goVersion:%v,prefix:%v \n", project, version, data, cfg, remoteCfg, network, la, lv, gv, prefix)
	return nil
}

func (m *MicroApp) version(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	ver := c.Args().Get(0)
	if len(ver) > 0 {
		m.setProjectVar(cfg, c, dockerVer, ver)
		fmt.Println("版本设置成功")
	} else {
		ver = m.getProjectVar(cfg, c, dockerVer)
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
		m.setProjectVar(cfg, c, dockerData, dir)
		fmt.Println("工作目录设置成功")
	} else {
		dir = m.getProjectVar(cfg, c, dockerData)
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
		m.setProjectVar(cfg, c, dockerPrefix, prefix)
		fmt.Println("运行前缀设置成功")
	} else {
		prefix = m.getProjectVar(cfg, c, dockerPrefix)
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
		m.setProjectVar(cfg, c, dockerConfig, configure)
		m.removeDockerFile()
		fmt.Println("工作目录下配置文件位置设置成功")
	} else {
		configure = m.getProjectVar(cfg, c, dockerConfig)
		if len(configure) == 0 {
			fmt.Println("工作目录下未设置配置文件位置")
		} else {
			fmt.Printf("当前工作目录下配置文件位置为:%v\n", configure)
		}
	}
	return nil
}

func (m *MicroApp) remoteConfig(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	configure := c.Args().Get(0)
	if len(configure) > 0 {
		m.setProjectVar(cfg, c, dockerRemoteConfig, configure)
		m.removeDockerFile()
		fmt.Println("远程配置文件位置设置成功")
	} else {
		configure = m.getProjectVar(cfg, c, dockerRemoteConfig)
		if len(configure) == 0 {
			fmt.Println("未设置远程配置文件位置")
		} else {
			fmt.Printf("当前远程配置文件位置为:%v\n", configure)
		}
	}
	return nil
}

func (m *MicroApp) projectSet(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	project := c.Args().Get(0)
	if len(project) > 0 {
		m.setGlobalVar(cfg, dockerProject, project)
		m.removeDockerFile()
		fmt.Println("项目设置成功")
	} else {
		project = m.getGlobalVar(cfg, dockerProject)
		if len(project) == 0 {
			project = defaultProject
		}
		fmt.Printf("当前项目为:%v\n", project)
	}
	return nil
}
func (m *MicroApp) labelAuthorSet(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	author := c.Args().Get(0)
	if len(author) > 0 {
		m.setGlobalVar(cfg, labelAuthor, author)
		m.removeDockerFile()
		fmt.Println("项目维护者设置成功")
	} else {
		author = m.getGlobalVar(cfg, labelAuthor)
		if len(author) == 0 {
			author = defaultLabelAuthor
		}
		fmt.Printf("当前项目维护者为:%v\n", author)
	}
	return nil
}
func (m *MicroApp) labelVersionSet(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	version := c.Args().Get(0)
	if len(version) > 0 {
		m.setGlobalVar(cfg, labelVersion, version)
		m.removeDockerFile()
		fmt.Println("docker版本设置成功")
	} else {
		version = m.getGlobalVar(cfg, labelVersion)
		if len(version) == 0 {
			version = defaultLabelVersion
		}
		fmt.Printf("当前docker版本为:%v\n", version)
	}
	return nil
}

func (m *MicroApp) goVersionSet(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	version := c.Args().Get(0)
	if len(version) > 0 {
		m.setGlobalVar(cfg, goVersion, version)
		fmt.Println("go版本设置成功")
	} else {
		version = m.getGlobalVar(cfg, goVersion)
		if len(version) == 0 {
			version = defaultGoVersion
		}
		fmt.Printf("当前go版本为:%v\n", version)
	}
	return nil
}

func (m *MicroApp) registry(c *cli.Context) error {
	cfg, err := ini.Load(storePath)
	if err != nil {
		return err
	}
	configure := c.Args().Get(0)
	if len(configure) > 0 {
		m.setGlobalVar(cfg, dockerRegistry, configure)
		fmt.Println("远程仓库地址设置成功")
	} else {
		configure = m.getGlobalVar(cfg, dockerRegistry)
		if len(configure) == 0 {
			fmt.Println("未设置远程仓库地址")
		} else {
			fmt.Printf("当前远程仓库地址为:%v\n", configure)
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
		m.setProjectVar(cfg, c, dockerNetwork, network)
		fmt.Println("运行网络设置成功")
	} else {
		network = m.getProjectVar(cfg, c, dockerNetwork)
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

func (m *MicroApp) runRemoteImage(depot, prefix, ver string, server *treaty.Server) string {
	return fmt.Sprintf("%v/%v", depot, m.runImage(prefix, ver, server))
}

func (m *MicroApp) saveName(prefix, ver string, list []string) string {
	item := strings.Join(list, "-")
	if len(prefix) == 0 {
		return fmt.Sprintf("%v_%v.tar", item, ver)
	}
	return fmt.Sprintf("%v_%v_%v.tar", prefix, item, ver)
}

func (m *MicroApp) buildServer(buildPath string, server *treaty.Server) ([]byte, error) {
	err := m.createServerDocker(server)
	if err != nil {
		return nil, err
	}
	defer m.removeServerDocker(server)
	fileName := m.serverDockerName(server)
	args := []string{"build", "-f", fileName, "--build-arg", fmt.Sprintf("run_server=%v", server.ServerId), "--build-arg", fmt.Sprintf(`client_port=%v`, server.ClientPort), "-t", m.runImage(m.prefix, m.ver, server), buildPath}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) runServer(mem, memSwap, memKernel, cpu, cpuSet string, server *treaty.Server) ([]byte, error) {
	args := []string{"run", "-d", "-v", fmt.Sprintf("%v:/data", m.data), "-p", fmt.Sprintf("%v:%v", server.ClientPort, server.ClientPort), fmt.Sprintf("--network=%s", m.network)}
	if len(mem) > 0 {
		args = append(args, fmt.Sprintf("--memory=%v", mem))
	}
	if len(memSwap) > 0 {
		args = append(args, fmt.Sprintf("--memory-swap=%v", memSwap))
	}
	if len(memKernel) > 0 {
		args = append(args, fmt.Sprintf("--kernel-memory=%v", memKernel))
	}
	if len(cpu) > 0 {
		args = append(args, fmt.Sprintf("--cpus=%v", cpu))
	}
	if len(cpuSet) > 0 {
		args = append(args, fmt.Sprintf("--cpuset-cpus=%v", cpuSet))
	}
	args = append(args, fmt.Sprintf("--name=%v", m.runName(m.prefix, m.ver, server)), m.runImage(m.prefix, m.ver, server))
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) startServer(server *treaty.Server) ([]byte, error) {
	args := []string{"start", m.runName(m.prefix, m.ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) stopServer(server *treaty.Server) ([]byte, error) {
	args := []string{"stop", m.runName(m.prefix, m.ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) rmServer(server *treaty.Server) ([]byte, error) {
	args := []string{"rm", m.runName(m.prefix, m.ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) rmiServer(server *treaty.Server) ([]byte, error) {
	args := []string{"rmi", m.runImage(m.prefix, m.ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) pushServer(server *treaty.Server) ([]byte, error) {
	bs, err := m.imageTag(m.depot, m.prefix, m.ver, server)
	if err != nil {
		return bs, err
	}
	bs, err = m.imagePush(m.depot, m.prefix, m.ver, server)
	if err != nil {
		return bs, err
	}
	return m.imageRemoteClear(m.depot, m.prefix, m.ver, server)
}

func (m *MicroApp) pullServer(server *treaty.Server) ([]byte, error) {
	bs, err := m.imagePull(m.depot, m.prefix, m.ver, server)
	if err != nil {
		return bs, err
	}
	bs, err = m.imageUnTag(m.depot, m.prefix, m.ver, server)
	if err != nil {
		return bs, err
	}
	return m.imageRemoteClear(m.depot, m.prefix, m.ver, server)
}

func (m *MicroApp) imageTag(depot, prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"tag", m.runImage(prefix, ver, server), m.runRemoteImage(depot, prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) imagePush(depot, prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"push", m.runRemoteImage(depot, prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) imagePull(depot, prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"pull", m.runRemoteImage(depot, prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) imageUnTag(depot, prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"tag", m.runRemoteImage(depot, prefix, ver, server), m.runImage(prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) imageRemoteClear(depot, prefix, ver string, server *treaty.Server) ([]byte, error) {
	args := []string{"rmi", m.runRemoteImage(depot, prefix, ver, server)}
	cmd := exec.Command("docker", args...)
	fmt.Println(cmd.String())
	return cmd.Output()
}

func (m *MicroApp) clearServer(server *treaty.Server) {
	bs, err := m.stopServer(server)
	fmt.Printf("stop server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
	bs, err = m.rmServer(server)
	fmt.Printf("rm server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
	bs, err = m.rmiServer(server)
	fmt.Printf("rmi server:%v, result,res:%v,err:%v \n", server.ServerId, string(bs), err)
}

func (m *MicroApp) getVersion(cfg *ini.File, c *cli.Context) string {
	version := m.getProjectVar(cfg, c, dockerVer)
	if ver := c.String("version"); len(ver) > 0 {
		version = ver
	}
	return version
}

func (m *MicroApp) getData(cfg *ini.File, c *cli.Context) string {
	dir := m.getProjectVar(cfg, c, dockerData)
	if tmp := c.String("data"); len(tmp) > 0 {
		dir = tmp
	}
	return dir
}

func (m *MicroApp) getConfig(cfg *ini.File, c *cli.Context) string {
	configure := m.getProjectVar(cfg, c, dockerConfig)
	if tmp := c.String("config"); len(tmp) > 0 {
		configure = tmp
	}
	return configure
}

func (m *MicroApp) getRemoteConfig(cfg *ini.File, c *cli.Context) string {
	configure := m.getProjectVar(cfg, c, dockerRemoteConfig)
	if tmp := c.String("remoteConfig"); len(tmp) > 0 {
		configure = tmp
	}
	if len(configure) == 0 {
		configure = m.getConfig(cfg, c)
	}
	return configure
}

func (m *MicroApp) getProject(cfg *ini.File, c *cli.Context) string {
	project := m.getGlobalVar(cfg, dockerProject)
	if tmp := c.String("project"); len(tmp) > 0 {
		project = tmp
	}
	if len(project) == 0 {
		project = defaultProject
	}
	return project
}

func (m *MicroApp) getLabelAuthor(cfg *ini.File, c *cli.Context) string {
	author := m.getGlobalVar(cfg, labelAuthor)
	if tmp := c.String("labelAuthor"); len(tmp) > 0 {
		author = tmp
	}
	if len(author) == 0 {
		author = defaultLabelAuthor
	}
	return author
}
func (m *MicroApp) getLabelVersion(cfg *ini.File, c *cli.Context) string {
	version := m.getGlobalVar(cfg, labelVersion)
	if tmp := c.String("labelVersion"); len(tmp) > 0 {
		version = tmp
	}
	if len(version) == 0 {
		version = defaultLabelVersion
	}
	return version
}

func (m *MicroApp) getGoVersion(cfg *ini.File, c *cli.Context) string {
	version := m.getGlobalVar(cfg, goVersion)
	if tmp := c.String("goVersion"); len(tmp) > 0 {
		version = tmp
	}
	if len(version) == 0 {
		version = defaultGoVersion
	}
	return version
}

func (m *MicroApp) getNetwork(cfg *ini.File, c *cli.Context) string {
	network := m.getProjectVar(cfg, c, dockerNetwork)
	if tmp := c.String("network"); len(tmp) > 0 {
		network = tmp
	}
	return network
}

func (m *MicroApp) getRegistry(cfg *ini.File, c *cli.Context) string {
	registry := m.getGlobalVar(cfg, dockerRegistry)
	if tmp := c.String("registry"); len(tmp) > 0 {
		registry = tmp
	}
	return registry
}

func (m *MicroApp) getPrefix(cfg *ini.File, c *cli.Context) string {
	prefix := m.getProjectVar(cfg, c, dockerPrefix)
	if tmp := c.String("prefix"); len(tmp) > 0 {
		prefix = tmp
		if prefix == Null {
			prefix = ""
		}
	}
	return prefix
}

func (m *MicroApp) getProjectVar(cfg *ini.File, c *cli.Context, key string) string {
	project := m.getProject(cfg, c)
	return cfg.Section(project).Key(key).String()
}

func (m *MicroApp) setProjectVar(cfg *ini.File, c *cli.Context, key, val string) {
	if val == Null {
		val = ""
	}
	project := m.getProject(cfg, c)
	cfg.Section(project).Key(key).SetValue(val)
	cfg.SaveTo(storePath)
}

func (m *MicroApp) getGlobalVar(cfg *ini.File, key string) string {
	return cfg.Section("").Key(key).String()
}

func (m *MicroApp) setGlobalVar(cfg *ini.File, key, val string) {
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
