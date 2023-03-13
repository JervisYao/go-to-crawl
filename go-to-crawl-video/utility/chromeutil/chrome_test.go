package chromeutil

import (
	"fmt"
	"testing"
)

func TestRandUA(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(GetRandomUA(false))
	}
}
