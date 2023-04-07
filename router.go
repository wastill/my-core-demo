package main

import (
	"github.com/wastill/my-core-demo/app/service1"
	"github.com/wastill/my-core-demo/framework"
)

func RegisterRouter(core *framework.Core) error {
	core.Get("foo", service1.FooControllerHandler)
	return nil
}
