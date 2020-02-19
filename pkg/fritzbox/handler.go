package fritzbox

import (
	"fmt"
	"github.com/rverst/miabdns/pkg/config"
	"github.com/rverst/miabdns/pkg/dnsupdate"
	"net/http"
)

type Handler struct {
	config config.AppConfig
}

func New(config config.AppConfig) *Handler {
	h := Handler{
		config: config,
	}
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	fmt.Printf("fritzbox - new request from %s\n", ip)

	mc, uc, error := dnsupdate.CheckRequest(h.config, r, []string{http.MethodGet})
	if error != nil {
		fmt.Println(fmt.Errorf("fritzbox - error in config.CheckRequest: %s (%d)", error.Err.Error(), error.Status))
		w.WriteHeader(error.Status)
		if error.ToClient {
			if _, e := w.Write([]byte(error.Err.Error())); e != nil {
				fmt.Println(e)
			}
		}
		return
	}

	s, err := dnsupdate.NewSetting(r.URL.Query().Get("ip4"), r.URL.Query().Get("ip6"), r.URL.Query().Get("domain"), (*uc).AllowedDomains())
	if err != nil {
		fmt.Println(fmt.Errorf("fritzbox - error creating settings: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ok4, ok6, err := dnsupdate.DoUpdate(*mc, *uc, s)

	if err != nil {
		fmt.Println(fmt.Errorf("fritzbox - error in dnsupdate.DoUpdate: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok4 && !ok6 {
		fmt.Printf("fritzbox - nothing updated for domains: %+v\n", s.Domains())
		w.WriteHeader(http.StatusNotModified)
		return
	}
}
