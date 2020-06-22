package main

import (
	"apigw/router"
	"apigw/serviceclient"
	"log"
)

func init() {
	serviceclient.RegisterService()
}
func main() {
	r := router.NewRouter()
	if err := r.Run("0.0.0.0:" + serviceclient.Port); err != nil {
		log.Print(err.Error())
	}
}
