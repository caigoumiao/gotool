package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/kylelemons/go-gypsy/yaml"
	"os/exec"
	"strings"
)

const (
	COMPILE = "COMPILE"
	DEPLOY = "DEPLOY"
	LIST = "LIST"
	ROLLBACK = "ROLLBACK"
)

func main() {
	// 获取命令行参数 : 操作类型
	var op string
	flag.Parse()
	op = flag.Arg(0)
	switch strings.ToUpper(op) {
		case COMPILE : compile()
		case DEPLOY : deploy()
		case LIST : list()
		case ROLLBACK : roolback()
		default :
			println("error")
	}
}

// 编译
func compile() {
	// 读取项目名、分支
	projectName := flag.Arg(1)
	branch := flag.Arg(2)
	if projectName == "" || branch == "" {
		println("args error !")
		return
	}

	// 读取配置文件
	config, err := yaml.ReadFile("conf.yaml")
	if err != nil {
		println("error : " + err.Error())
		return
	}

	gitUrl, err := config.Get("git.url")
	path, err := config.Get("path")

	if err != nil {
		println("error : " + err.Error())
		return
	}

	// 开始拉取git 代码
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(fmt.Sprintf("git clone %s %s", gitUrl, path))
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		println("cmd run error : " + err.Error())
	}

	println("cmd stdout = " + stdout.String())
}

// 部署
func deploy() {
	// 读取项目名、部署环境
	projectName := flag.Arg(1)
	env := flag.Arg(2)

	if projectName == "" || env == "" {
		println("args error !")
		return
	}
	// 读取配置文件

	println(projectName + ":" + env)
}

// 列出项目的之前版本
func list() {
	// 读取项目名
	projectName := flag.Arg(1)

	if projectName == "" {
		println("args error !")
		return
	}
}

// 回滚到指定版本
func roolback() {
	// 读取版本ID
	ver_id := flag.Arg(1)

	if ver_id == "" {
		println("args error !")
		return
	}
}