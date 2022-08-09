package rec

import (
	"context"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/utils"
)

type Options struct {
	ctx     context.Context
	http    string
	grpc    string
	name    string
	method  string
	timeout int
}

func WithCtx(ctx context.Context) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).ctx = getCtx(ctx)
	}
}

func WithHttp(s string) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).http = s
	}
}

func WithGrpc(s string) func(*Options) {
	return func(options *Options) {
		getOptionsOrSetDefault(options).grpc = s
	}
}

func WithName(s string) func(*Options) {
	return func(options *Options) {
		if s != "" {
			getOptionsOrSetDefault(options).name = s
		}
	}
}

func WithMethod(s string) func(*Options) {
	return func(options *Options) {
		if s != "" {
			getOptionsOrSetDefault(options).method = s
		}
	}
}

func WithTimeout(seconds int) func(*Options) {
	return func(options *Options) {
		if seconds > 0 {
			getOptionsOrSetDefault(options).timeout = seconds
		}
	}
}

func getOptionsOrSetDefault(options *Options) *Options {
	if options == nil {
		return &Options{
			ctx:     getCtx(nil),
			timeout: 10,
			http:    "http://127.0.0.1:2091",
			name:    "ocr",
			method:  "prediciton",
		}
	}
	return options
}

func getCtx(ctx context.Context) context.Context {
	if utils.InterfaceIsNil(ctx) {
		ctx = context.Background()
	}
	return context.WithValue(ctx, constant.LogSkipHelperCtxKey, false)
}
