# duckdns-go, a duckdns client in golang

A golang client to update, clear ip and records for [DuckDNS](https://www.duckdns.org/) domains.  
**Mod for get device IPv4 and IPv6, even IPv4 only.**

## Prerequisites

* [Go](https://golang.org/doc/)

## Installation

### From sources

You can download and build it from the sources. You have to retrieve the project sources by using one of the following way:
```bash
$ go get -u github.com/W0n9/duckdns-go
# or
$ git clone https://github.com/W0n9/duckdns-go.git
```

Then, build the binary:
```Powershell
go build -ldflags '-s -w' -o duckdns-go.exe .
```

## Client Usage

```Powershell
.\duckdns-go.exe -duckdns_token <token> -duckdns_domains <domain> -update-ip 
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
         Get device ipv4 and ipv6
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
  -ipv4-only
        Get device ipv4
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