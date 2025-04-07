/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package color_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/fatih/color"
	c "github.com/pvbouwel/sp/color"
)


func TestRotatingColor(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	// Given a 2 stride-length rotator colorer
	twoStride := c.NewFixedStrideLengthFunc(2)

	//Given a buffer to write into
	rb := new(bytes.Buffer)

	c1 := color.New(color.FgRed)
	c2 := color.New(color.FgGreen)
	//WHEN we create a writer with rotating colors
	w := c.NewRotatingColor(rb, []*color.Color{c1, c2}, twoStride)
	_, err := w.Write([]byte("Hello"))
	if err != nil {
		t.Errorf("Encountered error when writing msg: %s", err)
	}

	line, err := rb.ReadString(byte('\n'))
	if err != nil && err != io.EOF {
		t.Errorf("Encoutnered %s while reading buffer", err)
		t.FailNow()
	}
	expectedLine := "\x1b[31mHe\x1b[0m\x1b[32mll\x1b[0m\x1b[31mo\x1b[0m"
	if line != expectedLine {
		t.Errorf("\nExpected:%s\nGot     :%s", expectedLine, line)
	}
}


func TestRotatingColor2writes(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	// Given a 2 stride-length rotator colorer
	twoStride := c.NewFixedStrideLengthFunc(2)

	//Given a buffer to write into
	rb := new(bytes.Buffer)

	c1 := color.New(color.FgRed)
	c2 := color.New(color.FgGreen)
	//WHEN we create a writer with rotating colors
	w := c.NewRotatingColor(rb, []*color.Color{c1, c2}, twoStride)
	_, err := w.Write([]byte("Hel"))
	if err != nil {
		t.Errorf("Encountered error when writing msg: %s", err)
	}
	_, err = w.Write([]byte("lo"))
	if err != nil {
		t.Errorf("Encountered error when writing msg: %s", err)
	}

	line, err := rb.ReadString(byte('\n'))
	if err != nil && err != io.EOF {
		t.Errorf("Encoutnered %s while reading buffer", err)
		t.FailNow()
	}
	expectedLine := "\x1b[31mHe\x1b[0m\x1b[32ml\x1b[0m\x1b[32ml\x1b[0m\x1b[31mo\x1b[0m"
	if line != expectedLine {
		t.Errorf("\nExpected:%s\nGot     :%s", expectedLine, line)
	}
}
