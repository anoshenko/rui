package rui

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

// ViewStyle interface of the style of view
type ViewStyle interface {
	Properties

	// Transition returns the transition animation of the property. Returns nil is there is no transition animation.
	Transition(tag PropertyName) AnimationProperty

	// Transitions returns the map of transition animations. The result is always non-nil.
	Transitions() map[PropertyName]AnimationProperty

	// SetTransition sets the transition animation for the property if "animation" argument is not nil, and
	// removes the transition animation of the property if "animation" argument  is nil.
	// The "tag" argument is the property name.
	SetTransition(tag PropertyName, animation AnimationProperty)

	cssViewStyle(buffer cssBuilder, session Session)
}

type viewStyle struct {
	propertyList
	//transitions map[PropertyName]Animation
}

type stringWriter interface {
	writeString(buffer *strings.Builder, indent string)
}

func (style *viewStyle) init() {
	style.propertyList.init()
	style.normalize = normalizeViewStyleTag
}

// NewViewStyle create new ViewStyle object
func NewViewStyle(params Params) ViewStyle {
	style := new(viewStyle)
	style.init()
	for tag, value := range params {
		style.Set(tag, value)
	}
	return style
}

func textDecorationCSS(properties Properties, session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	noDecoration := false
	if strikethrough, ok := boolProperty(properties, Strikethrough, session); ok {
		if strikethrough {
			buffer.WriteString("line-through")
		}
		noDecoration = true
	}

	if overline, ok := boolProperty(properties, Overline, session); ok {
		if overline {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString("overline")
		}
		noDecoration = true
	}

	if underline, ok := boolProperty(properties, Underline, session); ok {
		if underline {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString("underline")
		}
		noDecoration = true
	}

	if buffer.Len() == 0 && noDecoration {
		return "none"
	}

	return buffer.String()
}

func split4Values(text string) []string {
	values := strings.Split(text, ",")
	count := len(values)
	switch count {
	case 1, 4:
		return values

	case 2:
		if strings.Trim(values[1], " \t\r\n") == "" {
			return values[:1]
		}

	case 5:
		if strings.Trim(values[4], " \t\r\n") != "" {
			return values[:4]
		}
	}
	return []string{}
}

func (style *viewStyle) cssViewStyle(builder cssBuilder, session Session) {

	if visibility, ok := enumProperty(style, Visibility, session, Visible); ok {
		switch visibility {
		case Invisible:
			builder.add(`visibility`, `hidden`)

		case Gone:
			builder.add(`display`, `none`)
		}
	}

	if margin, ok := getBounds(style, Margin, session); ok {
		margin.cssValue(Margin, builder, session)
	}

	if padding, ok := getBounds(style, Padding, session); ok {
		padding.cssValue(Padding, builder, session)
	}

	if border := getBorderProperty(style, Border); border != nil {
		border.cssStyle(builder, session)
		border.cssWidth(builder, session)
		border.cssColor(builder, session)
	}

	radius := getRadius(style, session)
	radius.cssValue(builder, session)

	if outline := getOutlineProperty(style); outline != nil {
		outline.ViewOutline(session).cssValue(builder, session)
	}

	for _, tag := range []PropertyName{ZIndex, Order} {
		if value, ok := intProperty(style, tag, session, 0); ok {
			builder.add(string(tag), strconv.Itoa(value))
		}
	}

	if opacity, ok := floatProperty(style, Opacity, session, 1.0); ok && opacity >= 0 && opacity <= 1 {
		builder.add(string(Opacity), strconv.FormatFloat(opacity, 'f', 3, 32))
	}

	for _, tag := range []PropertyName{ColumnCount, TabSize} {
		if value, ok := intProperty(style, tag, session, 0); ok && value > 0 {
			builder.add(string(tag), strconv.Itoa(value))
		}
	}

	for _, tag := range []PropertyName{
		Width, Height, MinWidth, MinHeight, MaxWidth, MaxHeight, Left, Right, Top, Bottom,
		TextSize, TextIndent, LetterSpacing, WordSpacing, LineHeight, TextLineThickness,
		ListRowGap, ListColumnGap, GridRowGap, GridColumnGap, ColumnGap, ColumnWidth, OutlineOffset} {

		if size, ok := sizeProperty(style, tag, session); ok && size.Type != Auto {
			cssTag, ok := sizeProperties[tag]
			if !ok {
				cssTag = string(tag)
			}
			builder.add(cssTag, size.cssString("", session))
		}
	}

	type propertyCss struct {
		property PropertyName
		cssTag   string
	}
	colorProperties := []propertyCss{
		//{BackgroundColor, string(BackgroundColor)},
		{TextColor, "color"},
		{TextLineColor, "text-decoration-color"},
		{CaretColor, string(CaretColor)},
		{AccentColor, string(AccentColor)},
	}
	for _, p := range colorProperties {
		if color, ok := colorProperty(style, p.property, session); ok && color != 0 {
			builder.add(p.cssTag, color.cssString())
		}
	}

	for _, tag := range []PropertyName{BackgroundClip, BackgroundOrigin, MaskClip, MaskOrigin} {
		if value, ok := enumProperty(style, tag, session, 0); ok {
			if data, ok := enumProperties[tag]; ok {
				builder.add(data.cssTag, data.cssValues[value])
			}
		}
	}

	if background := backgroundCSS(style, session); background != "" {
		builder.add("background", background)
	} else {
		backgroundColor, _ := colorProperty(style, BackgroundColor, session)
		if backgroundColor != 0 {
			builder.add("background-color", backgroundColor.cssString())
		}
	}

	if mask := maskCSS(style, session); mask != "" {
		builder.add("mask", mask)
	}

	if font, ok := stringProperty(style, FontName, session); ok && font != "" {
		builder.add(`font-family`, font)
	}

	writingMode := 0
	for _, tag := range []PropertyName{
		Overflow, TextAlign, TextTransform, TextWeight, TextLineStyle, WritingMode, TextDirection,
		VerticalTextOrientation, CellVerticalAlign, CellHorizontalAlign, GridAutoFlow, Cursor,
		WhiteSpace, WordBreak, TextOverflow, Float, TableVerticalAlign, Resize, MixBlendMode, BackgroundBlendMode} {

		if data, ok := enumProperties[tag]; ok {
			if tag != VerticalTextOrientation || (writingMode != VerticalLeftToRight && writingMode != VerticalRightToLeft) {
				if value, ok := enumProperty(style, tag, session, 0); ok {
					cssValue := data.values[value]
					if cssValue != "" {
						builder.add(data.cssTag, cssValue)
					}

					if tag == WritingMode {
						writingMode = value
					}
				}
			}
		}
	}

	type boolPropertyCss struct {
		tag             PropertyName
		cssTag, off, on string
	}
	for _, prop := range []boolPropertyCss{
		{tag: Italic, cssTag: "font-style", off: "normal", on: "italic"},
		{tag: SmallCaps, cssTag: "font-variant", off: "normal", on: "small-caps"},
	} {
		if flag, ok := boolProperty(style, prop.tag, session); ok {
			if flag {
				builder.add(prop.cssTag, prop.on)
			} else {
				builder.add(prop.cssTag, prop.off)
			}
		}
	}

	if text := textDecorationCSS(style, session); text != "" {
		builder.add("text-decoration", text)
	}

	if userSelect, ok := boolProperty(style, UserSelect, session); ok {
		if userSelect {
			builder.add("-webkit-user-select", "auto")
			builder.add("user-select", "auto")
		} else {
			builder.add("-webkit-user-select", "none")
			builder.add("user-select", "none")
		}
	}

	if css := shadowCSS(style, Shadow, session); css != "" {
		builder.add("box-shadow", css)
	}

	if css := shadowCSS(style, TextShadow, session); css != "" {
		builder.add("text-shadow", css)
	}

	if value, ok := style.properties[ColumnSeparator]; ok {
		if separator, ok := value.(ColumnSeparatorProperty); ok {
			if css := separator.cssValue(session); css != "" {
				builder.add("column-rule", css)
			}
		}
	}

	if avoid, ok := boolProperty(style, AvoidBreak, session); ok {
		if avoid {
			builder.add("break-inside", "avoid")
		} else {
			builder.add("break-inside", "auto")
		}
	}

	wrap, _ := enumProperty(style, ListWrap, session, 0)
	orientation, ok := valueToOrientation(style.Get(Orientation), session)
	if ok || wrap > 0 {
		cssText := enumProperties[Orientation].cssValues[orientation]
		switch wrap {
		case ListWrapOn:
			cssText += " wrap"

		case ListWrapReverse:
			cssText += " wrap-reverse"
		}
		builder.add(`flex-flow`, cssText)
	}

	rows := (orientation == StartToEndOrientation || orientation == EndToStartOrientation)

	var hAlignTag, vAlignTag string
	if rows {
		hAlignTag = `justify-content`
		vAlignTag = `align-items`
	} else {
		hAlignTag = `align-items`
		vAlignTag = `justify-content`
	}

	if align, ok := enumProperty(style, HorizontalAlign, session, LeftAlign); ok {
		switch align {
		case LeftAlign:
			if (!rows && wrap == ListWrapReverse) || orientation == EndToStartOrientation {
				builder.add(hAlignTag, `flex-end`)
			} else {
				builder.add(hAlignTag, `flex-start`)
			}
		case RightAlign:
			if (!rows && wrap == ListWrapReverse) || orientation == EndToStartOrientation {
				builder.add(hAlignTag, `flex-start`)
			} else {
				builder.add(hAlignTag, `flex-end`)
			}
		case CenterAlign:
			builder.add(hAlignTag, `center`)

		case StretchAlign:
			if rows {
				builder.add(hAlignTag, `space-between`)
			} else {
				builder.add(hAlignTag, `stretch`)
			}
		}
	}

	if align, ok := enumProperty(style, VerticalAlign, session, LeftAlign); ok {
		switch align {
		case TopAlign:
			if (rows && wrap == ListWrapReverse) || orientation == BottomUpOrientation {
				builder.add(vAlignTag, `flex-end`)
			} else {
				builder.add(vAlignTag, `flex-start`)
			}
		case BottomAlign:
			if (rows && wrap == ListWrapReverse) || orientation == BottomUpOrientation {
				builder.add(vAlignTag, `flex-start`)
			} else {
				builder.add(vAlignTag, `flex-end`)
			}
		case CenterAlign:
			builder.add(vAlignTag, `center`)

		case StretchAlign:
			if rows {
				builder.add(vAlignTag, `stretch`)
			} else {
				builder.add(vAlignTag, `space-between`)
			}
		}
	}

	if r, ok := rangeProperty(style, Row, session); ok {
		builder.add("grid-row", fmt.Sprintf("%d / %d", r.First+1, r.Last+2))
	}
	if r, ok := rangeProperty(style, Column, session); ok {
		builder.add("grid-column", fmt.Sprintf("%d / %d", r.First+1, r.Last+2))
	}
	if text := gridCellSizesCSS(style, CellWidth, session); text != "" {
		builder.add(`grid-template-columns`, text)
	}
	if text := gridCellSizesCSS(style, CellHeight, session); text != "" {
		builder.add(`grid-template-rows`, text)
	}

	style.writeViewTransformCSS(builder, session)

	if clip := getClipShapeProperty(style, Clip, session); clip != nil && clip.valid(session) {
		builder.add(`clip-path`, clip.cssStyle(session))
	}

	if clip := getClipShapeProperty(style, ShapeOutside, session); clip != nil && clip.valid(session) {
		builder.add(`shape-outside`, clip.cssStyle(session))
	}

	if value := style.getRaw(Filter); value != nil {
		if filter, ok := value.(FilterProperty); ok {
			if text := filter.cssStyle(session); text != "" {
				builder.add(string(Filter), text)
			}
		}
	}

	if value := style.getRaw(BackdropFilter); value != nil {
		if filter, ok := value.(FilterProperty); ok {
			if text := filter.cssStyle(session); text != "" {
				builder.add(`-webkit-backdrop-filter`, text)
				builder.add(string(BackdropFilter), text)
			}
		}
	}

	if transition := transitionCSS(style, session); transition != "" {
		builder.add(`transition`, transition)
	}

	if animation := animationCSS(style, session); animation != "" {
		builder.add(string(Animation), animation)
	}

	if pause, ok := boolProperty(style, AnimationPaused, session); ok {
		if pause {
			builder.add(`animation-play-state`, `paused`)
		} else {
			builder.add(`animation-play-state`, `running`)
		}
	}

	if spanAll, ok := boolProperty(style, ColumnSpanAll, session); ok {
		if spanAll {
			builder.add(`column-span`, `all`)
		} else {
			builder.add(`column-span`, `none`)
		}
	}
}

func valueToOrientation(value any, session Session) (int, bool) {
	if value != nil {
		switch value := value.(type) {
		case int:
			return value, true

		case string:
			text, ok := session.resolveConstants(value)
			if !ok {
				return 0, false
			}

			text = strings.ToLower(strings.Trim(text, " \t\n\r"))
			switch text {
			case "vertical":
				return TopDownOrientation, true

			case "horizontal":
				return StartToEndOrientation, true
			}

			if result, ok := enumStringToInt(text, enumProperties[Orientation].values, true); ok {
				return result, true
			}
		}
	}
	return 0, false
}

func normalizeViewStyleTag(tag PropertyName) PropertyName {
	tag = defaultNormalize(tag)
	switch tag {
	case "top-margin":
		return MarginTop

	case "right-margin":
		return MarginRight

	case "bottom-margin":
		return MarginBottom

	case "left-margin":
		return MarginLeft

	case "top-padding":
		return PaddingTop

	case "right-padding":
		return PaddingRight

	case "bottom-padding":
		return PaddingBottom

	case "left-padding":
		return PaddingLeft

	case "origin-x":
		return TransformOriginX

	case "origin-y":
		return TransformOriginY

	case "origin-z":
		return TransformOriginZ

	}
	return tag
}

func (style *viewStyle) Get(tag PropertyName) any {
	return viewStyleGet(style, normalizeViewStyleTag(tag))
}

func viewStyleGet(style Properties, tag PropertyName) any {
	switch tag {

	case BorderLeft, BorderRight, BorderTop, BorderBottom,
		BorderStyle, BorderLeftStyle, BorderRightStyle, BorderTopStyle, BorderBottomStyle,
		BorderColor, BorderLeftColor, BorderRightColor, BorderTopColor, BorderBottomColor,
		BorderWidth, BorderLeftWidth, BorderRightWidth, BorderTopWidth, BorderBottomWidth:
		if border := getBorderProperty(style, Border); border != nil {
			return border.Get(tag)
		}
		return nil

	case CellBorderLeft, CellBorderRight, CellBorderTop, CellBorderBottom,
		CellBorderStyle, CellBorderLeftStyle, CellBorderRightStyle, CellBorderTopStyle, CellBorderBottomStyle,
		CellBorderColor, CellBorderLeftColor, CellBorderRightColor, CellBorderTopColor, CellBorderBottomColor,
		CellBorderWidth, CellBorderLeftWidth, CellBorderRightWidth, CellBorderTopWidth, CellBorderBottomWidth:
		if border := getBorderProperty(style, CellBorder); border != nil {
			return border.Get(tag)
		}
		return nil

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		return getRadiusElement(style, tag)

	case ColumnSeparatorStyle, ColumnSeparatorWidth, ColumnSeparatorColor:
		if val := style.getRaw(ColumnSeparator); val != nil {
			separator := val.(ColumnSeparatorProperty)
			return separator.Get(tag)
		}
		return nil

	case RotateX, RotateY, RotateZ, Rotate, SkewX, SkewY, ScaleX, ScaleY, ScaleZ,
		TranslateX, TranslateY, TranslateZ:
		if transform := getTransformProperty(style, Transform); transform != nil {
			return transform.Get(tag)
		}
		return nil
	}

	return style.getRaw(tag)
}

func supportedPropertyValue(value any) bool {
	switch value := value.(type) {
	case string, bool, float32, float64, int, stringWriter, fmt.Stringer:
		return true

	case []string:
		return len(value) > 0

	case []ShadowProperty:
		return len(value) > 0

	case []View:
		return len(value) > 0

	case []any:
		return len(value) > 0

	case []BackgroundElement:
		return len(value) > 0

	case []BackgroundGradientPoint:
		return len(value) > 0

	case []BackgroundGradientAngle:
		return len(value) > 0

	case map[PropertyName]AnimationProperty:
		return len(value) > 0

	case []noArgListener[View]:
		return getNoArgBinding(value) != ""

	case []noArgListener[ImageView]:
		return getNoArgBinding(value) != ""

	case []noArgListener[MediaPlayer]:
		return getNoArgBinding(value) != ""

	case []oneArgListener[View, KeyEvent]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, MouseEvent]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, TouchEvent]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, PointerEvent]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, PropertyName]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, string]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, Frame]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[View, DragAndDropEvent]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[Checkbox, bool]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[FilePicker, []FileInfo]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[ListView, int]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[ListView, []int]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[MediaPlayer, float64]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[TableView, int]:
		return getOneArgBinding(value) != ""

	case []oneArgListener[TabsLayout, int]:
		return getOneArgBinding(value) != ""

	case []twoArgListener[ColorPicker, Color]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[DatePicker, time.Time]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[TimePicker, time.Time]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[DropDownList, int]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[EditView, string]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[NumberPicker, float64]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[TableView, int]:
		return getTwoArgBinding(value) != ""

	case []twoArgListener[TabsLayout, int]:
		return getTwoArgBinding(value) != ""

	case []mediaPlayerErrorListener:
		return getMediaPlayerErrorListenerBinding(value) != ""

	default:
		return false
	}
}

func writePropertyValue(buffer *strings.Builder, tag PropertyName, value any, indent string) {

	writeString := func(text string) {
		simple := (tag != Text && tag != Title && tag != Summary)
		if simple {
			if len(text) == 1 {
				simple = (text[0] >= '0' && text[0] <= '9') || (text[0] >= 'A' && text[0] <= 'Z') || (text[0] >= 'a' && text[0] <= 'z')
			} else {
				for _, ch := range text {
					if (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
						ch == '+' || ch == '-' || ch == '@' || ch == '/' || ch == '_' || ch == '.' ||
						ch == ':' || ch == '#' || ch == '%' || ch == 'π' || ch == '°' {
					} else {
						simple = false
						break
					}
				}
			}
		}

		if !simple {
			replace := []struct{ old, new string }{
				{old: "\\", new: `\\`},
				{old: "\t", new: `\t`},
				{old: "\r", new: `\r`},
				{old: "\n", new: `\n`},
				{old: "\"", new: `\"`},
			}
			for _, s := range replace {
				text = strings.Replace(text, s.old, s.new, -1)
			}
			buffer.WriteRune('"')
			buffer.WriteString(text)
			buffer.WriteRune('"')
		} else {
			buffer.WriteString(text)
		}
	}

	switch value := value.(type) {
	case string:
		writeString(value)

	case []string:
		if len(value) == 0 {
			buffer.WriteString("[]")
		} else {
			size := 0
			for _, text := range value {
				size += len(text) + 2
			}

			if size < 80 {
				lead := "["
				for _, text := range value {
					buffer.WriteString(lead)
					writeString(text)
					lead = ", "
				}
			} else {
				buffer.WriteString("[\n")
				for _, text := range value {
					buffer.WriteString(indent)
					buffer.WriteRune('\t')
					writeString(text)
					buffer.WriteString(",\n")
				}
				buffer.WriteString(indent)
			}
			buffer.WriteRune(']')
		}

	case bool:
		if value {
			buffer.WriteString("true")
		} else {
			buffer.WriteString("false")
		}

	case float32:
		buffer.WriteString(fmt.Sprintf("%g", float64(value)))

	case float64:
		buffer.WriteString(fmt.Sprintf("%g", value))

	case int:
		if prop, ok := enumProperties[tag]; ok && value >= 0 && value < len(prop.values) {
			buffer.WriteString(prop.values[value])
		} else {
			buffer.WriteString(strconv.Itoa(value))
		}

	case stringWriter:
		value.writeString(buffer, indent+"\t")

	case fmt.Stringer:
		writeString(value.String())

	case []ShadowProperty:
		switch len(value) {
		case 0:
			// do nothing

		case 1:
			value[0].writeString(buffer, indent)

		default:
			buffer.WriteString("[")
			indent2 := "\n" + indent + "\t"
			for _, shadow := range value {
				buffer.WriteString(indent2)
				shadow.writeString(buffer, indent)
				buffer.WriteRune(',')
			}
			buffer.WriteRune('\n')
			buffer.WriteString(indent)
			buffer.WriteRune(']')
		}

	case []View:
		switch len(value) {
		case 0:
			buffer.WriteString("[]")

		case 1:
			writeViewStyle(value[0].Tag(), value[0], buffer, indent, value[0].exscludeTags())

		default:
			buffer.WriteString("[\n")
			indent2 := indent + "\t"
			for _, v := range value {
				buffer.WriteString(indent2)
				writeViewStyle(v.Tag(), v, buffer, indent2, v.exscludeTags())
				buffer.WriteString(",\n")
			}

			buffer.WriteString(indent)
			buffer.WriteRune(']')
		}

	case []any:
		switch count := len(value); count {
		case 0:
			buffer.WriteString("[]")

		case 1:
			writePropertyValue(buffer, tag, value[0], indent)

		default:
			buffer.WriteString("[ ")
			comma := false
			for _, v := range value {
				if comma {
					buffer.WriteString(", ")
				}
				writePropertyValue(buffer, tag, v, indent)
				comma = true
			}
			buffer.WriteString(" ]")
		}

	case []BackgroundElement:
		switch len(value) {
		case 0:
			buffer.WriteString("[]\n")

		case 1:
			value[0].writeString(buffer, indent)

		default:
			buffer.WriteString("[\n")
			indent2 := indent + "\t"
			for _, element := range value {
				buffer.WriteString(indent2)
				element.writeString(buffer, indent2)
				buffer.WriteString(",\n")
			}

			buffer.WriteString(indent)
			buffer.WriteRune(']')
		}

	case []BackgroundGradientPoint:
		buffer.WriteRune('"')
		for i, point := range value {
			if i > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(point.String())
		}
		buffer.WriteRune('"')

	case []BackgroundGradientAngle:
		buffer.WriteRune('"')
		for i, point := range value {
			if i > 0 {
				buffer.WriteString(",")
			}
			buffer.WriteString(point.String())
		}
		buffer.WriteRune('"')

	case map[PropertyName]AnimationProperty:
		switch count := len(value); count {
		case 0:
			buffer.WriteString("[]")

		case 1:
			for tag, animation := range value {
				animation.writeTransitionString(tag, buffer)
				break
			}

		default:
			tags := make([]PropertyName, 0, len(value))
			for tag := range value {
				tags = append(tags, tag)
			}
			slices.Sort(tags)
			buffer.WriteString("[\n")
			indent2 := indent + "\t"
			for _, tag := range tags {
				if animation := value[tag]; animation != nil {
					buffer.WriteString(indent2)
					animation.writeTransitionString(tag, buffer)
					buffer.WriteString(",\n")
				}
			}
			buffer.WriteString(indent)
			buffer.WriteRune(']')
		}

	case []noArgListener[View]:
		buffer.WriteString(getNoArgBinding(value))

	case []noArgListener[ImageView]:
		buffer.WriteString(getNoArgBinding(value))

	case []noArgListener[MediaPlayer]:
		buffer.WriteString(getNoArgBinding(value))

	case []oneArgListener[View, KeyEvent]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, MouseEvent]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, TouchEvent]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, PointerEvent]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, PropertyName]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, string]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, Frame]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[View, DragAndDropEvent]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[Checkbox, bool]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[FilePicker, []FileInfo]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[ListView, int]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[ListView, []int]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[MediaPlayer, float64]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[TableView, int]:
		buffer.WriteString(getOneArgBinding(value))

	case []oneArgListener[TabsLayout, int]:
		buffer.WriteString(getOneArgBinding(value))

	case []twoArgListener[ColorPicker, Color]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[DatePicker, time.Time]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[TimePicker, time.Time]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[DropDownList, int]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[EditView, string]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[NumberPicker, float64]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[TableView, int]:
		buffer.WriteString(getTwoArgBinding(value))

	case []twoArgListener[TabsLayout, int]:
		buffer.WriteString(getTwoArgBinding(value))

	case []mediaPlayerErrorListener:
		buffer.WriteString(getMediaPlayerErrorListenerBinding(value))
	}
}

func writeViewStyle(name string, view Properties, buffer *strings.Builder, indent string, excludeTags []PropertyName) {
	buffer.WriteString(name)
	buffer.WriteString(" {\n")
	indent += "\t"

	writeProperty := func(tag PropertyName, value any) {
		if !slices.Contains(excludeTags, tag) {
			if supportedPropertyValue(value) {
				buffer.WriteString(indent)
				buffer.WriteString(string(tag))
				buffer.WriteString(" = ")
				writePropertyValue(buffer, tag, value, indent)
				buffer.WriteString(",\n")
			}
		}
	}

	tags := view.AllTags()
	removeTag := func(tag PropertyName) {
		for i, t := range tags {
			if t == tag {
				if i == 0 {
					tags = tags[1:]
				} else if i == len(tags)-1 {
					tags = tags[:i]
				} else {
					tags = append(tags[:i], tags[i+1:]...)
				}
				return
			}
		}
	}

	tagOrder := []PropertyName{
		ID, Row, Column, Top, Right, Bottom, Left, Semantics, Cursor, Visibility,
		Opacity, ZIndex, Width, Height, MinWidth, MinHeight, MaxWidth, MaxHeight,
		Margin, Padding, BackgroundColor, Background, BackgroundClip, BackgroundOrigin,
		Mask, MaskClip, MaskOrigin, Border, Radius, Outline, Shadow,
		Orientation, ListWrap, VerticalAlign, HorizontalAlign, CellWidth, CellHeight,
		CellVerticalAlign, CellHorizontalAlign, ListRowGap, ListColumnGap, GridRowGap, GridColumnGap,
		ColumnCount, ColumnWidth, ColumnSeparator, ColumnGap, AvoidBreak,
		Current, Expanded, Side, ResizeBorderWidth, EditViewType, MaxLength, Hint, Text, EditWrap,
		TextOverflow, FontName, TextSize, TextColor, TextWeight, Italic, SmallCaps,
		Strikethrough, Overline, Underline, TextLineStyle, TextLineThickness,
		TextLineColor, TextTransform, TextAlign, WhiteSpace, WordBreak, TextShadow, TextIndent,
		LetterSpacing, WordSpacing, LineHeight, TextDirection, WritingMode, VerticalTextOrientation,
	}

	for _, tag := range tagOrder {
		if value := view.Get(tag); value != nil {
			removeTag(tag)
			writeProperty(tag, value)
		}
	}

	finalTags := []PropertyName{
		PerspectiveOriginX, PerspectiveOriginY, BackfaceVisible,
		TransformOriginX, TransformOriginY, TransformOriginZ,
		Transform, Clip, Filter, BackdropFilter, Summary, Content, Transition}
	for _, tag := range finalTags {
		removeTag(tag)
	}

	for _, tag := range tags {
		if value := view.Get(tag); value != nil {
			writeProperty(tag, value)
		}
	}

	for _, tag := range finalTags {
		if value := view.Get(tag); value != nil {
			writeProperty(tag, value)
		}
	}

	indent = indent[:len(indent)-1]
	buffer.WriteString(indent)
	buffer.WriteString("}")
}

func runStringWriter(writer stringWriter) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	writer.writeString(buffer, "")
	return buffer.String()
}
