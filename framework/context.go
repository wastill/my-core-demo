package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	ctx            context.Context
	handler        ControllerHandler

	hasTimeout  bool
	writerMutex sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		ctx:            r.Context(),
		writerMutex:    sync.Mutex{},
	}
}

// ========================= base function =========================

func (ctx *Context) WriterMutex() *sync.Mutex {
	return &ctx.writerMutex
}

func (ctx *Context) SetTimeOut() {
	ctx.hasTimeout = true
}

func (ctx *Context) HasTimeOut() bool {
	return ctx.hasTimeout
}

func (ctx *Context) GetRequest() *http.Request {
	return ctx.request
}

func (ctx *Context) GetResponse() http.ResponseWriter {
	return ctx.responseWriter
}

// ========================= base function ========================= end

// ========================= baseContexter function =========================

func (ctx *Context) BaseContext() context.Context {
	return ctx.request.Context()
}

// ========================= baseContexter function ========================= end

// ========================= context.Context function =========================

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.BaseContext().Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.BaseContext().Done()
}

func (ctx *Context) Err() error {
	return ctx.BaseContext().Err()
}

func (ctx *Context) Value(key interface{}) interface{} {
	return ctx.BaseContext().Value(key)
}

// ========================= context.Context function ========================= end

// ========================= query url function =========================
// http:/.localhost:8080/foo?name=bar&age=18           // query string return name and age

func (ctx *Context) QueryInt(key string, def int) int {
	if ctx.request == nil {
		return def
	}
	if v := ctx.request.URL.Query().Get(key); v != "" {
		return cast.ToInt(v)
	}
	return def
}

func (ctx *Context) QueryString(key string, def string) string {
	if ctx.request == nil {
		return def
	}
	if v, ok := ctx.QueryAll()[key]; ok {
		return v[0]
	}
	return def
}

func (ctx *Context) QueryArray(key string, def []string) []string {
	if ctx.request == nil {
		return def
	}
	return ctx.request.URL.Query()[key]
}
func (ctx *Context) QueryAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}
	return cast.ToStringMapStringSlice(ctx.request.URL.Query())
}

// ========================= query url function ========================= end

// ========================= request form function =========================

func (ctx *Context) FormInt(key string, def int) int {
	if ctx.request == nil {
		return def
	}
	if v, ok := ctx.FormAll()[key]; ok {
		return cast.ToInt(v[0])
	}
	return def
}

func (ctx *Context) FormString(key string, def string) string {
	if ctx.request == nil {
		return def
	}
	if v, ok := ctx.FormAll()[key]; ok {
		return v[0]
	}
	return def
}

func (ctx *Context) FormArray(key string, def []string) []string {
	if ctx.request == nil {
		return def
	}
	if v, ok := ctx.request.PostForm[key]; ok {
		return v
	}
	return def
}

func (ctx *Context) FormAll() map[string][]string {
	if ctx.request == nil {
		return map[string][]string{}
	}
	return cast.ToStringMapStringSlice(ctx.request.ParseForm)
}

// ========================= request form function ========================= end

// ========================= request application/json post function =========================

func (ctx *Context) BindJson(obj interface{}) error {
	if ctx.request != nil {
		all, err := ioutil.ReadAll(ctx.request.Body)
		if err != nil {
			return err
		}
		// 重新设置body， request.Body是一个io.ReadCloser类型，只能读取一次，所以需要重新设置
		ctx.request.Body = ioutil.NopCloser(bytes.NewBuffer(all))

		if err := json.Unmarshal(all, obj); err != nil {
			return err
		}

	} else {
		return errors.New("request is nil")
	}
	return nil
}

// ========================= request application/json post function ========================= end

// ========================= response function =========================
func (ctx *Context) Json(status int, obj interface{}) error {
	if ctx.HasTimeOut() {
		return nil
	}
	ctx.responseWriter.Header().Set("Content-Type", "application/json")
	ctx.responseWriter.WriteHeader(status)
	byt, err := json.Marshal(obj)
	if err != nil {
		ctx.responseWriter.WriteHeader(500)
		return err
	}
	ctx.responseWriter.Write(byt)
	return nil
}

func (ctx *Context) HTML(status int, obj interface{}, template string) error {
	return nil
}

func (ctx *Context) Text(status int, obj string) error {
	return nil
}

// ========================= response function ========================= end
