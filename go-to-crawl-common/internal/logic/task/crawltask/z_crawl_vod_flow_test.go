package crawltask

import (
	"fmt"
	"go-to-crawl-common/internal/model/entity"
	"os/exec"

	"testing"
)

func TestNunuyy(t *testing.T) {
	doStartTmpSeed("https://www.nunuyy5.org/dongman/102123.html")
}

func TestNunuyyWithParams(t *testing.T) {
	seed := new(entity.CrawlQueue)
	seed.CrawlSeedUrl = "https://www.nunuyy5.org/dongman/102123.html"
	// 努努大部分资源都支持量子资源，且量子资源的m3u8未加密，因此videoItem需要切换到量子资源下的节目单来摘录到数据库
	seed.CrawlSeedParams = `{"videoItem":"第25集"}`
	DoStartCrawlVodFlow(seed)
}

// 测试通过，按需求对接数据库改成多线程下载就行
func TestBilibili(t *testing.T) {
	doStartTmpSeed("https://www.bilibili.com/video/BV1yg411T7Za?p=7&vd_source=e8bcede57b979b0eed49d9041d869a8e")
}

func TestCrawlUrlType1Task(t *testing.T) {
	CrawlTask.CrawlUrlType1Task(nil)
}

func doStartTmpSeed(url string) {
	seed := new(entity.CrawlQueue)
	seed.CrawlSeedUrl = url
	DoStartCrawlVodFlow(seed)
}

func TestName(t *testing.T) {
	cmd := exec.Command("curl", "https://www.baidu.com")

	out, err := cmd.Output()
	fmt.Println(out)
	fmt.Println(err)

	c := "curl https://www.baidu.com"

	cmd = exec.Command("sh", "-c", c)

	out, err = cmd.Output()
	fmt.Println(out)
	fmt.Println(err)
}
