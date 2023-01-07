package file

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"strings"
	"testing"
)

var (
	testUrl = "https://xlzycdn1.sy-precise.com:65/20220712/Ioe91OwL/2637kb/hls/index.m3u8"
	path    = "/app/docker/go-fastdfs/data/files/video/CN/11/22/org_index.m3u8"
)

func TestDownloadFile(t *testing.T) {
	proxyUrl := crawl.GetProxyByUrl(testUrl)
	err := DownloadFile(testUrl, proxyUrl, path, 2)
	fmt.Println(err)
}

func TestDownloadSeed(t *testing.T) {

	seed := new(model.CmsCrawlQueue)
	seed.CrawlSeedUrl = "https://dow.dowlz6.com/20221011/13088_57432b23/灵剑尊第322集.mp4"

	builder := CreateBuilder()
	builder.Url(seed.CrawlSeedUrl)

	path1 := "D:/cache2/"
	_ = gfile.Mkdir(path1)
	builder.SaveFile(path1 + "org.mp4")
	err2 := DownloadFileByBuilder(builder)

	fmt.Println(err2)
	fmt.Println("结束")
}

func TestSubStr(t *testing.T) {
	name := "111.jpg"
	str := gstr.SubStr(name, 0, strings.Index(name, "."))
	fmt.Println(str)
}
