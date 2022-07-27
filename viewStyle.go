package rui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ViewStyle interface of the style of view
type ViewStyle interface {
	Properties
	cssViewStyle(buffer cssBuilder, session Session)
}

type viewStyle struct {
	propertyList
	transitions map[string]Animation
}

// Range defines range limits. The First and Last value are included in the range
type Range struct {
	First, Last int
}

type stringWriter interface {
	writeString(buffer *strings.Builder, indent string)
}

// String returns a string representation of the Range struct
func (r Range) String() string {
	if r.First == r.Last {
		return fmt.Sprintf("%d", r.First)
	}
	return fmt.Sprintf("%d:%d", r.First, r.Last)
}

func (r *Range) setValue(value string) bool {
	var err error
	if strings.Contains(value, ":") {
		values := strings.Split(value, ":")
		if len(values) != 2 {
			ErrorLog("Invalid range value: " + value)
			return false
		}
		if r.First, err = strconv.Atoi(strings.Trim(values[0], " \t\n\r")); err != nil {
			ErrorLog(`Invalid first range value "` + value + `" (` + err.Error() + ")")
			return false
		}
		if r.Last, err = strconv.Atoi(strings.Trim(values[1], " \t\n\r")); err != nil {
			ErrorLog(`Invalid last range value "` + value + `" (` + err.Error() + ")")
			return false
		}
		return true
	}

	if r.First, err = strconv.Atoi(value); err != nil {
		ErrorLog(`Invalid range value "` + value + `" (` + err.Error() + ")")
		return false
	}
	r.Last = r.First
	return true
}

func (style *viewStyle) init() {
	style.propertyList.init()
	//style.shadows = []ViewShadow{}
	style.transitions = map[string]Animation{}
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

func (style *viewStyle) cssTextDecoration(session Session) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)

	noDecoration := false
	if strikethrough, ok := boolProperty(style, Strikethrough, session); ok {
		if strikethrough {
			buffer.WriteString("line-through")
		}
		noDecoration = true
	}

	if overline, ok := boolProperty(style, Overline, session); ok {
		if overline {
			if buffer.Len() > 0 {
				buffer.WriteRune(' ')
			}
			buffer.WriteString("overline")
		}
		noDecoration = true
	}

	if underline, ok := boolProperty(style, Underline, session); ok {
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

func (style *viewStyle) backgroundCSS(session Session) string {
	if value, ok := style.properties[Background]; ok {
		if backgrounds, ok := value.([]BackgroundElement); ok {
			buffer := allocStringBuilder()
			defer freeStringBuilder(buffer)

			for _, background := range backgrounds {
				if value := background.cssStyle(session); value != "" {
					if buffer.Len() > 0 {
						buffer.WriteString(", ")
					}
					buffer.WriteString(value)
				}
			}

			if buffer.Len() > 0 {
				return buffer.String()
			}
		}
	}
	return ""
}

func (style *viewStyle) cssViewStyle(builder cssBuilder, session Session) {

	if margin, ok := boundsProperty(style, Margin, session); ok {
		margin.cssValue(Margin, builder)
	}

	if padding, ok := boundsProperty(style, Padding, session); ok {
		padding.cssValue(Padding, builder)
	}

	if border := getBorder(style, Border); border != nil {
		border.cssStyle(builder, session)
		border.cssWidth(builder, session)
		border.cssColor(builder, session)
	}

	radius := getRadius(style, session)
	radius.cssValue(builder)

	if outline := getOutline(style); outline != nil {
		outline.ViewOutline(session).cssValue(builder)
	}

	if z, ok := intProperty(style, ZIndex, session, 0); ok {
		builder.add(ZIndex, strconv.Itoa(z))
	}

	if opacity, ok := floatProperty(style, Opacity, session, 1.0); ok && opacity >= 0 && opacity <= 1 {
		builder.add(Opacity, strconv.FormatFloat(opacity, 'f', 3, 32))
	}

	if n, ok := intProperty(style, ColumnCount, session, 0); ok && n > 0 {
		builder.add(ColumnCount, strconv.Itoa(n))
	}

	for _, tag := range []string{
		Width, Height, MinWidth, MinHeight, MaxWidth, MaxHeight, Left, Right, Top, Bottom,
		TextSize, TextIndent, LetterSpacing, WordSpacing, LineHeight, TextLineThickness,
		GridRowGap, GridColumnGap, ColumnGap, ColumnWidth} {

		if size, ok := sizeProperty(style, tag, session); ok && size.Type != Auto {
			cssTag, ok := sizeProperties[tag]
			if !ok {
				cssTag = tag
			}
			builder.add(cssTag, size.cssString(""))
		}
	}

	colorProperties := []struct{ property, cssTag string }{
		{BackgroundColor, BackgroundColor},
		{TextColor, "color"},
		{TextLineColor, "text-decoration-color"},
		{CaretColor, CaretColor},
	}
	for _, p := range colorProperties {
		if color, ok := colorProperty(style, p.property, session); ok && color != 0 {
			builder.add(p.cssTag, color.cssString())
		}
	}

	if value, ok := enumProperty(style, BackgroundClip, session, 0); ok {
		builder.add(BackgroundClip, enumProperties[BackgroundClip].values[value])
	}

	if background := style.backgroundCSS(session); background != "" {
		builder.add("background", background)
	}

	if font, ok := stringProperty(style, FontName, session); ok && font != "" {
		builder.add(`font-family`, font)
	}

	writingMode := 0
	for _, tag := range []string{
		Overflow, TextAlign, TextTransform, TextWeight, TextLineStyle, WritingMode, TextDirection,
		VerticalTextOrientation, CellVerticalAlign, CellHorizontalAlign, GridAutoFlow, Cursor,
		WhiteSpace, WordBreak, TextOverflow, Float, TableVerticalAlign, Resize} {

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

	for _, prop := range []struct{ tag, cssTag, off, on string }{
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

	if text := style.cssTextDecoration(session); text != "" {
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
		builder.add("grid-row-start", strconv.Itoa(r.First+1))
		builder.add("grid-row-end", strconv.Itoa(r.Last+2))
	}
	if r, ok := rangeProperty(style, Column, session); ok {
		builder.add("grid-column-start", strconv.Itoa(r.First+1))
		builder.add("grid-column-end", strconv.Itoa(r.Last+2))
	}
	if text := style.gridCellSizesCSS(CellWidth, session); text != "" {
		builder.add(`grid-template-columns`, text)
	}
	if text := style.gridCellSizesCSS(CellHeight, session); text != "" {
		builder.add(`grid-template-rows`, text)
	}

	style.writeViewTransformCSS(builder, session)

	if clip := getClipShape(style, Clip, session); clip != nil && clip.valid(session) {
		builder.add(`clip-path`, clip.cssStyle(session))
	}

	if clip := getClipShape(style, ShapeOutside, session); clip != nil && clip.valid(session) {
		builder.add(`shape-outside`, clip.cssStyle(session))
	}

	if value := style.getRaw(Filter); value != nil {
		if filter, ok := value.(ViewFilter); ok {
			if text := filter.cssStyle(session); text != "" {
				builder.add(Filter, text)
			}
		}
	}

	if value := style.getRaw(BackdropFilter); value != nil {
		if filter, ok := value.(ViewFilter); ok {
			if text := filter.cssStyle(session); text != "" {
				builder.add(`-webkit-backdrop-filter`, text)
				builder.add(BackdropFilter, text)
			}
		}
	}

	if transition := style.transitionCSS(session); transition != "" {
		builder.add(`transition`, transition)
	}

	if animation := style.animationCSS(session); animation != "" {
		builder.add(AnimationTag, animation)
	}

	if pause, ok := boolProperty(style, AnimationPaused, session); ok {
		if pause {
			builder.add(`animation-play-state`, `paused`)
		} else {
			builder.add(`animation-play-state`, `running`)
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

func (style *viewStyle) Get(tag string) any {
	return style.get(strings.ToLower(tag))
}

func (style *viewStyle) get(tag string) any {
	switch tag {
	case Border, CellBorder:
		return getBorder(&style.propertyList, tag)

	case BorderLeft, BorderRight, BorderTop, BorderBottom,
		BorderStyle, BorderLeftStyle, BorderRightStyle, BorderTopStyle, BorderBottomStyle,
		BorderColor, BorderLeftColor, BorderRightColor, BorderTopColor, BorderBottomColor,
		BorderWidth, BorderLeftWidth, BorderRightWidth, BorderTopWidth, BorderBottomWidth:
		if border := getBorder(style, Border); border != nil {
			return border.Get(tag)
		}
		return nil

	case CellBorderLeft, CellBorderRight, CellBorderTop, CellBorderBottom,
		CellBorderStyle, CellBorderLeftStyle, CellBorderRightStyle, CellBorderTopStyle, CellBorderBottomStyle,
		CellBorderColor, CellBorderLeftColor, CellBorderRightColor, CellBorderTopColor, CellBorderBottomColor,
		CellBorderWidth, CellBorderLeftWidth, CellBorderRightWidth, CellBorderTopWidth, CellBorderBottomWidth:
		if border := getBorder(style, CellBorder); border != nil {
			return border.Get(tag)
		}
		return nil

	case RadiusX, RadiusY, RadiusTopLeft, RadiusTopLeftX, RadiusTopLeftY,
		RadiusTopRight, RadiusTopRightX, RadiusTopRightY,
		RadiusBottomLeft, RadiusBottomLeftX, RadiusBottomLeftY,
		RadiusBottomRight, RadiusBottomRightX, RadiusBottomRightY:
		return getRadiusElement(style, tag)

	case ColumnSeparator:
		if val, ok := style.properties[ColumnSeparator]; ok {
			return val.(ColumnSeparatorProperty)
		}
		return nil

	case ColumnSeparatorStyle, ColumnSeparatorWidth, ColumnSeparatorColor:
		if val, ok := style.properties[ColumnSeparator]; ok {
			separator := val.(ColumnSeparatorProperty)
			return separator.Get(tag)
		}
		return nil

	case Transition:
		if len(style.transitions) == 0 {
			return nil
		}
		result := map[string]Animation{}
		for tag, animation := range style.transitions {
			result[tag] = animation
		}
		return result
	}

	return style.propertyList.getRaw(tag)
}

func (style *viewStyle) AllTags() []string {
	result := style.propertyList.AllTags()
	if len(style.transitions) > 0 {
		result = append(result, Transition)
	}
	return result
}

func supportedPropertyValue(value any) bool {
	switch value.(type) {
	case string:
	case []string:
	case bool:
	case float32:
	case float64:
	case int:
	case stringWriter:
	case fmt.Stringer:
	case []ViewShadow:
	case []View:
	case []any:
	case map[string]Animation:
	default:
		return false
	}
	return true
}

func writePropertyValue(buffer *strings.Builder, tag string, value any, indent string) {

	writeString := func(text string) {
		simple := (tag != Text && tag != Title && tag != Summary)
		if simple {
			if len(text) == 1 {
				simple = (text[0] >= '0' && text[0] <= '9') || (text[0] >= 'A' && text[0] <= 'Z') || (text[0] >= 'a' && text[0] <= 'z')
			} else {
				for _, ch := range text {
					if (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
						ch == '+' || ch == '-' || ch == '@' || ch == '/' || ch == '_' || ch == ':' {
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
			}
			buffer.WriteString(indent)
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
		buffer.WriteString(value.String())

	case []ViewShadow:
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
			}
			buffer.WriteRune('\n')
			buffer.WriteString(indent)
			buffer.WriteRune(']')
		}

	case []View:
		switch len(value) {
		case 0:
			buffer.WriteString("[]\n")

		case 1:
			writeViewStyle(value[0].Tag(), value[0], buffer, indent)

		default:
			buffer.WriteString("[\n")
			indent2 := indent + "\t"
			for _, v := range value {
				buffer.WriteString(indent2)
				writeViewStyle(v.Tag(), v, buffer, indent2)
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

	case map[string]Animation:
		switch count := len(value); count {
		case 0:
			buffer.WriteString("[]")

		case 1:
			for tag, animation := range value {
				animation.writeTransitionString(tag, buffer)
				break
			}

		default:
			tags := make([]string, 0, len(value))
			for tag := range value {
				tags = append(tags, tag)
			}
			sort.Strings(tags)
			buffer.WriteString("[\n")
			indent2 := indent + "\t"
			for _, tag := range tags {
				if animation := value[tag]; animation != nil {
					buffer.WriteString(indent2)
					animation.writeTransitionString(tag, buffer)
					buffer.WriteString("\n")
				}
			}
			buffer.WriteString(indent)
			buffer.WriteRune(']')
		}
	}
}

func writeViewStyle(name string, view ViewStyle, buffer *strings.Builder, indent string) {
	buffer.WriteString(name)
	buffer.WriteString(" {\n")
	indent += "\t"

	writeProperty := func(tag string, value any) {
		if supportedPropertyValue(value) {
			buffer.WriteString(indent)
			buffer.WriteString(tag)
			buffer.WriteString(" = ")
			writePropertyValue(buffer, tag, value, indent)
			buffer.WriteString(",\n")
		}
	}

	tags := view.AllTags()
	removeTag := func(tag string) {
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

	tagOrder := []string{
		ID, Row, Column, Top, Right, Bottom, Left, Semantics, Cursor, Visibility,
		Opacity, ZIndex, Width, Height, MinWidth, MinHeight, MaxWidth, MaxHeight,
		Margin, Padding, BackgroundClip, BackgroundColor, Background, Border, Radius, Outline, Shadow,
		Orientation, ListWrap, VerticalAlign, HorizontalAlign, CellWidth, CellHeight,
		CellVerticalAlign, CellHorizontalAlign, GridRowGap, GridColumnGap,
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

	finalTags := []string{
		Perspective, PerspectiveOriginX, PerspectiveOriginY, BackfaceVisible, OriginX, OriginY, OriginZ,
		TranslateX, TranslateY, TranslateZ, ScaleX, ScaleY, ScaleZ, Rotate, RotateX, RotateY, RotateZ,
		SkewX, SkewY, Clip, Filter, BackdropFilter, Summary, Content, Transition}
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

func getViewString(view View) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	writeViewStyle(view.Tag(), view, buffer, "")
	return buffer.String()

}

func runStringWriter(writer stringWriter) string {
	buffer := allocStringBuilder()
	defer freeStringBuilder(buffer)
	writer.writeString(buffer, "")
	return buffer.String()
}
