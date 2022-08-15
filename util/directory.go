package util

import (
	"os"
)

// PathExistsOrNot 判断所给路径的文件或文件夹是否存在
func PathExistsOrNot(path string) bool {
	//os.Stat获取文件信息
	_, err := os.Stat(path)
	/*
		当函数返回true时，其实文件并不一定存在。
		对目标path中的某一部分没有可读权限时，os.Lstat和syscall.Access同样会返回error，
		不过这个error不会让os.IsNotExist返回true。
		当文件不存在而你对文件所在的目录或者它的上层目录没有访问权限时，函数依旧会返回true，bug就在这时发生了。
		所以重要的一点是在判断文件是否存在前，应该先判断自己对文件及其路径是否有访问权限。
	*/
	if err != nil {
		//只有当错误为“文件或文件夹已存在”，才能返回true
		if os.IsExist(err) {
			return true
		}
		//否则返回false，可能是权限不够等问题，反正不是“文件或文件夹已存在”
		return false
	}
	return true
}
