package wait

import (
	"fmt"
	"time"
)

func spinner(delay time.Duration){
	for {
		for _, r := range `-\|/`{
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func Spinner(pat int, delay time.Duration){
	var load0 = []string{".", "..", "...", "....", "....."}
	var load1 = []string{"-", "\\", "|", "/"}
	var load2 = []string{"*", "**", "***"}
	var load []string

	switch  {
	case pat == 0: load = load0
	case pat == 1: load = load1
	case pat == 2: load = load2
	default: load = load1
	}

	for {
		for i := range load{
			fmt.Printf("\rPlease wait %s", load[i])
			time.Sleep(delay)
		}
	}
}

func WaitSeconds()  {
	for {
		for i := 0; i < (2 << 16) - 1; i++{
			fmt.Printf("\rSpent time: %d s", i)
			time.Sleep(1 * time.Second)
		}
	}
}
