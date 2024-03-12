package main

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"time"

	"go.bug.st/serial"
)

func main() {
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	port, err := serial.Open("/dev/ttyUSB0", mode)
	if err != nil {
		log.Fatal(err)
	}
	update := time.NewTicker(5 * time.Second)
	var prev float64
	for {
		time.Sleep(time.Millisecond * 5)
		var val float64
		n, err := fmt.Fscanf(port, "%f\n", &val)
		if err != nil {
			log.Println(err)
			continue
		}
		if n == 0 {
			fmt.Println("\nEOF")
			continue
		}
		if val == 0 || math.Abs(val-prev) < 1 {
			prev = val
			fmt.Println("not much difference, ignored")
			continue
		}
		prev = val
		select {
		case <-update.C:
			fmt.Printf("set brightness to %3.3f (raw)\n", val)
			for i := 1; i <= 3; i++ {
				br := int(val / 280 * maxForMonitor(i))
				fmt.Printf("mon %d val %d \n", i, br)
				if err := ddiset(i, br); err != nil {
					fmt.Printf("error changing brightness on mon %d (%d): %s\n", i, br, err)
				}
				time.Sleep(time.Millisecond * 20)
			}
		default:
		}
	}
}

// another hack
func maxForMonitor(mon int) float64 {
	switch mon {
	case 1: // NEC 24"
		return 7000
	case 2: // LG 38"
		return 100
	case 3: // NEC 27"
		return 100
	default:
		return 100
	}
}

// dirty hack based on ddcutil output
func ddiset(mon int, b int) error {
	cmd := exec.Command("/usr/bin/ddcutil", "-d", strconv.Itoa(mon), "setvcp", "10", strconv.Itoa(b))
	_, err := cmd.Output()
	return err
}
