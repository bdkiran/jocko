package nolan

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/bdkiran/nolan/protocol"
)

//Context is an interface that provides???
type Context struct {
	mu     sync.Mutex
	conn   io.ReadWriter
	err    error
	header *protocol.RequestHeader
	parent context.Context
	req    interface{}
	res    interface{}
	vals   map[interface{}]interface{}
}

//Request returns context request
func (ctx *Context) Request() interface{} {
	return ctx.req
}

//Response returns context response
func (ctx *Context) Response() interface{} {
	return ctx.res
}

//Header returns context header
func (ctx *Context) Header() *protocol.RequestHeader {
	return ctx.header
}

//Deadline returns the current time...
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

//Done returns nil..
func (ctx *Context) Done() <-chan struct{} {
	return nil
}

//Err returns the context error if any
func (ctx *Context) Err() error {
	return ctx.err
}

func (ctx *Context) String() string {
	return fmt.Sprintf("ctx: %s", ctx.header)
}

//Value returns a value for a provided key from vals within context
func (ctx *Context) Value(key interface{}) interface{} {
	ctx.mu.Lock()
	if ctx.vals == nil {
		ctx.vals = make(map[interface{}]interface{})
	}
	val := ctx.vals[key]
	if val == nil {
		val = ctx.parent.Value(key)
	}
	ctx.mu.Unlock()
	return val
}
