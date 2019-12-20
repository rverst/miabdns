package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type config struct {
	M []miab `json:"boxes"`
	U []user `json:"users"`
}

func NewJsonConfig(path string) (*config, error) {

	if path == "test" {
		c := config{
			M: []miab{
				{
					N: "test",
					A: "https://box.example.org",
					U: "testUser",
					P: "SuperSecret",
				},
			},
			U: []user{
				{
					M:   "test",
					U:   "user1",
					P:   "pwd1",
					Ip4: true,
					Ip6: false,
					D: []string{
						`.*\.?user1.example.org`,
					},
				},
			},
		}

		return &c, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var c config
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *config) MiabByName(name string) MiabConfig {

	for _, m := range c.M {
		if m.Name() == name {
			return &m
		}
	}
	return nil
}

func (c *config) UserByName(userName string) UserConfig {

	for _, u := range c.U {
		if u.UserName() == userName {
			return &u
		}
	}
	return nil
}

func (c *config) String() string {

	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("error marshalling config: %v", err)
	}
	return string(b)
}

type miab struct {
	N string `json:"name"`
	A string `json:"server"`
	U string `json:"username"`
	P string `json:"password"`
}

func (m *miab) Name() string {
	return m.N
}

func (m *miab) Address() string {
	return m.A
}

func (m *miab) UserName() string {
	return m.U
}

func (m *miab) Password() string {
	return m.P
}

type user struct {
	M   string   `json:"box"`
	U   string   `json:"username"`
	P   string   `json:"password"`
	Ip4 bool     `json:"allow_ip4"`
	Ip6 bool     `json:"allow_ip6"`
	D   []string `json:"allowed_domains"`
}

func (u *user) Miab() string {
	return u.M
}

func (u *user) UserName() string {
	return u.U
}

func (u *user) Password() string {
	return u.P
}

func (u *user) AllowedDomains() []string {
	return u.D
}

func (u *user) AllowIp4() bool {
	return u.Ip4
}
func (u *user) AllowIp6() bool {
	return u.Ip6
}
