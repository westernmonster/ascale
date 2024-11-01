package memcache

import (
	"ascale/pkg/conf/env"

	"context"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

const (
	_traceFamily        = "memcache"
	_traceSpanKind      = "client"
	_traceComponentName = "library/cache/memcache"
	_tracePeerService   = "memcache"
)

type traceConn struct {
	Conn
	ctx     context.Context
	address string
}

func (t *traceConn) setTrace(action, statement string) func(error) error {
	parent := trace.SpanFromContext(t.ctx)
	if !parent.IsRecording() {
		return func(err error) error { return err }
	}

	tracer := global.Tracer(env.AppID)

	startOpts := []trace.StartOption{
		trace.WithAttributes(
			label.String("peer.service", _tracePeerService),
			label.String("peer.address", t.address),
			label.String("db.stmt", action+" "+statement),
			label.String("component", _traceComponentName),
		),
		trace.WithSpanKind(trace.SpanKindClient),
	}

	ctx, span := tracer.Start(t.ctx, _traceFamily, startOpts...)

	t.ctx = ctx

	return func(err error) error {
		span.End()
		return err
	}
}

func (t *traceConn) WithContext(ctx context.Context) Conn {
	t.ctx = ctx
	t.Conn = t.Conn.WithContext(ctx)
	return t
}

func (t *traceConn) Add(item *Item) error {
	finishFn := t.setTrace("Add", item.Key)
	return finishFn(t.Conn.Add(item))
}

func (t *traceConn) Set(item *Item) error {
	finishFn := t.setTrace("Set", item.Key)
	return finishFn(t.Conn.Set(item))
}

func (t *traceConn) Replace(item *Item) error {
	finishFn := t.setTrace("Replace", item.Key)
	return finishFn(t.Conn.Replace(item))
}

func (t *traceConn) Get(key string) (*Item, error) {
	finishFn := t.setTrace("Get", key)
	item, err := t.Conn.Get(key)
	return item, finishFn(err)
}

func (t *traceConn) GetMulti(keys []string) (map[string]*Item, error) {
	finishFn := t.setTrace("GetMulti", strings.Join(keys, " "))
	items, err := t.Conn.GetMulti(keys)
	return items, finishFn(err)
}

func (t *traceConn) Delete(key string) error {
	finishFn := t.setTrace("Delete", key)
	return finishFn(t.Conn.Delete(key))
}

func (t *traceConn) Increment(key string, delta uint64) (newValue uint64, err error) {
	finishFn := t.setTrace("Increment", key+" "+strconv.FormatUint(delta, 10))
	newValue, err = t.Conn.Increment(key, delta)
	return newValue, finishFn(err)
}

func (t *traceConn) Decrement(key string, delta uint64) (newValue uint64, err error) {
	finishFn := t.setTrace("Decrement", key+" "+strconv.FormatUint(delta, 10))
	newValue, err = t.Conn.Decrement(key, delta)
	return newValue, finishFn(err)
}

func (t *traceConn) CompareAndSwap(item *Item) error {
	finishFn := t.setTrace("CompareAndSwap", item.Key)
	return finishFn(t.Conn.CompareAndSwap(item))
}

func (t *traceConn) Touch(key string, seconds int32) (err error) {
	finishFn := t.setTrace("Touch", key+" "+strconv.Itoa(int(seconds)))
	return finishFn(t.Conn.Touch(key, seconds))
}
