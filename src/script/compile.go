package script

import (
	"flag"
	"fmt"
	"github.com/labstack/gommon/log"
	"time"
	"util"
)

const COMPILE_PATH = "/home/miao/compile/"

func Compile() {
	// 读取项目名、分支
	projectName := flag.Arg(1)
	branch := flag.Arg(2)
	if projectName == "" || branch == "" {
		println("args error !")
		return
	}

	// 获取shell session
	session, in, out, err := util.Connect(util.GetConfig().Build.SSH.Host,
		util.GetConfig().Build.SSH.Port,
		util.GetConfig().Build.SSH.IdRsaPath)
	// 忽视登录产生的shell output
	<-out
	defer session.Close()

	if err != nil {
		log.Errorf("get shell session error : %s" + err.Error())
		panic(err.Error())
	}

	//session.Run("echo Hello,World!")

	/** 1. 开始拉取git 代码 */
	log.Infof("begin to pull code ...")
	projectDir := util.GetConfig().Build.Path + "/" + projectName
	gitCmd := fmt.Sprintf("git clone -b %s %s %s",
		branch,
		util.GetConfig().Build.GitUrl,
		projectDir)
	in <- gitCmd
	fmt.Println(<-out)

	/** 2. 下载依赖包 */
	// todo：可能会有某些依赖go get 失败，怎么处理？
	// todo: 可以分多session 去一起下载
	fmt.Println("begin to go get libs......")
	libs := util.GetConfig().Build.Lib
	for i, lib := range libs {
		in <- fmt.Sprintf("go get %s", lib)
		fmt.Printf("%d. go get %s ......\n", i+1, lib)
		println(<-out)
	}

	/** 3. 编译、上传 */
	fmt.Println("begin to build......")
	buildCmd := fmt.Sprintf("cd %s && GOPATH=$GOPATH:$(pwd) go build -v", projectDir)
	//ROOT=$(cd `dirname $0`; pwd)
	in <- buildCmd
	fmt.Println(<-out)

	// 编译结果打包为tar
	// todo: 检查编译目录是否存在

	fmt.Println("begin to tar " + projectDir)
	tarName := fmt.Sprintf("%s-%d.tar", projectName, time.Now().Unix())
	tarCmd := fmt.Sprintf("sudo tar -cvf %s %s && mv %s %s", tarName, projectDir, tarName, COMPILE_PATH+projectName)
	in <- tarCmd
	in <- "ceshi1234"
	<-out

	// tar 上传到oss
	ossPath := ""

	/** 4. 版本信息存储 */
	buildVersion := util.BuildVersion{
		tarName, "master", "xxx", ossPath,
	}
	err = util.AddVersion(buildVersion)
	if err != nil {
		panic(err.Error())
	}

	// 退出shell
	in <- "exit"
	session.Wait()
}
