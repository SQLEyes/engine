package engine

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/sqleyes/engine/abstract"
	"github.com/sqleyes/engine/pluginlog"
	"reflect"
	"runtime"
	"strings"
)

type Plugin struct {
	*pluginlog.PluginLog
	Name      string
	Version   string
	ptr       abstract.Plugin
	handle    *pcap.Handle
	status    abstract.Msg
	BPFFilter string
	Device    string
	DEBUG     bool
}

func (p *Plugin) setConfig(config any) {
	v := config.(abstract.Config)
	s := reflect.ValueOf(p.ptr)
	for key, value := range v {
		field := s.Elem().FieldByName(key)
		if field.IsValid() {
			if key == "BPFFilter" {
				p.BPFFilter = fmt.Sprintf("%s", value)
			}
			if key == "Device" {
				p.Device = fmt.Sprintf("%s", value)
			}
			field.Set(reflect.ValueOf(value))
		}
		if key == "DEBUG" {
			p.DEBUG = true
		}
	}
}

func (p *Plugin) startCap() {
	var handle *pcap.Handle
	var err error
	if p.DEBUG {
		handle, err = pcap.OpenOffline(p.Device)
	} else {
		handle, err = pcap.OpenLive(p.Device, 65535, false, pcap.BlockForever)
	}
	if err != nil {
		p.status.Code = 501
		p.status.Text = err.Error()
		p.ptr.React(abstract.ERROR(p.status))
		return
	}
	defer handle.Close()
	err = handle.SetBPFFilter(p.BPFFilter)
	if err != nil {
		p.status.Code = 501
		p.status.Text = fmt.Sprintf("can't parse filter expression on [BPFFilter:%s]", p.BPFFilter)
		p.ptr.React(abstract.ERROR(p.status))
		return
	}
	sources := gopacket.NewPacketSource(handle, handle.LinkType())
	p.Infof("%s is capture on %s", p.Name, p.BPFFilter)
	for {
		select {
		case packet := <-sources.Packets():
			if packet == nil || packet.NetworkLayer() == nil ||
				packet.TransportLayer() == nil ||
				packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				fmt.Println("ERR : Unknown Packet -_-")
				return
			}
			p.Broken(packet)
			p.Intact(packet)
		}
	}
}

func InstallPlugin(ptr abstract.Plugin) *Plugin {
	t := reflect.TypeOf(ptr).Elem()
	name := strings.TrimSuffix(t.Name(), "Config")
	plugin := &Plugin{
		Name:  name,
		ptr:   ptr,
		DEBUG: false,
	}
	//读取版本信息
	_, pluginFilePath, _, _ := runtime.Caller(1)
	plugin.Version = pluginFilePath
	plugin.PluginLog = pluginlog.NewPluginLog(name)
	log.Infof("plugin %s:%s installed", plugin.Name, plugin.Version)
	Plugins[name] = plugin
	plugin.status.Code = 200
	plugin.status.Text = "install success"
	return plugin
}
