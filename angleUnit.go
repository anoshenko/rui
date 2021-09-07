package rui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// AngleUnitType : type of enumerated constants for define a type of AngleUnit value.
// Can take the following values: Radian, Degree, Gradian, and Turn
type AngleUnitType uint8

const (
	// Radian - angle in radians
	Radian AngleUnitType = 0
	// Radian - angle in radians * π
	PiRadian AngleUnitType = 1
	// Degree - angle in degrees
	Degree AngleUnitType = 2
	// Gradian - angle in gradian (1⁄400 of a full circle)
	Gradian AngleUnitType = 3
	// Turn - angle in turns (1 turn = 360 degree)
	Turn AngleUnitType = 4
)

// AngleUnit describe a size (Value field) and size unit (Type field).
type AngleUnit struct {
	Type  AngleUnitType
	Value float64
}

// Deg creates AngleUnit with Degree type
func Deg(value float64) AngleUnit {
	return AngleUnit{Type: Degree, Value: value}
}

// Rad create AngleUnit with Radian type
func Rad(value float64) AngleUnit {
	return AngleUnit{Type: Radian, Value: value}
}

// PiRad create AngleUnit with PiRadian type
func PiRad(value float64) AngleUnit {
	return AngleUnit{Type: PiRadian, Value: value}
}

// Grad create AngleUnit with Gradian type
func Grad(value float64) AngleUnit {
	return AngleUnit{Type: Gradian, Value: value}
}

// Equal compare two AngleUnit. Return true if AngleUnit are equal
func (angle AngleUnit) Equal(size2 AngleUnit) bool {
	return angle.Type == size2.Type && angle.Value == size2.Value
}

func angleUnitSuffixes() map[AngleUnitType]string {
	return map[AngleUnitType]string{
		Degree:   "deg",
		Radian:   "rad",
		PiRadian: "pi",
		Gradian:  "grad",
		Turn:     "turn",
	}
}

// StringToAngleUnit converts the string argument to AngleUnit
func StringToAngleUnit(value string) (AngleUnit, bool) {
	var angle AngleUnit
	ok, err := angle.setValue(value)
	if !ok {
		ErrorLog(err)
	}
	return angle, ok
}

func (angle *AngleUnit) setValue(value string) (bool, string) {
	value = strings.ToLower(strings.Trim(value, " \t\n\r"))

	setValue := func(suffix string, unitType AngleUnitType) (bool, string) {
		val, err := strconv.ParseFloat(value[:len(value)-len(suffix)], 64)
		if err != nil {
			return false, `AngleUnit.SetValue("` + value + `") error: ` + err.Error()
		}
		angle.Value = val
		angle.Type = unitType
		return true, ""
	}

	if value == "π" {
		angle.Value = 1
		angle.Type = PiRadian
		return true, ""
	}

	if strings.HasSuffix(value, "π") {
		return setValue("π", PiRadian)
	}

	if strings.HasSuffix(value, "°") {
		return setValue("°", Degree)
	}

	for unitType, suffix := range angleUnitSuffixes() {
		if strings.HasSuffix(value, suffix) {
			return setValue(suffix, unitType)
		}
	}

	if val, err := strconv.ParseFloat(value, 64); err == nil {
		angle.Value = val
		angle.Type = Radian
		return true, ""
	}

	return false, `AngleUnit.SetValue("` + value + `") error: invalid argument`
}

// String - convert AngleUnit to string
func (angle AngleUnit) String() string {
	if suffix, ok := angleUnitSuffixes()[angle.Type]; ok {
		return fmt.Sprintf("%g%s", angle.Value, suffix)
	}

	return fmt.Sprintf("%g", angle.Value)
}

// cssString - convert AngleUnit to string
func (angle AngleUnit) cssString() string {
	if angle.Type == PiRadian {
		return fmt.Sprintf("%grad", angle.Value*math.Pi)
	}

	return angle.String()
}

// ToDegree returns the angle in radians
func (angle AngleUnit) ToRadian() AngleUnit {
	switch angle.Type {
	case PiRadian:
		return AngleUnit{Value: angle.Value * math.Pi, Type: Radian}

	case Degree:
		return AngleUnit{Value: angle.Value * math.Pi / 180, Type: Radian}

	case Gradian:
		return AngleUnit{Value: angle.Value * math.Pi / 200, Type: Radian}

	case Turn:
		return AngleUnit{Value: angle.Value * 2 * math.Pi, Type: Radian}
	}

	return angle
}

// ToDegree returns the angle in degrees
func (angle AngleUnit) ToDegree() AngleUnit {
	switch angle.Type {
	case Radian:
		return AngleUnit{Value: angle.Value * 180 / math.Pi, Type: Degree}

	case PiRadian:
		return AngleUnit{Value: angle.Value * 180, Type: Degree}

	case Gradian:
		return AngleUnit{Value: angle.Value * 360 / 400, Type: Degree}

	case Turn:
		return AngleUnit{Value: angle.Value * 360, Type: Degree}
	}

	return angle
}

// ToGradian returns the angle in gradians (1⁄400 of a full circle)
func (angle AngleUnit) ToGradian() AngleUnit {
	switch angle.Type {
	case Radian:
		return AngleUnit{Value: angle.Value * 200 / math.Pi, Type: Gradian}

	case PiRadian:
		return AngleUnit{Value: angle.Value * 200, Type: Gradian}

	case Degree:
		return AngleUnit{Value: angle.Value * 400 / 360, Type: Gradian}

	case Turn:
		return AngleUnit{Value: angle.Value * 400, Type: Gradian}
	}

	return angle
}

// ToTurn returns the angle in turns (1 turn = 360 degree)
func (angle AngleUnit) ToTurn() AngleUnit {
	switch angle.Type {
	case Radian:
		return AngleUnit{Value: angle.Value / (2 * math.Pi), Type: Turn}

	case PiRadian:
		return AngleUnit{Value: angle.Value / 2, Type: Turn}

	case Degree:
		return AngleUnit{Value: angle.Value / 360, Type: Turn}

	case Gradian:
		return AngleUnit{Value: angle.Value / 400, Type: Turn}
	}

	return angle
}
