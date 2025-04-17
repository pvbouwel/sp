/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package cmd

import (
	"io"
	"sync"
)

var syncedWriterMutex *sync.Mutex = &sync.Mutex{}

// syncedWriter is to make sure that if we have multiple output streams that we
// only write to one exclusively. This should help avoid streams being intertwined
// which are sent to a single text output (e.g. a terminal)
type syncedWriter struct {
	w io.Writer
}

func NewSyncedWriter(w io.Writer) io.Writer {
	return &syncedWriter{
		w: w,
	}
}

func (s *syncedWriter) Write(p []byte) (int, error) {
	syncedWriterMutex.Lock()
	defer syncedWriterMutex.Unlock()
	return s.w.Write(p)
}
