# Maxanon 
### REST API service for Maxmind Anonymous IP Database

[Maxmind database](https://www.maxmind.com/en/solutions/geoip2-enterprise-product-suite/anonymous-ip-database) consist ip which detected as proxies, VPNs, and other anonymizers.

There are two types of usage this service.
As for single requests or on stream with logstash processor.

## Deploy
Maxanon support databases:
1) redis
2) mongodb

### run

```
go run service/maxanon/main.go -file GeoIP2-Anonymous-IP-Block-IPv4.csv
```

### docker-compose
```sh
docker-compose up -d 
```

## Usage
Request to check information about ip address 192.168.1.1
```sh
curl /api/v1/info/192.168.1.1
```
Response 
```shell script
{"IP":"192.168.1.1","Anonymous":true,"AnonymousVPN":false,"IsHostingProvider":false,"IsPublicProxy":true,"IsTorExitNode":false}
```

Example request http module for logstash processor.
```sh
http {
    body => "%{source.ip}"
    target_body => "ip_flags"
    target_headers => "redis"
    url => "http://10.0.0.1:8000/api/v1/info/%{source.ip}"
    connect_timeout => 600
    request_timeout => 600
    socket_timeout => 600
    id => "redis"
  }
```

