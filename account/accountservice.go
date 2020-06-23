package main

import (
	"account/handler"
	"account/proto"
	"github.com/micro/go-micro/v2"
	_ "github.com/micro/go-plugins/registry/kubernetes/v2"
	"log"
)

func main() {
	service := micro.NewService(
		micro.Name("micro.service.account"),
	)
	service.Init()
	if err := proto.RegisterAccountServiceHandler(service.Server(), new(handler.AccountService)); err != nil {
		log.Print(err.Error())
		return
	}
	if err := service.Run(); err != nil {
		log.Print(err.Error())
		return
	}
}
