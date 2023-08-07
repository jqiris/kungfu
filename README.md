# kungfu

分布式可扩展可容器化部署游戏框架

## 基础组件部署
- 创建docker网络-docker network create kf-net
- 运行etcd服务- cd deploy/etcd && docker-compose up -d
- 运行nats服务- cd deploy/nats && docker-compose up -d

## 安装容器部署工具 
- go install github.com/jqiris/kungfu/v2@v2.10.8

## 常用命令
- kungfu build x1 ...x2  # 构建docker服务,会根据不同配置文件构建对应服务，例如account,hall等
- kungfu remote push x1 .. x2 # 将本地构建服务上传到远程仓库
- kungfu remote pull x1 .. x2 # 将远程仓库构建服务拉取到本地仓库
- kungfu run x1 .. x2 # 创建并运行本地构建服务
- kungfu stop x1 .. x2 # 停止本地构建服务
- kungfu start x1 .. x2 # 启动已经创建的本地构建服务
- kungfu restart x1 .. x2 # 重启已经创建的本地构建服务
- kungfu clear x1 .. x2 # 清理已经创建的本地构建服务并清除镜像