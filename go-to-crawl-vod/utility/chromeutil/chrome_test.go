package chromeutil

import (
	"fmt"
	"github.com/corpix/uarand"
	"testing"
)

func TestRandUA(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(uarand.GetRandom())
	}
}
