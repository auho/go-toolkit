package zone

import "github.com/miekg/dns"

type Records map[dns.RR]struct{}

func NewRecords() Records {
	return make(Records, 16)
}

func (rs Records) New(s string) error {
	rr, err := dns.NewRR(s)
	if err != nil {
		return err
	}

	rs.Add(rr)

	return nil
}

func (rs Records) Clear() {
	rs = NewRecords()
}

func (rs Records) Delete(s string) error {
	rr, err := dns.NewRR(s)
	if err != nil {
		return err
	}

	rs.Remove(rr)

	return nil
}

func (rs Records) Add(in dns.RR) {
	if !rs.exists(in) {
		rs[in] = struct{}{}
	}
}

func (rs Records) Remove(out dns.RR) {
	for rr := range rs {
		if dns.IsDuplicate(out, rr) {
			continue
		}

		delete(rs, rr)
	}
}

func (rs Records) ToSlice() []dns.RR {
	var out []dns.RR
	for rr := range rs {
		out = append(out, rr)
	}

	return out
}

func (rs Records) exists(rr dns.RR) bool {
	for record := range rs {
		if dns.IsDuplicate(record, rr) {
			return true
		}
	}

	return false
}
