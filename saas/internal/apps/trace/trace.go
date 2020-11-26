package trace

import (
	"errors"
	"fmt"
	"log"
	"net"
	e "saas/internal/debug/err"
	"syscall"
	"time"
)

type Opts struct {
	Port       int
	MaxHops    int
	FirstHop   int
	Timeout    int
	Retries    int
	PacketSize int
}

type TraceHop struct {
	N           int
	TTL         int
	Addr        [4]byte
	Success     bool
	ElapsedTime time.Duration
	Host        string
}

type TraceRes struct {
	DstAddr [4]byte
	Hops    []TraceHop
}

var opts = Opts{33434, 64, 1, 500, 3, 52}

func sockAddr() (ret [4]byte, err error) {
	addr, err := net.InterfaceAddrs()
	e.CheckErr(err)

	for i := range addr {
		if ip, ok := addr[i].(*net.IPNet); ok && ip.IP.IsLoopback() {
			if len(ip.IP.To4()) == net.IPv4len {
				copy(ret[:], ip.IP.To4())
				return ret, err
			}
		}
	}
	err = errors.New("msg")
	return ret, nil
}

func dstAddr(dst string) (dstAddr [4]byte, err error) {
	addrs, err := net.LookupHost(dst)
	e.CheckErr(err)
	addr := addrs[0]

	ipAddr, err := net.ResolveIPAddr("ip", addr)
	e.CheckErr(err)
	copy(dstAddr[:], ipAddr.IP.To4())
	return
}

func setConsts(c *Opts, port, maxHops, firstHop, timeout, retries, packetSize int) {
	c.Port = port
	c.MaxHops = maxHops
	c.FirstHop = firstHop
	c.Timeout = timeout
	c.Retries = retries
	c.PacketSize = packetSize
}

func getMaxHops(c *Opts) int {
	if c.MaxHops == 0 {
		c.MaxHops = opts.MaxHops
	}
	return c.MaxHops
}

func getPort(c *Opts) int {
	if c.Port == 0 {
		c.Port = opts.Port
	}
	return c.Port
}

func getFirstHop(c *Opts) int {
	if c.FirstHop == 0 {
		c.FirstHop = opts.FirstHop
	}
	return c.FirstHop
}

func getTimeout(c *Opts) int {
	if c.Timeout == 0 {
		c.Timeout = opts.Timeout
	}
	return c.Timeout
}

func getRetries(c *Opts) int {
	if c.Retries == 0 {
		c.Retries = opts.Retries
	}
	return c.Retries
}

func getPacketSize(c *Opts) int {
	if c.PacketSize == 0 {
		c.PacketSize = opts.PacketSize
	}
	return c.PacketSize
}

func hostAddr(h *TraceHop) string {
	ha := fmt.Sprintf("%v.%v.%v.%v", h.Addr[0], h.Addr[1], h.Addr[2], h.Addr[3])
	if h.Host != "" {
		ha = h.Host
	} else {
		log.Println("Error in func HostAddr")
	}
	return ha
}

func notifier(h TraceHop, ch []chan TraceHop) {
	for _, c := range ch {
		c <- h
	}
}

func closeNotifier(ch []chan TraceHop) {
	for _, c := range ch {
		close(c)
	}
}

func Trace(dst string, opts *Opts, c ...chan TraceHop) (TraceRes, error) {
	var res TraceRes
	res.Hops = []TraceHop{}
	dstAddr, err := dstAddr(dst)
	e.CheckErr(err)
	res.DstAddr = dstAddr
	sockAddr, err := sockAddr()
	e.CheckErr(err)

	timeout := int64(getTimeout(opts))
	ts := syscall.NsecToTimeval(100000 * timeout)

	ttl := getFirstHop(opts)
	retry := 0

	for {
		t1 := time.Now()
		recvSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
		if err != nil {
			return res, err
		}
		sendSock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
		if err != nil {
			return res, err
		}

		syscall.SetsockoptInt(sendSock, 0x0, syscall.IP_TTL, ttl)
		syscall.SetsockoptTimeval(recvSock, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &ts)
		defer syscall.Close(recvSock)
		defer syscall.Close(sendSock)

		syscall.Bind(recvSock, &syscall.SockaddrInet4{Port: getPort(opts), Addr: sockAddr})
		syscall.Sendto(sendSock, []byte{0x0}, 0, &syscall.SockaddrInet4{Port: getPort(opts), Addr: dstAddr})

		var p = make([]byte, opts.PacketSize)
		n, from, err := syscall.Recvfrom(recvSock, p, 0)
		el := time.Since(t1)
		if err == nil {
			curAddr := from.(*syscall.SockaddrInet4).Addr
			hop := TraceHop{Success: true, Addr: curAddr, N: n, ElapsedTime: el, TTL: ttl}
			curHost, err := net.LookupAddr(fmt.Sprintf("%v.%v.%v.%v", hop.Addr[0], hop.Addr[1], hop.Addr[2], hop.Addr[3]))
			if err == nil {
				hop.Host = curHost[0]
			}

			notifier(hop, c)
			res.Hops = append(res.Hops, hop)
			ttl++
			retry = 0

			if ttl > getMaxHops(opts) || curAddr == dstAddr {
				closeNotifier(c)
				return res, nil
			}
		} else {
			retry++
			if retry > getRetries(opts) {
				notifier(TraceHop{Success: false, TTL: ttl}, c)
				ttl++
				retry = 0
			}
			if ttl > getMaxHops(opts) {
				closeNotifier(c)
				return res, nil
			}
		}
	}

}
