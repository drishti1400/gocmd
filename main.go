package main

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"strconv"
	"strings"
)

func formatTime(time int) string {
	seconds := time / 100
	minutes := seconds / 60
	hours := minutes / 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes%60, seconds%60)
}

func main() {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			pid, err := strconv.Atoi(file.Name())
			if err != nil {
				continue
			}

			stat, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
			if err != nil {
				continue
			}

			fields := strings.Fields(string(stat))
			fmt.Println("Process ID: ", fields[0])
			fmt.Println("Executable name: ", fields[1])
			fmt.Println("Parent process id: ", fields[3])
			fmt.Println("Virtual memory size: ", fields[22])
			utime, _ := strconv.Atoi(fields[13])
			stime, _ := strconv.Atoi(fields[14])
			fmt.Println("utime: ", formatTime(utime))
			fmt.Println("stime: ", formatTime(stime))

			status, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
			if err != nil {
				continue
			}

			lines := strings.Split(string(status), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Uid:") {
					fields := strings.Fields(line)
					uid, _ := strconv.Atoi(fields[1])
					u, err := user.LookupId(strconv.Itoa(uid))
					if err == nil {
						fmt.Println("User: ", u.Username)
					}
				}
			}
			fmt.Println(" ")
		}
	}
}
