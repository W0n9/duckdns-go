package main

import (
	"context"
	"k8s.io/klog/v2"
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
	client *duckdns.Client
)

func main() {
	c := config.Load()
	clientConfig := &duckdns.Config{}
	clientConfig.Token = c.Token
	clientConfig.DomainNames = c.DomainNames
	clientConfig.Verbose = c.Verbose
	client = duckdns.NewClient(http.DefaultClient, clientConfig)

	if c.UpdateIP {
		UpdateIP(c.IPv4, c.IPv6)
		for range time.Tick(c.Interval) {
			UpdateIP(c.IPv4, c.IPv6)
		}
	} else if c.ClearIP {
		ClearIP()
	} else if c.UpdateRecord {
		UpdateRecord(c.Record)
	} else if c.GetRecord {
		GetRecord()
	} else if c.ClearRecord {
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

	klog.Infof("Got response %v", body)
	klog.Infof("IP has been updated at %v", time.Now())
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
