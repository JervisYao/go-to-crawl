package ffmpegutil

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
)

var (
	url  = "https://xlzycdn1.sy-precise.com:65/20220712/Ioe91OwL/2637kb/hls/index.m3u8"
	path = "D:\\cache2\\video\\CN\\2022\\12\\36\\org_index_bak2.m3u8"
)

func TestConvertM3U8(t *testing.T) {
	log := g.Log().Line()
	m3U8, _ := ConvertM3U8(url, "", path)
	log.Info(gctx.GetInitCtx(), m3U8)
}

func TestPaddingKey(t *testing.T) {
	ret := paddingKeyUrl("https://v.v1kd.com", "#EXT-X-KEY:METHOD=AES-128,URI=\"/20220510/EAlZ3rpV/2000kb/hls/key.key\",IV=0x9180b4da3f0c7e80975fad685f7f134e #EXTINF:6.416667")
	fmt.Println(ret)
}

func TestFormatExtName(t *testing.T) {
	imgMode := false
	name := formatExtName("http://xxx.com/a.jpg?sign=khi34h", &imgMode)
	fmt.Println(name)
	fmt.Println(imgMode)
}

func TestTruncateTS(t *testing.T) {
	truncateTS("D:\\下载\\media_34.ts", SrcTypeNormal, 308)
}

func TestIsPngType(t *testing.T) {
	pngType := IsPngType("D:\\cache2\\video\\YDJ\\2022\\1\\4\\segment-0001.ts")
	fmt.Println(pngType)
}

func TestRunFfmpegGenericMerge(t *testing.T) {
	err := RunFfmpegGenericMerge("D:\\cache2\\replay")
	fmt.Println(err)
}
