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
	Type     SizeUnitType
	Value    float64
	Function SizeFunc
}

// AutoSize creates SizeUnit with Auto type
func AutoSize() SizeUnit {
	return SizeUnit{Auto, 0, nil}
}

// Px creates SizeUnit with SizeInPixel type
func Px(value float64) SizeUnit {
	return SizeUnit{SizeInPixel, value, nil}
}

// Em creates SizeUnit with SizeInEM type
func Em(value float64) SizeUnit {
	return SizeUnit{SizeInEM, value, nil}
}

// Ex creates SizeUnit with SizeInEX type
func Ex(value float64) SizeUnit {
	return SizeUnit{SizeInEX, value, nil}
}

// Percent creates SizeUnit with SizeInDIP type
func Percent(value float64) SizeUnit {
	return SizeUnit{SizeInPercent, value, nil}
}

// Pt creates SizeUnit with SizeInPt type
func Pt(value float64) SizeUnit {
	return SizeUnit{SizeInPt, value, nil}
}

// Pc creates SizeUnit with SizeInPc type
func Pc(value float64) SizeUnit {
	return SizeUnit{SizeInPc, value, nil}
}

// Mm creates SizeUnit with SizeInMM type
func Mm(value float64) SizeUnit {
	return SizeUnit{SizeInMM, value, nil}
}

// Cm creates SizeUnit with SizeInCM type
func Cm(value float64) SizeUnit {
	return SizeUnit{SizeInCM, value, nil}
}

// Inch creates SizeUnit with SizeInInch type
func Inch(value float64) SizeUnit {
	return SizeUnit{SizeInInch, value, nil}
}

// Fr creates SizeUnit with SizeInFraction type
func Fr(value float64) SizeUnit {
	return SizeUnit{SizeInFraction, value, nil}
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
