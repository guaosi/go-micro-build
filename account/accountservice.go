package main

import (
	"account/handler"
	"account/proto"
	tracer "account/trace"
	"github.com/micro/go-micro/v2"
	"time"

	// 这里使用 kubernetes 是为了之后可以通过命令行指定注册中心用 kubernetes
	_ "github.com/micro/go-plugins/registry/kubernetes/v2"
	openTrace "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"log"
)

func main() {
	serviceName := "micro.service.account"
	t, io, err := tracer.NewTracer(serviceName, "127.0.0.1:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	service := micro.NewService(
		micro.Name(serviceName),
		micro.WrapHandler(openTrace.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.RegisterTTL(time.Second*30),      // 在注册中心上设置的过期时间
		micro.RegisterInterval(time.Second*20), // 本服务间隔自动重新注册时间
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
