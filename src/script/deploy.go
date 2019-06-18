package script

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"strconv"
	"util"
)

const DEPLOY_PATH_DIR = "/home/miao/deploy/"

func Deploy() {
	// 读取项目名、部署环境
	projectName := flag.Arg(1)
	env := flag.Arg(2)

	if projectName == "" || env == "" {
		println("args error !")
		return
	}

	m := util.GetVersion()
	if m[projectName] == nil {
		println("cannot found projectName: " + projectName)
		return
	}

	/** 1. 查找编译包（先本地，再oss） */
	// 默认获取最新编译包，也可指定

	var max int64 = -9999999
	for k, _ := range m[projectName] {
		v, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			panic(err.Error())
		}

		if v > max {
			max = v
		}
	}

	tarName := fmt.Sprintf("%s_%d", projectName, max)
	ossPath := m[projectName][tarName].OssPath
	// todo: 从oss 获取tar 包

	/** 2. 发送tar包到指定机器，并解压tar 包到指定路径 */

	// 获取shell session
	session, in, out, err := util.Connect(util.GetConfig().Deploy.SSH.Host,
		util.GetConfig().Deploy.SSH.Port,
		util.GetConfig().Deploy.SSH.IdRsaPath)
	// 忽视登录产生的shell output
	<-out
	defer session.Close()

	if err != nil {
		log.Errorf("get shell session error : %s" + err.Error())
		panic(err.Error())
	}

	// 复制tar 包到指定机器
	// todo: scp 命令需要输入 yes/no
	scpCmd := fmt.Sprintf("scp %s %s@%s:%s",
		tarName,
		util.GetConfig().Deploy.SSH.UserName,
		util.GetConfig().Deploy.SSH.Host,
		DEPLOY_PATH_DIR)
	fmt.Printf("begin to mv %s to %s ......", tarName, util.GetConfig().Deploy.SSH.Host)
	in <- scpCmd
	println(<-out)

	// 解压tar 包到指定路径: DEPLOY_PATH/env/projectName
	// todo: 项目路径需要提前创建
	deployPath := DEPLOY_PATH_DIR + "/" + env + "/" + projectName
	tarCmd := fmt.Sprintf("cd %s && tar -zxvf %s -C %s",
		DEPLOY_PATH_DIR,
		tarName,
		deployPath)
	fmt.Printf("begin to untar %s to dir:%s ......", tarName, deployPath)
	in <- tarCmd
	<-out // 解包的日志不输出

	/** 3. 建立screen 环境 */
	// todo: 在screen 环境部署，暂不做

	/** 4. 在指定环境部署 （路径举例：/home/miao/deploy/prod） ./main --env=prod */
	deployCmd := fmt.Sprintf("cd %s && ./%s --env=%s",
		deployPath,
		util.GetConfig().Deploy.MainFunc,
		env)
	fmt.Printf("begin to deploy %s ......", projectName)
	in <- deployCmd
	println(<-out)

	// 退出shell
	in <- "exit"
	session.Wait()
}
