package handler

import (
	"account/proto"
	"context"
	"fmt"
)

type AccountService struct {
}

func (a *AccountService) AccountRegister(c context.Context, req *proto.ReqAccountRegister, res *proto.ResAccountRegister) error {
	fmt.Println("hit here")
	if req.Username == "guaosi" && req.Password == "guaosi" {
		res.Code = 0
		res.Message = ""
		return nil
	}
	res.Code = -1
	res.Message = "账号或者密码不正确"
	return nil
}
