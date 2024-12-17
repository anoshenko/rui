package rui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SizeUnitType : type of enumerated constants for define a type of SizeUnit value.
//
// Can take the following values: Auto, SizeInPixel, SizeInPercent,
// SizeInDIP, SizeInPt, SizeInInch, SizeInMM, SizeInFraction
type SizeUnitType uint8

// Constants which represent values of a [SizeUnitType]
const (
	// Auto is the SizeUnit type: default value.
	Auto SizeUnitType = 0
	// SizeInPixel is the SizeUnit type: the Value field specifies the size in pixels.
	SizeInPixel SizeUnitType = 1
	// SizeInEM is the SizeUnit type: the Value field specifies the size in em.
	SizeInEM SizeUnitType = 2
	// SizeInEX is the SizeUnit type: the Value field specifies the size in em.
	SizeInEX SizeUnitType = 3
	// SizeInPercent is the SizeUnit type: the Value field specifies the size in percents of the parent size.
	SizeInPercent SizeUnitType = 4
	// SizeInPt is the SizeUnit type: the Value field specifies the size in pt (1/72 inch).
	SizeInPt SizeUnitType = 5
	// SizeInPc is the SizeUnit type: the Value field specifies the size in pc (1pc = 12pt).
	SizeInPc SizeUnitType = 6
	// SizeInInch is the SizeUnit type: the Value field specifies the size in inches.
	SizeInInch SizeUnitType = 7
	// SizeInMM is the SizeUnit type: the Value field specifies the size in millimeters.
	SizeInMM SizeUnitType = 8
	// SizeInCM is the SizeUnit type: the Value field specifies the size in centimeters.
	SizeInCM SizeUnitType = 9
	// SizeInFraction is the SizeUnit type: the Value field specifies the size in fraction.
	// Used only for "cell-width" and "cell-height" property.
	SizeInFraction SizeUnitType = 10
	// SizeFunction is the SizeUnit type: the Function field specifies the size function.
	// "min", "max", "clamp", "sum", "sub", "mul", and "div" functions are available.
	SizeFunction = 11
)

// SizeUnit describe a size (Value field) and size unit (Type field).
type SizeUnit struct {
	// Type or dimension of the value
	Type SizeUnitType

	// Value of the size in Type units
	Value float64

	// Function representation of a size unit.
	// When setting this value type should be set to SizeFunction
	Function SizeFunc
}

// AutoSize creates SizeUnit with Auto type
func AutoSize() SizeUnit {
	return SizeUnit{Auto, 0, nil}
}

// Px creates SizeUnit with SizeInPixel type
func Px[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInPixel, Value: float64(value), Function: nil}
}

// Em creates SizeUnit with SizeInEM type
func Em[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInEM, Value: float64(value), Function: nil}
}

// Ex creates SizeUnit with SizeInEX type
func Ex[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInEX, Value: float64(value), Function: nil}
}

// Percent creates SizeUnit with SizeInDIP type
func Percent[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInPercent, Value: float64(value), Function: nil}
}

// Pt creates SizeUnit with SizeInPt type
func Pt[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInPt, Value: float64(value), Function: nil}
}

// Pc creates SizeUnit with SizeInPc type
func Pc[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInPc, Value: float64(value), Function: nil}
}

// Mm creates SizeUnit with SizeInMM type
func Mm[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInMM, Value: float64(value), Function: nil}
}

// Cm creates SizeUnit with SizeInCM type
func Cm[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInCM, Value: float64(value), Function: nil}
}

// Inch creates SizeUnit with SizeInInch type
func Inch[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{Type: SizeInInch, Value: float64(value), Function: nil}
}

// Fr creates SizeUnit with SizeInFraction type
func Fr[T float64 | float32 | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](value T) SizeUnit {
	return SizeUnit{SizeInFraction, float64(value), nil}
}

// Equal compare two SizeUnit. Return true if SizeUnit are equal
func (size SizeUnit) Equal(size2 SizeUnit) bool {
	return size.Type == size2.Type && (size.Type == Auto || size.Value == size2.Value)
}

func sizeUnitSuffixes() map[SizeUnitType]string {
	return map[SizeUnitType]string{
		SizeInPixel:    "px",
		SizeInPercent:  "%",
		SizeInEM:       "em",
		SizeInEX:       "ex",
		SizeInPt:       "pt",
		SizeInPc:       "pc",
		SizeInInch:     "in",
		SizeInMM:       "mm",
		SizeInCM:       "cm",
		SizeInFraction: "fr",
	}
}

// StringToSizeUnit converts the string argument to SizeUnit
func StringToSizeUnit(value string) (SizeUnit, bool) {
	size, err := stringToSizeUnit(value)
	if err != nil {
		ErrorLog(err.Error())
		return size, false
	}
	return size, true
}

func stringToSizeUnit(value string) (SizeUnit, error) {
	value = strings.Trim(value, " \t\n\r")

	switch value {
	case "auto", "none", "":
		return SizeUnit{Type: Auto, Value: 0}, nil

	case "0":
		return SizeUnit{Type: SizeInPixel, Value: 0}, nil
	}

	suffixes := sizeUnitSuffixes()
	for unitType, suffix := range suffixes {
		if strings.HasSuffix(value, suffix) {
			var err error
			var val float64
			if val, err = strconv.ParseFloat(value[:len(value)-len(suffix)], 64); err != nil {
				return SizeUnit{Type: Auto, Value: 0}, err
			}
			return SizeUnit{Type: unitType, Value: val}, nil
		}
	}

	if val, err := strconv.ParseFloat(value, 64); err == nil {
		return SizeUnit{Type: SizeInPixel, Value: val}, nil
	}

	return SizeUnit{Type: Auto, Value: 0}, errors.New(`Invalid SizeUnit value: "` + value + `"`)
}

// String - convert SizeUnit to string
func (size SizeUnit) String() string {
	switch size.Type {
	case Auto:
		return "auto"

	case SizeFunction:
		if size.Function == nil {
			return "auto"
		}
		return size.Function.String()
	}
	if suffix, ok := sizeUnitSuffixes()[size.Type]; ok {
		return fmt.Sprintf("%g%s", size.Value, suffix)
	}
	return strconv.FormatFloat(size.Value, 'g', -1, 64)
}

// cssString - convert SizeUnit to string
func (size SizeUnit) cssString(textForAuto string, session Session) string {
	switch size.Type {
	case Auto:
		return textForAuto

	case SizeInEM:
		return fmt.Sprintf("%grem", size.Value)

	case SizeFunction:
		if size.Function == nil {
			return textForAuto
		}
		return size.Function.cssString(session)
	}

	if size.Value == 0 {
		return "0"
	}

	return size.String()
}
