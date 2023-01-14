package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "go-to-crawl-vod/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"go-to-crawl-vod/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
