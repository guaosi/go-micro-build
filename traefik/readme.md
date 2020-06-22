创建docker网络
```
docker network create gomicro --driver bridge
```

docker-compose启动
```
# 因为启动的目录名都是deploy docker-compose 认为冲突 所以需要特定
docker-compose -p go_micro_traefik up -d
```