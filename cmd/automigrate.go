package main

import (
	"fmt"
	"web-base/global"
	"web-base/internal/model"
	"web-base/pkg/setting"
)

func init() {
	setting, err := setting.NewSetting()
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		fmt.Println(err)
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		fmt.Println(err)
	}
	err = setupDBEngine()
	if err != nil {
		fmt.Println(err)
	}
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	global.DBEngine.AutoMigrate(&model.Article{}, &model.Tag{}, &model.Auth{})
	fmt.Println("自动迁移数据库完成")
}
