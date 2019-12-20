package main

import (
	"fmt"
	"github.com/integrii/flaggy"
	"github.com/rverst/miabdns/pkg/config"
	"github.com/rverst/miabdns/pkg/fritzbox"
	"github.com/rverst/miabdns/pkg/health"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	Version    = "0.0.0"
	CommitHash = ""
	BuildDate  = time.Now().UTC().Format("2006-01-02 15:04:05")

	flPort   uint16 = 8080
	flBind          = "0.0.0.0"
	flConfig        = "config.json"
)

func init() {
	flaggy.SetName("miabdns")
	flaggy.SetDescription("miabdns - super simple dynamic DNS service for Mail-in-a-Box")
	flaggy.SetVersion(fmt.Sprintf("%s - %s (%s)", Version, CommitHash, BuildDate))

	flaggy.String(&flConfig, "c", "config", "path to the configuration file (default \"config.json\")")
	flaggy.String(&flBind, "b", "bind", "interface to which the server will bind (default \"0.0.0.0\")")
	flaggy.UInt16(&flPort, "p", "port", " port on which the server will listen (default 8080)")
	flaggy.Parse()

	if os.Getenv("MIABDNS_PORT") != "" {
		if up, err := strconv.ParseUint(os.Getenv("MIABDNS_PORT"), 10, 16); err != nil {
			flPort = uint16(up)
		}
	}

	if os.Getenv("MIABDNS_BIND") != "" {
		flConfig = os.Getenv("MIABDNS_BIND")
	}

	if os.Getenv("MIABDNS_CONFIG") != "" {
		flConfig = os.Getenv("MIABDNS_CONFIG")
	}
}

func main() {

	cfg, err := config.NewJsonConfig(flConfig)
	if err != nil {
		fmt.Printf("unable to load cfg: %v", err)
		os.Exit(1)
	}

	fmt.Printf("cfg: %+v\n\n", cfg)

	mux := http.NewServeMux()
	mux.Handle("/health/", health.New())
	mux.Handle("/dnsupdate/fritz/", fritzbox.New(cfg))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("unhandled route:\t\t %s\n", r.URL)
		http.NotFound(w, r)
		return
	})

	address := fmt.Sprintf("%s:%d", flBind, flPort)
	l, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("error creating listener: %v", err)
		os.Exit(1)
	}
	fmt.Printf("listening on: %s\n", address)

	defer func() {
		err := l.Close()
		if err != nil {
			fmt.Printf("error closing listener: %v", err)
			os.Exit(1)
		}
	}()

	for {
		err := http.Serve(l, mux)
		if err != nil {
			fmt.Printf("error starting server: %v", err)
			os.Exit(1)
		}
	}
}
