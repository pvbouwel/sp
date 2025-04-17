/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package color

import (
	"io"

	"github.com/fatih/color"
)

type defaultColor struct {
	wrapped io.Writer

	c color.Color
}

func NewDefaultColor(w io.Writer, c color.Color) io.Writer {
	return &defaultColor{
		wrapped: w,
		c:       c,
	}
}

func (dc *defaultColor) Write(p []byte) (n int, err error) {
	dc.c.SetWriter(dc.wrapped)
	_, err = dc.wrapped.Write(p)
	if err != nil {
		return
	}
	dc.c.UnsetWriter(dc.wrapped)
	n = len(p)

	return n, err
}
