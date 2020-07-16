package serviceclient

import (
	"apigw/handler"
	proto "apigw/proto/account"
	hystrix_go "github.com/afex/hystrix-go/hystrix"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	// 这里使用 kubernetes 是为了之后可以通过命令行指定注册中心用 kubernetes
	_ "github.com/micro/go-plugins/registry/kubernetes/v2"
)

var Port string

func RegisterService() {
	// 连接服务注册中心
	service := micro.NewService(
		micro.Flags(
			&cli.StringFlag{
				Name:  "p",
				Usage: "port",
			},
		),
	)
	// 解析命令行参数
	// 我们希望可以使用 -p 参数来手动指定我们HTTP服务对外提供服务时的端口
	service.Init(
		micro.Action(func(c *cli.Context) error {
			Port = c.String("p")
			if len(Port) == 0 {
				Port = "8091"
			}
			return nil
		},
		),
	)
	// 复用服务注册的客户端
	hystrix_go.DefaultVolumeThreshold = 1
	hystrix_go.DefaultErrorPercentThreshold = 1
	cli := client.DefaultClient

	// 重试 只有当某一个服务的其中某一个实例被请求到时，该实例正好出现问题时，才会触发重试
	//cli.Init(
	//	client.Retries(3),
	//	//为了调试看log方便，始终返回true, nil，即会一直重试直至重试次数用尽
	//	client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
	//		log.Print(req.Method(), "retry count:", retryCount, ", client request service failed ,client retry")
	//		return true, nil
	//	}),
	//)

	// 获取在服务注册中心上 micro.service.account 的客户端
	handler.AccountServiceClient = proto.NewAccountService("micro.service.account", cli)
}
