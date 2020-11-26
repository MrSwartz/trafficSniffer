package err

import (
	"fmt"
	"log"
	"os"
)

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func FatalErr(msg string, err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func RetErr(device string, err error, code int) {
	if err != nil {
		fmt.Printf("Error  opening device %s: %v", device, err)
		os.Exit(code)
	}
}
