package publisher

import (
	"log"
	"net"
)

type interfaceIP struct {
	ipv4Addrs []string
	ipv6Addrs []string
}

func (i *interfaceIP) ipv4Addr() string {
	if len(i.ipv4Addrs) <= 0 {
		return ""
	}

	return i.ipv4Addrs[0]
}

func (i *interfaceIP) ipv6Addr() string {
	if len(i.ipv6Addrs) <= 0 {
		return ""
	}

	return i.ipv6Addrs[0]
}

type interfacesIP map[string]interfaceIP

func newInterfacesIp() (interfacesIP, error) {
	ifsIP := make(interfacesIP, 16)

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, _i := range ifs {
		if _i.Flags&net.FlagLoopback != 0 {
			continue
		}

		iIP, err := ifsIP.getInterfaceIPAddresses(_i)
		if err != nil {
			log.Fatal(err)
		}

		ifsIP[_i.Name] = iIP
	}

	return ifsIP, nil
}

func (ii interfacesIP) getInterfaceIPAddresses(i net.Interface) (interfaceIP, error) {
	var iIP interfaceIP

	addrs, err := i.Addrs()
	if err != nil {
		return iIP, err
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil {
			continue
		}

		if ip.To4() != nil {
			iIP.ipv4Addrs = append(iIP.ipv4Addrs, ip.String())
		} else if ip.To16() != nil {
			iIP.ipv6Addrs = append(iIP.ipv6Addrs, ip.String())
		}
	}

	return iIP, nil
}
