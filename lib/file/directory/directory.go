package directory

import (
	"os"
	"path/filepath"
)


/**获取工作目录*/
func GetWorkDir() (string, error) {
	wd, err := os.Getwd()
	return wd, err
}


/**获取执行目录*/
func GetExecDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	ed := filepath.Dir(execPath)
	return ed, nil
}

/**判断是否为绝对路径**/
func IsAbs(path string) bool {
	return len(path) > 0 && path[0] == '/'
}