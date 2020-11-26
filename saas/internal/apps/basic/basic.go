package basic

import (
	"fmt"
	"net"
	e "saas/internal/debug/err"
	"sync"
)

func LookIP(arg string) {
	fmt.Println("Looking up IP addresses for hostname: " + arg)

	ips, err := net.LookupHost(arg)
	e.CheckErr(err)
	for _, ip := range ips {
		fmt.Println(ip)
	}
}

func LookMX(arg string) {
	fmt.Println("Looking up MX records for : " + arg)

	mxRecords, err := net.LookupMX(arg)
	e.CheckErr(err)

	for _, mxRecord := range mxRecords {
		fmt.Printf("Host: %s\tPreference: %d\n", mxRecord.Host, mxRecord.Pref)
	}
}

func LookName(arg string) {
	fmt.Println("Looking up nameservers for : " + arg)

	nameservers, err := net.LookupNS(arg)
	e.CheckErr(err)

	for _, nameserver := range nameservers {
		fmt.Println(nameserver.Host)
	}
}

func IsPortOpen(addr string, port int) (flag bool) {
	flag = false
	_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	e.CheckErr(err)
	flag = true
	return
}

func ScanSysPorts(addr string) []int {
	var open []int

	var wg sync.WaitGroup
	for i := 1; i <= 1024; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
			if err != nil {
				return
			}
			conn.Close()
			open = append(open, port)
		}(i)
	}
	wg.Wait()
	return open
}

func ScanPortsInRange(addr string, r1 int, r2 int) [] int {
	var open []int

	var wg sync.WaitGroup
	for i := r1; r1 <= r2; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
			if err != nil {
				return
			}
			conn.Close()
			open = append(open, port)
		}(i)
	}
	wg.Wait()
	return open
}
