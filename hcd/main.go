package main

import (
	"fmt"
	"hcd/sensor"
	"math"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	sensor.Open()
	var prev float64
	update := time.NewTicker(5 * time.Second)
	for {
		val := sensor.Read()
		if val == 0 || math.Abs(val-prev) < 1 {
			prev = val
			continue
		}
		prev = val
		select {
		case <-update.C:
			fmt.Printf("set brightness to %3.3f (raw)\n", val)
			for i := 1; i <= 3; i++ {
				m, max := tuneForMonitor(i)
				br := int(val * m / 280 * max)
				if br > int(max) {
					br = int(max)
				}
				fmt.Printf("mon %d val %d \n", i, br)
				if err := ddiset(i, br); err != nil {
					fmt.Printf("error changing brightness on mon %d (%d): %s\n", i, br, err)
				}
				time.Sleep(time.Millisecond * 20)
			}
		}
	}
}

// another hack
func tuneForMonitor(mon int) (mult, max float64) {
	switch mon {
	case 1: // NEC 24"
		return 1.5, 7000
	case 2: // LG 38"
		return 1, 100
	case 3: // NEC 27"
		return 1, 100
	default:
		return 1, 100
	}
}

// dirty hack based on ddcutil output
func ddiset(mon int, b int) error {
	cmd := exec.Command("/usr/bin/ddcutil", "-d", strconv.Itoa(mon), "setvcp", "10", strconv.Itoa(b))
	_, err := cmd.Output()
	return err
}
