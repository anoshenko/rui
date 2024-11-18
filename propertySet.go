package rui

import (
	"math"
	"strconv"
	"strings"
)

var colorProperties = []PropertyName{
	ColorTag,
	BackgroundColor,
	TextColor,
	CaretColor,
	BorderColor,
	BorderLeftColor,
	BorderRightColor,
	BorderTopColor,
	BorderBottomColor,
	OutlineColor,
	TextLineColor,
	ColorPickerValue,
	AccentColor,
}

func isPropertyInList(tag PropertyName, list []PropertyName) bool {
	for _, prop := range list {
		if prop == tag {
			return true
		}
	}
	return false
}

var angleProperties = []PropertyName{
	From,
}

var boolProperties = []PropertyName{
	Disabled,
	Focusable,
	Inset,
	BackfaceVisible,
	ReadOnly,
	EditWrap,
	Spellcheck,
	CloseButton,
	OutsideClose,
	Italic,
	SmallCaps,
	Strikethrough,
	Overline,
	Underline,
	Expanded,
	AvoidBreak,
	NotTranslate,
	Controls,
	Loop,
	Muted,
	AnimationPaused,
	Multiple,
	TabCloseButton,
	Repeating,
	UserSelect,
	ColumnSpanAll,
}

var intProperties = []PropertyName{
	ZIndex,
	TabSize,
	HeadHeight,
	FootHeight,
	RowSpan,
	ColumnSpan,
	ColumnCount,
	Order,
	TabIndex,
	MaxLength,
	NumberPickerPrecision,
}

var floatProperties = map[PropertyName]struct{ min, max float64 }{
	Opacity:           {min: 0, max: 1},
	NumberPickerMax:   {min: -math.MaxFloat64, max: math.MaxFloat64},
	NumberPickerMin:   {min: -math.MaxFloat64, max: math.MaxFloat64},
	NumberPickerStep:  {min: -math.MaxFloat64, max: math.MaxFloat64},
	NumberPickerValue: {min: -math.MaxFloat64, max: math.MaxFloat64},
	ProgressBarMax:    {min: 0, max: math.MaxFloat64},
	ProgressBarValue:  {min: 0, max: math.MaxFloat64},
	VideoWidth:        {min: 0, max: 10000},
	VideoHeight:       {min: 0, max: 10000},
}

var sizeProperties = map[PropertyName]string{
	Width:              string(Width),
	Height:             string(Height),
	MinWidth:           string(MinWidth),
	MinHeight:          string(MinHeight),
	MaxWidth:           string(MaxWidth),
	MaxHeight:          string(MaxHeight),
	Left:               string(Left),
	Right:              string(Right),
	Top:                string(Top),
	Bottom:             string(Bottom),
	TextSize:           "font-size",
	TextIndent:         string(TextIndent),
	LetterSpacing:      string(LetterSpacing),
	WordSpacing:        string(WordSpacing),
	LineHeight:         string(LineHeight),
	TextLineThickness:  "text-decoration-thickness",
	ListRowGap:         "row-gap",
	ListColumnGap:      "column-gap",
	GridRowGap:         string(GridRowGap),
	GridColumnGap:      string(GridColumnGap),
	ColumnWidth:        string(ColumnWidth),
	ColumnGap:          string(ColumnGap),
	Gap:                string(Gap),
	Margin:             string(Margin),
	MarginLeft:         string(MarginLeft),
	MarginRight:        string(MarginRight),
	MarginTop:          string(MarginTop),
	MarginBottom:       string(MarginBottom),
	Padding:            string(Padding),
	PaddingLeft:        string(PaddingLeft),
	PaddingRight:       string(PaddingRight),
	PaddingTop:         string(PaddingTop),
	PaddingBottom:      string(PaddingBottom),
	BorderWidth:        string(BorderWidth),
	BorderLeftWidth:    string(BorderLeftWidth),
	BorderRightWidth:   string(BorderRightWidth),
	BorderTopWidth:     string(BorderTopWidth),
	BorderBottomWidth:  string(BorderBottomWidth),
	OutlineWidth:       string(OutlineWidth),
	OutlineOffset:      string(OutlineOffset),
	XOffset:            string(XOffset),
	YOffset:            string(YOffset),
	BlurRadius:         string(BlurRadius),
	SpreadRadius:       string(SpreadRadius),
	Perspective:        string(Perspective),
	PerspectiveOriginX: string(PerspectiveOriginX),
	PerspectiveOriginY: string(PerspectiveOriginY),
	TransformOriginX:   string(TransformOriginX),
	TransformOriginY:   string(TransformOriginY),
	TransformOriginZ:   string(TransformOriginZ),
	Radius:             string(Radius),
	RadiusX:            string(RadiusX),
	RadiusY:            string(RadiusY),
	RadiusTopLeft:      string(RadiusTopLeft),
	RadiusTopLeftX:     string(RadiusTopLeftX),
	RadiusTopLeftY:     string(RadiusTopLeftY),
	RadiusTopRight:     string(RadiusTopRight),
	RadiusTopRightX:    string(RadiusTopRightX),
	RadiusTopRightY:    string(RadiusTopRightY),
	RadiusBottomLeft:   string(RadiusBottomLeft),
	RadiusBottomLeftX:  string(RadiusBottomLeftX),
	RadiusBottomLeftY:  string(RadiusBottomLeftY),
	RadiusBottomRight:  string(RadiusBottomRight),
	RadiusBottomRightX: string(RadiusBottomRightX),
	RadiusBottomRightY: string(RadiusBottomRightY),
	ItemWidth:          string(ItemWidth),
	ItemHeight:         string(ItemHeight),
	CenterX:            string(CenterX),
	CenterY:            string(CenterX),
}

var enumProperties = map[PropertyName]struct {
	values    []string
	cssTag    string
	cssValues []string
}{
	Semantics: {
		[]string{"default", "article", "section", "aside", "header", "main", "footer", "navigation", "figure", "figure-caption", "button", "p", "h1", "h2", "h3", "h4", "h5", "h6", "blockquote", "code"},
		"",
		[]string{"div", "article", "section", "aside", "header", "main", "footer", "nav", "figure", "figcaption", "button", "p", "h1", "h2", "h3", "h4", "h5", "h6", "blockquote", "code"},
	},
	Visibility: {
		[]string{"visible", "invisible", "gone"},
		"",
		[]string{"visible", "invisible", "gone"},
	},
	Overflow: {
		[]string{"hidden", "visible", "scroll", "auto"},
		string(Overflow),
		[]string{"hidden", "visible", "scroll", "auto"},
	},
	TextAlign: {
		[]string{"left", "right", "center", "justify"},
		string(TextAlign),
		[]string{"left", "right", "center", "justify"},
	},
	TextTransform: {
		[]string{"none", "capitalize", "lowercase", "uppercase"},
		string(TextTransform),
		[]string{"none", "capitalize", "lowercase", "uppercase"},
	},
	TextWeight: {
		[]string{"inherit", "thin", "extra-light", "light", "normal", "medium", "semi-bold", "bold", "extra-bold", "black"},
		"font-weight",
		[]string{"inherit", "100", "200", "300", "normal", "500", "600", "bold", "800", "900"},
	},
	WhiteSpace: {
		[]string{"normal", "nowrap", "pre", "pre-wrap", "pre-line", "break-spaces"},
		string(WhiteSpace),
		[]string{"normal", "nowrap", "pre", "pre-wrap", "pre-line", "break-spaces"},
	},
	WordBreak: {
		[]string{"normal", "break-all", "keep-all", "break-word"},
		string(WordBreak),
		[]string{"normal", "break-all", "keep-all", "break-word"},
	},
	TextOverflow: {
		[]string{"clip", "ellipsis"},
		string(TextOverflow),
		[]string{"clip", "ellipsis"},
	},
	TextWrap: {
		[]string{"wrap", "nowrap", "balance"},
		string(TextWrap),
		[]string{"wrap", "nowrap", "balance"},
	},
	WritingMode: {
		[]string{"horizontal-top-to-bottom", "horizontal-bottom-to-top", "vertical-right-to-left", "vertical-left-to-right"},
		string(WritingMode),
		[]string{"horizontal-tb", "horizontal-bt", "vertical-rl", "vertical-lr"},
	},
	TextDirection: {
		[]string{"system", "left-to-right", "right-to-left"},
		"direction",
		[]string{"", "ltr", "rtl"},
	},
	VerticalTextOrientation: {
		[]string{"mixed", "upright"},
		"text-orientation",
		[]string{"mixed", "upright"},
	},
	TextLineStyle: {
		[]string{"inherit", "solid", "dashed", "dotted", "double", "wavy"},
		"text-decoration-style",
		[]string{"inherit", "solid", "dashed", "dotted", "double", "wavy"},
	},
	BorderStyle: {
		[]string{"none", "solid", "dashed", "dotted", "double"},
		string(BorderStyle),
		[]string{"none", "solid", "dashed", "dotted", "double"},
	},
	TopStyle: {
		[]string{"none", "solid", "dashed", "dotted", "double"},
		"",
		[]string{"none", "solid", "dashed", "dotted", "double"},
	},
	RightStyle: {
		[]string{"none", "solid", "dashed", "dotted", "double"},
		"",
		[]string{"none", "solid", "dashed", "dotted", "double"},
	},
	BottomStyle: {
		[]string{"none", "solid", "dashed", "dotted", "double"},
		"",
		[]string{"none", "solid", "dashed", "dotted", "double"},
	},
	LeftStyle: {
		[]string{"none", "solid", "dashed", "dotted", "double"},
		"",
		[]string{"none", "solid", "dashed", "dotted", "double"},
	},
	OutlineStyle: {
		[]string{"none", "solid", "dashed", "dotted", "double"},
		string(OutlineStyle),
		[]string{"none", "solid", "dashed", "dotted", "double"},
	},
	Tabs: {
		[]string{"top", "bottom", "left", "right", "left-list", "right-list", "hidden"},
		"",
		[]string{"top", "bottom", "left", "right", "left-list", "right-list", "hidden"},
	},
	NumberPickerType: {
		[]string{"editor", "slider"},
		"",
		[]string{"editor", "slider"},
	},
	EditViewType: {
		[]string{"text", "password", "email", "emails", "url", "phone", "multiline"},
		"",
		[]string{"text", "password", "email", "emails", "url", "phone", "multiline"},
	},
	Orientation: {
		[]string{"up-down", "start-to-end", "bottom-up", "end-to-start"},
		"",
		[]string{"column", "row", "column-reverse", "row-reverse"},
	},
	ListWrap: {
		[]string{"off", "on", "reverse"},
		"",
		[]string{"nowrap", "wrap", "wrap-reverse"},
	},
	"list-orientation": {
		[]string{"vertical", "horizontal"},
		"",
		[]string{"vertical", "horizontal"},
	},
	VerticalAlign: {
		[]string{"top", "bottom", "center", "stretch"},
		"",
		[]string{"top", "bottom", "center", "stretch"},
	},
	HorizontalAlign: {
		[]string{"left", "right", "center", "stretch"},
		"",
		[]string{"left", "right", "center", "stretch"},
	},
	ButtonsAlign: {
		[]string{"left", "right", "center", "stretch"},
		"",
		[]string{"left", "right", "center", "stretch"},
	},
	ArrowAlign: {
		[]string{"left", "right", "center"},
		"",
		[]string{"left", "right", "center"},
	},
	CellVerticalAlign: {
		[]string{"top", "bottom", "center", "stretch"},
		"align-items",
		[]string{"start", "end", "center", "stretch"},
	},
	CellHorizontalAlign: {
		[]string{"left", "right", "center", "stretch"},
		"justify-items",
		[]string{"start", "end", "center", "stretch"},
	},
	CellVerticalSelfAlign: {
		[]string{"top", "bottom", "center", "stretch"},
		"align-self",
		[]string{"start", "end", "center", "stretch"},
	},
	CellHorizontalSelfAlign: {
		[]string{"left", "right", "center", "stretch"},
		"justify-self",
		[]string{"start", "end", "center", "stretch"},
	},
	GridAutoFlow: {
		[]string{"row", "column", "row-dense", "column-dense"},
		string(GridAutoFlow),
		[]string{"row", "column", "row dense", "column dense"},
	},
	ImageVerticalAlign: {
		[]string{"top", "bottom", "center"},
		"",
		[]string{"top", "bottom", "center"},
	},
	ImageHorizontalAlign: {
		[]string{"left", "right", "center"},
		"",
		[]string{"left", "right", "center"},
	},
	ItemVerticalAlign: {
		[]string{"top", "bottom", "center", "stretch"},
		"",
		[]string{"start", "end", "center", "stretch"},
	},
	ItemHorizontalAlign: {
		[]string{"left", "right", "center", "stretch"},
		"",
		[]string{"start", "end", "center", "stretch"},
	},
	CheckboxVerticalAlign: {
		[]string{"top", "bottom", "center"},
		"",
		[]string{"start", "end", "center"},
	},
	CheckboxHorizontalAlign: {
		[]string{"left", "right", "center"},
		"",
		[]string{"start", "end", "center"},
	},
	TableVerticalAlign: {
		[]string{"top", "bottom", "center", "stretch", "baseline"},
		"vertical-align",
		[]string{"top", "bottom", "middle", "baseline", "baseline"},
	},
	Cursor: {
		[]string{"auto", "default", "none", "context-menu", "help", "pointer", "progress", "wait", "cell", "crosshair", "text", "vertical-text", "alias", "copy", "move", "no-drop", "not-allowed", "e-resize", "n-resize", "ne-resize", "nw-resize", "s-resize", "se-resize", "sw-resize", "w-resize", "ew-resize", "ns-resize", "nesw-resize", "nwse-resize", "col-resize", "row-resize", "all-scroll", "zoom-in", "zoom-out", "grab", "grabbing"},
		string(Cursor),
		[]string{"auto", "default", "none", "context-menu", "help", "pointer", "progress", "wait", "cell", "crosshair", "text", "vertical-text", "alias", "copy", "move", "no-drop", "not-allowed", "e-resize", "n-resize", "ne-resize", "nw-resize", "s-resize", "se-resize", "sw-resize", "w-resize", "ew-resize", "ns-resize", "nesw-resize", "nwse-resize", "col-resize", "row-resize", "all-scroll", "zoom-in", "zoom-out", "grab", "grabbing"},
	},
	Fit: {
		[]string{"none", "contain", "cover", "fill", "scale-down"},
		"object-fit",
		[]string{"none", "contain", "cover", "fill", "scale-down"},
	},
	backgroundFit: {
		[]string{"none", "contain", "cover"},
		"",
		[]string{"none", "contain", "cover"},
	},
	Repeat: {
		[]string{"no-repeat", "repeat", "repeat-x", "repeat-y", "round", "space"},
		"",
		[]string{"no-repeat", "repeat", "repeat-x", "repeat-y", "round", "space"},
	},
	Attachment: {
		[]string{"scroll", "fixed", "local"},
		"",
		[]string{"scroll", "fixed", "local"},
	},
	BackgroundClip: {
		[]string{"border-box", "padding-box", "content-box"}, // "text"},
		"background-clip",
		[]string{"border-box", "padding-box", "content-box"}, // "text"},
	},
	Direction: {
		[]string{"to-top", "to-right-top", "to-right", "to-right-bottom", "to-bottom", "to-left-bottom", "to-left", "to-left-top"},
		"",
		[]string{"to top", "to right top", "to right", "to right bottom", "to bottom", "to left bottom", "to left", "to left top"},
	},
	AnimationDirection: {
		[]string{"normal", "reverse", "alternate", "alternate-reverse"},
		"",
		[]string{"normal", "reverse", "alternate", "alternate-reverse"},
	},
	RadialGradientShape: {
		[]string{"ellipse", "circle"},
		"",
		[]string{"ellipse", "circle"},
	},
	RadialGradientRadius: {
		[]string{"closest-side", "closest-corner", "farthest-side", "farthest-corner"},
		"",
		[]string{"closest-side", "closest-corner", "farthest-side", "farthest-corner"},
	},
	ItemCheckbox: {
		[]string{"none", "single", "multiple"},
		"",
		[]string{"none", "single", "multiple"},
	},
	Float: {
		[]string{"none", "left", "right"},
		"float",
		[]string{"none", "left", "right"},
	},
	Preload: {
		[]string{"none", "metadata", "auto"},
		"",
		[]string{"none", "metadata", "auto"},
	},
	SelectionMode: {
		[]string{"none", "cell", "row"},
		"",
		[]string{"none", "cell", "row"},
	},
	Resize: {
		[]string{"none", "both", "horizontal", "vertical"},
		"resize",
		[]string{"none", "both", "horizontal", "vertical"},
	},
	Arrow: {
		[]string{"none", "top", "right", "bottom", "left"},
		"",
		[]string{"none", "top", "right", "bottom", "left"},
	},
	MixBlendMode: {
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
		string(MixBlendMode),
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
	},
	BackgroundBlendMode: {
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
		string(BackgroundBlendMode),
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
	},
	ColumnFill: {
		[]string{"balance", "auto"},
		string(ColumnFill),
		[]string{"balance", "auto"},
	},
}

func notCompatibleType(tag PropertyName, value any) {
	ErrorLogF(`"%T" type not compatible with "%s" property`, value, string(tag))
}

func invalidPropertyValue(tag PropertyName, value any) {
	ErrorLogF(`Invalid value "%v" of "%s" property`, value, string(tag))
}

func isConstantName(text string) bool {
	len := len(text)
	if len <= 1 || text[0] != '@' {
		return false
	}

	if len > 2 {
		last := len - 1
		if (text[1] == '`' && text[last] == '`') ||
			(text[1] == '"' && text[last] == '"') ||
			(text[1] == '\'' && text[last] == '\'') {
			return true
		}
	}

	return !strings.ContainsAny(text, ",;|\"'`+(){}[]<>/\\*&%! \t\n\r")
}

func isInt(value any) (int, bool) {
	var n int
	switch value := value.(type) {
	case int:
		n = value

	case int8:
		n = int(value)

	case int16:
		n = int(value)

	case int32:
		n = int(value)

	case int64:
		n = int(value)

	case uint:
		n = int(value)

	case uint8:
		n = int(value)

	case uint16:
		n = int(value)

	case uint32:
		n = int(value)

	case uint64:
		n = int(value)

	default:
		return 0, false
	}

	return n, true
}

func setSimpleProperty(properties Properties, tag PropertyName, value any) bool {
	if value == nil {
		properties.setRaw(tag, nil)
		return true
	} else if text, ok := value.(string); ok {
		text = strings.Trim(text, " \t\n\r")
		if text == "" {
			properties.setRaw(tag, nil)
			return true
		}
		if isConstantName(text) {
			properties.setRaw(tag, text)
			return true
		}
	}
	return false
}

func setStringPropertyValue(properties Properties, tag PropertyName, text any) []PropertyName {
	if text != "" {
		properties.setRaw(tag, text)
	} else if properties.getRaw(tag) != nil {
		properties.setRaw(tag, nil)
	} else {
		return []PropertyName{}
	}
	return []PropertyName{tag}
}

func setArrayPropertyValue[T any](properties Properties, tag PropertyName, value []T) []PropertyName {
	if len(value) > 0 {
		properties.setRaw(tag, value)
	} else if properties.getRaw(tag) != nil {
		properties.setRaw(tag, nil)
	} else {
		return []PropertyName{}
	}
	return []PropertyName{tag}
}

func setSizeProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		var size SizeUnit
		switch value := value.(type) {
		case string:
			var ok bool
			if fn := parseSizeFunc(value); fn != nil {
				size.Type = SizeFunction
				size.Function = fn
			} else if size, ok = StringToSizeUnit(value); !ok {
				invalidPropertyValue(tag, value)
				return nil
			}
		case SizeUnit:
			size = value

		case SizeFunc:
			size.Type = SizeFunction
			size.Function = value

		case float32:
			size.Type = SizeInPixel
			size.Value = float64(value)

		case float64:
			size.Type = SizeInPixel
			size.Value = value

		default:
			if n, ok := isInt(value); ok {
				size.Type = SizeInPixel
				size.Value = float64(n)
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}

		if size.Type == Auto {
			properties.setRaw(tag, nil)
		} else {
			properties.setRaw(tag, size)
		}
	}

	return []PropertyName{tag}
}

func setAngleProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		var angle AngleUnit
		switch value := value.(type) {
		case string:
			var ok bool
			if angle, ok = StringToAngleUnit(value); !ok {
				invalidPropertyValue(tag, value)
				return nil
			}
		case AngleUnit:
			angle = value

		case float32:
			angle = Rad(float64(value))

		case float64:
			angle = Rad(value)

		default:
			if n, ok := isInt(value); ok {
				angle = Rad(float64(n))
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}
		properties.setRaw(tag, angle)
	}

	return []PropertyName{tag}
}

func setColorProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		var result Color
		switch value := value.(type) {
		case string:
			var err error
			if result, err = stringToColor(value); err != nil {
				invalidPropertyValue(tag, value)
				return nil
			}
		case Color:
			result = value

		default:
			if color, ok := isInt(value); ok {
				result = Color(color)
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}

		properties.setRaw(tag, result)
	}

	return []PropertyName{tag}
}

func setEnumProperty(properties Properties, tag PropertyName, value any, values []string) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		var n int
		if text, ok := value.(string); ok {
			if n, ok = enumStringToInt(text, values, false); !ok {
				invalidPropertyValue(tag, value)
				return nil
			}
		} else if i, ok := isInt(value); ok {
			if i < 0 || i >= len(values) {
				invalidPropertyValue(tag, value)
				return nil
			}
			n = i
		} else {
			notCompatibleType(tag, value)
			return nil
		}

		properties.setRaw(tag, n)
	}

	return []PropertyName{tag}
}

func setBoolProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		if text, ok := value.(string); ok {
			switch strings.ToLower(strings.Trim(text, " \t")) {
			case "true", "yes", "on", "1":
				properties.setRaw(tag, true)

			case "false", "no", "off", "0":
				properties.setRaw(tag, false)

			default:
				invalidPropertyValue(tag, value)
				return nil
			}
		} else if n, ok := isInt(value); ok {
			switch n {
			case 1:
				properties.setRaw(tag, true)

			case 0:
				properties.setRaw(tag, false)

			default:
				invalidPropertyValue(tag, value)
				return nil
			}
		} else if b, ok := value.(bool); ok {
			properties.setRaw(tag, b)
		} else {
			notCompatibleType(tag, value)
			return nil
		}
	}

	return []PropertyName{tag}
}

func setIntProperty(properties Properties, tag PropertyName, value any) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		if text, ok := value.(string); ok {
			n, err := strconv.Atoi(strings.Trim(text, " \t"))
			if err != nil {
				invalidPropertyValue(tag, value)
				ErrorLog(err.Error())
				return nil
			}
			properties.setRaw(tag, n)
		} else if n, ok := isInt(value); ok {
			properties.setRaw(tag, n)
		} else {
			notCompatibleType(tag, value)
			return nil
		}
	}

	return []PropertyName{tag}
}

func setFloatProperty(properties Properties, tag PropertyName, value any, min, max float64) []PropertyName {
	if !setSimpleProperty(properties, tag, value) {
		f := float64(0)
		switch value := value.(type) {
		case string:
			var err error
			if f, err = strconv.ParseFloat(strings.Trim(value, " \t"), 64); err != nil {
				invalidPropertyValue(tag, value)
				ErrorLog(err.Error())
				return nil
			}
			if f < min || f > max {
				ErrorLogF(`"%T" out of range of "%s" property`, value, tag)
				return nil
			}
			properties.setRaw(tag, value)
			return nil

		case float32:
			f = float64(value)

		case float64:
			f = value

		default:
			if n, ok := isInt(value); ok {
				f = float64(n)
			} else {
				notCompatibleType(tag, value)
				return nil
			}
		}

		if f >= min && f <= max {
			properties.setRaw(tag, f)
		} else {
			ErrorLogF(`"%T" out of range of "%s" property`, value, tag)
			return nil
		}
	}

	return []PropertyName{tag}
}

func propertiesSet(properties Properties, tag PropertyName, value any) []PropertyName {
	if _, ok := sizeProperties[tag]; ok {
		return setSizeProperty(properties, tag, value)
	}

	if valuesData, ok := enumProperties[tag]; ok {
		return setEnumProperty(properties, tag, value, valuesData.values)
	}

	if limits, ok := floatProperties[tag]; ok {
		return setFloatProperty(properties, tag, value, limits.min, limits.max)
	}

	if isPropertyInList(tag, colorProperties) {
		return setColorProperty(properties, tag, value)
	}

	if isPropertyInList(tag, angleProperties) {
		return setAngleProperty(properties, tag, value)
	}

	if isPropertyInList(tag, boolProperties) {
		return setBoolProperty(properties, tag, value)
	}

	if isPropertyInList(tag, intProperties) {
		return setIntProperty(properties, tag, value)
	}

	if text, ok := value.(string); ok {
		properties.setRaw(tag, text)
		return []PropertyName{tag}
	}

	notCompatibleType(tag, value)
	return nil
}

/*
func (properties *propertyList) Set(tag PropertyName, value any) bool {
	tag = properties.normalize(tag)
	if value == nil {
		properties.remove(properties, tag)
		return true
	}

	return properties.set(properties, tag, value) != nil
}
*/

func (data *dataProperty) Set(tag PropertyName, value any) bool {
	if value == nil {
		data.Remove(tag)
		return true
	}

	tag = data.normalize(tag)
	for _, supported := range data.supportedProperties {
		if tag == supported {
			return data.set(data, tag, value) != nil
		}
	}

	ErrorLogF(`"%s" property is not supported`, string(tag))
	return false
}
