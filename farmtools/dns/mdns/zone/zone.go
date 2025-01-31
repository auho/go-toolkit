package zone

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

type Zone struct {
	records map[dns.RR]struct{}
	ifs     []net.Interface
	ipvs    []Ipv
}

func NewZone(c Config) (*Zone, error) {
	z := &Zone{}
	err := z.multicastInterfaces()
	if err != nil {
		return nil, err
	}

	if c.EnableIpv4 {
		ipv4, err := NewIpv4(z.ifs)
		if err != nil {
			return nil, fmt.Errorf("ipv4: %v", err)
		}

		z.ipvs = append(z.ipvs, ipv4)
	}

	return z, nil
}

func (z *Zone) BroadcastEntries(rs Records) error {
	if len(rs) <= 0 {
		return fmt.Errorf("no record found")
	}

	entries := rs.ToSlice()

	msg := &dns.Msg{}
	msg.MsgHdr.Response = true
	msg.Answer = entries

	for _, ipv := range z.ipvs {
		err := ipv.multicastResponse(msg)
		if err != nil {
			z.tryJoinMulticast(ipv)
		}
	}

	return nil
}

func (z *Zone) tryJoinMulticast(ipv Ipv) {
	var retry = time.Second
	for {
		err := z.multicastInterfaces()
		if err == nil {
			err = ipv.joinMulticast(z.ifs)
			if err == nil {
				return
			}
		}

		<-time.NewTimer(retry).C
		if retry < time.Second*20 {
			retry *= 2
		}
	}
}

func (z *Zone) multicastInterfaces() error {
	z.ifs = z.listMulticastInterfaces()
	if len(z.ifs) <= 0 {
		return fmt.Errorf("no interfaces found")
	}

	return nil
}

func (z *Zone) listMulticastInterfaces() []net.Interface {
	var out []net.Interface
	ifs, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, ifi := range ifs {
		if (ifi.Flags & net.FlagUp) == 0 {
			continue
		}
		if (ifi.Flags & net.FlagMulticast) > 0 {
			out = append(out, ifi)
		}
	}

	return out
}
