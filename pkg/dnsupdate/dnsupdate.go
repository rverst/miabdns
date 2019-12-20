package dnsupdate

import (
	"fmt"
	"github.com/rverst/go-miab/miab"
	"github.com/rverst/miabdns/pkg/config"
)

type UpdateResult int

func DoUpdate(miabConfig config.MiabConfig, userConfig config.UserConfig, settings config.UpdateSettings) (ip4Ok, ip6Ok bool, err error) {

	mc, err := miab.NewConfig(miabConfig.UserName(), miabConfig.Password(), miabConfig.Address())
	if err != nil {
		return false, false, err
	}

	var ok4 bool
	var ok6 bool

	if userConfig.AllowIp4() && len(settings.IP4()) > 0 && len(settings.Domains()) > 0 {

		for _, d := range settings.Domains() {
			if ok, err := miab.UpdateDns4(mc, d, settings.IP4().String()); err != nil {
				fmt.Printf("dnsupdate - error updating A record to value '%s' for: %s\n            err: %s\n", settings.IP4(), d, err)
			} else if ok {
				ok4 = true
				fmt.Printf("dnsupdate - IPv4 updated, set A record to value '%s' for: %s\n", settings.IP4(), d)
			}
		}
	}

	if userConfig.AllowIp6() && len(settings.IP6()) > 0 && len(settings.Domains()) > 0 {
		for _, d := range settings.Domains() {
			if ok, err := miab.UpdateDns6(mc, d, settings.IP6().String()); err != nil {
				fmt.Printf("dnsupdate - error updating AAAA record to value '%s' for: %s\n            err: %s\n", settings.IP6(), d, err)
			} else if ok {
				ok6 = true
				fmt.Printf("dnsupdate - IPv6 updated, set AAAA record to value '%s' for: %s\n", settings.IP6(), d)
			}
		}
	}

	return ok4, ok6, nil
}
