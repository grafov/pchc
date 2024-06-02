package tuner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

var (
	mons []chan float64
	m    sync.Mutex
)

func Setup() {
	for i := range 3 {
		mon := make(chan float64, 60)
		mons = append(mons, mon)
		go func(mon chan float64, i int) {
			var br, count int
			var prevBr int
			mul, max := tuneForMonitor(i)
			for {
				select {
				case val := <-mon:
					// получаем все накопившиеся значения
					br = int(val * mul / 280 * max)
					if br > int(max) {
						br = int(max)
					}
					if prevBr == br {
						continue
					}
					prevBr = br
					count = 0 // получение нового значения останавливает ретраи старого
				default:
					if checkRunningLocker() {
						continue
					}
					// когда нет свежих значений в канале, обновляем яркость
					switch {
					case count > 50:
						count = 0
						br = 0
					case br != 0:
						if err := ddiset(i, br); err != nil {
							fmt.Fprintf(os.Stderr, "error on mon %d (%d): %s\n", i, br, err)
							count++
							time.Sleep(time.Duration(count) * time.Second)
						} else {
							br = 0
							count = 0
						}
					}
					time.Sleep(5 * time.Second) // установка яркости не чаще этого периода
				}
			}
		}(mon, i+1)
	}
}

func Set(val float64) {
	for _, m := range mons {
		m <- val
	}
}

// dirty hack based on ddcutil output
func ddiset(mon int, b int) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	m.Lock() // для стабильности вызывать ddcutil последовательно для мониторов
	cmd := exec.CommandContext(ctx, "/usr/bin/ddcutil", "-d", strconv.Itoa(mon), "setvcp", "10", strconv.Itoa(b))
	_, err := cmd.Output()
	m.Unlock()
	return err
}

// another hack
func tuneForMonitor(mon int) (mult, max float64) {
	switch mon {
	case 1: // NEC 24"
		return 4.5, 7000
	case 2: // LG 38"
		return 1.2, 100
	case 3: // NEC 27"
		return 1.2, 100
	default:
		return 1, 100
	}
}
