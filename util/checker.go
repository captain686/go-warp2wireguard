package util

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
)

func ConfigFileCheck(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Error(fmt.Sprintf("%s does not exist", filename))
		return false
	}

	// 判断文件是否为空
	fileInfo, err := os.Stat(filename)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to obtain %s file information: %E", filename, err))
		return false
	}
	if fileInfo.Size() == 0 {
		//log.Error(fmt.Sprintf("%s 为空\n", filename))
		return false
	} else {
		//log.Info(fmt.Sprintf("文件 %s 不为空\n", filename))
		return true
	}
}
