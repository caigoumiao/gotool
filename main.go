package main

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"strings"
	"time"
	"util"
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
		compile()
	case DEPLOY:
		deploy()
	case LIST:
		list()
	case ROLLBACK:
		roolback()
	default:
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

	// 获取shell session
	session, in, out, err := util.Connect(util.GetConfig().SSH.Host, util.GetConfig().SSH.Port, util.GetConfig().SSH.IdRsaPath)
	// 忽视登录产生的shell output
	<-out
	defer session.Close()

	if err != nil {
		log.Errorf("get shell session error : %s" + err.Error())
		panic(err.Error())
	}

	//session.Run("echo Hello,World!")

	// 开始拉取git 代码
	// 脚本执行过程中的日志输出

	// todo : 加上分支
	log.Infof("begin to pull code ...")
	projectDir := util.GetConfig().Build.Path + "/" + projectName
	gitCmd := fmt.Sprintf("git clone -b %s %s %s",
		branch,
		util.GetConfig().Build.GitUrl,
		projectDir)
	in <- gitCmd
	fmt.Println(<-out)

	// 下载依赖包
	// todo：可能会有某些依赖go get 失败，怎么处理？
	// todo: 可以分多session 去一起下载
	fmt.Println("begin to go get libs......")
	libs := util.GetConfig().Build.Lib
	for i, lib := range libs {
		in <- fmt.Sprintf("go get %s", lib)
		fmt.Printf("%d. go get %s ......\n", i+1, lib)
		println(<-out)
	}

	// 开始build 吗？
	fmt.Println("begin to build......")
	buildCmd := fmt.Sprintf("cd %s && go build -v main.go", projectDir)
	in <- buildCmd
	fmt.Println(<-out)

	// 编译结果打包为tar
	fmt.Println("begin to tar " + projectDir)
	tarName := fmt.Sprintf("%s-%d.tar", projectName, time.Now().Unix())
	tarCmd := fmt.Sprintf("sudo tar -cvf %s %s", tarName, projectDir)
	in <- tarCmd
	in <- "ceshi1234"
	<-out

	// tar 上传到oss

	// 退出shell
	in <- "exit"
	session.Wait()
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
