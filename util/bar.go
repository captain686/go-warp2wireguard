package util

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"time"
)

func SleepBar(describe string, duration int) {
	bar := pb.Full.Start(duration)

	// 设置进度条样式
	bar.SetTemplateString(fmt.Sprintf(`{{ "%s ... :" }} {{percent . }} {{ bar . "|" "=" ( cycle . ">" ) "_" "|"}} {{counters . }} {{"s"}}`, describe))

	// 每秒更新进度条
	for i := 0; i < duration; i++ {
		time.Sleep(time.Second)
		bar.Increment()
	}
	bar.Finish()
}
