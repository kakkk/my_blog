package config

import (
	"fmt"
	"os"
)

var root string

func InitConfig() error {
	err := InitApplicationConfig()
	if err != nil {
		return err
	}
	err = InitCachexConfig()
	if err != nil {
		return err
	}
	return nil
}

func MustInit() {
	rootPath := os.Getenv("ROOT_PATH")
	if rootPath == "" {
		panic(fmt.Errorf("ROOT_PATH not found"))
	}
	root = rootPath
	err := InitConfig()
	if err != nil {
		panic(err)
	}
}
