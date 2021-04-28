- docker调用客户端操作
````
docker exec etcd1 /bin/sh -c "ETCDCTL_API=3 /usr/local/bin/etcdctl get --prefix /server"
````