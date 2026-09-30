package serve

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// clockTicks is the kernel's USER_HZ, which /proc counts a process's start
// in: 100 on every architecture Linux has shipped for years.
const clockTicks = 100

// process is pid's session and when it started, from /proc: the start is
// the boot time in /proc/stat plus the clock ticks since boot in
// /proc/<pid>/stat, and zero when the boot time cannot be read.
func process(pid int) (sid int, started time.Time, err error) {
	stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return 0, time.Time{}, err
	}
	// the command's name, in parentheses, may hold spaces: count from after
	// it, where state is field 3, session field 6, and starttime field 22
	i := strings.LastIndexByte(string(stat), ')')
	fields := strings.Fields(string(stat[i+1:]))
	if i < 0 || len(fields) < 20 {
		return 0, time.Time{}, errors.New("/proc/" + strconv.Itoa(pid) + "/stat has fewer fields than it should")
	}
	if sid, err = strconv.Atoi(fields[3]); err != nil {
		return 0, time.Time{}, err
	}
	ticks, err := strconv.ParseInt(fields[19], 10, 64)
	if err != nil {
		return 0, time.Time{}, err
	}
	if boot, ok := bootTime(); ok {
		started = boot.Add(time.Duration(ticks) * time.Second / clockTicks)
	}
	return sid, started, nil
}

// bootTime is when the system booted, from btime in /proc/stat.
func bootTime() (time.Time, bool) {
	stat, err := os.ReadFile("/proc/stat")
	if err != nil {
		return time.Time{}, false
	}
	for line := range strings.SplitSeq(string(stat), "\n") {
		if v, ok := strings.CutPrefix(line, "btime "); ok {
			secs, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
			if err != nil {
				return time.Time{}, false
			}
			return time.Unix(secs, 0), true
		}
	}
	return time.Time{}, false
}
