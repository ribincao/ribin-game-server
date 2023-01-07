package utils

import (
	"context"
	"runtime"

	logger "github.com/ribincao/ribin-game-server/logger"
	"go.uber.org/zap"
)

func PrintStack(r interface{}) {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	logger.Error("[ENGINE_SERVER_PANIC]", zap.Any("err", r), zap.String("Stack", string(buf[:n])))
}

func RunWithRecover(f func()) {
	defer func() {
		if r := recover(); r != nil {
			PrintStack(r)
		}
	}()

	f()
}

func GoWithRecover(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				PrintStack(r)
			}
		}()

		f()
	}()
}

func GoCtxWithRecover(ctx context.Context, f func(context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				PrintStack(r)
			}
		}()

		f(ctx)
	}()
}
