package pluginlog

import (
	"github.com/husanpao/logrus-easy-formatter"
	"github.com/husanpao/timewriter"
	"github.com/sirupsen/logrus"
)

type PluginLog struct {
	*logrus.Logger
}

func NewPluginLog(args ...string) *PluginLog {
	p := &PluginLog{}
	name := "engine"
	if len(args) > 0 {
		name = args[0]
	}
	log := logrus.New()
	log.SetReportCaller(true)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%lvl% - [" + name + "]: %time% - %msg%\n",
	})
	log.SetLevel(logrus.DebugLevel)
	timeWriter := &timewriter.TimeWriter{
		Dir:           "./logs",
		Compress:      true,
		ReserveDay:    30,
		Screen:        true,
		LogFilePrefix: name, // default is process name
	}
	log.Out = timeWriter
	p.Logger = log
	return p
}
