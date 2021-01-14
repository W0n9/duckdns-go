package main

import (
	"context"
	"k8s.io/klog"
	"net/http"
	"strings"
	"time"

	"github.com/ebrianne/duckdns-go/config"
	"github.com/ebrianne/duckdns-go/duckdns"
)

const (
	name = "duckdns-client"
)

var (
	c      *config.ClientConfig
	client *duckdns.Client
)

func main() {
	c = config.Load()
	config := &duckdns.Config{}
	config.Token = c.Token
	config.DomainNames = c.DomainNames
	config.Verbose = c.Verbose
	client = duckdns.NewClient(http.DefaultClient, config)

	if c.UpdateIP {
		UpdateIP(c.IPv4, c.IPv6)
		for range time.Tick(c.Interval) {
			UpdateIP(c.IPv4, c.IPv6)
		}
	} else if c.ClearIP {
		ClearIP()
	} else if c.UpdateRecord {
		if c.Record == "" {
			klog.Error("Provided TXT record empty... It needs to be provided with -record string to update the txt record")
			return
		}
		UpdateRecord(c.Record)
	} else if c.GetRecord {
		GetRecord()
	} else if c.ClearRecord {
		if c.Record == "" {
			klog.Error("Provided TXT record empty... It needs to be provided with -record string to clear the txt record")
			return
		}
		ClearRecord(c.Record)
	} else {
		klog.Error("CLI option provided unknown...")
	}
}

func UpdateIP(ipv4, ipv6 string) {
	var body string

	if ipv4 == "" && ipv6 == "" {
		resp, err := client.UpdateIP(context.Background())
		if err != nil {
			klog.Fatal("UpdateIP() returned error: ", err)
		}
		body = SplitAndJoin(resp.Data)
	} else {
		resp, err := client.UpdateIPWithValues(context.Background(), ipv4, ipv6)
		if err != nil {
			klog.Fatal("UpdateIPWithValues() returned error: ", err)
		}
		body = SplitAndJoin(resp.Data)
	}

	if strings.Contains(body, "KO") {
		klog.Errorf("Got response containing KO, verify the provided arguments, will try again in %v", c.Interval)
	} else {
		klog.Infof("Got response %v", body)
		klog.Infof("IP has been updated at %v", time.Now())
	}
}

func ClearIP() {
	resp, err := client.ClearIP(context.Background())
	if err != nil {
		klog.Fatal("ClearIP() returned error: ", err)
	}
	klog.Infof("Got response %v", resp.Data)
	klog.Infof("IP has been cleared at %v", time.Now())
}

func UpdateRecord(record string) {
	resp, err := client.UpdateRecord(context.Background(), record)
	if err != nil {
		klog.Fatal("UpdateRecord() returned error: ", err)
	}
	klog.Infof("Got response %v", resp.Data)
	klog.Infof("TXT Record has been update with %v at %v", record, time.Now())
}

func GetRecord() {
	record, err := client.GetRecord()
	if err != nil {
		klog.Fatal("GetRecord() returned error: ", err)
	}
	klog.Infof("TXT Record is %q", record)
}

func ClearRecord(record string) {
	resp, err := client.ClearRecord(context.Background(), record)
	if err != nil {
		klog.Fatal("ClearRecord() returned error: ", err)
	}
	klog.Infof("Got response %v", resp.Data)
	klog.Infof("TXT Record has been cleared at %v", time.Now())
}

func SplitAndJoin(data string) string {
	s := strings.Split(data, "\n")
	body := strings.Join(s, ", ")
	return body
}
