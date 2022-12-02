package pluginlog

import (
	"fmt"
	"github.com/husanpao/logrus-easy-formatter"
	"github.com/husanpao/timewriter"
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

type PluginLog struct {
	*logrus.Logger
}

func (plug *PluginLog) Black(format string, args ...interface{}) {
	plug.Info(aurora.Black(fmt.Sprintf(format, args...)))
}
func (plug *PluginLog) Red(format string, args ...interface{}) {
	plug.Info(aurora.Red(fmt.Sprintf(format, args...)))
}
func (plug *PluginLog) Green(format string, args ...interface{}) {
	plug.Info(aurora.Green(fmt.Sprintf(format, args...)))
}
func (plug *PluginLog) Yellow(format string, args ...interface{}) {
	plug.Info(aurora.Yellow(fmt.Sprintf(format, args...)))
}

func (plug *PluginLog) Blue(format string, args ...interface{}) {
	plug.Info(aurora.Blue(fmt.Sprintf(format, args...)))
}
func (plug *PluginLog) Magenta(format string, args ...interface{}) {
	plug.Info(aurora.Magenta(fmt.Sprintf(format, args...)))
}
func (plug *PluginLog) Cyan(format string, args ...interface{}) {
	plug.Info(aurora.Cyan(fmt.Sprintf(format, args...)))
}
func (plug *PluginLog) White(format string, args ...interface{}) {
	plug.Info(aurora.White(fmt.Sprintf(format, args...)))
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
