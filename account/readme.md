# 运行

默认127地址
```
docker run -e PARAMS="--registry etcd" --name="account" -d registry.cn-shenzhen.aliyuncs.com/go_micro/account:v1.0
```

指定地址
```
docker run -e PARAMS="--registry etcd --registry_address 172.17.0.5:2379" --name="account" -d registry.cn-shenzhen.aliyuncs.com/go_micro/account:v1.0
```

创建docker网络
```
docker network create gomicro --driver bridge
```

docker-compose启动
```
# 因为启动的目录名都是deploy docker-compose 认为冲突 所以需要特定
docker-compose -p go_micro_account up -d
```