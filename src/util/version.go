package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type BuildVersion struct {
	Name     string `json:"name"`
	Branch   string `json:"branch"`
	CommitId string `json:"commitId"`
	OssPath  string `json:"ossPath"`
}

/**
 * 添加一个版本信息
 */
func AddVersion(buildVersion BuildVersion) error {
	projectName := strings.Split(buildVersion.Name, "_")[0]
	m := GetVersion()
	m[projectName][buildVersion.Name] = buildVersion

	// 重新写入version.json 文件
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("version.json", bytes, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 通过包名：tarName 取得此版本信息：BuildVersion
 */
func GetVersionByName(name string) BuildVersion {
	// tarName = projectName_time
	// todo: 不规范的包名格式需要检查
	projectName := strings.Split(name, "_")[0]
	m := GetVersionList(projectName)
	return m[name]
}

/**
 * 通过项目名：projectName 取得该项目全部版本信息：map[string]BuildVersion （tarName -> BuildVersion）
 */
func GetVersionList(projectName string) map[string]BuildVersion {
	m := GetVersion()
	return m[projectName]
}

/**
 * 从version.json 中读取全部版本信息
 */
func GetVersion() (m map[string](map[string]BuildVersion)) {
	bytes, err := ioutil.ReadFile("version.json")
	m = make(map[string](map[string]BuildVersion))
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(bytes, &m)
	if err != nil {
		panic(err.Error())
	}
	return
}
