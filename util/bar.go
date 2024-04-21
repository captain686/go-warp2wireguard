package util

import (
	"github.com/cheggaaa/pb/v3"
	"time"
)

func SleepBar(duration int) {
	bar := pb.Full.Start(duration)

	// 设置进度条样式
	bar.SetTemplateString(`{{ "Sleep Countdown ... :" }} {{percent . }} {{ bar . "|" "=" ( cycle . ">" ) "_" "|"}} {{counters . }} {{"s"}}`)

	// 每秒更新进度条
	for i := 0; i < duration; i++ {
		time.Sleep(time.Second)
		bar.Increment()
	}
	bar.Finish()
}

func RetryBar(duration int) {
	bar := pb.Full.Start(duration)

	// 设置进度条样式
	bar.SetTemplateString(`{{ "Waiting for retry ... :" }} {{percent . }} {{ bar . "|" "=" ( cycle . ">" ) "_" "|"}} {{counters . }} {{"s"}}`)

	// 每秒更新进度条
	for i := 0; i < duration; i++ {
		time.Sleep(time.Second)
		bar.Increment()
	}
	bar.Finish()
}
