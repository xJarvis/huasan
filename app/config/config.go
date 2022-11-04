package config

import (
	"errors"
	"gopkg.in/ini.v1"
	"strings"
)

var cfg *ini.File

/**初始化配置文件导入*/
func InitConfig(confFile string) {
	var err error
	cfg, err = ini.Load(confFile)
	if err != nil {
		panic(err)
	}
}


/**获取配置文件指定value*/
func GetValue(key string) (string,error) {
	if cfg == nil {
		return "",errors.New("config not init")
	}
	keys := strings.Split(key,".")
	if len(keys) != 2 {
		return "",errors.New("config key failed")
	}
	session := cfg.Section(keys[0])
	if !session.HasKey(keys[1]) {
		return "",errors.New("config key not exist")
	}
	value := session.Key(keys[1]).Value()
	return value,nil
}


/**获取配置文件指定value[若不存在返回空]*/
func Get(key string) string {
	value,err := GetValue(key)
	if err != nil {
		return ""
	}
	return value
}


/**获取配置文件指定int[若不存在返回0]*/
func GetInt(key string) int {
	if cfg == nil {
		return 0
	}
	keys := strings.Split(key,".")
	if len(keys) != 2 {
		return 0
	}
	session := cfg.Section(keys[0])
	if !session.HasKey(keys[1]) {
		return 0
	}
	value,err := session.Key(keys[1]).Int()
	if err != nil {
		return 0
	}
	return value
}



/**新参数设定*/
func SetValue(key string,value string) error {
	if cfg == nil {
		return errors.New("config not init")
	}
	keys := strings.Split(key,".")
	if len(keys) != 2 {
		return errors.New("config key failed")
	}
	session := cfg.Section(keys[0])
	_,err := session.NewKey(keys[1],value)
	return err
}
