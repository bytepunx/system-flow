package serve

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// process is pid's session and its start in clock ticks since boot, from
// /proc/<pid>/stat. The start is fixed for the life of the process and
// counted from boot, not the wall clock, so it does not move when the
// clock does, as a WSL host's does after a sleep.
func process(pid int) (sid int, start int64, err error) {
	stat, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/stat")
	if err != nil {
		return 0, 0, err
	}
	// the command's name, in parentheses, may hold spaces: count from after
	// it, where state is field 3, session field 6, and starttime field 22
	i := strings.LastIndexByte(string(stat), ')')
	fields := strings.Fields(string(stat[i+1:]))
	if i < 0 || len(fields) < 20 {
		return 0, 0, errors.New("/proc/" + strconv.Itoa(pid) + "/stat has fewer fields than it should")
	}
	if sid, err = strconv.Atoi(fields[3]); err != nil {
		return 0, 0, err
	}
	if start, err = strconv.ParseInt(fields[19], 10, 64); err != nil {
		return 0, 0, err
	}
	return sid, start, nil
}
