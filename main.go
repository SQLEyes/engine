package engine

import (
	"fmt"
	"github.com/sqleyes/engine/abstract"
	"github.com/sqleyes/engine/pluginlog"
	"gopkg.in/yaml.v3"
	"moul.io/banner"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

var (
	log       = pluginlog.NewPluginLog()
	Plugins   = make(map[string]*Plugin)
	ExecPath  = os.Args[0]
	ExecDir   = filepath.Dir(ExecPath)
	ConfigRaw []byte
)

func Run(configFile string) (err error) {
	if ConfigRaw, err = os.ReadFile(configFile); err != nil {
		log.Error("read config file error:", err.Error())
		return
	}
	var cg abstract.Config
	if ConfigRaw != nil {
		if err = yaml.Unmarshal(ConfigRaw, &cg); err != nil {
			log.Error("parsing yml error:", err)
			return
		}
	}
	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println(banner.Inline("sqleyes"))
	log.Infof("sqleyes@1.0 %s", " start success")
	for _, plugin := range Plugins {
		//配置文件注入到每个插件
		plugin.setConfig(cg[strings.ToLower(plugin.Name)])
		//通知用户插件安装完成了
		switch plugin.ptr.React(abstract.Installed(plugin.status)) {
		case abstract.Start:
			go plugin.startCap()
		default:
			log.Warnf("%s installed but not in use", plugin.Name)
		}
	}
	<-c
	log.Infof("Bye!")
	return
}
