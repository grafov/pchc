package sensor

import (
	"fmt"
	"os"
	"time"

	"go.bug.st/serial"
)

type arduino struct {
	port serial.Port
}

var dev arduino

// Open TTY on Arduino device. Take into account that TTY number may
// changed in time.
func Open() {
	mode := &serial.Mode{
		BaudRate: 9600,
	}
	for {
		for i := range 9 {
			if port, err := serial.Open(fmt.Sprintf("/dev/ttyUSB%d", i), mode); err == nil {
				dev.port = port
				return
			}
		}
		fmt.Fprintf(os.Stderr, "can't open any TTY device, trying again in 30 sec...")
		time.Sleep(30 * time.Second)
	}
}

// Read value from sensor. Try to reopen TTY if the read attempt
// failed.
func Read() float64 {
	var val float64
	for {
		n, err := fmt.Fscanf(dev.port, "%f\n", &val)
		if err != nil {
			dev.port.Close()
			time.Sleep(1 * time.Second)
			Open()
			continue
		}
		if n == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	return val
}
