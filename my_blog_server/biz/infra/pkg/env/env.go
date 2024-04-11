package env

import (
	"fmt"
	"os"
)

var rootPath string

func MustInitEnv() {
	rootPath = os.Getenv("ROOT_PATH")
	if rootPath == "" {
		panic("get env ROOT_PATH fail")
	}
	fmt.Printf("[ENV] ROOT_PATH=%v\n", rootPath)
}

func RootPath() string {
	return rootPath
}
