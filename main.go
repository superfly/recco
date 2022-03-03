package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {

	ctx, cancel := newContext()
	defer cancel()
	run(ctx)
}

func newContext() (context.Context, context.CancelFunc) {
	// NOTE: when signal.Notify is called for os.Interrupt it traps both
	// ^C (Control-C) and ^BREAK (Control-Break) on Windows.

	signals := []os.Signal{os.Interrupt}
	if runtime.GOOS != "windows" {
		signals = append(signals, syscall.SIGTERM)
	}

	return signal.NotifyContext(context.Background(), signals...)
}

type DetectedSource struct {
	Runtime          string
	Framework        string
	RuntimeVersion   string
	FrameworkVersion string
}

func run(ctx context.Context) {
	fmt.Println("implement me!")
	// See README for implementation ideas
}
