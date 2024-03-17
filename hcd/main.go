package main

import (
	"fmt"
	"hcd/sensor"
	"hcd/tuner"
	"time"
)

func main() {
	sensor.Open()
	tuner.Setup()
	update := time.NewTicker(500 * time.Millisecond)
	for {
		val := sensor.Read()
		select {
		case <-update.C:
			fmt.Printf("brightness is %3.3f (raw)\n", val)
			tuner.Set(val)
		}
	}
}
