package config

import (
	"net"
)

type AppConfig interface {
	MiabByName(name string) MiabConfig
	UserByName(userName string) UserConfig
}

type MiabConfig interface {
	Name() string
	Address() string
	UserName() string
	Password() string
}

type UserConfig interface {
	Miab() string
	UserName() string
	Password() string
	AllowedDomains() []string
	AllowIp4() bool
	AllowIp6() bool
}

type UpdateSettings interface {
	IP4() net.IP
	IP6() net.IP
	Domains() []string
}
