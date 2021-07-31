package setting

import (
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath("configs/")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	// vp.WatchConfig()
	// vp.OnConfigChange(func(in fsnotify.Event) {
	// 	fmt.Println("配置文件已修改")
	// })
	return &Setting{vp}, nil
}
