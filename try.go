package try

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type Context struct {
	catched atomic.Bool
	parent  context.Context
	cancel  context.CancelCauseFunc
}

func New(parent context.Context) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(parent)
	return &Context{parent: ctx, cancel: cancel}, func() {
		cancel(nil)
	}
}

func (ctx *Context) Try(fn func() error) {
	defer func() {
		panicValue := recover()
		if panicValue != nil {
			ctx.cancel(fmt.Errorf("panic: %v", panicValue))
		}
	}()
	if err := ctx.Err(); err != nil {
		return
	}
	if err := fn(); err != nil {
		ctx.cancel(errors.WithStack(err))
	}
}

func (ctx *Context) Catch(fn func(error)) {
	if ctx.catched.Load() {
		return
	}

	if err := context.Cause(ctx); err != nil {
		ctx.catched.Store(true)
		fn(err)
	}
}

func (ctx *Context) CatchError(target error, fn func(error)) {
	if ctx.catched.Load() {
		return
	}

	if err := context.Cause(ctx); errors.Is(err, target) {
		ctx.catched.Store(true)
		fn(err)
	}
}

func (ctx *Context) PassError(target error) {
	if err := context.Cause(ctx); errors.Is(err, target) {
		ctx.catched.Store(true)
	}
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.parent.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.parent.Done()
}

func (ctx *Context) Err() error {
	return ctx.parent.Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.parent.Value(key)
}
