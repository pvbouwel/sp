/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package jsonwriter_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/fatih/color"
	jsonwriter "github.com/pvbouwel/sp/json"
)

func getTestMapBasedColourDeciders() []jsonwriter.ColourDecider {
	var trafficDecider = jsonwriter.NewMapBasedColourDecider(
		false,
		jsonwriter.JSONColor{Key: "level", Value: "info", Color: []*color.Color{color.RGB(0, 255, 0), color.RGB(0, 205, 0)}},
		jsonwriter.JSONColor{Key: "level", Value: "warning", Color: []*color.Color{color.RGB(255, 128, 0)}},
		jsonwriter.JSONColor{Key: "level", Value: "error", Color: []*color.Color{color.RGB(255, 0, 0)}},
	)

	var trafficDeciderIgnoreCase = jsonwriter.NewMapBasedColourDecider(
		true,
		jsonwriter.JSONColor{Key: "level", Value: "InfO", Color: []*color.Color{color.RGB(0, 255, 0), color.RGB(0, 205, 0)}},
		jsonwriter.JSONColor{Key: "level", Value: "WARNING", Color: []*color.Color{color.RGB(255, 128, 0)}},
		jsonwriter.JSONColor{Key: "level", Value: "eRRor", Color: []*color.Color{color.RGB(255, 0, 0)}},
	)
	return []jsonwriter.ColourDecider{
		&trafficDecider,
		&trafficDeciderIgnoreCase,
	}
}

func TestJSONTrafficSimple(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	for i, td := range getTestMapBasedColourDeciders() {
		//Given a buffer to write into
		rb := new(bytes.Buffer)

		//WHEN we create a writer with the decider
		w := jsonwriter.NewJSONWriter(rb, td)

		_, err := w.Write([]byte("{\"level\": \"info\"}{\"level\": \"warning\"}\n{\"level\":\"error\"}"))
		if err != nil {
			t.Errorf("%d: Encountered error when writing msg: %s", i, err)
		}

		line, err := rb.ReadString(byte('\n'))
		if err != nil {
			t.Errorf("%d: Encountered %s while reading buffer", i, err)
			t.FailNow()
		}
		expectedLine := "\x1b[38;2;0;255;0m{\"level\": \"info\"}\x1b[0m\x1b[38;2;255;128;0m{\"level\": \"warning\"}\x1b[0m\n"
		if line != expectedLine {
			t.Errorf("\n%d: Expected:%s\nGot     :%s", i, expectedLine, line)
		}
		line, err = rb.ReadString(byte('\n'))
		if err != io.EOF {
			t.Errorf("%d: Encountered %s while reading buffer", i, err)
			t.FailNow()
		}
		expectedLine = "\x1b[38;2;255;0;0m{\"level\":\"error\"}\x1b[0m"
		if line != expectedLine {
			t.Errorf("\n%d: Expected:%s\nGot     :%s", i, expectedLine, line)
		}
	}
}

func TestJSONTrafficSimpleBanding(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	for i, td := range getTestMapBasedColourDeciders() {
		//Given a buffer to write into
		rb := new(bytes.Buffer)

		//WHEN we create a writer with the decider
		w := jsonwriter.NewJSONWriter(rb, td)

		_, err := w.Write([]byte("{\"level\": \"info\"}\n{\"level\": \"info\"}\n{\"level\": \"info\"}\n"))
		if err != nil {
			t.Errorf("%d: Encountered error when writing msg: %s", i, err)
		}

		line, err := rb.ReadString(byte('\n'))
		if err != nil {
			t.Errorf("%d: Encountered %s while reading buffer", i, err)
			t.FailNow()
		}
		expectedLine := "\x1b[38;2;0;255;0m{\"level\": \"info\"}\x1b[0m\n"
		if line != expectedLine {
			t.Errorf("\n%d.A: Expected:%s\nGot     :%s", i, expectedLine, line)
		}
		line, _ = rb.ReadString(byte('\n'))
		expectedLineAlternate := "\x1b[38;2;0;205;0m{\"level\": \"info\"}\x1b[0m\n"
		if line != expectedLineAlternate {
			t.Errorf("\n%d.B: Expected:%s\nGot     :%s", i, expectedLineAlternate, line)
		}
		line, _ = rb.ReadString(byte('\n'))
		if line != expectedLine {
			t.Errorf("\n%d: Expected:%s\nGot     :%s", i, expectedLine, line)
		}
	}
}

func TestJSONTrafficWithIgnorableClosingBrace(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	for i, td := range getTestMapBasedColourDeciders() {

		//Given a buffer to write into
		rb := new(bytes.Buffer)

		//WHEN we create a writer with the decider
		w := jsonwriter.NewJSONWriter(rb, td)

		_, err := w.Write([]byte("{\"level\": \"info\", \"msg\": \"}\"}{\"level\": \"warning\"}\n{\"level\":\"error\"}"))
		if err != nil {
			t.Errorf("%d: Encountered error when writing msg: %s", i, err)
		}

		line, err := rb.ReadString(byte('\n'))
		if err != nil && err != io.EOF {
			t.Errorf("%d: Encountered %s while reading buffer", i, err)
			t.FailNow()
		}
		expectedLine := "\x1b[38;2;0;255;0m{\"level\": \"info\", \"msg\": \"}\"}\x1b[0m\x1b[38;2;255;128;0m{\"level\": \"warning\"}\x1b[0m\n"
		if line != expectedLine {
			t.Errorf("\n%d: Expected:%s\nGot     :%s", i, expectedLine, line)
		}
	}
}

func TestJSONTrafficWithIgnorableClosingBraceAndEscapedLiteral(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	for i, td := range getTestMapBasedColourDeciders() {

		//Given a buffer to write into
		rb := new(bytes.Buffer)

		//WHEN we create a writer with the decider
		w := jsonwriter.NewJSONWriter(rb, td)

		_, err := w.Write([]byte("{\"level\": \"info\", \"msg\": \"\\\"}\"}{\"level\": \"warning\"}\n{\"level\":\"error\"}"))
		if err != nil {
			t.Errorf("%d: Encountered error when writing msg: %s", i, err)
		}

		line, err := rb.ReadString(byte('\n'))
		if err != nil && err != io.EOF {
			t.Errorf("%d: Encountered %s while reading buffer", i, err)
			t.FailNow()
		}
		expectedLine := "\x1b[38;2;0;255;0m{\"level\": \"info\", \"msg\": \"\\\"}\"}\x1b[0m\x1b[38;2;255;128;0m{\"level\": \"warning\"}\x1b[0m\n"
		if line != expectedLine {
			t.Errorf("\n%d: Expected:%s\nGot     :%s", i, expectedLine, line)
		}
	}
}

func TestJSONTrafficFunkeyValues(t *testing.T) {
	//Given color is to be done
	color.NoColor = false

	//Given a buffer to write into
	rb := new(bytes.Buffer)

	var trafficDeciderFunkyValues = jsonwriter.NewMapBasedColourDecider(
		false,
		jsonwriter.JSONColor{Key: "level", Value: "i.nfo", Color: []*color.Color{color.RGB(0, 255, 0)}},
		jsonwriter.JSONColor{Key: "level", Value: "w.a.r.n.i.n.g", Color: []*color.Color{color.RGB(255, 128, 0)}},
		jsonwriter.JSONColor{Key: "level", Value: "error.", Color: []*color.Color{color.RGB(255, 0, 0)}},
	)

	//WHEN we create a writer with the decider
	w := jsonwriter.NewJSONWriter(rb, &trafficDeciderFunkyValues)

	_, err := w.Write([]byte("{\"level\": \"i.nfo\"}{\"level\": \"w.a.r.n.i.n.g\"}\n{\"level\":\"error.\"}"))
	if err != nil {
		t.Errorf("Encountered error when writing msg: %s", err)
	}

	line, err := rb.ReadString(byte('\n'))
	if err != nil {
		t.Errorf("Encountered %s while reading buffer", err)
		t.FailNow()
	}
	expectedLine := "\x1b[38;2;0;255;0m{\"level\": \"i.nfo\"}\x1b[0m\x1b[38;2;255;128;0m{\"level\": \"w.a.r.n.i.n.g\"}\x1b[0m\n"
	if line != expectedLine {
		t.Errorf("\nExpected:%s\nGot     :%s", expectedLine, line)
	}
	line, err = rb.ReadString(byte('\n'))
	if err != io.EOF {
		t.Errorf("Encountered %s while reading buffer", err)
		t.FailNow()
	}
	expectedLine = "\x1b[38;2;255;0;0m{\"level\":\"error.\"}\x1b[0m"
	if line != expectedLine {
		t.Errorf("\nExpected:%s\nGot     :%s", expectedLine, line)
	}
}
