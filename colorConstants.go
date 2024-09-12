package rui

import "sort"

// A set of predefined colors used in the library
const (
	// Black color constant
	Black Color = 0xff000000
	// Silver color constant
	Silver Color = 0xffc0c0c0
	// Gray color constant
	Gray Color = 0xff808080
	// White color constant
	White Color = 0xffffffff
	// Maroon color constant
	Maroon Color = 0xff800000
	// Red color constant
	Red Color = 0xffff0000
	// Purple color constant
	Purple Color = 0xff800080
	// Fuchsia color constant
	Fuchsia Color = 0xffff00ff
	// Green color constant
	Green Color = 0xff008000
	// Lime color constant
	Lime Color = 0xff00ff00
	// Olive color constant
	Olive Color = 0xff808000
	// Yellow color constant
	Yellow Color = 0xffffff00
	// Navy color constant
	Navy Color = 0xff000080
	// Blue color constant
	Blue Color = 0xff0000ff
	// Teal color constant
	Teal Color = 0xff008080
	// Aqua color constant
	Aqua Color = 0xff00ffff
	// Orange color constant
	Orange Color = 0xffffa500
	// AliceBlue color constant
	AliceBlue Color = 0xfff0f8ff
	// AntiqueWhite color constant
	AntiqueWhite Color = 0xfffaebd7
	// Aquamarine color constant
	Aquamarine Color = 0xff7fffd4
	// Azure color constant
	Azure Color = 0xfff0ffff
	// Beige color constant
	Beige Color = 0xfff5f5dc
	// Bisque color constant
	Bisque Color = 0xffffe4c4
	// BlanchedAlmond color constant
	BlanchedAlmond Color = 0xffffebcd
	// BlueViolet color constant
	BlueViolet Color = 0xff8a2be2
	// Brown color constant
	Brown Color = 0xffa52a2a
	// BurlyWood color constant
	BurlyWood Color = 0xffdeb887
	// CadetBlue color constant
	CadetBlue Color = 0xff5f9ea0
	// Chartreuse color constant
	Chartreuse Color = 0xff7fff00
	// Chocolate color constant
	Chocolate Color = 0xffd2691e
	// Coral color constant
	Coral Color = 0xffff7f50
	// CornflowerBlue color constant
	CornflowerBlue Color = 0xff6495ed
	// CornSilk color constant
	CornSilk Color = 0xfffff8dc
	// Crimson color constant
	Crimson Color = 0xffdc143c
	// Cyan color constant
	Cyan Color = 0xff00ffff
	// DarkBlue color constant
	DarkBlue Color = 0xff00008b
	// DarkCyan color constant
	DarkCyan Color = 0xff008b8b
	// DarkGoldenRod color constant
	DarkGoldenRod Color = 0xffb8860b
	// DarkGray color constant
	DarkGray Color = 0xffa9a9a9
	// DarkGreen color constant
	DarkGreen Color = 0xff006400
	// DarkGrey color constant
	DarkGrey Color = 0xffa9a9a9
	// DarkKhaki color constant
	DarkKhaki Color = 0xffbdb76b
	// DarkMagenta color constant
	DarkMagenta Color = 0xff8b008b
	// DarkOliveGreen color constant
	DarkOliveGreen Color = 0xff556b2f
	// DarkOrange color constant
	DarkOrange Color = 0xffff8c00
	// DarkOrchid color constant
	DarkOrchid Color = 0xff9932cc
	// DarkRed color constant
	DarkRed Color = 0xff8b0000
	// DarkSalmon color constant
	DarkSalmon Color = 0xffe9967a
	// DarkSeaGreen color constant
	DarkSeaGreen Color = 0xff8fbc8f
	// DarkSlateBlue color constant
	DarkSlateBlue Color = 0xff483d8b
	// DarkSlateGray color constant
	DarkSlateGray Color = 0xff2f4f4f
	// DarkSlateGrey color constant
	DarkSlateGrey Color = 0xff2f4f4f
	// DarkTurquoise color constant
	DarkTurquoise Color = 0xff00ced1
	// DarkViolet color constant
	DarkViolet Color = 0xff9400d3
	// DeepPink color constant
	DeepPink Color = 0xffff1493
	// DeepSkyBlue color constant
	DeepSkyBlue Color = 0xff00bfff
	// DimGray color constant
	DimGray Color = 0xff696969
	// DimGrey color constant
	DimGrey Color = 0xff696969
	// DodgerBlue color constant
	DodgerBlue Color = 0xff1e90ff
	// FireBrick color constant
	FireBrick Color = 0xffb22222
	// FloralWhite color constant
	FloralWhite Color = 0xfffffaf0
	// ForestGreen color constant
	ForestGreen Color = 0xff228b22
	// Gainsboro color constant
	Gainsboro Color = 0xffdcdcdc
	// GhostWhite color constant
	GhostWhite Color = 0xfff8f8ff
	// Gold color constant
	Gold Color = 0xffffd700
	// GoldenRod color constant
	GoldenRod Color = 0xffdaa520
	// GreenYellow color constant
	GreenYellow Color = 0xffadff2f
	// Grey color constant
	Grey Color = 0xff808080
	// Honeydew color constant
	Honeydew Color = 0xfff0fff0
	// HotPink color constant
	HotPink Color = 0xffff69b4
	// IndianRed color constant
	IndianRed Color = 0xffcd5c5c
	// Indigo color constant
	Indigo Color = 0xff4b0082
	// Ivory color constant
	Ivory Color = 0xfffffff0
	// Khaki color constant
	Khaki Color = 0xfff0e68c
	// Lavender color constant
	Lavender Color = 0xffe6e6fa
	// LavenderBlush color constant
	LavenderBlush Color = 0xfffff0f5
	// LawnGreen color constant
	LawnGreen Color = 0xff7cfc00
	// LemonChiffon color constant
	LemonChiffon Color = 0xfffffacd
	// LightBlue color constant
	LightBlue Color = 0xffadd8e6
	// LightCoral color constant
	LightCoral Color = 0xfff08080
	// LightCyan color constant
	LightCyan Color = 0xffe0ffff
	// LightGoldenrodYellow color constant
	LightGoldenRodYellow Color = 0xfffafad2
	// LightGray color constant
	LightGray Color = 0xffd3d3d3
	// LightGreen color constant
	LightGreen Color = 0xff90ee90
	// LightGrey color constant
	LightGrey Color = 0xffd3d3d3
	// LightPink color constant
	LightPink Color = 0xffffb6c1
	// LightSalmon color constant
	LightSalmon Color = 0xffffa07a
	// LightSeaGreen color constant
	LightSeaGreen Color = 0xff20b2aa
	// LightSkyBlue color constant
	LightSkyBlue Color = 0xff87cefa
	// LightSlateGray color constant
	LightSlateGray Color = 0xff778899
	// LightSlateGrey color constant
	LightSlateGrey Color = 0xff778899
	// LightSteelBlue color constant
	LightSteelBlue Color = 0xffb0c4de
	// LightYellow color constant
	LightYellow Color = 0xffffffe0
	// LimeGreen color constant
	LimeGreen Color = 0xff32cd32
	// Linen color constant
	Linen Color = 0xfffaf0e6
	// Magenta color constant
	Magenta Color = 0xffff00ff
	// MediumAquamarine color constant
	MediumAquamarine Color = 0xff66cdaa
	// MediumBlue color constant
	MediumBlue Color = 0xff0000cd
	// MediumOrchid color constant
	MediumOrchid Color = 0xffba55d3
	// MediumPurple color constant
	MediumPurple Color = 0xff9370db
	// MediumSeaGreen color constant
	MediumSeaGreen Color = 0xff3cb371
	// MediumSlateBlue color constant
	MediumSlateBlue Color = 0xff7b68ee
	// MediumSpringGreen color constant
	MediumSpringGreen Color = 0xff00fa9a
	// MediumTurquoise color constant
	MediumTurquoise Color = 0xff48d1cc
	// MediumVioletRed color constant
	MediumVioletRed Color = 0xffc71585
	// MidnightBlue color constant
	MidnightBlue Color = 0xff191970
	// MintCream color constant
	MintCream Color = 0xfff5fffa
	// MistyRose color constant
	MistyRose Color = 0xffffe4e1
	// Moccasin color constant
	Moccasin Color = 0xffffe4b5
	// NavajoWhite color constant
	NavajoWhite Color = 0xffffdead
	// OldLace color constant
	OldLace Color = 0xfffdf5e6
	// OliveDrab color constant
	OliveDrab Color = 0xff6b8e23
	// OrangeRed color constant
	OrangeRed Color = 0xffff4500
	// Orchid color constant
	Orchid Color = 0xffda70d6
	// PaleGoldenrod color constant
	PaleGoldenrod Color = 0xffeee8aa
	// PaleGreen color constant
	PaleGreen Color = 0xff98fb98
	// PaleTurquoise color constant
	PaleTurquoise Color = 0xffafeeee
	// PaleVioletRed color constant
	PaleVioletRed Color = 0xffdb7093
	// PapayaWhip color constant
	PapayaWhip Color = 0xffffefd5
	// PeachPuff color constant
	PeachPuff Color = 0xffffdab9
	// Peru color constant
	Peru Color = 0xffcd853f
	// Pink color constant
	Pink Color = 0xffffc0cb
	// Plum color constant
	Plum Color = 0xffdda0dd
	// PowderBlue color constant
	PowderBlue Color = 0xffb0e0e6
	// RosyBrown color constant
	RosyBrown Color = 0xffbc8f8f
	// RoyalBlue color constant
	RoyalBlue Color = 0xff4169e1
	// SaddleBrown color constant
	SaddleBrown Color = 0xff8b4513
	// Salmon color constant
	Salmon Color = 0xfffa8072
	// SandyBrown color constant
	SandyBrown Color = 0xfff4a460
	// SeaGreen color constant
	SeaGreen Color = 0xff2e8b57
	// SeaShell color constant
	SeaShell Color = 0xfffff5ee
	// Sienna color constant
	Sienna Color = 0xffa0522d
	// SkyBlue color constant
	SkyBlue Color = 0xff87ceeb
	// SlateBlue color constant
	SlateBlue Color = 0xff6a5acd
	// SlateGray color constant
	SlateGray Color = 0xff708090
	// SlateGrey color constant
	SlateGrey Color = 0xff708090
	// Snow color constant
	Snow Color = 0xfffffafa
	// SpringGreen color constant
	SpringGreen Color = 0xff00ff7f
	// SteelBlue color constant
	SteelBlue Color = 0xff4682b4
	// Tan color constant
	Tan Color = 0xffd2b48c
	// Thistle color constant
	Thistle Color = 0xffd8bfd8
	// Tomato color constant
	Tomato Color = 0xffff6347
	// Turquoise color constant
	Turquoise Color = 0xff40e0d0
	// Violet color constant
	Violet Color = 0xffee82ee
	// Wheat color constant
	Wheat Color = 0xfff5deb3
	// WhiteSmoke color constant
	WhiteSmoke Color = 0xfff5f5f5
	// YellowGreen color constant
	YellowGreen Color = 0xff9acd32
)

var colorConstants = map[string]Color{
	"black":                0xff000000,
	"silver":               0xffc0c0c0,
	"gray":                 0xff808080,
	"white":                0xffffffff,
	"maroon":               0xff800000,
	"red":                  0xffff0000,
	"purple":               0xff800080,
	"fuchsia":              0xffff00ff,
	"green":                0xff008000,
	"lime":                 0xff00ff00,
	"olive":                0xff808000,
	"yellow":               0xffffff00,
	"navy":                 0xff000080,
	"blue":                 0xff0000ff,
	"teal":                 0xff008080,
	"aqua":                 0xff00ffff,
	"orange":               0xffffa500,
	"aliceblue":            0xfff0f8ff,
	"antiquewhite":         0xfffaebd7,
	"aquamarine":           0xff7fffd4,
	"azure":                0xfff0ffff,
	"beige":                0xfff5f5dc,
	"bisque":               0xffffe4c4,
	"blanchedalmond":       0xffffebcd,
	"blueviolet":           0xff8a2be2,
	"brown":                0xffa52a2a,
	"burlywood":            0xffdeb887,
	"cadetblue":            0xff5f9ea0,
	"chartreuse":           0xff7fff00,
	"chocolate":            0xffd2691e,
	"coral":                0xffff7f50,
	"cornflowerblue":       0xff6495ed,
	"cornsilk":             0xfffff8dc,
	"crimson":              0xffdc143c,
	"cyan":                 0xff00ffff,
	"darkblue":             0xff00008b,
	"darkcyan":             0xff008b8b,
	"darkgoldenrod":        0xffb8860b,
	"darkgray":             0xffa9a9a9,
	"darkgreen":            0xff006400,
	"darkgrey":             0xffa9a9a9,
	"darkkhaki":            0xffbdb76b,
	"darkmagenta":          0xff8b008b,
	"darkolivegreen":       0xff556b2f,
	"darkorange":           0xffff8c00,
	"darkorchid":           0xff9932cc,
	"darkred":              0xff8b0000,
	"darksalmon":           0xffe9967a,
	"darkseagreen":         0xff8fbc8f,
	"darkslateblue":        0xff483d8b,
	"darkslategray":        0xff2f4f4f,
	"darkslategrey":        0xff2f4f4f,
	"darkturquoise":        0xff00ced1,
	"darkviolet":           0xff9400d3,
	"deeppink":             0xffff1493,
	"deepskyblue":          0xff00bfff,
	"dimgray":              0xff696969,
	"dimgrey":              0xff696969,
	"dodgerblue":           0xff1e90ff,
	"firebrick":            0xffb22222,
	"floralwhite":          0xfffffaf0,
	"forestgreen":          0xff228b22,
	"gainsboro":            0xffdcdcdc,
	"ghostwhite":           0xfff8f8ff,
	"gold":                 0xffffd700,
	"goldenrod":            0xffdaa520,
	"greenyellow":          0xffadff2f,
	"grey":                 0xff808080,
	"honeydew":             0xfff0fff0,
	"hotpink":              0xffff69b4,
	"indianred":            0xffcd5c5c,
	"indigo":               0xff4b0082,
	"ivory":                0xfffffff0,
	"khaki":                0xfff0e68c,
	"lavender":             0xffe6e6fa,
	"lavenderblush":        0xfffff0f5,
	"lawngreen":            0xff7cfc00,
	"lemonchiffon":         0xfffffacd,
	"lightblue":            0xffadd8e6,
	"lightcoral":           0xfff08080,
	"lightcyan":            0xffe0ffff,
	"lightgoldenrodyellow": 0xfffafad2,
	"lightgray":            0xffd3d3d3,
	"lightgreen":           0xff90ee90,
	"lightgrey":            0xffd3d3d3,
	"lightpink":            0xffffb6c1,
	"lightsalmon":          0xffffa07a,
	"lightseagreen":        0xff20b2aa,
	"lightskyblue":         0xff87cefa,
	"lightslategray":       0xff778899,
	"lightslategrey":       0xff778899,
	"lightsteelblue":       0xffb0c4de,
	"lightyellow":          0xffffffe0,
	"limegreen":            0xff32cd32,
	"linen":                0xfffaf0e6,
	"magenta":              0xffff00ff,
	"mediumaquamarine":     0xff66cdaa,
	"mediumblue":           0xff0000cd,
	"mediumorchid":         0xffba55d3,
	"mediumpurple":         0xff9370db,
	"mediumseagreen":       0xff3cb371,
	"mediumslateblue":      0xff7b68ee,
	"mediumspringgreen":    0xff00fa9a,
	"mediumturquoise":      0xff48d1cc,
	"mediumvioletred":      0xffc71585,
	"midnightblue":         0xff191970,
	"mintcream":            0xfff5fffa,
	"mistyrose":            0xffffe4e1,
	"moccasin":             0xffffe4b5,
	"navajowhite":          0xffffdead,
	"oldlace":              0xfffdf5e6,
	"olivedrab":            0xff6b8e23,
	"orangered":            0xffff4500,
	"orchid":               0xffda70d6,
	"palegoldenrod":        0xffeee8aa,
	"palegreen":            0xff98fb98,
	"paleturquoise":        0xffafeeee,
	"palevioletred":        0xffdb7093,
	"papayawhip":           0xffffefd5,
	"peachpuff":            0xffffdab9,
	"peru":                 0xffcd853f,
	"pink":                 0xffffc0cb,
	"plum":                 0xffdda0dd,
	"powderblue":           0xffb0e0e6,
	"rosybrown":            0xffbc8f8f,
	"royalblue":            0xff4169e1,
	"saddlebrown":          0xff8b4513,
	"salmon":               0xfffa8072,
	"sandybrown":           0xfff4a460,
	"seagreen":             0xff2e8b57,
	"seashell":             0xfffff5ee,
	"sienna":               0xffa0522d,
	"skyblue":              0xff87ceeb,
	"slateblue":            0xff6a5acd,
	"slategray":            0xff708090,
	"slategrey":            0xff708090,
	"snow":                 0xfffffafa,
	"springgreen":          0xff00ff7f,
	"steelblue":            0xff4682b4,
	"tan":                  0xffd2b48c,
	"thistle":              0xffd8bfd8,
	"tomato":               0xffff6347,
	"turquoise":            0xff40e0d0,
	"violet":               0xffee82ee,
	"wheat":                0xfff5deb3,
	"whitesmoke":           0xfff5f5f5,
	"yellowgreen":          0xff9acd32,
}

// NamedColor make a relation between color and its name
type NamedColor struct {
	// Name of a color
	Name string

	// Color value
	Color Color
}

// NamedColors returns the list of named colors
func NamedColors() []NamedColor {
	count := len(colorConstants)
	result := make([]NamedColor, 0, count)
	for name, color := range colorConstants {
		result = append(result, NamedColor{Name: name, Color: color})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}
