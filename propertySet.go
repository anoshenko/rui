package rui

import (
	"math"
	"strconv"
	"strings"
)

var colorProperties = []string{
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
}

func isPropertyInList(tag string, list []string) bool {
	for _, prop := range list {
		if prop == tag {
			return true
		}
	}
	return false
}

var angleProperties = []string{
	Rotate,
	SkewX,
	SkewY,
	From,
}

var boolProperties = []string{
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
}

var intProperties = []string{
	ZIndex,
	TabSize,
	HeadHeight,
	FootHeight,
	RowSpan,
	ColumnSpan,
	ColumnCount,
	Order,
	TabIndex,
}

var floatProperties = map[string]struct{ min, max float64 }{
	Opacity:           {min: 0, max: 1},
	ScaleX:            {min: -math.MaxFloat64, max: math.MaxFloat64},
	ScaleY:            {min: -math.MaxFloat64, max: math.MaxFloat64},
	ScaleZ:            {min: -math.MaxFloat64, max: math.MaxFloat64},
	RotateX:           {min: 0, max: 1},
	RotateY:           {min: 0, max: 1},
	RotateZ:           {min: 0, max: 1},
	NumberPickerMax:   {min: -math.MaxFloat64, max: math.MaxFloat64},
	NumberPickerMin:   {min: -math.MaxFloat64, max: math.MaxFloat64},
	NumberPickerStep:  {min: -math.MaxFloat64, max: math.MaxFloat64},
	NumberPickerValue: {min: -math.MaxFloat64, max: math.MaxFloat64},
	ProgressBarMax:    {min: 0, max: math.MaxFloat64},
	ProgressBarValue:  {min: 0, max: math.MaxFloat64},
	VideoWidth:        {min: 0, max: 10000},
	VideoHeight:       {min: 0, max: 10000},
}

var sizeProperties = map[string]string{
	Width:              Width,
	Height:             Height,
	MinWidth:           MinWidth,
	MinHeight:          MinHeight,
	MaxWidth:           MaxWidth,
	MaxHeight:          MaxHeight,
	Left:               Left,
	Right:              Right,
	Top:                Top,
	Bottom:             Bottom,
	TextSize:           "font-size",
	TextIndent:         TextIndent,
	LetterSpacing:      LetterSpacing,
	WordSpacing:        WordSpacing,
	LineHeight:         LineHeight,
	TextLineThickness:  "text-decoration-thickness",
	ListRowGap:         "row-gap",
	ListColumnGap:      "column-gap",
	GridRowGap:         GridRowGap,
	GridColumnGap:      GridColumnGap,
	ColumnWidth:        ColumnWidth,
	ColumnGap:          ColumnGap,
	Gap:                Gap,
	Margin:             Margin,
	MarginLeft:         MarginLeft,
	MarginRight:        MarginRight,
	MarginTop:          MarginTop,
	MarginBottom:       MarginBottom,
	Padding:            Padding,
	PaddingLeft:        PaddingLeft,
	PaddingRight:       PaddingRight,
	PaddingTop:         PaddingTop,
	PaddingBottom:      PaddingBottom,
	BorderWidth:        BorderWidth,
	BorderLeftWidth:    BorderLeftWidth,
	BorderRightWidth:   BorderRightWidth,
	BorderTopWidth:     BorderTopWidth,
	BorderBottomWidth:  BorderBottomWidth,
	OutlineWidth:       OutlineWidth,
	XOffset:            XOffset,
	YOffset:            YOffset,
	BlurRadius:         BlurRadius,
	SpreadRadius:       SpreadRadius,
	Perspective:        Perspective,
	PerspectiveOriginX: PerspectiveOriginX,
	PerspectiveOriginY: PerspectiveOriginY,
	OriginX:            OriginX,
	OriginY:            OriginY,
	OriginZ:            OriginZ,
	TranslateX:         TranslateX,
	TranslateY:         TranslateY,
	TranslateZ:         TranslateZ,
	Radius:             Radius,
	RadiusX:            RadiusX,
	RadiusY:            RadiusY,
	RadiusTopLeft:      RadiusTopLeft,
	RadiusTopLeftX:     RadiusTopLeftX,
	RadiusTopLeftY:     RadiusTopLeftY,
	RadiusTopRight:     RadiusTopRight,
	RadiusTopRightX:    RadiusTopRightX,
	RadiusTopRightY:    RadiusTopRightY,
	RadiusBottomLeft:   RadiusBottomLeft,
	RadiusBottomLeftX:  RadiusBottomLeftX,
	RadiusBottomLeftY:  RadiusBottomLeftY,
	RadiusBottomRight:  RadiusBottomRight,
	RadiusBottomRightX: RadiusBottomRightX,
	RadiusBottomRightY: RadiusBottomRightY,
	ItemWidth:          ItemWidth,
	ItemHeight:         ItemHeight,
	CenterX:            CenterX,
	CenterY:            CenterX,
}

var enumProperties = map[string]struct {
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
		Overflow,
		[]string{"hidden", "visible", "scroll", "auto"},
	},
	TextAlign: {
		[]string{"left", "right", "center", "justify"},
		TextAlign,
		[]string{"left", "right", "center", "justify"},
	},
	TextTransform: {
		[]string{"none", "capitalize", "lowercase", "uppercase"},
		TextTransform,
		[]string{"none", "capitalize", "lowercase", "uppercase"},
	},
	TextWeight: {
		[]string{"inherit", "thin", "extra-light", "light", "normal", "medium", "semi-bold", "bold", "extra-bold", "black"},
		"font-weight",
		[]string{"inherit", "100", "200", "300", "normal", "500", "600", "bold", "800", "900"},
	},
	WhiteSpace: {
		[]string{"normal", "nowrap", "pre", "pre-wrap", "pre-line", "break-spaces"},
		WhiteSpace,
		[]string{"normal", "nowrap", "pre", "pre-wrap", "pre-line", "break-spaces"},
	},
	WordBreak: {
		[]string{"normal", "break-all", "keep-all", "break-word"},
		WordBreak,
		[]string{"normal", "break-all", "keep-all", "break-word"},
	},
	TextOverflow: {
		[]string{"clip", "ellipsis"},
		TextOverflow,
		[]string{"clip", "ellipsis"},
	},
	WritingMode: {
		[]string{"horizontal-top-to-bottom", "horizontal-bottom-to-top", "vertical-right-to-left", "vertical-left-to-right"},
		WritingMode,
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
		BorderStyle,
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
		OutlineStyle,
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
	GridAutoFlow: {
		[]string{"row", "column", "row-dense", "column-dense"},
		GridAutoFlow,
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
		Cursor,
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
		MixBlendMode,
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
	},
	BackgroundBlendMode: {
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
		BackgroundBlendMode,
		[]string{"normal", "multiply", "screen", "overlay", "darken", "lighten", "color-dodge", "color-burn", "hard-light", "soft-light", "difference", "exclusion", "hue", "saturation", "color", "luminosity"},
	},
	ColumnFill: {
		[]string{"balance", "auto"},
		ColumnFill,
		[]string{"balance", "auto"},
	},
}

func notCompatibleType(tag string, value any) {
	ErrorLogF(`"%T" type not compatible with "%s" property`, value, tag)
}

func invalidPropertyValue(tag string, value any) {
	ErrorLogF(`Invalid value "%v" of "%s" property`, value, tag)
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

func (properties *propertyList) setSimpleProperty(tag string, value any) bool {
	if value == nil {
		delete(properties.properties, tag)
		return true
	} else if text, ok := value.(string); ok {
		text = strings.Trim(text, " \t\n\r")
		if text == "" {
			delete(properties.properties, tag)
			return true
		}
		if isConstantName(text) {
			properties.properties[tag] = text
			return true
		}
	}
	return false
}

func (properties *propertyList) setSizeProperty(tag string, value any) bool {
	if !properties.setSimpleProperty(tag, value) {
		var size SizeUnit
		switch value := value.(type) {
		case string:
			var ok bool
			if fn := parseSizeFunc(value); fn != nil {
				size.Type = SizeFunction
				size.Function = fn
			} else if size, ok = StringToSizeUnit(value); !ok {
				invalidPropertyValue(tag, value)
				return false
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
				return false
			}
		}

		if size.Type == Auto {
			delete(properties.properties, tag)
		} else {
			properties.properties[tag] = size
		}
	}

	return true
}

func (properties *propertyList) setAngleProperty(tag string, value any) bool {
	if !properties.setSimpleProperty(tag, value) {
		var angle AngleUnit
		switch value := value.(type) {
		case string:
			var ok bool
			if angle, ok = StringToAngleUnit(value); !ok {
				invalidPropertyValue(tag, value)
				return false
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
				return false
			}
		}
		properties.properties[tag] = angle
	}

	return true
}

func (properties *propertyList) setColorProperty(tag string, value any) bool {
	if !properties.setSimpleProperty(tag, value) {
		var result Color
		switch value := value.(type) {
		case string:
			var err error
			if result, err = stringToColor(value); err != nil {
				invalidPropertyValue(tag, value)
				return false
			}
		case Color:
			result = value

		default:
			if color, ok := isInt(value); ok {
				result = Color(color)
			} else {
				notCompatibleType(tag, value)
				return false
			}
		}

		if result == 0 {
			delete(properties.properties, tag)
		} else {
			properties.properties[tag] = result
		}
	}

	return true
}

func (properties *propertyList) setEnumProperty(tag string, value any, values []string) bool {
	if !properties.setSimpleProperty(tag, value) {
		var n int
		if text, ok := value.(string); ok {
			if n, ok = enumStringToInt(text, values, false); !ok {
				invalidPropertyValue(tag, value)
				return false
			}
		} else if i, ok := isInt(value); ok {
			if i < 0 || i >= len(values) {
				invalidPropertyValue(tag, value)
				return false
			}
			n = i
		} else {
			notCompatibleType(tag, value)
			return false
		}

		properties.properties[tag] = n
	}

	return true
}

func (properties *propertyList) setBoolProperty(tag string, value any) bool {
	if !properties.setSimpleProperty(tag, value) {
		if text, ok := value.(string); ok {
			switch strings.ToLower(strings.Trim(text, " \t")) {
			case "true", "yes", "on", "1":
				properties.properties[tag] = true

			case "false", "no", "off", "0":
				properties.properties[tag] = false

			default:
				invalidPropertyValue(tag, value)
				return false
			}
		} else if n, ok := isInt(value); ok {
			switch n {
			case 1:
				properties.properties[tag] = true

			case 0:
				properties.properties[tag] = false

			default:
				invalidPropertyValue(tag, value)
				return false
			}
		} else if b, ok := value.(bool); ok {
			properties.properties[tag] = b
		} else {
			notCompatibleType(tag, value)
			return false
		}
	}

	return true
}

func (properties *propertyList) setIntProperty(tag string, value any) bool {
	if !properties.setSimpleProperty(tag, value) {
		if text, ok := value.(string); ok {
			n, err := strconv.Atoi(strings.Trim(text, " \t"))
			if err != nil {
				invalidPropertyValue(tag, value)
				ErrorLog(err.Error())
				return false
			}
			properties.properties[tag] = n
		} else if n, ok := isInt(value); ok {
			properties.properties[tag] = n
		} else {
			notCompatibleType(tag, value)
			return false
		}
	}

	return true
}

func (properties *propertyList) setFloatProperty(tag string, value any, min, max float64) bool {
	if !properties.setSimpleProperty(tag, value) {
		f := float64(0)
		switch value := value.(type) {
		case string:
			var err error
			if f, err = strconv.ParseFloat(strings.Trim(value, " \t"), 64); err != nil {
				invalidPropertyValue(tag, value)
				ErrorLog(err.Error())
				return false
			}
			if f < min || f > max {
				ErrorLogF(`"%T" out of range of "%s" property`, value, tag)
				return false
			}
			properties.properties[tag] = value
			return true

		case float32:
			f = float64(value)

		case float64:
			f = value

		default:
			if n, ok := isInt(value); ok {
				f = float64(n)
			} else {
				notCompatibleType(tag, value)
				return false
			}
		}

		if f >= min && f <= max {
			properties.properties[tag] = f
		} else {
			ErrorLogF(`"%T" out of range of "%s" property`, value, tag)
			return false
		}
	}

	return true
}

func (properties *propertyList) Set(tag string, value any) bool {
	return properties.set(strings.ToLower(tag), value)
}

func (properties *propertyList) set(tag string, value any) bool {
	if value == nil {
		delete(properties.properties, tag)
		return true
	}

	if _, ok := sizeProperties[tag]; ok {
		return properties.setSizeProperty(tag, value)
	}

	if valuesData, ok := enumProperties[tag]; ok {
		return properties.setEnumProperty(tag, value, valuesData.values)
	}

	if limits, ok := floatProperties[tag]; ok {
		return properties.setFloatProperty(tag, value, limits.min, limits.max)
	}

	if isPropertyInList(tag, colorProperties) {
		return properties.setColorProperty(tag, value)
	}

	if isPropertyInList(tag, angleProperties) {
		return properties.setAngleProperty(tag, value)
	}

	if isPropertyInList(tag, boolProperties) {
		return properties.setBoolProperty(tag, value)
	}

	if isPropertyInList(tag, intProperties) {
		return properties.setIntProperty(tag, value)
	}

	if text, ok := value.(string); ok {
		properties.properties[tag] = text
		return true
	}

	notCompatibleType(tag, value)
	return false
}
