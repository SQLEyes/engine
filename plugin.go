package engine

import (
	"engine/config"
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
	ptr       config.Plugin
	handle    *pcap.Handle
	status    config.Msg
	BPFFilter string
	Device    string
}

func (p *Plugin) startCap() {
	handle, err := pcap.OpenLive(p.Device, 65535, false, pcap.BlockForever)
	if err != nil {
		p.status.Code = 501
		p.status.Text = err.Error()
		p.ptr.React(config.ERROR(p.status))
		return
	}
	defer handle.Close()
	err = handle.SetBPFFilter(p.BPFFilter)
	if err != nil {
		p.status.Code = 501
		p.status.Text = fmt.Sprintf("can't parse filter expression on [BPFFilter:%s]", p.BPFFilter)
		p.ptr.React(config.ERROR(p.status))
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

func InstallPlugin(ptr config.Plugin) *Plugin {
	t := reflect.TypeOf(ptr).Elem()
	name := strings.TrimSuffix(t.Name(), "Config")
	plugin := &Plugin{
		Name: name,
		ptr:  ptr,
	}
	v := reflect.ValueOf(ptr).Elem()
	bpf := v.FieldByName("BPFFilter")
	device := v.FieldByName("Device")
	_, pluginFilePath, _, _ := runtime.Caller(1)
	plugin.Version = pluginFilePath
	plugin.PluginLog = pluginlog.NewPluginLog(name)
	log.Infof("plugin %s:%s installed", plugin.Name, plugin.Version)
	Plugins[name] = plugin
	plugin.status = CheckParameter(bpf, "BPFFilter")
	if plugin.status.Code != 200 {
		return plugin
	}
	plugin.BPFFilter = bpf.String()
	plugin.status = CheckParameter(device, "Device")
	if plugin.status.Code == 200 {
		plugin.Device = device.String()
	}
	return plugin
}
func CheckParameter(v reflect.Value, field string) config.Msg {
	if !v.IsValid() || len(v.String()) == 0 {
		return config.Msg{
			Code: 500,
			Text: "install failure caused by: " + field + " IS NULL",
		}
	} else {
		return config.Msg{
			Code: 200,
			Text: "install success",
		}
	}
}
