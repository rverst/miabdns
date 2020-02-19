package dnsupdate

import (
	"errors"
	"fmt"
	"github.com/rverst/miabdns/pkg/config"
	"github.com/rverst/miabdns/pkg/models"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type setting struct {
	ip4     net.IP
	ip6     net.IP
	domains []string
}

// Creates a new setting model based on the query parameters of the request
func NewSetting(ip4, ip6, domains string, allowedDomains []string) (*setting, error) {

	v6 := net.IP{}
	v4 := net.IP{}

	if ip4 != "" {
		if err := v4.UnmarshalText([]byte(ip4)); err != nil {
			return nil, err
		}
	}

	if ip6 != "" {
		if err := v6.UnmarshalText([]byte(ip6)); err != nil {
			return nil, err
		}
	}

	split := strings.Split(domains, ",")
	d := make([]string, 0)
	for _, s := range split {
		if s == "" {
			continue
		}

		if len(allowedDomains) == 0 {
			d = append(d, s)
			continue
		}

		for _, a := range allowedDomains {

			if reg, err := regexp.Compile(a); err == nil {
				if reg.MatchString(s) {
					d = append(d, s)
					continue
				}
			}
		}
	}

	return &setting{
		ip4:     v4,
		ip6:     v6,
		domains: d,
	}, nil
}

func (s *setting) IP4() net.IP {
	return s.ip4
}

func (s *setting) IP6() net.IP {
	return s.ip6
}

func (s *setting) Domains() []string {
	return s.domains
}

// Checks the request for username/password and a configuration
func CheckRequest(c config.AppConfig, r *http.Request, allowedMethods []string) (*config.MiabConfig, *config.UserConfig, *models.Error) {

	if allowedMethods != nil {
		allowed := false
		for _, am := range allowedMethods {
			if am == r.Method {
				allowed = true
			}
		}
		if !allowed {
			return nil, nil, models.NewError(http.StatusMethodNotAllowed, fmt.Errorf("method not allowd: %s", r.Method), false)
		}
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		return nil, nil, models.NewError(http.StatusUnauthorized, errors.New("no basic auth header"), true)
	}

	user := c.UserByName(username)
	if user == nil {
		return nil, nil, models.NewError(http.StatusForbidden, fmt.Errorf("no user with name '%s' found", username), false)
	}

	if user.Password() != password {
		return nil, nil, models.NewError(http.StatusForbidden, fmt.Errorf("password for user '%s' did not match", username), false)
	}

	miab := c.MiabByName(user.Miab())
	if miab == nil {
		return nil, nil, models.NewError(http.StatusInternalServerError, fmt.Errorf("no miab config found with name '%s' ", user.Miab()), false)
	}

	return &miab, &user, nil
}
