package settings

import (
	"os"
	"strings"
)

/*
  @Author : Mustang Kong
*/

// ObjectPath 获取当前项目路径
func ObjectPath() (projectPath string) {
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if strings.HasSuffix(projectPath, "tmp") {
		sliceTmp := strings.Split(projectPath, "/")
		projectPath = strings.Join(sliceTmp[:len(sliceTmp)-1], "/")
	}

	return
}
