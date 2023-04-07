package service1

import (
	"context"
	"fmt"
	"github.com/wastill/my-core-demo/framework"
	"log"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	// 结束信号
	finish := make(chan struct{}, 1)
	// panic信号
	panicChan := make(chan interface{}, 1)

	// 设置超时时间
	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	defer cancel()

	// mu := sync.Mutex{}
	// 执行业务逻辑
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// Do real action
		//time.Sleep(10 * time.Second)
		c.Json(200, "ok1")

		finish <- struct{}{}
	}()

	// 监听panic、程序运行结束、程序超时的信号
	select {
	case p := <-panicChan:
		c.WriterMutex().Lock()
		defer c.WriterMutex().Unlock()
		log.Println(p)
		c.Json(500, "panic")
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		c.WriterMutex().Lock()
		defer c.WriterMutex().Unlock()
		c.Json(500, "time out")
		c.SetTimeOut()
	}
	return nil
}
