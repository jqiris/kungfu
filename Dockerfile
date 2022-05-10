
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
	