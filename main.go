package main

import (
	"flag"
	"script"
	"strings"
)

const (
	COMPILE  = "COMPILE"
	DEPLOY   = "DEPLOY"
	LIST     = "LIST"
	ROLLBACK = "ROLLBACK"
)

func main() {
	// 获取命令行参数 : 操作类型
	var op string
	flag.Parse()
	op = flag.Arg(0)
	switch strings.ToUpper(op) {
	case COMPILE:
		script.Compile()
	case DEPLOY:
		script.Deploy()
	case LIST:
		script.List()
	case ROLLBACK:
		script.Rollback()
	default:
		println("error")
	}
}
