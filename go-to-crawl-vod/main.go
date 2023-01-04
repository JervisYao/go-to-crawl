package main

import (
	_ "go-to-crawl-vod/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"go-to-crawl-vod/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
