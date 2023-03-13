package timeutil

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"
	"testing"
	"time"
)

func TestTimeAdd(t *testing.T) {
	waterMark := gtime.Now().Add(time.Duration(30) * -time.Minute)
	fmt.Println(waterMark)
}
