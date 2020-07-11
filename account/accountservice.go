package main

import (
	"account/handler"
	"account/proto"
	"github.com/micro/go-micro/v2"
	// 这里使用 kubernetes 是为了之后可以通过命令行指定注册中心用 kubernetes
	_ "github.com/micro/go-plugins/registry/kubernetes/v2"
	"log"
)

func main() {
	service := micro.NewService(
		micro.Name("micro.service.account"),
	)
	// 初始化相关操作
	service.Init()
	// 注册实现了服务的handler
	if err := proto.RegisterAccountServiceHandler(service.Server(), new(handler.AccountService)); err != nil {
		log.Print(err.Error())
		return
	}
	// 运行server
	if err := service.Run(); err != nil {
		log.Print(err.Error())
		return
	}
}
