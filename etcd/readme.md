ETCD单例 
```
docker run -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --name etcd01 \
  quay.io/coreos/etcd:v3.3.8 \
  /usr/local/bin/etcd \
  --name s1 \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379
```

docker-compose启动
```
# 因为启动的目录名都是deploy docker-compose 认为冲突 所以需要特定
docker-compose -p go_micro_etcd up -d
```

设置docker中环境变量
```
export ETCDCTL_API=3
```
切换至etcdctl
```
cd /bin
```
