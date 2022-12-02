package engine

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sqleyes/engine/abstract"
)

func (p *Plugin) Intact(packet gopacket.Packet) {
	p.ptr.React(packet)
}
func (p *Plugin) Broken(packet gopacket.Packet) {
	broken := abstract.Broken{}
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		broken.SrcIP = ip.SrcIP.String()
		broken.DstIP = ip.DstIP.String()
	}
	tcp := packet.TransportLayer().(*layers.TCP)
	broken.SrcPort = int(tcp.SrcPort)
	broken.DstPort = int(tcp.DstPort)
	if !tcp.SYN && !tcp.FIN && !tcp.RST && len(tcp.LayerPayload()) == 0 {
		return
	}
	if len(tcp.Payload) == 0 {
		return
	}
	broken.Payload = tcp.Payload
	p.ptr.React(broken)
}
