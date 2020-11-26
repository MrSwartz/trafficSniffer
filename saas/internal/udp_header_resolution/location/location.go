package location

/*
This package contains IP location functions, parsing from a csv file,
downloaded from https://lite.ip2location.com/
*/

import (
	e "saas/internal/debug/err"
	"strconv"
)

type IPv4 struct {
	p uint8
	q uint8
	r uint8
	s uint8
}

type IPv6 struct {
	p uint64
	q uint64
	r uint64
	s uint64
	t uint64
	x uint64
	y uint64
	z uint64
}

type Coordinates struct {
	latitude  float32
	longitude float32
}

type IPv4v6Host struct {
	IPv4From      IPv4
	IPv4To        IPv4
	IPv6From	  IPv6
	IPv6To		  IPv6
	Location    Coordinates
	CountryCode string
	CountryName string
	RegionName  string
	City        string
	Timezone    string
}

type IPv4v6Proxy struct {
	IPv4From      IPv4
	IPv4To        IPv4
	IPv6From      IPv6
	IPv6To        IPv6
	ASN         uint64
	ProxyType   string
	CountryCode string
	CountryName string
	RegionName  string
	City        string
	ISP         string
	Domain      string
	UsageType   string
	AS          string
}

func ConvertToIPv4(IP string) IPv4 {
	var ip IPv4
	code, err := strconv.ParseUint(IP, 10, 64)
	e.CheckErr(err)
	ip.p = uint8((code >> 24) & 0xFF)
	ip.q = uint8((code >> 16) & 0xFF)
	ip.r = uint8((code >> 8) & 0xFF)
	ip.s = uint8(code & 0xFF)
	return ip
}

func ConvertToIPv6(IP string) IPv6{
	var ip IPv6
	code, err := strconv.ParseUint(IP, 16, 64)
	e.CheckErr(err)
	ip.p = (code / (65536 << 7)) % 65536
	ip.q = (code / (65536 << 6)) % 65536
	ip.r = (code / (65536 << 5)) % 65536
	ip.s = (code / (65536 << 4)) % 65536
	ip.t = (code / (65536 << 3)) % 65536
	ip.x = (code / (65536 << 2)) % 65536
	ip.y = (code / 65536) % 65536
	ip.z = code % 65536
	return ip
}
