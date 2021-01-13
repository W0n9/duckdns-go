package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"reflect"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

// Config is the exporter CLI configuration.
type ClientConfig struct {
	Token       string        `config:"duckdns_token"`
	DomainNames []string      `config:"duckdns_domains"`
	Record      string        `config:"record"`
	IPv4        string        `config:"ipv4"`
	IPv6        string        `config:"ipv6"`
	Interval    time.Duration `config:"update_interval"`

	Verbose      bool `config:"verbose"`
	AutoIP       bool `config:"auto-ip"`
	UpdateIP     bool `config:"update-ip"`
	ClearIP      bool `config:"clear-ip"`
	UpdateRecord bool `config:"update-record"`
	GetRecord    bool `config:"get-record"`
	ClearRecord  bool `config:"clear-record"`
}

func getDefaultConfig() *ClientConfig {
	return &ClientConfig{
		Token:        "",
		DomainNames:  []string{},
		Record:       "",
		IPv4:         "",
		IPv6:         "",
		Interval:     60 * time.Minute,
		Verbose:      false,
		AutoIP:       false,
		UpdateIP:     true,
		ClearIP:      false,
		UpdateRecord: false,
		GetRecord:    false,
		ClearRecord:  false,
	}
}

// Load method loads the configuration by using both flag or environment variables.
func Load() *ClientConfig {
	loaders := []backend.Backend{
		env.NewBackend(),
		flags.NewBackend(),
	}

	loader := confita.NewLoader(loaders...)

	cfg := getDefaultConfig()
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		klog.Fatal("Could not load the configuration...")
	}

	if cfg.AutoIP {
		cfg.getPublicIPv4()
		cfg.getPublicIPv6()
	}

	if cfg.Interval < 10*time.Minute {
		klog.Infof("A time interval below 10 mins is not recommanded. Setting it to 10 mins.")
		cfg.Interval = 10 * time.Minute
	}

	cfg.show()

	return cfg
}

func (c *ClientConfig) show() {
	val := reflect.ValueOf(c).Elem()
	klog.Info("---------------------------------------")
	klog.Info("- DuckDNS client configuration -")
	klog.Info("---------------------------------------")
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		if valueField.Interface() != false {
			klog.Info(fmt.Sprintf("%s : %v", typeField.Name, valueField.Interface()))
		}
	}
	klog.Info("---------------------------------------")
}

func (c *ClientConfig) getPublicIPv4() {
	url := "http://ipv4bot.whatismyipaddress.com"
	resp, err := http.Get(url)
	if err != nil {
		klog.Error(err)
		return
	}
	defer resp.Body.Close()
	ipv4, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Infof("Got IPv4 %v", string(ipv4))
	c.IPv4 = string(ipv4)
}

func (c *ClientConfig) getPublicIPv6() {
	url := "http://ipv6bot.whatismyipaddress.com"
	resp, err := http.Get(url)
	if err != nil {
		klog.Error(err)
		return
	}
	defer resp.Body.Close()
	ipv6, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Fatal(err)
	}

	klog.Infof("Got IPv6 %v", string(ipv6))
	c.IPv6 = string(ipv6)
}
