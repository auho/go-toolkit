package certificate

import (
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"time"
)

func QuickCheckExpired(domains []string, alarmDays int) {

	expiredDomains := make(map[string]int)

	var nowTime = time.Now()
	for _, domain := range domains {
		conn, _ := tls.Dial("tcp", domain+":443", nil)
		cert := conn.ConnectionState().PeerCertificates[0]

		color.Black(fmt.Sprintf("Subject: %v\n", cert.Subject))
		color.Black(fmt.Sprintf("Expired: %s => %s\n", cert.NotBefore.Format(time.DateTime), cert.NotAfter.Format(time.DateTime)))
		color.Magenta(fmt.Sprintf("Expired Days: %f 天", cert.NotAfter.Sub(nowTime).Seconds()/86400))

		if alarmDays <= 10 {
			alarmDays = 10
		}

		_expiredDays := int(cert.NotAfter.Sub(nowTime).Seconds() / 86400)

		if _expiredDays <= alarmDays {
			expiredDomains[domain] = _expiredDays
		}

		fmt.Println()
	}

	fmt.Println("Approaching Expired:")
	for _d, _i := range expiredDomains {
		color.Red(fmt.Sprintf("%s: %d 天", _d, _i))
	}
}
