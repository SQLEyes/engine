package engine

import (
	"engine/config"
	"engine/pluginlog"
	"fmt"
	"moul.io/banner"
	"os"
	"os/signal"
)

var (
	log     = pluginlog.NewPluginLog()
	Plugins = make(map[string]*Plugin)
)

func Run() {
	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println(banner.Inline("sqleyes"))
	log.Infof("sqleyes@1.0 %s", " start success")
	for _, plugin := range Plugins {
		plugin.ptr.React(config.Installed(plugin.status))
		if plugin.status.Code == 200 {
			go plugin.startCap()
		}
	}
	<-c
	log.Infof("Bye!")
}
