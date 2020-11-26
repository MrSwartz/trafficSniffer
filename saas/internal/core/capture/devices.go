package packets

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	e "saas/internal/debug/err"
)

func ParseDevs(DEVICES int) []Device {
	if devices, err := pcap.FindAllDevs(); err != nil {
		e.CheckErr(err)
		return nil
	} else {
		list := make([]Device, DEVICES)
		temp := Device{}

		for _, v := range devices {
			temp.Name = v.Name
			temp.Flags = v.Flags
			temp.Description = v.Description
			temp.Addresses = v.Addresses
			list = append(list, temp)
		}
		return list
	}
}

func PrintDevs(DEVICES int) {
	if devs := ParseDevs(DEVICES); devs != nil {
		for i := range devs {
			fmt.Printf("Index : %d\n", i)
			fmt.Printf("Name %s\n", devs[i].Name)
			fmt.Printf("Flags %d\n", devs[i].Flags)
			fmt.Printf("Description %s\n", devs[i].Description)
			for ind := range devs[i].Addresses {
				fmt.Printf("\tBroadaddr %v\n", devs[i].Addresses[ind].Broadaddr)
				fmt.Printf("\tIP %v\n", devs[i].Addresses[ind].IP)
				fmt.Printf("\tNetmask %v\n", devs[i].Addresses[ind].Netmask)
				fmt.Printf("\tP2P %v\n", devs[i].Addresses[ind].P2P)
			}
			fmt.Println()
		}
	} else {
		fmt.Println("error in devices")
	}
}