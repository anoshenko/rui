package rui

import (
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
	// Auto - default value.
	Auto SizeUnitType = 0
	// SizeInPixel - size in pixels.
	SizeInPixel SizeUnitType = 1
	// SizeInEM - size in em.
	SizeInEM SizeUnitType = 2
	// SizeInEX - size in em.
	SizeInEX SizeUnitType = 3
	// SizeInPercent - size in percents of a parant size.
	SizeInPercent SizeUnitType = 4
	// SizeInPt - size in pt (1/72 inch).
	SizeInPt SizeUnitType = 5
	// SizeInPc - size in pc (1pc = 12pt).
	SizeInPc SizeUnitType = 6
	// SizeInInch - size in inches.
	SizeInInch SizeUnitType = 7
	// SizeInMM - size in millimeters.
	SizeInMM SizeUnitType = 8
	// SizeInCM - size in centimeters.
	SizeInCM SizeUnitType = 9
	// SizeInFraction - size in fraction. Used only for "cell-width" and "cell-height" property
	SizeInFraction SizeUnitType = 10
)

// SizeUnit describe a size (Value field) and size unit (Type field).
type SizeUnit struct {
	Type  SizeUnitType
	Value float64
}

// AutoSize creates SizeUnit with Auto type
func AutoSize() SizeUnit {
	return SizeUnit{Auto, 0}
}

// Px creates SizeUnit with SizeInPixel type
func Px(value float64) SizeUnit {
	return SizeUnit{SizeInPixel, value}
}

// Em creates SizeUnit with SizeInEM type
func Em(value float64) SizeUnit {
	return SizeUnit{SizeInEM, value}
}

// Ex creates SizeUnit with SizeInEX type
func Ex(value float64) SizeUnit {
	return SizeUnit{SizeInEX, value}
}

// Percent creates SizeUnit with SizeInDIP type
func Percent(value float64) SizeUnit {
	return SizeUnit{SizeInPercent, value}
}

// Pt creates SizeUnit with SizeInPt type
func Pt(value float64) SizeUnit {
	return SizeUnit{SizeInPt, value}
}

// Pc creates SizeUnit with SizeInPc type
func Pc(value float64) SizeUnit {
	return SizeUnit{SizeInPc, value}
}

// Mm creates SizeUnit with SizeInMM type
func Mm(value float64) SizeUnit {
	return SizeUnit{SizeInMM, value}
}

// Cm creates SizeUnit with SizeInCM type
func Cm(value float64) SizeUnit {
	return SizeUnit{SizeInCM, value}
}

// Inch creates SizeUnit with SizeInInch type
func Inch(value float64) SizeUnit {
	return SizeUnit{SizeInInch, value}
}

// Fr creates SizeUnit with SizeInFraction type
func Fr(value float64) SizeUnit {
	return SizeUnit{SizeInFraction, value}
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

	value = strings.Trim(value, " \t\n\r")

	switch value {
	case "auto", "none", "":
		return SizeUnit{Type: Auto, Value: 0}, true

	case "0":
		return SizeUnit{Type: SizeInPixel, Value: 0}, true
	}

	suffixes := sizeUnitSuffixes()
	for unitType, suffix := range suffixes {
		if strings.HasSuffix(value, suffix) {
			var err error
			var val float64
			if val, err = strconv.ParseFloat(value[:len(value)-len(suffix)], 64); err != nil {
				ErrorLog(err.Error())
				return SizeUnit{Type: Auto, Value: 0}, false
			}
			return SizeUnit{Type: unitType, Value: val}, true
		}
	}

	ErrorLog(`Invalid SizeUnit value: "` + value + `"`)
	return SizeUnit{Type: Auto, Value: 0}, false
}

// String - convert SizeUnit to string
func (size SizeUnit) String() string {
	if size.Type == Auto {
		return "auto"
	}
	if suffix, ok := sizeUnitSuffixes()[size.Type]; ok {
		return fmt.Sprintf("%g%s", size.Value, suffix)
	}
	return strconv.FormatFloat(size.Value, 'g', -1, 64)
}

// cssString - convert SizeUnit to string
func (size SizeUnit) cssString(textForAuto string) string {
	switch size.Type {
	case Auto:
		return textForAuto

	case SizeInEM:
		return fmt.Sprintf("%grem", size.Value)
	}

	if size.Value == 0 {
		return "0"
	}

	return size.String()
}
