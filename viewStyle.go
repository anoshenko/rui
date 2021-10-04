package rui

import (
	"fmt"
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
		TextAlign, TextTransform, TextWeight, TextLineStyle, WritingMode, TextDirection,
		VerticalTextOrientation, CellVerticalAlign, CellHorizontalAlign, Cursor, WhiteSpace,
		WordBreak, TextOverflow, Float, TableVerticalAlign} {

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

	wrap, _ := enumProperty(style, Wrap, session, 0)
	orientation, ok := getOrientation(style, session)
	if ok || wrap > 0 {
		cssText := enumProperties[Orientation].cssValues[orientation]
		switch wrap {
		case WrapOn:
			cssText += " wrap"

		case WrapReverse:
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
			if (!rows && wrap == WrapReverse) || orientation == EndToStartOrientation {
				builder.add(hAlignTag, `flex-end`)
			} else {
				builder.add(hAlignTag, `flex-start`)
			}
		case RightAlign:
			if (!rows && wrap == WrapReverse) || orientation == EndToStartOrientation {
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
			if (rows && wrap == WrapReverse) || orientation == BottomUpOrientation {
				builder.add(vAlignTag, `flex-end`)
			} else {
				builder.add(vAlignTag, `flex-start`)
			}
		case BottomAlign:
			if (rows && wrap == WrapReverse) || orientation == BottomUpOrientation {
				builder.add(vAlignTag, `flex-start`)
			} else {
				builder.add(vAlignTag, `flex-end`)
			}
		case CenterAlign:
			builder.add(vAlignTag, `center`)

		case StretchAlign:
			if rows {
				builder.add(hAlignTag, `stretch`)
			} else {
				builder.add(hAlignTag, `space-between`)
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
				builder.add(`filter`, text)
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
