/*
Copyright © 2025 Peter Van Bouwel <https://github.com/pvbouwel>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	c "github.com/pvbouwel/sp/color"
	jsonwriter "github.com/pvbouwel/sp/json"
	"github.com/spf13/cobra"
)

// colorCmd represents the color command
var colorCmd = &cobra.Command{
	Use:   "color",
	Short: "Add color to a stream",
	Long: `Add color to a stream.
	
	Example 1 : rainbow colors with widths between 2 and 10
	sp color --color-type rotating --rotating-type random --stride-length 2-10
	
	Example 2 : Have json color-coded based on key named level
	sp color --force --color-type JSON --json-key level --colors info.0.255.0,warning.255.128.0,error.255.0.0
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		force, err := cmd.Flags().GetBool("force")
		if err == nil && force {
			color.NoColor = false
		}
		stdoutWriter, err = getWriter(cmd, stdout)
		if err != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("Encountered error: %s", err)))
		}
		stderrWriter, err = getWriter(cmd, stderr)
		if err != nil {
			os.Stderr.Write([]byte(fmt.Sprintf("Encountered error: %s", err)))
		}
	},
}

var colourTable map [string]*color.Color = map [string]*color.Color{
	"red": color.New(color.FgRed),
	"green": color.New(color.FgGreen),
	"white": color.New(color.FgWhite),
}

func availableColors() ([]string){
	var availableColors []string = make([]string, len(colourTable))
	var i int = 0
	for c := range(colourTable) {
		availableColors[i] = c
		i += 1
	}
	return availableColors
}

func getWriter(cmd *cobra.Command, outputType outputType) (io.Writer, error) {
	getFlag := getFlagNameFunc(outputType)

	colorType, err := cmd.Flags().GetString(getFlag(fColorType))
	if err != nil {
		return nil, err
	}

	baseWriter := getBaseWriter(outputType)

	switch colorType {
	case fColorTypeSingle:
		//Logic to get writer for single color
		textColor, err := cmd.Flags().GetString(getFlag(fTextColor))
		if err !=nil {
			return nil, err
		}
		fgColor, err := getColor(textColor)
		if err != nil {
			return nil, err
		}
		return c.NewDefaultColor(baseWriter, *fgColor), nil
	case fColorTypeRotating:
		colors, err := cmd.Flags().GetString(getFlag(fColors))
		if err !=nil {
			return nil, err
		}
		strings.Split(colors, ",")
		var rotColorStrings []string = strings.Split(colors, ",")
		rotColors, err := RGBStringsToColors(rotColorStrings)
		if err != nil {
			return nil, err
		}
		rt, err := cmd.Flags().GetString(getFlag(fRotatingType))
		if err != nil {
			return nil, err
		}
		lengthStr, err := cmd.Flags().GetString(getFlag(fRotatingStrideLength))
		if err != nil {
			return nil, err
		}
		var strideLen func() int
		switch rt{
		case fRotatingFixed:
			i64, err := strconv.ParseInt(lengthStr, 0, 32)
			if err != nil {
				return nil, err
			}
			strideLen = c.NewFixedStrideLengthFunc(int(i64))
		case fRotatingRandom:
			lengthStrParts := strings.Split(lengthStr, "-")
			if len(lengthStrParts) != 2 {
				return nil, fmt.Errorf("when using flag %s you should specify length as min-max", fRotatingRandom)
			}
			i64min, err := strconv.ParseInt(lengthStrParts[0], 0, 32)
			if err != nil {
				return nil, err
			}
			i64max, err := strconv.ParseInt(lengthStrParts[1], 0, 32)
			if err != nil {
				return nil, err
			}
			strideLen = c.NewRandomStrideLengthFunc(int(i64min), int(i64max))
		}
		return c.NewRotatingColor(baseWriter, rotColors, strideLen), nil
	case fColorTypeJSON:
		colors, err := cmd.Flags().GetString(getFlag(fColors))
		if err !=nil {
			return nil, err
		}
		jsonKey, err := cmd.Flags().GetString(getFlag(fJSONKey))
		if err != nil {
			return nil, err
		}
		var colorStrings []string = strings.Split(colors, ",")
		var jColors []jsonwriter.JSONColor = make([]jsonwriter.JSONColor, len(colorStrings))
		for i, colorString := range colorStrings {
			colorStringParts := strings.Split(colorString, ".")
			colorDotParts := len(colorStringParts)
			if colorDotParts < 4 {
				return nil, fmt.Errorf("invalid JSON color string should be value.R.G.B got %s", colorString)
			}
			c, err := RGBValuesToColor(colorStringParts[colorDotParts-3:colorDotParts])
			if err != nil {
				return nil, fmt.Errorf("invalid JSON color string RGB value got %v from %s", colorStringParts[colorDotParts-3:colorDotParts] ,colorString)	
			}
			value := strings.Join(colorStringParts[0:colorDotParts-3],".")
			jColors[i] = jsonwriter.JSONColor{
				Key: jsonKey,
				Value: value,
				Color: c,

			}
		}
		return jsonwriter.NewJSONWriter(baseWriter, jsonwriter.NewColourDecider(jColors...)), nil
		
	default:
		return nil, fmt.Errorf("unknown color type: %s", colorType)
	}
}

func RGBValuesToColor(rgbValues []string) (*color.Color, error) {
	if len(rgbValues) != 3 {
		return nil, fmt.Errorf("invRGBValuesToColor requires 3 .-separated color values (got %d)", len(rgbValues))
	}
	r, err := strconv.ParseInt(rgbValues[0], 0, 24)
	if err != nil {
		return nil, fmt.Errorf("invalid R value: %s", rgbValues[0])
	}
	g, err := strconv.ParseInt(rgbValues[1], 0, 24)
	if err != nil {
		return nil, fmt.Errorf("invalid R value: %s", rgbValues[1])
	}
	b, err := strconv.ParseInt(rgbValues[2], 0, 24)
	if err != nil {
		return nil, fmt.Errorf("invalid R value: %s", rgbValues[2])
	}
	return color.RGB(int(r), int(g), int(b)), nil
}

func RGBStringsToColors(icolors []string) ([]*color.Color, error) {
	var result []*color.Color = make([]*color.Color, len(icolors))
	for i, clr := range icolors {
		rgbValues := strings.Split(clr, ".")
		c, err := RGBValuesToColor((rgbValues))
		if err != nil {
			return nil, fmt.Errorf("invalid RGB value: %s: %s", clr, err)
		}
		result[i] = c
	}
	return result, nil
}

func getColor(colorString string) (*color.Color, error) {
	c, ok := colourTable[colorString]
	if ok {
		return c, nil
	}
	return nil, fmt.Errorf("unknown color: %s", colorString)
}

func getBaseWriter(outputType outputType) io.Writer {
	switch outputType {
	case stdout:
		return os.Stdout
	case stderr:
		return os.Stderr
	}
	panic(fmt.Sprintf("Invalid outputType %s", outputType))
}




const fColorType = "color-type"
const fColorTypeSingle = "single"
const fColorTypeRotating = "rotating"
const fColorTypeJSON = "JSON"
const fColors = "colors"
const fColorsRainbow = "230.42.42,255.128.0,250.235.54,121.195.20,72.125.231,75.54.157,112.54.157"
const fRotatingType = "rotating-type"
const fRotatingFixed = "fixed"
const fRotatingRandom = "random"
const fRotatingStrideLength = "stride-length"
const fJSONKey = "json-key"

var fRotatingTypes = []string {
	fRotatingFixed,
	fRotatingRandom,
}

var fColorTypes = []string{
	fColorTypeSingle,
	fColorTypeRotating,
}
const fTextColor = "text-color"


var colorFlags []outErrStringFlag = []outErrStringFlag{
	{
		Name: fColorType,
		OutDefault: fColorTypeSingle,
		ErrDefault: fColorTypeSingle,
		Usage: fmt.Sprintf("The color type to use for text [%s].", strings.Join(fColorTypes, ", ")),
	},
	{
		Name: fTextColor,
		OutDefault: "white",
		ErrDefault: "red",
		Usage: fmt.Sprintf("The default color to use for text [%s].", strings.Join(availableColors(), ", ")),
	},
	{
		Name: fColors,
		OutDefault: fColorsRainbow,
		ErrDefault: fColorsRainbow,
		Usage: "The colors to use for color types with multiple colors. comma separated R.G.B values (0-255)",
	},
	{
		Name: fRotatingType,
		OutDefault: fRotatingFixed,
		ErrDefault: fRotatingFixed,
		Usage: fmt.Sprintf("The rotating type [%s].", strings.Join(fRotatingTypes, ", ")),
	},
	{
		Name: fRotatingStrideLength,
		OutDefault: "2",
		ErrDefault: "2",
		Usage: "The length used for strides of colors",
	},
	{
		Name: fJSONKey,
		OutDefault: "level",
		ErrDefault: "level",
		Usage: "The key of the JSON field that decides the color (default: level)",
	},
}

func getOutFlagName(flagName string) string {
	return flagName
}

func getErrFlagName(flagName string) string {
	return fmt.Sprintf("err-%s", flagName)
}

func getFlagNameFunc(outputType outputType) func (string) string{
	switch outputType {
	case stdout:
		return getOutFlagName
	case stderr:
		return getErrFlagName
	}
	panic(fmt.Sprintf("Unsupported outputType %s", outputType))
}

func init() {
	rootCmd.AddCommand(colorCmd)

	for _, colorFlag := range colorFlags {
		colorCmd.Flags().String(getOutFlagName(colorFlag.Name), colorFlag.OutDefault, fmt.Sprintf("%s for stdout", colorFlag.Usage))
		colorCmd.Flags().String(getErrFlagName(colorFlag.Name), colorFlag.ErrDefault, fmt.Sprintf("%s for stderr", colorFlag.Usage))

	}

	colorCmd.Flags().Bool("force", false, "Whether to force coloring regardless of type of outputstream.")

}
