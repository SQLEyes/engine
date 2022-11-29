package engine

import (
	"engine/abstract"
	"engine/pluginlog"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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
	}
}

func (p *Plugin) startCap() {
	handle, err := pcap.OpenLive(p.Device, 65535, false, pcap.BlockForever)
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
		Name: name,
		ptr:  ptr,
	}
	//v := reflect.ValueOf(ptr).Elem()
	//bpf := v.FieldByName("BPFFilter")
	//device := v.FieldByName("Device")
	//读取版本信息
	_, pluginFilePath, _, _ := runtime.Caller(1)
	plugin.Version = pluginFilePath
	plugin.PluginLog = pluginlog.NewPluginLog(name)
	log.Infof("plugin %s:%s installed", plugin.Name, plugin.Version)
	Plugins[name] = plugin
	plugin.status.Code = 200
	plugin.status.Text = "install success"
	//plugin.status = CheckParameter(bpf, "BPFFilter")
	//if plugin.status.Code != 200 {
	//	return plugin
	//}
	//plugin.BPFFilter = bpf.String()
	//plugin.status = CheckParameter(device, "Device")
	//if plugin.status.Code == 200 {
	//	plugin.Device = device.String()
	//}
	return plugin
}
