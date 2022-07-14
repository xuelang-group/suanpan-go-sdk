package component

import (
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/app"
	"github.com/xuelang-group/suanpan-go-sdk/suanpan/v1/stream"
)

type Component interface {
	InitHandler()
	CallHandler(stream.Request)
	SioHandler()
}

func Run(c Component)  {
	c.InitHandler()
	c.SioHandler()
	app.Run(c.CallHandler)
}