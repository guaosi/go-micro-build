package main

import (
	"apigw/router"
	"apigw/serviceclient"
	tracer "apigw/trace"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/opentracing/opentracing-go"
	"log"
	"net/http"
)

func main() {
	t, io, err := tracer.NewTracer("apigateway", "127.0.0.1:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	serviceclient.RegisterService()
	r := router.NewRouter()

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go func() {
		err := http.ListenAndServe(":81", hystrixStreamHandler)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := r.Run("0.0.0.0:" + serviceclient.Port); err != nil {
		log.Print(err.Error())
	}
}
