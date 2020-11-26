package manufacture

import (
	"io/ioutil"
	e "saas/internal/debug/err"
)

func LoadMacList() []string {
	fl, err := ioutil.ReadFile("manufacture/mac.txt")
	e.CheckErr(err)
	list := string(fl)

	var count int
	var buf string
	mac := make([]string, 38999)

	for _, v := range list{
		switch {
		case v != '\n': buf += string(v)
		case v == '\n': mac[count] = buf
			buf = ""
			count++
		}
	}
	return mac
}		// загружает в память мак адреса (использовать в main)

func FindManufacture(mac string, macList []string)(manuf string){
	b1 := string(mac[0]) + string(mac[1])
	b2 := string(mac[3]) + string(mac[4])
	b3 := string(mac[6]) + string(mac[7])

	for i := range macList{
		temp := macList[i]
		if (string(temp[0]) + string(temp[1])) == b1 && (string(temp[3]) + string(temp[4])) == b2  && (string(temp[6]) + string(temp[7])) == b3{
			manuf = temp
			break
		}
	}
	return
}