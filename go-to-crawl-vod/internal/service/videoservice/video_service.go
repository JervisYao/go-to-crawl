package videoservice

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"go-to-crawl-vod/internal/model/entity"
	"go-to-crawl-vod/internal/service/crawlservice"
)

func GetVideoDir(countryCode string, videoYear int, videoCollId, videoItemId int64) string {
	S := gfile.Separator
	dfsRootPathVar, _ := g.Cfg().Get(gctx.GetInitCtx(), "dfs.rootPath")
	dfsRootPath := dfsRootPathVar.String()
	videoRootPath := fmt.Sprintf("%s%s%s%s", dfsRootPath, S, "video", S)
	finalFilePath := fmt.Sprintf("%s%s%s%d%s%d%s%d%s", videoRootPath, countryCode, S, videoYear, S, videoCollId, S, videoItemId, S)
	return finalFilePath
}

func UpdateDownloadStatus(seed *entity.CrawlQueue, err error) {
	if err != nil {
		g.Log().Line().Error(gctx.GetInitCtx(), err)
		seed.ErrorMsg = err.Error()
		crawlservice.UpdateStatus(seed, crawlservice.M3U8Err)
	} else {
		crawlservice.UpdateStatus(seed, crawlservice.M3U8Parsed)
	}
}
