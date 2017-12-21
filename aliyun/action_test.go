package aliyun

import (
	"fmt"
	"testing"
)

func TestDescribeLiveStreamPublishList(t *testing.T) {
	start := "2017-12-10T09:56:39Z"
	end := "2017-12-12T09:56:39Z"
	rsp := DescribeLiveStreamPublishList(start, end)
	if rsp == "" {
		t.Errorf("DescribeLiveStreamPublishList failed")
	} else {
		fmt.Printf("DescribeLiveStreamPublishList rsp:%s", rsp)
	}
}
