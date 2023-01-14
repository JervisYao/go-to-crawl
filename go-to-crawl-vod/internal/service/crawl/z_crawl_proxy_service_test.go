package crawl

import (
	"fmt"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"testing"
)

func TestGetProxyByUrl(t *testing.T) {
	fmt.Println(GetProxyByUrl("https://www.nunuyy5.org/dongman/102123.html"))
}
