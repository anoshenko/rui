package rui

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Color - represent color in argb format
type Color uint32

// ARGB - return alpha, red, green and blue components of the color
func (color Color) ARGB() (uint8, uint8, uint8, uint8) {
	return uint8(color >> 24),
		uint8((color >> 16) & 0xFF),
		uint8((color >> 8) & 0xFF),
		uint8(color & 0xFF)
}

// Alpha - return the alpha component of the color
func (color Color) Alpha() int {
	return int((color >> 24) & 0xFF)
}

// Red - return the red component of the color
func (color Color) Red() int {
	return int((color >> 16) & 0xFF)
}

// Green - return the green component of the color
func (color Color) Green() int {
	return int((color >> 8) & 0xFF)
}

// Blue - return the blue component of the color
func (color Color) Blue() int {
	return int(color & 0xFF)
}

// String get a text representation of the color
func (color Color) String() string {
	return fmt.Sprintf("#%08X", int(color))
}

func (color Color) rgbString() string {
	return fmt.Sprintf("#%06X", int(color&0xFFFFFF))
}

// writeData write a text representation of the color to the buffer
func (color Color) writeData(buffer *bytes.Buffer) {
	buffer.WriteString(color.String())
}

// cssString get the text representation of the color in CSS format
func (color Color) cssString() string {
	red := color.Red()
	green := color.Green()
	blue := color.Blue()

	if alpha := color.Alpha(); alpha < 255 {
		aText := fmt.Sprintf("%.2f", float64(alpha)/255.0)
		if len(aText) > 1 {
			aText = aText[1:]
		}
		return fmt.Sprintf("rgba(%d,%d,%d,%s)", red, green, blue, aText)
	}

	return fmt.Sprintf("rgb(%d,%d,%d)", red, green, blue)
}

// StringToColor converts the string argument to Color value
func StringToColor(text string) (Color, bool) {

	text = strings.Trim(text, " \t\r\n")
	if text == "" {
		ErrorLog(`Invalid color value: ""`)
		return 0, false
	}

	if text[0] == '#' {
		c, err := strconv.ParseUint(text[1:], 16, 32)
		if err != nil {
			ErrorLog("Set color value error: " + err.Error())
			return 0, false
		}

		switch len(text) - 1 {
		case 8:
			return Color(c), true

		case 6:
			return Color(c | 0xFF000000), true

		case 4:
			a := (c >> 12) & 0xF
			r := (c >> 8) & 0xF
			g := (c >> 4) & 0xF
			b := c & 0xF
			return Color((a << 28) | (a << 24) | (r << 20) | (r << 16) | (g << 12) | (g << 8) | (b << 4) | b), true

		case 3:
			r := (c >> 8) & 0xF
			g := (c >> 4) & 0xF
			b := c & 0xF
			return Color(0xFF000000 | (r << 20) | (r << 16) | (g << 12) | (g << 8) | (b << 4) | b), true
		}

		ErrorLog(`Invalid color format: "` + text + `". Valid formats: #AARRGGBB, #RRGGBB, #ARGB, #RGB`)
		return 0, false
	}

	parseRGB := func(args string) []int {
		args = strings.Trim(args, " \t")
		count := len(args)
		if count < 3 || args[0] != '(' || args[count-1] != ')' {
			return []int{}
		}

		arg := strings.Split(args[1:count-1], ",")
		result := make([]int, len(arg))
		for i, val := range arg {
			val = strings.Trim(val, " \t")
			size := len(val)
			if size == 0 {
				return []int{}
			}
			if val[size-1] == '%' {
				if n, err := strconv.Atoi(val[:size-1]); err == nil && n >= 0 && n <= 100 {
					result[i] = n * 255 / 100
				} else {
					return []int{}
				}
			} else if strings.ContainsRune(val, '.') {
				if val[0] == '.' {
					val = "0" + val
				}
				if f, err := strconv.ParseFloat(val, 32); err == nil && f >= 0 && f <= 1 {
					result[i] = int(f * 255)
				} else {
					return []int{}
				}
			} else {
				if n, err := strconv.Atoi(val); err == nil && n >= 0 && n <= 255 {
					result[i] = n
				} else {
					return []int{}
				}
			}
		}
		return result
	}

	text = strings.ToLower(text)
	if strings.HasPrefix(text, "rgba") {
		args := parseRGB(text[4:])
		if len(args) == 4 {
			return Color((args[3] << 24) | (args[0] << 16) | (args[1] << 8) | args[2]), true
		}
	}

	if strings.HasPrefix(text, "rgb") {
		args := parseRGB(text[3:])
		if len(args) == 3 {
			return Color(0xFF000000 | (args[0] << 16) | (args[1] << 8) | args[2]), true
		}
	}

	// TODO hsl(360,100%,50%), hsla(360,100%,50%,.5)

	if color, ok := colorConstants[text]; ok {
		return color, true
	}

	ErrorLog(`Invalid color format: "` + text + `"`)
	return 0, false
}
