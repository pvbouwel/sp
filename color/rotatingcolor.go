/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package color

import (
	"io"
	"math/rand"

	"github.com/fatih/color"
)

type rotatingColor struct {
	wrapped io.Writer

	rotColors []*color.Color
	colorIdx  int

	currentStrideRemaining int

	//A function to decide how much characters of the current color need to be printed
	//The value can occasionally be 0 but cannot be negative
	getNextStrideLength func() int
}

func NewFixedStrideLengthFunc(l int) func() int {
	return func() int {
		return l
	}
}

func NewRandomStrideLengthFunc(lMin, lMax int) func() int {
	return func() int {
		return rand.Intn(lMax-lMin) + lMin
	}
}

func NewRotatingColor(w io.Writer, rotColors []*color.Color, getNextStrideLength func() int) io.Writer {
	return &rotatingColor{
		wrapped:   w,
		rotColors: rotColors,
		colorIdx:  -1,

		currentStrideRemaining: 0,
		getNextStrideLength:    getNextStrideLength,
	}
}

func (c *rotatingColor) Write(bytes []byte) (int, error) {
	if c.colorIdx >= 0 {
		//Continue where we left off
		c.rotColors[c.colorIdx].SetWriter(c.wrapped)
	}

	bytesToWrite := len(bytes)
	var bytesWritten int

	for bytesWritten = 0; bytesWritten < bytesToWrite; {
		if c.currentStrideRemaining == 0 {
			// Will have to go to the next color
			if c.colorIdx >= 0 {
				//Undo last color
				c.rotColors[c.colorIdx].UnsetWriter(c.wrapped)
			}
			c.colorIdx = (c.colorIdx + 1) % len(c.rotColors)
			c.rotColors[c.colorIdx].SetWriter(c.wrapped)
			c.currentStrideRemaining = c.getNextStrideLength()
			if c.currentStrideRemaining < 0 {
				panic("Invalid getNextStrideLength function returned negative number")
			}
		} else {
			remaining := min(bytesToWrite-bytesWritten, c.currentStrideRemaining)
			writtenChunk, errChunk := c.wrapped.Write(bytes[bytesWritten : bytesWritten+remaining])
			if errChunk != nil {
				return bytesWritten, errChunk
			}
			bytesWritten += writtenChunk
			c.currentStrideRemaining -= writtenChunk
		}
	}
	c.rotColors[c.colorIdx].UnsetWriter(c.wrapped)
	return bytesWritten, nil
}
