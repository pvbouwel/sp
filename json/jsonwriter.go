/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package jsonwriter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
)

type enclosedWriter struct {
	//Underlying writer to which everything ends up
	wrapped io.Writer

	//Bytes used for enclosing opening and closing (e.g. {})
	braceBytes []byte

	//Byte used for declaring literals (e.g. `"`)
	literalByte byte

	//Byte used for escaping upcoming character (e.g. `\`)
	escapeByte byte

	//writer that is called with slices that are possible json dicts
	embracedWriter io.Writer
}

type colourDecider interface {
	decide(m map[string]any) *color.Color
}

type mapBasedColourDecider struct {
	m map[string]map[any][]*color.Color

	i int
}

func NewMapBasedColourDecider(c ...JSONColor) mapBasedColourDecider {
	var d = mapBasedColourDecider{}
	d.m = map[string]map[any][]*color.Color{}
	for _, ci := range c {
		colourMap, ok := d.m[ci.Key]
		if !ok {
			cm := map[any][]*color.Color{}
			cm[ci.Value] = ci.Color
			d.m[ci.Key] = cm
		} else {
			colourMap[ci.Value] = append(colourMap[ci.Value], ci.Color...)
		}
	}
	return d
}

func (d *mapBasedColourDecider) decide(m map[string]any) *color.Color {
	for key, value := range m {
		mapToColour, ok := (d.m)[key]
		if ok {
			c, ok := mapToColour[value]
			if ok {
				chosenColor := c[d.i%len(c)]
				d.i = (d.i + 1) % 100000
				return chosenColor
			}

		}
	}
	return nil
}

type possibleJSONWriter struct {
	wrapped io.Writer

	//A mapping to decide the colour. it maps key names to a mapping to decide on colour
	//if there are multiple matches there is no guarantee on a winner.
	colourDecider colourDecider
}

func (j *possibleJSONWriter) Write(p []byte) (n int, err error) {
	var decoded map[string]any
	err = json.Unmarshal(p, &decoded)
	if err != nil {
		//Unsupported JSON let's not fail
		return j.wrapped.Write(p)
	}
	c := j.colourDecider.decide(decoded)
	if c == nil {
		return j.wrapped.Write(p)
	}
	c.SetWriter(j.wrapped)
	n, err = j.wrapped.Write(p)
	c.UnsetWriter(j.wrapped)

	return n, err
}

type JSONColor struct {
	Key   string
	Value any
	Color []*color.Color
}

func NewJSONWriter(w io.Writer, c mapBasedColourDecider) io.Writer {
	return &enclosedWriter{
		wrapped:     w,
		braceBytes:  []byte{byte('{'), byte('}')},
		literalByte: byte('"'),
		escapeByte:  byte('\\'),
		embracedWriter: &possibleJSONWriter{
			wrapped:       w,
			colourDecider: &c,
		},
	}
}

func (j *enclosedWriter) Write(p []byte) (n int, err error) {
	n = 0
	var idxOpeningBrace int
	curlyBraceDepth := 0
	var inLiteral bool

	for i := 0; i < len(p); i++ {
		switch p[i] {
		case j.braceBytes[0]:
			if inLiteral {
				continue
			}
			if curlyBraceDepth == 0 && i != 0 {
				//First curly brace what comes before must be written unprocessed
				ni, err := j.wrapped.Write(p[n:i])
				n += ni
				if err != nil {
					return n, err
				}
			}
			curlyBraceDepth += 1
			if curlyBraceDepth == 1 {
				idxOpeningBrace = i
			}
		case j.braceBytes[1]:
			if inLiteral {
				continue
			}
			curlyBraceDepth -= 1
			if curlyBraceDepth == 0 {
				ni, err := j.embracedWriter.Write(p[idxOpeningBrace : i+1])
				n += ni
				if err != nil {
					return n, err
				}
			}
		case j.literalByte:
			inLiteral = !inLiteral
		case j.escapeByte:
			i += 1
		default:
		}
	}
	if n < len(p) {
		ni, err := j.wrapped.Write(p[n:])
		n += ni
		if err != nil {
			return n, err
		}
	}
	if n != len(p) {
		err = fmt.Errorf("expected to write %d bytes written %d", len(p), n)
	}
	return n, err
}
