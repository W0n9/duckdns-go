package main

import (
	"context"
	"k8s.io/klog/v2"
	"net/http"
	"time"

	"github.com/ebrianne/duckdns-go/duckdns"
)

func main() {
	config := &duckdns.Config{}
	config.Token = ""
	config.DomainNames = []string{""}
	client := duckdns.NewClient(http.DefaultClient, config)

	resp, err := client.UpdateIP(context.Background())
	if err != nil {
		klog.Fatal("UpdateIP() returned error: ", err)
	}

	klog.Info("Got response ", resp.Data)
	klog.Info("IP has been updated at ", time.Now())
}
