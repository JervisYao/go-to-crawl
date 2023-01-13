package bilibili

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
	"go-to-crawl-vod/internal/logic/task/dto"
	"go-to-crawl-vod/utility/fileutil"
	"strconv"
)

type BilibiliCrawl struct {
	*dto.AbstractCrawlVodFlow
}

func (crawl *BilibiliCrawl) UseBrowser() bool {
	return false
}

func (crawl *BilibiliCrawl) FillTargetRequest(ctx *dto.BrowserContext) {

	coll := colly.NewCollector()
	coll.OnResponse(func(response *colly.Response) {
		bodyStr := string(response.Body)

		var wd Data
		dataJsons, err := gregex.MatchString(`__INITIAL_STATE__=(.*);\(function\(\)`, bodyStr)
		if err != nil {
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}

		err = gconv.Struct(dataJsons[1], &wd)
		if err != nil {
			ctx.Log.Error(gctx.GetInitCtx(), err)
			return
		}

		var pid string
		pidS, err := gregex.MatchString(`\?p=(\d+)`, ctx.CrawlQueueSeed.CrawlSeedUrl)
		if err != nil || len(pidS) != 2 {
			pid = "1"
		} else {
			pid = pidS[1]
		}

		var cid int
		for _, page := range wd.VideoData.Pages {
			atoi, _ := strconv.Atoi(pid)
			if page.Page == atoi {
				cid = page.Cid
				break
			}
		}

		var video Video
		videoJsons, err := gregex.MatchString(`window\.__playinfo__=(.*)\</script\>\<script\>window\.__INITIAL_STATE__=`, bodyStr)
		if err != nil {
			return
		}
		err = gconv.Struct(videoJsons[1], &video)
		if err != nil {
			return
		}

		if video.Code != 0 {
			return
		}

		innerColl := colly.NewCollector()
		innerColl.OnResponse(func(response *colly.Response) {
			jsonObj := gjson.New(response.Body)
			flvUrl := jsonObj.GetString("data.durl.0.url")
			ctx.Log.Infof(gctx.GetInitCtx(), "flv url = %s", flvUrl)

			downloadBuilder := fileutil.CreateBuilder().Url(flvUrl).SaveFile("D:\\刘星\\bilibili.flv")
			downloadBuilder.Header("Referer", ctx.CrawlQueueSeed.CrawlSeedUrl)
			fileutil.DownloadFileByBuilder(downloadBuilder)
		})
		url := "https://api.bilibili.com/x/player/playurl"
		err = innerColl.Visit(fmt.Sprintf("%s?qn=80&avid=%d&cid=%d", url, wd.Aid, cid))

	})

	err := coll.Visit(ctx.CrawlQueueSeed.CrawlSeedUrl)
	if err != nil {
		ctx.Log.Error(gctx.GetInitCtx(), err)
		return
	}

}

type Data struct {
	Aid       int `json:"aid"`
	VideoData struct {
		Title string `json:"title"`
		Pages []struct {
			Cid       int    `json:"cid"`
			Page      int    `json:"page"`
			Part      string `json:"part"`
			Duration  int    `json:"duration"`
			Dimension struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"dimension"`
		} `json:"pages"`
	} `json:"videoData"`
}

type Video struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Result `json:"data"`
	Result  Result `json:"result"`
}

type Result struct {
	Quality    int    `json:"quality"`
	Format     string `json:"format"`
	Timelength int    `json:"timelength"` // ms
	Durl       []struct {
		URL    string `json:"url"`
		Order  int    `json:"order"`
		Length int    `json:"length"`
		Size   int    `json:"size"`
	} `json:"durl"`
	SupportFormats []struct {
		Quality        int    `json:"quality"`
		Format         string `json:"format"`
		NewDescription string `json:"new_description"`
	} `json:"support_formats"`
}
