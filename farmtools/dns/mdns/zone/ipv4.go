package zone

import (
	"fmt"
	"net"

	"github.com/miekg/dns"
	"golang.org/x/net/ipv4"
)

var (
	_ Ipv = (*Ipv4)(nil)

	// Multicast groups used by mDNS
	mdnsGroupIPv4 = net.IPv4(224, 0, 0, 251)

	// mDNS wildcard addresses
	mdnsWildcardAddrIPv4 = &net.UDPAddr{
		IP:   net.ParseIP("224.0.0.0"),
		Port: 5353,
	}

	// mDNS endpoint addresses
	ipv4Addr = &net.UDPAddr{
		IP:   mdnsGroupIPv4,
		Port: 5353,
	}
)

type Ipv4 struct {
	net4 *ipv4.PacketConn
	ifs  []net.Interface
}

func NewIpv4(ifs []net.Interface) (*Ipv4, error) {
	i := &Ipv4{}

	err := i.joinMulticast(ifs)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Ipv4) multicastResponse(msg *dns.Msg) error {
	buf, err := msg.Pack()
	if err != nil {
		return err
	}

	if i.net4 != nil {
		var wcm ipv4.ControlMessage
		var ifs []net.Interface

		for _, _i := range i.ifs {
			wcm.IfIndex = _i.Index
			_, err = i.net4.WriteTo(buf, &wcm, ipv4Addr)
			if err != nil {
				ifs = append(ifs, _i)
			}
		}

		if len(ifs) <= 0 {
			return fmt.Errorf("no interfaces found")
		}

		i.ifs = ifs
	}

	return nil
}

func (i *Ipv4) joinMulticast(ifs []net.Interface) error {
	var err error
	i.net4, err = i.joinUDPMulticast(ifs)
	if err != nil {
		return err
	}

	return nil
}

func (i *Ipv4) joinUDPMulticast(ifs []net.Interface) (*ipv4.PacketConn, error) {
	if len(ifs) == 0 {
		return nil, fmt.Errorf("no interfaces found")
	}

	udpConn, err := net.ListenUDP("udp4", mdnsWildcardAddrIPv4)
	if err != nil {
		// log.Printf("[ERR] bonjour: Failed to bind to udp4 multicast: %v", err)
		return nil, err
	}

	// Join multicast groups to receive announcements
	pkConn := ipv4.NewPacketConn(udpConn)
	err = pkConn.SetControlMessage(ipv4.FlagInterface, true)
	if err != nil {
		return nil, err
	}

	var failedJoins int
	for _, _i := range ifs {
		if err := pkConn.JoinGroup(&_i, &net.UDPAddr{IP: mdnsGroupIPv4}); err == nil {
			i.ifs = append(i.ifs, _i)
		} else {
			// log.Println("Udp4 JoinGroup failed for interface ", _i)
			failedJoins++
		}
	}

	if failedJoins == len(ifs) {
		_ = pkConn.Close()
		return nil, fmt.Errorf("udp4: failed to join any of these interfaces: %v", ifs)
	}

	return pkConn, nil
}
