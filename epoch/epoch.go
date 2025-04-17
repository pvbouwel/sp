/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package epoch

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

func bytesToInt64(bytes []byte) int64 {
	i, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		fmt.Printf("Could not parse int64 %s due to %s", string(bytes), err)
	}
	return i
}

func replaceIfEpoch(bytes []byte) string {
	r, _ := regexp.Compile(`\.`)
	loc := r.FindIndex(bytes)
	var t time.Time
	var sec, nsec int64
	var secEndIdx int
	if loc == nil {
		// No dot
		nsec = 0
		secEndIdx = len(bytes)
	} else {
		secEndIdx = loc[0]
		nsec = bytesToInt64(bytes[loc[1]:])
	}
	sec = bytesToInt64(bytes[0:secEndIdx])

	t = time.Unix(sec, nsec)

	return t.UTC().Format("2006-01-02T15:04:05Z")
}

func replaceEpochs(line []byte) []byte {
	r, _ := regexp.Compile(`[0-9]{10}`)
	loc := r.FindIndex(line)
	if loc == nil {
		return line
	} else {
		var rewrite = make([]byte, 0)

		rewrite = append(rewrite, line[0:loc[0]]...)
		rewrite = append(rewrite, []byte(replaceIfEpoch(line[loc[0]:loc[1]]))...)
		rewrite = append(rewrite, line[loc[1]:]...)
		return rewrite
	}
}

type epoch struct {
	wrapped io.Writer
}

func (e epoch) Write(b []byte) (n int, err error) {
	n = len(b)

	_, err = e.wrapped.Write(replaceEpochs(b))
	return
}

func NewEpoch(w io.Writer) *epoch {
	return &epoch{
		wrapped: w,
	}
}
