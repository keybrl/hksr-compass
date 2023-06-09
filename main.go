package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/keybrl/hksr-compass/pkg/commands"
)

var (
	version = "0.0.0-dev"
)

func main() {
	// 将中断信号绑定到上下文
	ctx, cancel := notifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	// 设置版本
	commands.Cmd.Version = version
	// 执行命令
	if err := commands.Cmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}

// notifyContext 将信号绑定到上下文
func notifyContext(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, signals...)
	if ctx.Err() == nil {
		go func() {
			// 第一次取消上下文
			select {
			case <-ctx.Done():
			case <-ch:
				cancel()
			}
			// 第二次直接退出
			select {
			case <-ctx.Done():
			case <-ch:
				os.Exit(1)
			}
		}()
	}
	return ctx, cancel
}
