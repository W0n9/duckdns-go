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

func main() {
	c := config.Load()
	clientConfig := &duckdns.Config{}
	clientConfig.Token = c.Token
	clientConfig.DomainNames = c.DomainNames
	clientConfig.Verbose = c.Verbose

	if c.UpdateIP {
		for range time.Tick(c.Interval) {
			client := duckdns.NewClient(http.DefaultClient, clientConfig)

			if c.IPv4 == "" && c.IPv4 == "" {
				resp, err := client.UpdateIP(context.Background())
				if err != nil {
					klog.Fatal("UpdateIP() returned error: ", err)
				}
				s := strings.Split(resp.Data, "\n")
				body := strings.Join(s, ", ")
				klog.Infof("Got response %v", body)
			} else {
				resp, err := client.UpdateIPWithValues(context.Background(), c.IPv4, c.IPv6)
				if err != nil {
					klog.Fatal("UpdateIPWithValues() returned error: ", err)
				}
				s := strings.Split(resp.Data, "\n")
				body := strings.Join(s, ", ")
				klog.Infof("Got response %v", body)
			}

			klog.Infof("IP has been cleared at %v", time.Now())
		}
	}

	if c.ClearIP {
		client := duckdns.NewClient(http.DefaultClient, clientConfig)
		resp, err := client.ClearIP(context.Background())
		if err != nil {
			klog.Fatal("ClearIP() returned error: ", err)
		}
		klog.Infof("Got response %v", resp.Data)
		klog.Infof("IP has been cleared at %v", time.Now())
	}

	if c.UpdateRecord {
		client := duckdns.NewClient(http.DefaultClient, clientConfig)
		resp, err := client.UpdateRecord(context.Background(), c.Record)
		if err != nil {
			klog.Fatal("UpdateRecord() returned error: ", err)
		}
		klog.Infof("Got response %v", resp.Data)
		klog.Infof("TXT Record has been update with %v at %v", c.Record, time.Now())
	}

	if c.GetRecord {
		client := duckdns.NewClient(http.DefaultClient, clientConfig)
		record, err := client.GetRecord()
		if err != nil {
			klog.Fatal("GetRecord() returned error: ", err)
		}
		klog.Infof("TXT Record is %q", record)
	}

	if c.ClearRecord {
		client := duckdns.NewClient(http.DefaultClient, clientConfig)
		resp, err := client.ClearRecord(context.Background(), c.Record)
		if err != nil {
			klog.Fatal("ClearRecord() returned error: ", err)
		}
		klog.Infof("Got response %v", resp.Data)
		klog.Infof("TXT Record has been cleared at %v", time.Now())
	}
}
