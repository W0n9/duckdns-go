package config

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strings"
	"time"

	"k8s.io/klog/v2"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/flags"
)

// Config is the exporter CLI configuration.
type ClientConfig struct {
	Token       string        `config:"duckdns_token,description=DuckDNS Token (mandatory)"`
	DomainNames []string      `config:"duckdns_domains,description=List of duckdns domains to update, needs to be comma separated (mandatory)"`
	Record      string        `config:"record,description=TXT record (mandatory with -update-record/-clear-record flags)"`
	IPv4        string        `config:"ipv4,description=IPv4 address (optional)"`
	IPv6        string        `config:"ipv6,description=IPv6 address (optional)"`
	Interval    time.Duration `config:"update_interval,description=Interval between IP updates (min 10 mins)"`

	Verbose      bool `config:"verbose,description=Verbose flag for duckdns response"`
	AutoIP       bool `config:"auto-ip,description=Get device ipv4 and ipv6"`
	IPv4Only     bool `config:"ipv4-only,description=Get device ipv4"`
	UpdateIP     bool `config:"update-ip,description=Update IP routine"`
	ClearIP      bool `config:"clear-ip,description=Clear ip in duckdns with clear=true`
	UpdateRecord bool `config:"update-record,description=Update TXT record routine"`
	GetRecord    bool `config:"get-record,description=Get txt record"`
	ClearRecord  bool `config:"clear-record,description=Clear txt record in duckdns with clear=true"`
}

func getDefaultConfig() *ClientConfig {
	return &ClientConfig{
		Token:        "",
		DomainNames:  nil,
		Record:       "",
		IPv4:         "",
		IPv6:         "",
		Interval:     60 * time.Minute,
		Verbose:      false,
		AutoIP:       false,
		UpdateIP:     false,
		ClearIP:      false,
		UpdateRecord: false,
		GetRecord:    false,
		ClearRecord:  false,
	}
}

// Load method loads the configuration by using both flag or environment variables.
func Load() *ClientConfig {
	cfg := getDefaultConfig()

	loader := confita.NewLoader(env.NewBackend(), flags.NewBackend())
	err := loader.Load(context.Background(), cfg)
	if err != nil {
		klog.Fatal("Could not load the configuration...")
	}

	if cfg.AutoIP {
		// cfg.getPublicIPv4()
		// cfg.getPublicIPv6()
		cfg.getDeviceIPv4()
		cfg.getDeviceIPv6()
	}

	if cfg.IPv4Only {
		cfg.getDeviceIPv4()
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

// func (c *ClientConfig) getPublicIPv4() {
// 	url := "http://ipv4bot.whatismyipaddress.com"
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		klog.Error(err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	ipv4, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		klog.Fatal(err)
// 	}

// 	klog.Infof("Got IPv4 %v", string(ipv4))
// 	c.IPv4 = string(ipv4)
// }

// func (c *ClientConfig) getPublicIPv6() {
// 	url := "http://ipv6bot.whatismyipaddress.com"
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		klog.Error(err)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	ipv6, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		klog.Fatal(err)
// 	}

// 	klog.Infof("Got IPv6 %v", string(ipv6))
// 	c.IPv6 = string(ipv6)
// }

func (c *ClientConfig) getDeviceIPv4() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		klog.Error(err)
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				klog.Infof("Got IPv4 %v", ipnet.IP.String())
				c.IPv4 = ipnet.IP.String()
			}
		}
	}
}

func (c *ClientConfig) getDeviceIPv6() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		klog.Error(err)
		return
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To16() != nil {
				i := regexp.MustCompile(`(\w+:){7}\w+`).FindString(ipnet.IP.String())
				if strings.Count(i, ":") == 7 {
					klog.Infof("Got IPv6 %v", i)
					c.IPv6 = i
				}
			}
		}
	}
}
