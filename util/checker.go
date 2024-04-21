package util

import (
	"fmt"
	"github.com/charmbracelet/log"
	"os"
)

func FileCheck(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Error(fmt.Sprintf("%s does not exist", filepath))
		return false
	}

	// 判断文件是否为空
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to obtain %s file information: %E", filepath, err))
		return false
	}
	if fileInfo.Size() == 0 {
		//log.Error(fmt.Sprintf("%s 为空\n", filepath))
		return false
	} else {
		//log.Info(fmt.Sprintf("文件 %s 不为空\n", filepath))
		return true
	}
}
