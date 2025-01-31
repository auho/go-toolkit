package zone

import (
	"net"

	"github.com/miekg/dns"
)

type Ipv interface {
	joinMulticast(ifs []net.Interface) error
	multicastResponse(msg *dns.Msg) error
}
