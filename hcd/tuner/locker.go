package tuner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const lockerName = "i3lock-fancy"

func checkRunningLocker() bool {
	files, err := os.ReadDir("/proc")
	if err != nil {
		fmt.Fprintln(os.Stderr, "unable to read /proc directory")
		return false
	}
	for _, f := range files {
		f, err := f.Info()
		if err != nil {
			continue
		}
		if pid := findName(f); pid > 0 {
			return true
		}
	}
	return false
}

func findName(file os.FileInfo) int {
	if !file.IsDir() {
		return 0
	}
	// Our directory name should convert to integer
	// if it's a PID
	pid, err := strconv.Atoi(file.Name())
	if err != nil {
		return 0
	}
	// Open the /proc/xxx/stat file to read the name
	f, err := os.Open(file.Name() + "/stat")
	if err != nil {
		// fmt.Fprintf(os.Stderr, "unable to stat for: %s\n", file.Name())
		return 0
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), lockerName) {
			return pid
		}
	}
	return 0
}
