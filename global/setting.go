package global

import (
	"web-base/pkg/logger"
	"web-base/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	JWTSetting      *setting.JWTSettingS

	Logger *logger.Logger
)
