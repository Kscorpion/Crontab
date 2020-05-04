package master

import (
	"encoding/json"
	"io/ioutil"
)

//程序配置
type Config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeOut  int      `json:"apiReadTimeOut"`
	ApiWriteTimeOut int      `json:"apiWriteTimeOut"`
	EtcdEndPoints   []string `json:"etcdEndPoints"`
	etcdDialTimeOut int      `json:"etcdDialTimeOut"`
}

var (
	//单例
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	//把配置文件读进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	//json反序列化
	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = &conf
	return
}
