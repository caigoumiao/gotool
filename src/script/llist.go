package script

import (
	"flag"
	"fmt"
	"util"
)

/**
 * 列出项目的全部版本
 */
func List() {
	// 读取项目名
	projectName := flag.Arg(1)

	if projectName == "" {
		println("args error !")
		return
	}

	m := util.GetVersion()
	if m[projectName] == nil {
		println("cannot find project:" + projectName)
		return
	}

	fmt.Printf("Project:%s has following versions :", projectName)
	for k, _ := range m[projectName] {
		fmt.Println(k)
	}
}
