package zone

import (
	"testing"
)

func TestZone(t *testing.T) {
	z, err := NewZone(Config{
		EnableIpv4: true,
		EnableIpv6: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	var rs = Records{}

	// Publish an A record
	newRecord(t, rs, "test-zone.local. 120 IN A 192.168.37.137")
	newRecord(t, rs, "137.37.168.192.in-addr.arpa. 60 IN PTR test-zone.local.")

	// Publish a PTR record for the _ssh._tcp DNS-SD type
	newRecord(t, rs, "_ssh._tcp.local. 60 IN PTR test-zone._ssh._tcp.local.")

	// Publish a SRV record tying the _ssh._tcp record to an A record and a port.
	newRecord(t, rs, "test-zone._ssh._tcp.local. 60 IN SRV 0 0 22 test-zone.local.")

	// Most mDNS browsing tools expect a TXT record for the service even if there
	// are not records defined by RFC 2782.
	newRecord(t, rs, `test-zone._ssh._tcp.local. 60 IN TXT ""`)

	// Bind this service into the list of registered services for dns-sd.
	newRecord(t, rs, "_services._dns-sd._udp.local. 60 IN PTR _ssh._tcp.local.")

	err = z.BroadcastEntries(rs)
	if err != nil {
		t.Fatal(err)
	}
}

func newRecord(t *testing.T, rs Records, s string) {
	err := rs.New(s)
	if err != nil {
		t.Fatal(err)
	}
}
