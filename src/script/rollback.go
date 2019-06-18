package script

import "flag"

/**
 * 滚到指定版本
 */
func Rollback() {
	// 读取版本ID
	ver_id := flag.Arg(1)

	if ver_id == "" {
		println("args error !")
		return
	}
}
