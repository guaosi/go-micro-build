# 手摸手教你从0到1搭建部署Go微服务

本仓库配合 [手摸手教你从开发到部署(CI/CD)GO微服务系列](https://www.guaosi.com/2020/07/05/go-microservice-series-from-development-to-deployment-introduction-contents/) 食用最佳噢~

从`原生搭建`、`容器搭建`、`Docker-Compose搭建`、`Kubernetes搭建`这四个过程，从0到1体验基于GO的微服务搭建部署的全过程。

## 涉及

|     技术      |           使用            |      版本       |
| :-----------: | :-----------------------: | :-------------: |
|     语言      |          Golang           |     1.14.1      |
| Web框架(网关) |            Gin            |     v1.6.3      |
|   通讯格式    |         Protobuf          |     v3.12.1      |
|  微服务框架   |         Go-micro          |     v2.9.0      |
|   反向代理    |          Traefik          |     v2.2.1      |
| 服务注册中心  |      Etcd/Kubernetes      | v3.3.8/v1.16.5  |
|     容器      |          Docker           |    v19.03.8     |
|   编排工具    | Docker-Compose/Kubernetes | v1.25.5/v1.16.5 |


## 目录讲解

```
- account 基于go-micro的微服务

- apigateway 基于go-micro的网关,负责调用微服务

- deploy 初始化部署调用,k8s搭建时使用到

- etcd 服务注册使用,在原生搭建跟docker-compose搭建时使用到

- traefik 反向代理
```

## 原生搭建

### 确保go环境
由于`原生搭建`是直接跑`go`的代码,所以请先确保有`go`的运行环境。同时,请设置好`go mod`代理.

> 不会设置请看 https://goproxy.io/zh/

### etcd

使用`etcd`来作为服务注册组件,`go-micro`通过`etcd`来做服务注册的存储。

```
# etcd单一节点

docker run -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --name etcd1 \
  quay.io/coreos/etcd:v3.3.8 \
  /usr/local/bin/etcd \
  --name s1 \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379
```

如果想在etcd容器中使用cli
```
# 进入容器
docker exec -it etcd1 sh

# 设置docker中环境变量
export ETCDCTL_API=3

# 切换至etcdctl目录
cd /bin
```

### account

切换到 `account` 目录下
```
cd account
```

下载依赖

```
go mod download
```

执行

```
go run . --registry etcd
```

> 如果想`account`高可用,多几个进程执行`go run .`即可

### apigateway

切换到 `apigateway` 目录下
```
cd apigateway
```

下载依赖

```
go mod download
```

执行

```
go run . -p 8091 --registry etcd
```

### 验证网关以及微服务

```
curl -X POST -d "username=guaosi&password=guaosi" http://127.0.0.1:8091/account/register
```
如果返回 `{"code":0,"message":""}` 则证明成功。

## 容器搭建

### etcd
在`原生搭建`中使用的就是`etcd`的镜像创建的容器,这里可以跳过

获取etcd在网络里的内网IP地址

```
docker inspect etcd1 --format "{{.NetworkSettings.IPAddress}}"
```

### account

```
# 此时寻找etcd的地址为容器内的127.0.0.1
docker run -e PARAMS="--registry etcd" --name="account" -d registry.cn-shenzhen.aliyuncs.com/go_micro/account:v1.0

# 此时寻找etcd的地址为容器内的172.17.0.5(实际多少请根据上面etcd查询IP方法获取到)
docker run -e PARAMS="--registry etcd --registry_address 172.17.0.5:2379" --name="account" -d registry.cn-shenzhen.aliyuncs.com/go_micro/account:v1.0
```

### apigateway

```
# 此时寻找etcd的地址为容器内的127.0.0.1
docker run -e PARAMS="--registry etcd" --name="account" -d registry.cn-shenzhen.aliyuncs.com/go_micro/account:v1.0

# 此时寻找etcd的地址为容器内的172.17.0.5(实际多少请根据上面etcd查询IP方法获取到)
docker run -e PARAMS="-p 8091 --registry etcd --registry_address 172.17.0.5:2379" -p 8091:8091 --name="apigw" -d registry.cn-shenzhen.aliyuncs.com/go_micro/apigw:v1.0
```

### 验证网关以及微服务

```
curl -X POST -d "username=guaosi&password=guaosi" http://127.0.0.1:8091/account/register
```
如果返回 `{"code":0,"message":""}` 则证明成功。

## Docker Compose下搭建

### 创建专属网络

为了隔离和管理,所以先创建一个专属网络
```
docker network create gomicro --driver bridge
```

### etcd

切换到 `etcd` 目录下.
```
cd etcd
```

因为docker-compose启动时会根据当前的所在目录名取名,这样目录名称下执行`docker-compose`会被认为冲突 所以需要特定名称
```
docker-compose -p go_micro_etcd up -d
```

### account

切换到 `account/deploy` 目录下.
```
cd account/deploy
```

因为docker-compose启动时会根据当前的所在目录名取名,这样目录名称下执行`docker-compose`会被认为冲突 所以需要特定名称
```
docker-compose -p go_micro_account up -d
```

### apigateway

切换到 `apigateway/deploy` 目录下.
```
cd apigateway/deploy
```

因为docker-compose启动时会根据当前的所在目录名取名,这样目录名称下执行`docker-compose`会被认为冲突 所以需要特定名称
```
docker-compose -p go_micro_apigw up -d
```

### traefik

切换到 `traefik` 目录下.
```
cd traefik
```

因为docker-compose启动时会根据当前的所在目录名取名,这样目录名称下执行`docker-compose`会被认为冲突 所以需要特定名称
```
docker-compose -p go_micro_traefik up -d
```

### 配置 Host 文件

客户端想通过域名访问服务,必须要进行`DNS` 解析,由于这里没有`DNS`服务器进行域名解析,所以修改`hosts`文件将`apigw`所在节点服务器的`IP`和自定义`Host`绑定,。打开电脑的`Hosts`配置文件,往其加入下面配置：

```
127.0.0.1  apigw.guaosi.com
```

### 验证Traefik Dashboard

打开浏览器输入地址：http://127.0.0.1:8080,即可打开`Traefik Dashboard`

### 验证网关以及微服务

```
curl -X POST -d "username=guaosi&password=guaosi" http://apigw.guaosi.com/account/register
```
如果返回 `{"code":0,"message":""}` 则证明成功。

## Kubernetes下搭建

当前kubernetes集群是kubernetes单节点,即master上进行任务调度和pod创建。

```
# 通过污点 设置master允许pod创建
kubectl taint node k8s-master node-role.kubernetes.io/master-

# 通过污点 设置master禁止pod创建
kubectl taint node k8s-master node-role.kubernetes.io/master=""
```

当如果想下载自己的私有镜像时,记得需要现在kubernetes里创建secret认证,否则没有权限进行下载。

**我的镜像已经开放了公有权限,可以直接下载,不需要进行验证。如果你想使用自己的镜像,请按照下面进行操作。**

下面以阿里云为例

### 下载私有阿里云镜像

#### 登陆认证阿里云

```
docker login --username=<your-name> registry.cn-shenzhen.aliyuncs.com
# <your-name> 是你在阿里云上的登陆名
```

#### 在集群中创建保存授权令牌的 Secret

```
# 创建 Secret,命名为 regcred：

kubectl create secret docker-registry regcred --docker-server=<your-registry-server> --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email> -n <your-namespace>

# <your-registry-server> 是你的私有 Docker 仓库全限定域名（FQDN）
# <your-name> 是你的 Docker 用户名。
# <your-pword> 是你的 Docker 密码。
# <your-email> 是你的 Docker 邮箱。
# <your-namespace> 创建后的该密钥属于k8s集群中哪一个namespace的
```
举个栗子
```
kubectl create secret docker-registry regcred --docker-server=registry.cn-shenzhen.aliyuncs.com --docker-username=guaosi@vip.qq.com --docker-password=a123654 --docker-email=guaosi@vip.qq.com -n go-micro

# 注意: -n 后面是我想这个密钥归属于哪个namespace,即哪个namespace可以使用
```
#### 在pod清单中加入使用密钥

```
apiVersion: v1
kind: Pod
metadata:
  name: private-reg
  namespace: <your-namespace>
spec:
  containers:
  - name: private-reg-container
    image: <your-private-image>
  imagePullSecrets: # 使用指定的密钥
  - name: regcred # 与上面创建的secret名称相同
```

### 初始化

创建名字为go-micro的namespace
```
kubectl create -f deploy/k8s/k8s-namespace.yml
```

给pod创建RBAC权限
```
kubectl create -f deploy/k8s/k8s-pod-rbac.yml
```

### account

创建account的pod
```
kubectl create -f account/deploy/k8s/k8s-pod-account.yml
```

创建account的services
```
kubectl create -f account/deploy/k8s/k8s-svc-account.yml
```

> 关于`account`的负载均衡,只需要扩容`deploy`的`replicas`即可  kubectl scale --replicas=3 deployment/svc-account -n go-micro

### apigateway

创建apigateway的pod
```
kubectl create -f apigateway/deploy/k8s/k8s-pod-apigw.yml
```

创建apigateway的services
```
kubectl create -f apigateway/deploy/k8s/k8s-svc-apigw.yml
```

此时,进行验证
```
curl -X POST -d "username=guaosi&password=guaosi" http://127.0.0.1:30088/account/register
```
如果返回 `{"code":0,"message":""}` 则证明成功。可以继续下面使用`traefik`来进行反向代理了。

> 关于`apigateway`的负载均衡,只需要扩容`deploy`的`replicas`即可  kubectl scale --replicas=3 deployment/svc-apigw -n go-micro

### traefik

> 该部分参考 [Kubernetes 部署 Ingress 控制器 Traefik v2.2](http://www.mydlq.club/article/72/) 后运行成功后,总结得出。

#### 创建 CRD 资源

> 在`Traefik v2.0`版本后,开始使用`CRD`（Custom Resource Definition）来完成路由配置等,所以需要提前创建`CRD`资源。
```
kubectl create -f traefik/k8s/k8s-crd-traefik.yaml
```

#### 创建 RBAC 权限

```
kubectl create -f traefik/k8s/k8s-rbac-traefik.yml
```

#### 创建 Traefik 配置文件

下面配置中可以通过配置`kubernetesCRD`与`kubernetesIngress`两项参数,让`Traefik`支持`CRD`与`Ingress`两种路由方式。
```
kubectl create -f traefik/k8s/k8s-config-traefik.yml
```

#### Kubernetes 部署 Traefik

下面将用`DaemonSet`方式部署`Traefik`,便于在多服务器间扩展,用`hostport`方式绑定服务器`80`、`443`端口,方便流量通过物理机进入`Kubernetes`内部。

> 即docker的host模式,但是mac系统无法支持host模式。只能通过nodePort

```
# Linux系统执行
kubectl create -f traefik/k8s/k8s-pod-traefik.yml

# Mac系统执行
kubectl create -f traefik/k8s/mac-k8s-pod-traefik.yml
```

##### 注意

`DaemonSet`保证在每个`Node`都运行一个`Pod`,如果 新增一个`Node`,这个`Pod`也会运行在新增的`Node`上,如果删除这个`DadmonSet`,就会清除它所创建的`Pod`。

如果想指定只能在哪些`node`上创建`traefik`,则需要提前指定`Label`,这样当程序部署时会自动调度到设置`Label`的节点上。

```
# 格式：kubectl label nodes [节点名] [key=value]
kubectl label nodes docker-desktop IngressProxy=true

# 查看节点是否设置 Label 成功
kubectl get nodes --show-labels
```

同时,需要修改 `traefik/k8s/k8s-pod-traefik.yml`
```
      volumes:
        - name: config
          configMap:
            name: traefik-config 
      # 以上是已经存在的设置,下面为要添加的设置
      tolerations:              # 设置容忍所有污点,防止节点被设置污点
        - operator: "Exists"
      nodeSelector:             # 设置node筛选器,在特定label的节点上启动
        IngressProxy: "true"
```

#### 配置路由规则

Traefik 应用已经部署完成,但是想让外部访问`Kubernetes`内部服务,还需要配置路由规则,上面部署`Traefik`时开启了`Traefik Dashboard`,这是`Traefik`提供的视图看板,所以,首先配置基于`HTTP`的`Traefik Dashboard`路由规则,使外部能够访问 `Traefik Dashboard`。这里使用`CRD`进行演示。
> 想使用Ingress或者加上HTTPS认证,请参考[这篇文章](http://www.mydlq.club/article/72/)

> 使用 CRD 方式创建路由规则可言参考 Traefik 文档 [Kubernetes IngressRoute](https://docs.traefik.io/v2.2/routing/providers/kubernetes-crd/)

```
kubectl create -f traefik/k8s/k8s-crd-router-traefik.yml
```

#### 配置 Host 文件

客户端想通过域名访问服务,必须要进行`DNS` 解析,由于这里没有`DNS`服务器进行域名解析,所以修改`hosts`文件将`Traefik`、`apigw`所在节点服务器的`IP`和自定义`Host`绑定,。打开电脑的`Hosts`配置文件,往其加入下面配置：

```
127.0.0.1  traefik.guaosi.com
127.0.0.1  apigw.guaosi.com
```

##### 验证Traefik Dashboard

打开浏览器输入地址：http://traefik.guaosi.com  ,即可打开`Traefik Dashboard`

> mac请打开 http://traefik.guaosi.com:30180 , 因为mac无法使用host模式,只能使用nodePort来映射80端口到宿主机上的30180。

##### 验证网关以及微服务

```
# linux系统请执行
curl -X POST -d "username=guaosi&password=guaosi" http://apigw.guaosi.com/account/register

# mac系统请执行
curl -X POST -d "username=guaosi&password=guaosi" http://apigw.guaosi.com:30180/account/register
```

如果返回 `{"code":0,"message":""}` 则证明成功。