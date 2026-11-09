package serve

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// process is pid's session, its start in clock ticks since boot, and whether
// it has ended and waits only to be reaped, from /proc/<pid>/stat. The start
// is fixed for the life of the process and counted from boot, not the wall
// clock, so it does not move when the clock does, as a WSL host's does after
// a sleep.
func process(pid int) (sid int, start int64, ended bool, err error) {
	path := "/proc/" + strconv.Itoa(pid) + "/stat"
	stat, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, false, err
	}
	return parseStat(path, stat)
}

// parseStat reads the session, the start, and the state from the contents of
// a /proc/<pid>/stat file read from path. A zombie (Z) or a dead process (X)
// has ended: signal 0 still reaches a zombie until its parent reaps it, so
// only the state says it no longer runs (I-0299).
func parseStat(path string, stat []byte) (sid int, start int64, ended bool, err error) {
	// the command's name, in parentheses, may hold spaces and parentheses:
	// count from after the last, where state is field 3, session field 6,
	// and starttime field 22
	i := strings.LastIndexByte(string(stat), ')')
	if i < 0 {
		return 0, 0, false, errors.New(path + " has no command name in parentheses")
	}
	fields := strings.Fields(string(stat[i+1:]))
	if len(fields) < 20 {
		return 0, 0, false, errors.New(path + " has fewer fields than it should")
	}
	if sid, err = strconv.Atoi(fields[3]); err != nil {
		return 0, 0, false, err
	}
	if start, err = strconv.ParseInt(fields[19], 10, 64); err != nil {
		return 0, 0, false, err
	}
	return sid, start, fields[0] == "Z" || fields[0] == "X", nil
}
