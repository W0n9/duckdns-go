# duckdns-go, a duckdns client in golang

![Build/Push (master)](https://github.com/ebrianne/duckdns-go/workflows/Build/Push%20(master)/badge.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/ebrianne/duckdns-go?status.png)](https://godoc.org/github.com/ebrianne/duckdns-go)
[![GoReportCard](https://goreportcard.com/badge/github.com/ebrianne/duckdns-go)](https://goreportcard.com/report/github.com/ebrianne/duckdns-go)
[![Known Vulnerabilities](https://snyk.io/test/github/ebrianne/duckdns-go/badge.svg?targetFile=Dockerfile)](https://snyk.io/test/github/ebrianne/duckdns-go?targetFile=Dockerfile)

A golang client to update, clear ip and records for [DuckDNS](https://www.duckdns.org/) domains.

## Prerequisites

* [Go](https://golang.org/doc/)

## Installation

### Download binary

You can download the latest version of the binary built for your architecture here:

* Architecture **i386** [
[Darwin](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-darwin-386) /
[FreeBSD](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-freebsd-386) /
[Linux](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-linux-386) /
[Windows](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-windows-386.exe)
]
* Architecture **amd64** [
[Darwin](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-darwin-amd64) /
[FreeBSD](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-freebsd-amd64) /
[Linux](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-linux-amd64) /
[Windows](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-windows-amd64.exe)
]
* Architecture **arm** [
[Linux](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-linux-arm)
]
* Architecture **arm64** [
[Linux](https://github.com/ebrianne/duckdns-go/releases/latest/download/duckdns-go-linux-arm64)
]

### From sources

You can download and build it from the sources. You have to retrieve the project sources by using one of the following way:
```bash
$ go get -u github.com/ebrianne/duckdns-go
# or
$ git clone https://github.com/ebrianne/duckdns-go.git
```

Install the needed vendors:

```
$ GO111MODULE=on go mod vendor
```

Then, build the binary (here, an example to run on Raspberry PI ARM architecture):
```bash
$ GOOS=linux GOARCH=arm GOARM=7 go build -o duckdns-go .
```
## Using Docker

The client has been made available as a docker image. You can simply run it by the following command and pass the configuration with environment variables. 
By default it executes `./duckdns-go -update-ip`.

```bash
docker run \
-e 'DUCKDNS_TOKEN=<token>' \
-e 'DUCKDNS_DOMAINS=<domains>' \
ebrianne/duckdns-go
```

You can also provide the command to run

```bash
docker run \
-e 'DUCKDNS_TOKEN=<token>' \
-e 'DUCKDNS_DOMAINS=<domains>' \
ebrianne/duckdns-go ./duckdns-go [ARG]
```

## Client Usage

```bash
$ ./duckdns-go -duckdns_token <token> -duckdns_domains <domain> -update-ip 
```

```bash
I0113 11:17:15.063439  426646 configuration.go:86] ---------------------------------------
I0113 11:17:15.063895  426646 configuration.go:87] - DuckDNS client configuration -
I0113 11:17:15.064026  426646 configuration.go:88] ---------------------------------------
I0113 11:17:15.064115  426646 configuration.go:94] Token : **************
I0113 11:17:15.064135  426646 configuration.go:94] DomainNames : [******]
I0113 11:17:15.064146  426646 configuration.go:94] Record : 
I0113 11:17:15.064151  426646 configuration.go:94] IPv4 : 
I0113 11:17:15.064166  426646 configuration.go:94] IPv6 : 
I0113 11:17:15.064177  426646 configuration.go:94] Interval : 1h0m0s
I0113 11:17:15.064187  426646 configuration.go:94] UpdateIP : true
I0113 11:17:15.064220  426646 configuration.go:97] ---------------------------------------
I0113 11:17:15.064242  426646 client.go:96] Sending request to https://www.duckdns.org/update?domains=******&token=**************&ip=
I0113 11:17:15.940591  426646 main.go:71] Got response OK
I0113 11:17:15.940629  426646 main.go:72] IP has been updated at 2021-01-13 11:17:15.940624102 +0100 CET m=+0.877805589
```
## Available CLI options

```bash
Usage of ./duckdns-go:
  -auto-ip
        Get public ipv4 and ipv6 via whatismyipaddress.com
  -clear-record
        Clear txt record in duckdns with clear=true
  -duckdns_domains value
        List of duckdns domains to update (default duckdns_domains)
  -duckdns_token string
        DuckDNS Token (mandatory)
  -get-record
        Get txt record
  -ipv4 string
        IPv4 address (optional)
  -ipv6 string
        IPv6 address (optional)
  -record string
        TXT record (mandatory with -update-record/-clear-record flags)
  -update-ip
        Update IP routine
  -update-record
        Update TXT record routine
  -update_interval duration
        Interval between IP updates (min 10 mins) (default 1h0m0s)
  -verbose
        Verbose flag for duckdns response
  ```

### Environment Variables

All CLI commands can be specified as an environment variable such as:

```bash
export DUCKDNS_TOKEN="<your token>"
export DUCKDNS_DOMAINS="domain1,domain2" #use space comma separated names
duckdns
```