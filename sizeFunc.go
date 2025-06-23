package rui

import (
	"fmt"
	"strconv"
	"strings"
)

// SizeFunc describes a function that calculates the SizeUnit size.
//
// Used as the value of the SizeUnit properties.
//
// "min", "max", "clamp", "sum", "sub", "mul", "div", mod,
// "round", "round-up", "round-down" and "round-to-zero" functions are available.
type SizeFunc interface {
	fmt.Stringer
	// Name() returns the function name: "min", "max", "clamp", "sum", "sub", "mul",
	// "div", "mod", "rem", "round", "round-up", "round-down" or "round-to-zero"
	Name() string
	// Args() returns a list of function arguments
	Args() []any
	cssString(session Session) string
	writeCSS(topFunc string, buffer *strings.Builder, session Session)
	writeString(topFunc string, buffer *strings.Builder)
}

type sizeFuncData struct {
	tag  string
	args []any
}

func parseSizeFunc(text string) SizeFunc {
	text = strings.Trim(text, " ")

	for _, tag := range []string{
		"min", "max", "sum", "sub", "mul", "div", "mod", "rem", "clamp",
		"round-up", "round-down", "round-to-zero", "round"} {
		if strings.HasPrefix(text, tag) {
			text = strings.Trim(strings.TrimPrefix(text, tag), " ")
			last := len(text) - 1
			if text[0] == '(' && text[last] == ')' {
				text = text[1:last]
				bracket := 0
				start := 0
				args := []any{}
				for i, ch := range text {
					switch ch {
					case ',':
						if bracket == 0 {
							args = append(args, text[start:i])
							start = i + 1
						}

					case '(':
						bracket++

					case ')':
						bracket--
					}
				}
				if bracket != 0 {
					ErrorLogF(`Invalid "%s" function`, tag)
					return nil
				}

				args = append(args, text[start:])
				switch tag {
				case "sub", "mul", "div", "mod", "rem", "round-up", "round-down", "round-to-zero", "round":
					if len(args) != 2 {
						ErrorLogF(`"%s" function needs 2 arguments`, tag)
						return nil
					}
				case "clamp":
					if len(args) != 3 {
						ErrorLog(`"clamp" function needs 3 arguments`)
						return nil
					}
				}

				data := new(sizeFuncData)
				data.tag = tag
				if data.parseArgs(args, tag == "mul" || tag == "div" || tag == "mod" ||
					tag == "rem" || tag == "round-up" || tag == "round-down" ||
					tag == "round-to-zero" || tag == "round") {
					return data
				}
			}

			ErrorLogF(`Invalid "%s" function`, tag)
			return nil
		}
	}

	return nil
}

func (data *sizeFuncData) parseArgs(args []any, allowNumber bool) bool {
	data.args = []any{}

	numberArg := func(index int, value float64) bool {
		if allowNumber {
			if index == 1 {
				if value == 0 {
					if data.tag == "div" || data.tag == "mod" {
						ErrorLogF(`Division by 0 in "%s" function`, data.tag)
						return false
					}
					if data.tag == "round" || data.tag == "round-up" ||
						data.tag == "round-down" || data.tag == "round-to-zero" {
						ErrorLogF(`The rounding interval is 0 in "%s" function`, data.tag)
						return false
					}
				}
				data.args = append(data.args, value)
				return true
			} else {
				ErrorLogF(`Only the second %s function argument can be a number`, data.tag)
			}
		} else {
			ErrorLogF(`The %s function argument can't be a number`, data.tag)
		}
		return false
	}

	for i, arg := range args {
		switch arg := arg.(type) {
		case string:
			if arg = strings.Trim(arg, " \t\n"); arg == "" {
				ErrorLogF(`Unsupported %s function argument #%d: ""`, data.tag, i)
				return false
			}

			if arg[0] == '@' {
				data.args = append(data.args, arg)
			} else if val, err := strconv.ParseFloat(arg, 64); err == nil {
				return numberArg(i, val)
			} else if fn := parseSizeFunc(arg); fn != nil {
				data.args = append(data.args, fn)
			} else if size, err := stringToSizeUnit(arg); err == nil {
				data.args = append(data.args, size)
			} else {
				ErrorLogF(`Unsupported %s function argument #%d: "%s"`, data.tag, i, arg)
				return false
			}

		case SizeFunc:
			data.args = append(data.args, arg)

		case SizeUnit:
			if arg.Type == Auto {
				ErrorLogF(`Unsupported %s function argument #%d: "auto"`, data.tag, i)
			}
			data.args = append(data.args, arg)

		case float64:
			return numberArg(i, arg)

		case float32:
			return numberArg(i, float64(arg))

		default:
			if n, ok := isInt(arg); ok {
				return numberArg(i, float64(n))
			}
			ErrorLogF(`Unsupported %s function argument #%d: %v`, data.tag, i, arg)
			return false
		}
	}
	return true
}

func (data *sizeFuncData) String() string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	data.writeString("", buffer)
	return buffer.String()
}

func (data *sizeFuncData) Name() string {
	return data.tag
}

func (data *sizeFuncData) Args() []any {
	args := make([]any, len(data.args))
	copy(args, data.args)
	return args
}

func (data *sizeFuncData) writeString(topFunc string, buffer *strings.Builder) {
	buffer.WriteString(data.tag)
	buffer.WriteRune('(')
	for i, arg := range data.args {
		if i > 0 {
			buffer.WriteString(", ")
		}
		switch arg := arg.(type) {
		case string:
			buffer.WriteString(arg)

		case SizeFunc:
			arg.writeString(data.tag, buffer)

		case SizeUnit:
			buffer.WriteString(arg.String())

		case fmt.Stringer:
			buffer.WriteString(arg.String())

		case float64:
			fmt.Fprintf(buffer, "%g", arg)
		}

	}
	buffer.WriteRune(')')
}

func (data *sizeFuncData) cssString(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	data.writeCSS("", buffer, session)
	return buffer.String()
}

func (data *sizeFuncData) writeCSS(topFunc string, buffer *strings.Builder, session Session) {
	bracket := true
	sep := ", "

	mathFunc := func(s string) {
		sep = s
		switch topFunc {
		case "":
			buffer.WriteString("calc(")

		case "min", "max", "clamp", "mod", "rem", "round", "round-up", "round-down", "round-to-zero":
			bracket = false

		default:
			buffer.WriteRune('(')
		}
	}

	switch data.tag {
	case "min", "max", "clamp", "mod", "rem":
		buffer.WriteString(data.tag)
		buffer.WriteRune('(')

	case "round":
		buffer.WriteString("round(nearest, ")

	case "round-up":
		buffer.WriteString("round(up, ")

	case "round-down":
		buffer.WriteString("round(down, ")

	case "round-to-zero":
		buffer.WriteString("round(to-zero, ")

	case "sum":
		mathFunc(" + ")

	case "sub":
		mathFunc(" - ")

	case "mul":
		mathFunc(" * ")

	case "div":
		mathFunc(" / ")

	default:
		return
	}

	for i, arg := range data.args {
		if i > 0 {
			buffer.WriteString(sep)
		}
		switch arg := arg.(type) {
		case string:
			if arg, ok := session.resolveConstants(arg); ok {
				if fn := parseSizeFunc(arg); fn != nil {
					fn.writeCSS(data.tag, buffer, session)
				} else if size, err := stringToSizeUnit(arg); err == nil {
					buffer.WriteString(size.cssString("0", session))
				} else {
					buffer.WriteString("0")
				}
			} else {
				buffer.WriteString("0")
			}

		case SizeFunc:
			arg.writeCSS(data.tag, buffer, session)

		case SizeUnit:
			buffer.WriteString(arg.cssString("0", session))

		case fmt.Stringer:
			buffer.WriteString(arg.String())

		case float64:
			fmt.Fprintf(buffer, "%g", arg)
		}

	}

	if bracket {
		buffer.WriteRune(')')
	}
}

// MaxSize creates a SizeUnit function that calculates the maximum argument.
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc
func MaxSize(arg0, arg1 any, args ...any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "max"
	if !data.parseArgs(append([]any{arg0, arg1}, args...), false) {
		return nil
	}
	return data
}

// MinSize creates a SizeUnit function that calculates the minimum argument.
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
func MinSize(arg0, arg1 any, args ...any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "min"
	if !data.parseArgs(append([]any{arg0, arg1}, args...), false) {
		return nil
	}
	return data
}

// SumSize creates a SizeUnit function that calculates the sum of arguments.
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
func SumSize(arg0, arg1 any, args ...any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "sum"
	if !data.parseArgs(append([]any{arg0, arg1}, args...), false) {
		return nil
	}
	return data
}

// SumSize creates a SizeUnit function that calculates the result of subtracting the arguments (arg1 - arg2).
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
func SubSize(arg0, arg1 any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "sub"
	if !data.parseArgs([]any{arg0, arg1}, false) {
		return nil
	}
	return data
}

// MulSize creates a SizeUnit function that calculates the result of multiplying the arguments (arg1 * arg2).
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func MulSize(arg0, arg1 any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "mul"
	if !data.parseArgs([]any{arg0, arg1}, true) {
		return nil
	}
	return data
}

// DivSize creates a SizeUnit function that calculates the result of dividing the arguments (arg1 / arg2).
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func DivSize(arg0, arg1 any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "div"
	if !data.parseArgs([]any{arg0, arg1}, true) {
		return nil
	}
	return data
}

// RemSize creates a SizeUnit function that calculates the remainder of a division operation
// with the same sign as the dividend (arg1 % arg2).
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func RemSize(arg0, arg1 any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "rem"
	if !data.parseArgs([]any{arg0, arg1}, true) {
		return nil
	}
	return data
}

// ModSize creates a SizeUnit function that calculates the remainder of a division operation
// with the same sign as the divisor (arg1 % arg2).
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func ModSize(arg0, arg1 any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "mod"
	if !data.parseArgs([]any{arg0, arg1}, true) {
		return nil
	}
	return data
}

// RoundSize creates a SizeUnit function that calculates a rounded number.
// The function rounds valueToRound (first argument) to the nearest integer multiple
// of roundingInterval (second argument), which may be either above or below the value.
// If the valueToRound is half way between the rounding targets above and below (neither is "nearest"), it will be rounded up.
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func RoundSize(valueToRound, roundingInterval any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "round"
	if !data.parseArgs([]any{valueToRound, roundingInterval}, true) {
		return nil
	}
	return data
}

// RoundUpSize creates a SizeUnit function that calculates a rounded number.
// The function rounds valueToRound (first argument) up to the nearest integer multiple
// of roundingInterval (second argument) (if the value is negative, it will become "more positive").
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func RoundUpSize(valueToRound, roundingInterval any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "round-up"
	if !data.parseArgs([]any{valueToRound, roundingInterval}, true) {
		return nil
	}
	return data
}

// RoundDownSize creates a SizeUnit function that calculates a rounded number.
// The function rounds valueToRound (first argument) down to the nearest integer multiple
// of roundingInterval (second argument) (if the value is negative, it will become "more negative").
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func RoundDownSize(valueToRound, roundingInterval any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "round-down"
	if !data.parseArgs([]any{valueToRound, roundingInterval}, true) {
		return nil
	}
	return data
}

// RoundToZeroSize creates a SizeUnit function that calculates a rounded number.
// The function rounds valueToRound (first argument) to the nearest integer multiple
// of roundingInterval (second argument), which may be either above or below the value.
// If the valueToRound is half way between the rounding targets above and below.
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
// The second argument can also be a number (float32, float32, int, int8...int64, uint, uint8...unit64)
// or a string which is a text representation of a number.
func RoundToZeroSize(valueToRound, roundingInterval any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "round-to-zero"
	if !data.parseArgs([]any{valueToRound, roundingInterval}, true) {
		return nil
	}
	return data
}

// ClampSize creates a SizeUnit function whose the result is calculated as follows:
//
//	min ≤ value ≤ max -> value;
//	value < min -> min;
//	max < value -> max;
//
// Valid arguments types are SizeUnit, SizeFunc and a string which is a text description of SizeUnit or SizeFunc.
func ClampSize(min, value, max any) SizeFunc {
	data := new(sizeFuncData)
	data.tag = "clamp"
	if !data.parseArgs([]any{min, value, max}, false) {
		return nil
	}
	return data
}
