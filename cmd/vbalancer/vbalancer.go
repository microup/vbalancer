package main

import (
	"context"
	"runtime"
	"vbalancer/internal/app"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	app.Run(ctx)
}
